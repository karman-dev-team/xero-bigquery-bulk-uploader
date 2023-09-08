package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"xero-bigquery-bulk-uploader/models"

	"cloud.google.com/go/bigquery"
)

func uploadInvoices(invoices []models.XeroInvoice, company string, lookup []models.RevenueLineCSV, jobIDLookup []models.JobIDCSV) error {
	bqInvoices, err := convertToBQInvoice(invoices, company, lookup, jobIDLookup)
	if err != nil {
		return err
	}
	batchSize := 1000
	batches := splitIntoBatches(bqInvoices, batchSize)
	err = uploadToBQ(batches)
	if err != nil {
		return err
	}
	return nil
}

func uploadToBQ(batches [][]models.BQInvoice) error {
	maxRetries := 10
	retryInterval := 5 * time.Second
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, "reporting-393509")
	if err != nil {
		log.Printf("Failed to create BigQuery client: %v", err)
		return err
	}
	defer client.Close()

	dataset := client.Dataset("internal_reporting")
	table := dataset.Table("xero_invoices")
	uploader := table.Uploader()

	for i, batch := range batches {
		var retryCount int
		for retryCount < maxRetries {
			err := uploader.Put(ctx, batch)
			if err == nil {
				fmt.Printf("Uploaded Batch %d...\n", i+1)
				break
			}
			retryCount++
			log.Printf("Failed to insert data for Batch %d: %v. Retrying", i+1, err)
			if retryCount < maxRetries {
				time.Sleep(retryInterval)
			}
		}
		if retryCount == maxRetries {
			fmt.Printf("Exceeded maximum retries for Batch %d, giving up.\n", i+1)
		}
	}
	return nil
}
