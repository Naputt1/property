package queue

import (
	"backend/internal/models"
	"backend/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"encoding/csv"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type MigrationHandler struct {
	svcs *services.Services
}

func NewMigrationHandler(svcs *services.Services) *MigrationHandler {
	return &MigrationHandler{
		svcs: svcs,
	}
}

func (h *MigrationHandler) HandleCSVMigrateTask(ctx context.Context, t *asynq.Task) error {
	var payload models.CSVConfigPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	jobID := payload.JobID
	slog.Info("Starting CSV migration task", "file", payload.FilePath, "job_id", jobID)

	// Update job status to RUNNING
	if jobID != "" {
		_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusRunning, "Processing CSV file")
	}

	// Open CSV file
	file, err := os.Open(payload.FilePath)
	if err != nil {
		slog.Error("Failed to open CSV file", "error", err)
		if jobID != "" {
			_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to open file: %v", err))
		}
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// UK Land Registry dataset doesn't have headers by default, but it's configurable
	if payload.HasHeader {
		if _, err := reader.Read(); err != nil {
			slog.Error("Failed to read header", "error", err)
			if jobID != "" {
				_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusFailed, fmt.Sprintf("Failed to read header: %v", err))
			}
			return err
		}
	}

	var batch []models.Property
	batchSize := 1000
	totalProcessed := 0

	for {
		// Respect context cancellation
		select {
		case <-ctx.Done():
			slog.Warn("Task cancelled by context")
			if len(batch) > 0 {
				_ = h.svcs.Property.CreateBatch(ctx, batch, len(batch)) 
			}
			if jobID != "" {
				_ = h.svcs.Job.UpdateJobStatus(context.Background(), jobID, models.JobStatusFailed, "Task cancelled")
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

		// Very basic mapping based on dataset.csv
		// Ensure record length is sufficient
		if len(record) < 16 {
			slog.Warn("Skipping malformed record", "record", record)
			continue
		}

		price, _ := strconv.ParseInt(record[1], 10, 64)
		dateOfTransfer, _ := time.Parse("2006-01-02 15:04", record[2])

		property := models.Property{
			ID:               uuid.New().String(), // Generate new UUID for primary key
			Price:            price,
			DateOfTransfer:   dateOfTransfer,
			Postcode:         record[3],
			PropertyType:     record[4],
			OldNew:           record[5],
			Duration:         record[6],
			PAON:             record[7],
			SAON:             record[8],
			Street:           record[9],
			Locality:         record[10],
			TownCity:         record[11],
			District:         record[12],
			County:           record[13],
			PPDCategoryType:  record[14],
			RecordStatus:     record[15],
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
			
			slog.Info("Processed batch", "total", totalProcessed)
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

	// Invalidate analytics cache as new data has been added
	_ = h.svcs.Analytics.ClearCache(ctx)

	// Update job status to SUCCESS
	if jobID != "" {
		_ = h.svcs.Job.UpdateJobStatus(ctx, jobID, models.JobStatusSuccess, fmt.Sprintf("Successfully processed %d records", totalProcessed))
	}

	slog.Info("CSV migration task completed successfully", "totalProcessed", totalProcessed, "job_id", jobID)
	return nil
}
