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
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

// cleanUUID removes curly braces and normalizes the TUI to a standard UUID string format
func cleanUUID(id string) string {
	return strings.ToLower(strings.Trim(id, "{}"))
}

func toTitleCase(s string) string {
	if s == "" {
		return ""
	}
	words := strings.Fields(strings.ToLower(s))
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
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
func getCountyByPostcode(postcodeOutward string) string {
	if postcodeOutward == "" {
		return ""
	}

	// Extract Postcode Area (first 1-2 letters)
	area := ""
	for _, char := range postcodeOutward {
		if (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
			area += string(char)
		} else {
			break
		}
	}
	area = strings.ToUpper(area)

	// Postcode Area to Ceremonial County mapping (Comprehensive Title Case)
	areaToCounty := map[string]string{
		"AB": "Aberdeenshire",
		"AL": "Hertfordshire",
		"B":  "West Midlands",
		"BA": "Somerset",
		"BB": "Lancashire",
		"BD": "West Yorkshire",
		"BH": "Dorset",
		"BL": "Greater Manchester",
		"BN": "East Sussex",
		"BR": "Greater London",
		"BS": "Bristol",
		"BT": "Northern Ireland",
		"CA": "Cumbria",
		"CB": "Cambridgeshire",
		"CF": "South Glamorgan",
		"CH": "Cheshire",
		"CM": "Essex",
		"CO": "Essex",
		"CR": "Greater London",
		"CT": "Kent",
		"CV": "Warwickshire",
		"CW": "Cheshire",
		"DA": "Kent",
		"DD": "Angus",
		"DE": "Derbyshire",
		"DG": "Dumfries and Galloway",
		"DH": "Durham",
		"DL": "North Yorkshire",
		"DN": "South Yorkshire",
		"DT": "Dorset",
		"DY": "West Midlands",
		"E":  "Greater London",
		"EC": "Greater London",
		"EH": "City of Edinburgh",
		"EN": "Greater London",
		"EX": "Devon",
		"FK": "Stirling",
		"FY": "Lancashire",
		"G":  "Glasgow",
		"GL": "Gloucestershire",
		"GU": "Surrey",
		"HA": "Greater London",
		"HD": "West Yorkshire",
		"HG": "North Yorkshire",
		"HP": "Buckinghamshire",
		"HR": "Herefordshire",
		"HS": "Western Isles",
		"HU": "East Riding of Yorkshire",
		"HX": "West Yorkshire",
		"IG": "Greater London",
		"IP": "Suffolk",
		"IV": "Highland",
		"KA": "Ayrshire",
		"KT": "Surrey",
		"KW": "Highland",
		"KY": "Fife",
		"L":  "Merseyside",
		"LA": "Lancashire",
		"LD": "Powys",
		"LE": "Leicestershire",
		"LL": "Clwyd",
		"LN": "Lincolnshire",
		"LS": "West Yorkshire",
		"LU": "Bedfordshire",
		"M":  "Greater Manchester",
		"ME": "Kent",
		"MK": "Buckinghamshire",
		"ML": "Lanarkshire",
		"N":  "Greater London",
		"NE": "Tyne and Wear",
		"NG": "Nottinghamshire",
		"NN": "Northamptonshire",
		"NP": "Gwent",
		"NR": "Norfolk",
		"NW": "Greater London",
		"OL": "Greater Manchester",
		"OX": "Oxfordshire",
		"PA": "Renfrewshire",
		"PE": "Cambridgeshire",
		"PH": "Perthshire",
		"PL": "Devon",
		"PO": "Hampshire",
		"PR": "Lancashire",
		"RG": "Berkshire",
		"RH": "West Sussex",
		"RM": "Greater London",
		"S":  "South Yorkshire",
		"SA": "Dyfed",
		"SE": "Greater London",
		"SG": "Hertfordshire",
		"SK": "Cheshire",
		"SL": "Berkshire",
		"SM": "Greater London",
		"SN": "Wiltshire",
		"SO": "Hampshire",
		"SP": "Wiltshire",
		"SR": "Tyne and Wear",
		"SS": "Essex",
		"ST": "Staffordshire",
		"SW": "Greater London",
		"SY": "Shropshire",
		"TA": "Somerset",
		"TD": "Roxburghshire",
		"TF": "Shropshire",
		"TN": "Kent",
		"TQ": "Devon",
		"TR": "Cornwall",
		"TS": "Durham",
		"TW": "Greater London",
		"UB": "Greater London",
		"W":  "Greater London",
		"WA": "Cheshire",
		"WC": "Greater London",
		"WD": "Hertfordshire",
		"WF": "West Yorkshire",
		"WN": "Greater Manchester",
		"WR": "Worcestershire",
		"WS": "Staffordshire",
		"WV": "West Midlands",
		"YO": "North Yorkshire",
		"ZE": "Shetland",
	}

	if county, ok := areaToCounty[area]; ok {
		return county
	}

	return ""
}

func splitPostcode(raw string) (outward, inward string) {
	raw = strings.ToUpper(strings.ReplaceAll(raw, " ", ""))
	if len(raw) < 3 {
		return raw, ""
	}
	// UK Postcodes have inward code as last 3 characters
	inwardStart := len(raw) - 3
	return raw[:inwardStart], raw[inwardStart:]
}

func interpretCounty(ppdCounty, district, town, postcodeOutward string) string {
	// 1. Try deriving from postcode (Source of Truth for ceremonial mapping)
	if derived := getCountyByPostcode(postcodeOutward); derived != "" {
		return derived
	}

	// 2. Fallback to existing logic if postcode fails
	c := toTitleCase(ppdCounty)
	d := toTitleCase(district)

	// Common PPD Unitary Authority to Ceremonial County mappings
	mappings := map[string]string{
		"Brighton And Hove":                   "East Sussex",
		"Bath And North East Somerset":        "Somerset",
		"North East Lincolnshire":             "Lincolnshire",
		"North Lincolnshire":                  "Lincolnshire",
		"Bournemouth, Christchurch And Poole": "Dorset",
		"Bournemouth":                         "Dorset",
		"Poole":                               "Dorset",
		"West Berkshire":                      "Berkshire",
		"Windsor And Maidenhead":              "Berkshire",
		"Wokingham":                           "Berkshire",
		"Bracknell Forest":                    "Berkshire",
		"Reading":                             "Berkshire",
		"Slough":                              "Berkshire",
		"Wokingham Borough":                   "Berkshire",
		"West Northamptonshire":               "Northamptonshire",
		"North Northamptonshire":              "Northamptonshire",
		"Buckinghamshire":                     "Buckinghamshire",
		"Milton Keynes":                       "Buckinghamshire",
		"Central Bedfordshire":                "Bedfordshire",
		"Bedford":                             "Bedfordshire",
		"Cheshire East":                       "Cheshire",
		"Cheshire West And Chester":           "Cheshire",
		"Halton":                              "Cheshire",
		"Warrington":                          "Cheshire",
		"Cumberland":                          "Cumbria",
		"Westmorland And Furness":             "Cumbria",
		"East Riding Of Yorkshire":            "East Riding of Yorkshire",
		"Hull":                                "East Riding of Yorkshire",
		"City Of Kingston Upon Hull":          "East Riding of Yorkshire",
		"York":                                "North Yorkshire",
		"North Yorkshire":                     "North Yorkshire",
		"Middlesbrough":                       "North Yorkshire",
		"Redcar And Cleveland":                "North Yorkshire",
		"Stockton-on-tees":                    "Durham",
		"Darlington":                          "Durham",
		"Hartlepool":                          "Durham",
		"Stoke-on-trent":                      "Staffordshire",
		"Telford And Wrekin":                  "Shropshire",
		"Herefordshire":                       "Herefordshire",
		"City Of Hereford":                    "Herefordshire",
		"Isle Of Wight":                       "Isle of Wight",
		"Medway":                              "Kent",
		"Southampton":                         "Hampshire",
		"Portsmouth":                          "Hampshire",
		"Plymouth":                            "Devon",
		"Torbay":                              "Devon",
		"Swindon":                             "Wiltshire",
		"Bristol":                             "Bristol",
		"City Of Bristol":                     "Bristol",
		"Leicester":                           "Leicestershire",
		"City Of Leicester":                   "Leicestershire",
		"Rutland":                             "Rutland",
		"Nottingham":                          "Nottinghamshire",
		"City Of Nottingham":                  "Nottinghamshire",
		"Derby":                               "Derbyshire",
		"City Of Derby":                       "Derbyshire",
		"Peterborough":                        "Cambridgeshire",
		"City Of Peterborough":                "Cambridgeshire",
		"Thurrock":                            "Essex",
		"Southend-on-sea":                     "Essex",
		"Greater Manchester":                  "Greater Manchester",
		"Merseyside":                          "Merseyside",
		"South Yorkshire":                     "South Yorkshire",
		"West Yorkshire":                      "West Yorkshire",
		"Tyne And Wear":                       "Tyne and Wear",
		"West Midlands":                       "West Midlands",
	}

	if ceremonial, ok := mappings[c]; ok {
		return ceremonial
	}

	// Fallback to district check
	if ceremonial, ok := mappings[d]; ok {
		return ceremonial
	}

	return c
}

func formatAddress(paon, saon, street, locality string) string {
	var parts []string
	for _, s := range []string{saon, paon, street, locality} {
		s = strings.TrimSpace(s)
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}

func mapRecordToProperty(record []string) (*models.Property, error) {
	if len(record) < 16 {
		return nil, fmt.Errorf("malformed record: not enough columns")
	}

	price, _ := strconv.ParseInt(record[1], 10, 64)
	dateOfTransfer, _ := time.Parse("2006-01-02 15:04", record[2])

	id, err := uuid.Parse(cleanUUID(record[0]))
	if err != nil {
		id = uuid.New()
	}

	var outward, inward string
	var propertyType, oldNew, duration, paon, saon, street, locality, townCity, district, county, ppdCat, recordStatus string

	if len(record) >= 17 {
		// Split postcode format (17+ columns)
		// 0:ID, 1:Price, 2:Date, 3:PostOut, 4:PostIn, 5:Type, 6:OldNew, 7:Dur, 8:PAON, 9:SAON, 10:Street, 11:Locality, 12:Town, 13:Dist, 14:County, 15:Cat, 16:Stat
		outward = record[3]
		inward = record[4]
		propertyType = record[5]
		oldNew = record[6]
		duration = record[7]
		paon = record[8]
		saon = record[9]
		street = record[10]
		locality = record[11]
		townCity = record[12]
		district = record[13]
		county = record[14]
		ppdCat = record[15]
		recordStatus = record[16]
	} else {
		// Standard 16-column format
		// 0:ID, 1:Price, 2:Date, 3:Postcode, 4:Type, 5:OldNew, 6:Dur, 7:PAON, 8:SAON, 9:Street, 10:Locality, 11:Town, 12:Dist, 13:County, 14:Cat, 15:Stat
		outward, inward = splitPostcode(record[3])
		propertyType = record[4]
		oldNew = record[5]
		duration = record[6]
		paon = record[7]
		saon = record[8]
		street = record[9]
		locality = record[10]
		townCity = record[11]
		district = record[12]
		county = record[13]
		ppdCat = record[14]
		recordStatus = record[15]
	}

	return &models.Property{
		ID:              id,
		Price:           price,
		DateOfTransfer:  dateOfTransfer,
		PostcodeOutward: outward,
		PostcodeInward:  inward,
		PropertyType:    propertyType,
		OldNew:          oldNew,
		Duration:        duration,
		Address:         formatAddress(paon, saon, street, locality),
		TownCity:        toTitleCase(townCity),
		District:        toTitleCase(district),
		County:          interpretCounty(county, district, townCity, outward),
		PPDCategoryType: ppdCat,
		RecordStatus:    recordStatus,
	}, nil
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

		// Use the new mapping function that handles both 16 and 17 column formats
		property, err := mapRecordToProperty(record)
		if err != nil {
			slog.Warn("Skipping invalid record", "error", err, "record", record)
			continue
		}

		batch = append(batch, *property)
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

	// Refresh materialized view for analytics via queue with short delay
	_ = h.svcs.Job.EnqueueAnalyticsRefresh(context.Background(), 5*time.Second)

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
