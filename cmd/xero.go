package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"xero-bigquery-bulk-uploader/models"

	"golang.org/x/oauth2"
)

func getAllInvoices(token *oauth2.Token, tenantID string) ([]models.XeroInvoice, error) {
	invoices := []models.XeroInvoice{}

	page := 1
	for {
		invoice := models.InvoiceBody{}
		invoicesBytes, err := getInvoices(token, page, tenantID)
		if err != nil {
			fmt.Println("Error getting invoices", err)
			log.Fatal(err)
		}
		err = json.Unmarshal(invoicesBytes, &invoice)
		if err != nil {
			fmt.Println("Error unmarshalling", err)
			log.Fatal(err)
		}
		invoices = append(invoices, invoice.Invoices...)
		if len(invoice.Invoices) < 100 {
			break
		}
		page++
	}

	return invoices, nil
}

func getInvoices(token *oauth2.Token, page int, tenantID string) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/Invoices", nil)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("where", `Status!="DELETED" AND Status!="VOIDED" AND Type=="ACCREC"`)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("xero-tenant-id", tenantID)
	req.Header.Add("Accept", "application/json")
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("body: %s\n", body)
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return body, nil
}
