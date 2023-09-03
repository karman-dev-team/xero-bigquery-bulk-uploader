package models

import "time"

type XeroCompany struct {
	ID      string
	Company string
}

type XeroInvoice struct {
	Type           string  `json:"Type"`
	InvoiceID      string  `json:"InvoiceID"`
	InvoiceNumber  string  `json:"InvoiceNumber"`
	Reference      string  `json:"Reference"`
	AmountDue      float32 `json:"AmountDue"`
	AmountPaid     float32 `json:"AmountPaid"`
	AmountCredited float32 `json:"AmountCredited"`
	SentToContact  bool    `json:"SentToContact"`
	CurrencyRate   float32 `json:"CurrencyRate"`
	IsDiscounted   bool    `json:"IsDiscounted"`
	HasAttachments bool    `json:"HasAttachments"`
	HasErrors      bool    `json:"HasErrors"`
	Contact        struct {
		ContactID           string `json:"ContactID"`
		ContactNumber       string `json:"ContactNumber"`
		Name                string `json:"Name"`
		Addresses           []any  `json:"Addresses"`
		Phones              []any  `json:"Phones"`
		ContactGroups       []any  `json:"ContactGroups"`
		ContactPersons      []any  `json:"ContactPersons"`
		HasValidationErrors bool   `json:"HasValidationErrors"`
	} `json:"Contact"`
	DateString      string `json:"DateString"`
	DueDateString   string `json:"DueDateString"`
	Status          string `json:"Status"`
	LineAmountTypes string `json:"LineAmountTypes"`
	LineItems       []struct {
		Description string  `json:"Description"`
		UnitAmount  float32 `json:"UnitAmount"`
		TaxType     string  `json:"TaxType"`
		TaxAmount   float32 `json:"TaxAmount"`
		LineAmount  float32 `json:"LineAmount"`
		AccountCode string  `json:"AccountCode"`
		Quantity    float32 `json:"Quantity"`
		LineItemID  string  `json:"LineItemID"`
		AccountID   string  `json:"AccountID"`
	} `json:"LineItems"`
	SubTotal     float32 `json:"SubTotal"`
	TotalTax     float32 `json:"TotalTax"`
	Total        float32 `json:"Total"`
	CurrencyCode string  `json:"CurrencyCode"`
}

type InvoiceBody struct {
	Invoices []XeroInvoice `json:"Invoices"`
}

type BQInvoice struct {
	InvoiceID   string    `bigquery:"invoice_id"`
	ContactID   string    `bigquery:"contact_id"`
	ContactName string    `bigquery:"contact_name"`
	InvoiceDate time.Time `bigquery:"invoice_date"`
	DueDate     time.Time `bigquery:"due_date"`
	TotalPreTax float32   `bigquery:"total_pre_tax"`
	TotalTax    float32   `bigquery:"total_tax"`
	Total       float32   `bigquery:"total"`
	Company     string    `bigquery:"company"`
	Status      string    `bigquery:"status"`
	Reference   string    `bigquery:"reference"`
	Type        string    `bigquery:"type"`
}
