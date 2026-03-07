package queue

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

// cleanUUID removes curly braces and normalizes the TUI to a standard UUID string format
func cleanUUID(id string) string {
	return strings.ToLower(strings.Trim(id, "{}"))
}

type MigrationHandler struct {
	svcs   *services.Services
	bucket repository.BucketService
}

func NewMigrationHandler(svcs *services.Services, bucket repository.BucketService) *MigrationHandler {
	return &MigrationHandler{
		svcs:   svcs,
		bucket: bucket,
	}
}

func (h *MigrationHandler) HandleCSVMigrateTask(ctx context.Context, t *asynq.Task) (err error) {
	defer func() {
		if err != nil {
			slog.Info("HandleCSVMigrateTask failed", "error", err.Error())
		}
	}()

	var payload models.CSVConfigPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	jobID := payload.JobID
	slog.Info("Starting CSV migration task", "bucket_key", payload.BucketKey, "job_id", jobID)

	// Fetch current progress to allow resumption
	processedCount := 0
	if jobID != "" {
		jobDetail, err := h.svcs.Job.GetJobByID(ctx, jobID)
		if err == nil && jobDetail != nil {
			if jobDetail.Status == models.JobStatusSuccess {
				slog.Info("Job already completed successfully", "job_id", jobID)
				return nil
			}
			processedCount = jobDetail.Total
			if processedCount > 0 {
				slog.Info("Resuming migration from row", "job_id", jobID, "start_row", processedCount)
			}
		}
	}

	// Update job status to RUNNING
	if jobID != "" {
		_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusRunning, "Processing CSV file from bucket")
	}

	var object io.ReadCloser
	var totalSize int64

	object, totalSize, err = h.bucket.GetObject(ctx, payload.BucketKey)
	if err != nil {
		slog.Error("Failed to get object from bucket", "error", err)
		if jobID != "" {
			_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to open bucket object: %v", err))
		}
		return err
	}

	defer object.Close()

	reader := csv.NewReader(object)

	// Skip header if needed
	if payload.HasHeader {
		if _, err := reader.Read(); err != nil {
			slog.Error("Failed to read header", "error", err)
			if jobID != "" {
				_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to read header: %v", err))
			}
			return err
		}
	}

	// Skip already processed rows
	if processedCount > 0 {
		slog.Info("Skipping already processed rows", "count", processedCount)
		for i := 0; i < processedCount; i++ {
			if _, err := reader.Read(); err != nil {
				if err == io.EOF {
					slog.Info("Reached EOF while skipping", "skipped", i)
					break
				}
				slog.Error("Error while skipping rows", "error", err, "at", i)
				break
			}
		}
	}

	var batch []models.Property
	batchSize := 1000
	totalProcessed := processedCount
	lastProgressUpdate := time.Now()
	// We can't easily track bytesRead accurately when skipping rows without a custom reader
	// but we can estimate or just start from 0 for the remainder.
	// For better accuracy, we'll just use row count for progress if we resumed.
	bytesRead := int64(0)

	for {
		// Respect context cancellation
		select {
		case <-ctx.Done():
			slog.Warn("Task cancelled by context")
			if len(batch) > 0 {
				_ = h.svcs.Property.CreateBatch(context.Background(), batch, len(batch))
			}
			if jobID != "" {
				_ = h.svcs.Job.UpdateJobStatus(context.Background(), jobID, models.JobStatusFailed, "Task cancelled or timed out")
			}
			return ctx.Err()
		default:
		}

		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Failed to read CSV record", "error", err)
			continue
		}

		// Estimate bytes read (very rough estimate based on record content)
		for _, s := range record {
			bytesRead += int64(len(s)) + 1 // +1 for comma
		}
		bytesRead += 1 // for newline

		// Very basic mapping based on dataset.csv
		if len(record) < 16 {
			slog.Warn("Skipping malformed record", "record", record)
			continue
		}

		price, _ := strconv.ParseInt(record[1], 10, 64)
		dateOfTransfer, _ := time.Parse("2006-01-02 15:04", record[2])

		property := models.Property{
			ID:              cleanUUID(record[0]), // Use TUI as ID for idempotency
			Price:           price,
			DateOfTransfer:  dateOfTransfer,
			Postcode:        record[3],
			PropertyType:    record[4],
			OldNew:          record[5],
			Duration:        record[6],
			PAON:            record[7],
			SAON:            record[8],
			Street:          record[9],
			Locality:        record[10],
			TownCity:        record[11],
			District:        record[12],
			County:          record[13],
			PPDCategoryType: record[14],
			RecordStatus:    record[15],
		}

		batch = append(batch, property)
		totalProcessed++

		if len(batch) >= batchSize {
			if err := h.svcs.Property.CreateBatch(ctx, batch, batchSize); err != nil {
				slog.Error("Failed to process batch", "error", err)
				if jobID != "" {
					_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to save batch: %v", err))
				}
				return err
			}
			batch = batch[:0] // clear batch

			// Update progress every 5 seconds
			if time.Since(lastProgressUpdate) >= 5*time.Second {
				progress := 0
				if totalSize > 0 {
					// Note: bytesRead is only accurate if we didn't resume.
					// If we resumed, we might want a better way to calculate progress.
					// For now, let's just use it as is.
					progress = int((float64(bytesRead) / float64(totalSize)) * 100)
					if progress > 100 {
						progress = 99 // Cap at 99 until finish
					}
				}

				if jobID != "" {
					_ = h.svcs.Job.UpdateJobProgress(ctx, jobID, progress, totalProcessed)
					h.svcs.Socket.Broadcast(gin.H{
						"type": config.WsMessageTypeJobUpdate,
						"data": gin.H{
							"id":       jobID,
							"status":   models.JobStatusRunning,
							"progress": progress,
							"total":    totalProcessed,
						},
					})
				}
				lastProgressUpdate = time.Now()
				slog.Info("Processed progress", "total", totalProcessed, "progress", progress)
			}
		}
	}

	// Final batch
	if len(batch) > 0 {
		if err := h.svcs.Property.CreateBatch(ctx, batch, len(batch)); err != nil {
			slog.Error("Failed to process final batch", "error", err)
			if jobID != "" {
				_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to save final batch: %v", err))
			}
			return err
		}
	}

	// Invalidate analytics cache
	_ = h.svcs.Analytics.ClearCache(context.Background())

	// Precompute analytics cache in background
	go func() {
		// Use a fresh context for background work
		_ = h.svcs.Analytics.PrecomputeCache(context.Background())
	}()

	// Update job status to SUCCESS
	if jobID != "" {
		_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusSuccess, fmt.Sprintf("Successfully processed %d records", totalProcessed))
		_ = h.svcs.Job.UpdateJobProgress(ctx, jobID, 100, totalProcessed)
		h.svcs.Socket.Broadcast(gin.H{
			"type": config.WsMessageTypeJobUpdate,
			"data": gin.H{
				"id":       jobID,
				"status":   models.JobStatusSuccess,
				"progress": 100,
				"total":    totalProcessed,
			},
		})
	}

	// Optional: Delete from bucket after successful migration
	// _ = h.bucket.RemoveObject(ctx, h.bucketName, payload.BucketKey, minio.RemoveObjectOptions{})

	slog.Info("CSV migration task completed successfully", "totalProcessed", totalProcessed, "job_id", jobID)
	return nil
}
