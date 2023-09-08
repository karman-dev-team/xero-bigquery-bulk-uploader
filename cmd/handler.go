package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"xero-bigquery-bulk-uploader/models"

	"github.com/gocarina/gocsv"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	authURL := oauth2Config.AuthCodeURL("")
	http.Redirect(w, r, authURL, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	lookup, err := loadCSV()
	if err != nil {
		log.Fatal(err)
		return
	}
	jobIDLookup, err := loadJobIdCSV()
	if err != nil {
		log.Fatal(err)
		return
	}
	tenantID := []models.XeroCompany{
		{ID: os.Getenv("CF_TENANT_ID"), Company: "CF"},
		{ID: os.Getenv("KD_TENANT_ID"), Company: "KD"}}

	for _, tenant := range tenantID {
		invoices, err := getAllInvoices(token, tenant.ID)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = uploadInvoices(invoices, tenant.Company, lookup, jobIDLookup)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func loadCSV() ([]models.RevenueLineCSV, error) {
	filenames := []string{"cf_revenue_line.csv", "kd_revenue_line.csv"}
	var lookup []models.RevenueLineCSV
	for _, filename := range filenames {
		var csvLookup []models.RevenueLineCSV
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Println("Error opening CSV:", err)
			return lookup, err
		}
		defer file.Close()

		if err := gocsv.UnmarshalFile(file, &csvLookup); err != nil {
			fmt.Println("Error unmarshaling CSV:", err)
			return lookup, err
		}
		lookup = append(lookup, csvLookup...)
	}
	return lookup, nil
}

func loadJobIdCSV() ([]models.JobIDCSV, error) {
	fileName := "job_ids.csv"
	var lookup []models.JobIDCSV
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Error opening CSV:", err)
		return lookup, err
	}
	defer file.Close()
	if err := gocsv.UnmarshalFile(file, &lookup); err != nil {
		fmt.Println("Error unmarshaling CSV:", err)
		return lookup, err
	}
	return lookup, nil
}
