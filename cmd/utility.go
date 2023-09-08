package main

import (
	"time"
	"xero-bigquery-bulk-uploader/models"
)

func splitIntoBatches(slice []models.BQInvoice, batchSize int) [][]models.BQInvoice {
	var batches [][]models.BQInvoice

	for batchSize < len(slice) {
		slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
	}
	batches = append(batches, slice)

	return batches
}

func convertToBQInvoice(invoices []models.XeroInvoice, company string, lookup []models.RevenueLineCSV, jobIDLookup []models.JobIDCSV) ([]models.BQInvoice, error) {
	layout := "2006-01-02T15:04:05"
	bqInvoices := []models.BQInvoice{}
	for _, invoice := range invoices {
		invoiceDate, err := time.Parse(layout, invoice.DateString)
		if err != nil {
			return bqInvoices, err
		}
		dueDate, err := time.Parse(layout, invoice.DueDateString)
		if err != nil {
			return bqInvoices, err
		}
		bqInvoice := models.BQInvoice{
			InvoiceID:   invoice.InvoiceID,
			ContactID:   invoice.Contact.ContactID,
			ContactName: invoice.Contact.Name,
			InvoiceDate: invoiceDate,
			DueDate:     dueDate,
			TotalPreTax: invoice.SubTotal,
			TotalTax:    invoice.TotalTax,
			Total:       invoice.Total,
			Company:     company,
			Status:      invoice.Status,
			Reference:   invoice.Reference,
			Type:        invoice.Type,
			Description: invoice.LineItems[0].Description,
			RevenueLine: findRevenueLineById(lookup, jobIDLookup, invoice.LineItems[0].AccountCode, invoice.Reference),
		}
		bqInvoices = append(bqInvoices, bqInvoice)
	}
	return bqInvoices, nil
}

func findRevenueLineById(lookup []models.RevenueLineCSV, jobIDLookup []models.JobIDCSV, accountCode string, reference string) string {
	revenueLine := ""
	for i := range lookup {
		if lookup[i].XeroRevenueCode == accountCode {
			revenueLine = lookup[i].HubspotRevenueLine
		}
	}
	if revenueLine == "" {
		for i := range jobIDLookup {
			if jobIDLookup[i].JobID == reference {
				revenueLine = jobIDLookup[i].RevenueLine
			}
		}
	}

	return revenueLine
}
