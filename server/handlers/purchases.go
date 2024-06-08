package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

type Purchase struct {
	Date          string        `json:"date" validate:"required"` // TODO: /* time.Time */
	VendorDetails VendorDetails `json:"vendorDetails" validate:"required"`
	LineItems     []Item        `json:"lineItems"`
}

type Item struct {
	SlNo     int    `json:"slNo" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Quantity string `json:"quantity"`
	Discount string `json:"discount"`
	Net      string `json:"net"`
	RPU      string `json:"rpu"`
	Unit     string `json:"unit"` // TODO convey
}

type VendorDetails struct {
	Name      string `json:"vendorName" validate:"required"`
	Address   string `json:"address"`
	GstNo     string `json:"gstNo"`
	InvoiceNo string `json:"invoiceNo"`
	TinNo     string `json:"tinNo"`
	ContactNo string `json:"contactNo"`
}

var newPurchases = make(map[string]Purchase)

func NewPurchase(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Registering new purchase...")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Only HTTP POST request is allowed"))
		log.Default().Fatal("Illegal HTTP method used for registering new purchase")
		return
	}

	r.ParseMultipartForm(200 << 20) // Maximum of 200MB file allowed

	file, handler, err := r.FormFile("file")
	if err != nil {
		errStr := fmt.Sprintf("Error in reading the file %s\n", err)
		fmt.Println(errStr)
		fmt.Fprintf(w, errStr)
		return
	}

	_, err = saveFile(file, handler)

	if err != nil {
		// Error handling here
		log.Default().Fatal("Error while saving file to temp location: ", err.Error())
		return
	}

	var purchase Purchase
	purDataJson := r.FormValue("data")
	fmt.Println("Purchase data is: ", purDataJson)
	err = json.Unmarshal([]byte(purDataJson), &purchase)
	if err != nil {
		errStr := fmt.Sprintf("Error in Unmarshalling purchase data %s\n", err)
		fmt.Println(errStr)
		fmt.Fprintf(w, errStr)
		return
	}
	log.Default().Print("New Purchase Created Successfully: ", purchase.VendorDetails.InvoiceNo)
	newPurchases[purchase.VendorDetails.InvoiceNo] = purchase
	w.WriteHeader(http.StatusCreated)
}

func saveFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	defer file.Close()
	log.Default().Println("File name to be uploaded: ", handler.Filename)
	filebytes, err := io.ReadAll(file)
	if err != nil {
		errStr := fmt.Sprintf("Error in reading the file buffer %s\n", err)
		fmt.Println(errStr)
		return errStr, err
	}
	err = os.WriteFile(handler.Filename, filebytes, os.ModeAppend)

	if err != nil {
		log.Default().Fatal("File upload failed: ", err.Error())
		return "", err
	}
	return "Successfully uploaded\n", nil
}

func SearchPurchase(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	invoice := q["invoiceNumber"]
	log.Default().Println("Invoice number to be searched: ", invoice[0])
	purchase := newPurchases[strings.Trim(invoice[0], " ")]
	data, err := json.Marshal(purchase)
	if err != nil {
		log.Fatal("Marshalling searched purchase failed: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Marshalling searched purchase failed: ", err.Error())
		return
	}
	if purchase.VendorDetails.InvoiceNo == "" {
		log.Println("Requested search purhcase not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
