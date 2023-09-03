package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"xero-bigquery-bulk-uploader/models"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	oauth2Config = oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/callback",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.xero.com/identity/connect/authorize",
			TokenURL: "https://identity.xero.com/connect/token",
		},
		Scopes: []string{"offline_access openid profile email accounting.transactions"},
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/callback", handleCallback)
	http.ListenAndServe(":8080", nil)

}

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
	w.Write([]byte("Import started"))
	invoices, err := bulkImport(token)
	if err != nil {
		log.Fatal(err)
		return
	}
	body, err := json.Marshal(invoices)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(body)
}

func bulkImport(token *oauth2.Token) ([]models.XeroInvoice, error) {
	invoices := []models.XeroInvoice{}
	tenantID := []models.XeroCompany{
		{ID: os.Getenv("CF_TENANT_ID"), Company: "CF"},
		{ID: os.Getenv("KD_TENANT_ID"), Company: "KD"}}
	for _, tenant := range tenantID {
		page := 1
		for {
			invoice := models.InvoiceBody{}
			invoicesBytes, err := getInvoices(token, page, tenant)
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
	}
	fmt.Println(invoices)

	return invoices, nil
}

func getInvoices(token *oauth2.Token, page int, tenant models.XeroCompany) ([]byte, error) {
	req, err := http.NewRequest("GET", "https://api.xero.com/api.xro/2.0/Invoices", nil)
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("where", `Status!="DELETED" AND Status!="VOIDED" AND Type=="ACCREC"`)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("xero-tenant-id", tenant.ID)
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
