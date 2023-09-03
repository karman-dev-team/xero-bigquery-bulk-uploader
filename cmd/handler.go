package main

import (
	"log"
	"net/http"
	"os"
	"xero-bigquery-bulk-uploader/models"
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
	tenantID := []models.XeroCompany{
		{ID: os.Getenv("CF_TENANT_ID"), Company: "CF"},
		{ID: os.Getenv("KD_TENANT_ID"), Company: "KD"}}

	for _, tenant := range tenantID {
		invoices, err := getAllInvoices(token, tenant.ID)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = uploadInvoices(invoices, tenant.Company)
		if err != nil {
			log.Fatal(err)
			return
		}

	}

}
