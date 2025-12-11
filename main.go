package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// --- 1. JSON (INPUT) STRUKTŪRA ---
// Atitinka jūsų turimą input.json failą

type InputInvoice struct {
	ID                string     `json:"id"`
	IssueDate         string     `json:"issue_date"`
	DueDate           string     `json:"due_date"`
	Currency          string     `json:"currency"`
	Supplier          Company    `json:"supplier"`
	Customer          Company    `json:"customer"`
	Lines             []LineItem `json:"lines"`
	TaxAmount         float64    `json:"tax_amount"`
	TaxSubtotalAmount float64    `json:"tax_subtotal_amount"`
	TaxPercent        float64    `json:"tax_percent"`
	NetAmount         float64    `json:"net_amount"`
	PayableAmount     float64    `json:"payable_amount"`
}

type Company struct {
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
	VatID     string `json:"vat_id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Country   string `json:"country"`
}

type LineItem struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	Price       float64 `json:"price"`
	Amount      float64 `json:"amount"`
}

// --- 2. UBL XML (OUTPUT) STRUKTŪRA ---
// Naudojame tiesioginius prefiksus "cbc:" ir "cac:" taguose

type UBLInvoice struct {
	XMLName xml.Name `xml:"Invoice"`

	// Namespace deklaracijos (būtinos šaknyje)
	Xmlns    string `xml:"xmlns,attr"`
	XmlnsCac string `xml:"xmlns:cac,attr"`
	XmlnsCbc string `xml:"xmlns:cbc,attr"`

	UBLVersionID     string `xml:"cbc:UBLVersionID"`
	CustomizationID  string `xml:"cbc:CustomizationID"`
	ProfileID        string `xml:"cbc:ProfileID"`
	ID               string `xml:"cbc:ID"`
	IssueDate        string `xml:"cbc:IssueDate"`
	DueDate          string `xml:"cbc:DueDate,omitempty"` // omitempty, jei nėra
	InvoiceTypeCode  string `xml:"cbc:InvoiceTypeCode"`
	DocumentCurrency string `xml:"cbc:DocumentCurrencyCode"`

	// --- Šalys (Parties) ---
	AccountingSupplierParty SupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty CustomerParty `xml:"cac:AccountingCustomerParty"`

	// --- Mokesčiai ir Sumos ---
	TaxTotal           TaxTotal      `xml:"cac:TaxTotal"`
	LegalMonetaryTotal MonetaryTotal `xml:"cac:LegalMonetaryTotal"`

	// --- Eilutės ---
	InvoiceLine []InvoiceLine `xml:"cac:InvoiceLine"`
}

// --- Pagalbinės XML struktūros ---

type SupplierParty struct {
	Party Party `xml:"cac:Party"`
}

type CustomerParty struct {
	Party Party `xml:"cac:Party"`
}

type Party struct {
	PartyName        PartyName         `xml:"cac:PartyName"`
	PostalAddress    *PostalAddress    `xml:"cac:PostalAddress,omitempty"`
	PartyTaxScheme   *PartyTaxScheme   `xml:"cac:PartyTaxScheme,omitempty"`
	PartyLegalEntity *PartyLegalEntity `xml:"cac:PartyLegalEntity,omitempty"`
}

type PartyName struct {
	Name string `xml:"cbc:Name"`
}

type PostalAddress struct {
	StreetName string      `xml:"cbc:StreetName"`
	CityName   string      `xml:"cbc:CityName"`
	Country    CountryType `xml:"cac:Country"`
}

type CountryType struct {
	IdentificationCode string `xml:"cbc:IdentificationCode"`
}

type PartyTaxScheme struct {
	CompanyID string    `xml:"cbc:CompanyID"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

type TaxScheme struct {
	ID string `xml:"cbc:ID"`
}

type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
	CompanyID        string `xml:"cbc:CompanyID"`
}

// --- Mokesčių struktūros ---

type TaxTotal struct {
	TaxAmount   Amount        `xml:"cbc:TaxAmount"`
	TaxSubtotal []TaxSubtotal `xml:"cac:TaxSubtotal,omitempty"`
}

type TaxSubtotal struct {
	TaxableAmount Amount      `xml:"cbc:TaxableAmount"`
	TaxAmount     Amount      `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategory `xml:"cac:TaxCategory"`
}

type TaxCategory struct {
	ID        string    `xml:"cbc:ID"`
	Percent   float64   `xml:"cbc:Percent"`
	TaxScheme TaxScheme `xml:"cac:TaxScheme"`
}

// --- Galutinės sumos ---

type MonetaryTotal struct {
	LineExtensionAmount Amount `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount  Amount `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount  Amount `xml:"cbc:TaxInclusiveAmount"`
	PayableAmount       Amount `xml:"cbc:PayableAmount"`
}

// --- Eilutės ---

type InvoiceLine struct {
	ID                  string   `xml:"cbc:ID"`
	InvoicedQuantity    Quantity `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount Amount   `xml:"cbc:LineExtensionAmount"`
	Item                Item     `xml:"cac:Item"`
	Price               *Price   `xml:"cac:Price,omitempty"`
}

type Item struct {
	Name string `xml:"cbc:Name"`
	// Galima pridėti ClassifiedTaxCategory čia, jei reikia
}

type Price struct {
	PriceAmount Amount `xml:"cbc:PriceAmount"`
}

// --- Baziniai tipai su atributais ---

type Amount struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Value      float64 `xml:",chardata"`
}

type Quantity struct {
	UnitCode string  `xml:"unitCode,attr"`
	Value    float64 `xml:",chardata"`
}

// --- MAIN FUNKCIJA ---

func main() {
	// 1. Nuskaitome failą
	file, err := os.Open("input.json")
	if err != nil {
		fmt.Println("Klaida: Nėra input.json failo")
		return
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var input InputInvoice
	json.Unmarshal(byteValue, &input)

	currency := input.Currency

	// 2. Kuriame UBL struktūrą
	ubl := UBLInvoice{
		Xmlns:            "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		XmlnsCac:         "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		XmlnsCbc:         "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		UBLVersionID:     "2.1",
		CustomizationID:  "urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0",
		ProfileID:        "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0",
		ID:               input.ID,
		IssueDate:        input.IssueDate,
		DueDate:          input.DueDate,
		InvoiceTypeCode:  "380",
		DocumentCurrency: currency,

		// PARDAVĖJAS
		AccountingSupplierParty: SupplierParty{
			Party: Party{
				PartyName: PartyName{Name: input.Supplier.Name},
				PostalAddress: &PostalAddress{
					StreetName: input.Supplier.Street,
					CityName:   input.Supplier.City,
					Country:    CountryType{IdentificationCode: input.Supplier.Country},
				},
				PartyTaxScheme: &PartyTaxScheme{
					CompanyID: input.Supplier.VatID,
					TaxScheme: TaxScheme{ID: "VAT"},
				},
				PartyLegalEntity: &PartyLegalEntity{
					RegistrationName: input.Supplier.Name,
					CompanyID:        input.Supplier.CompanyID,
				},
			},
		},

		// PIRKĖJAS
		AccountingCustomerParty: CustomerParty{
			Party: Party{
				PartyName: PartyName{Name: input.Customer.Name},
				PostalAddress: &PostalAddress{
					StreetName: input.Customer.Street,
					CityName:   input.Customer.City,
					Country:    CountryType{IdentificationCode: input.Customer.Country},
				},
				PartyTaxScheme: &PartyTaxScheme{
					CompanyID: input.Customer.VatID,
					TaxScheme: TaxScheme{ID: "VAT"},
				},
				PartyLegalEntity: &PartyLegalEntity{
					RegistrationName: input.Customer.Name,
					CompanyID:        input.Customer.CompanyID,
				},
			},
		},

		// MOKESČIAI
		TaxTotal: TaxTotal{
			TaxAmount: Amount{CurrencyID: currency, Value: input.TaxAmount},
			TaxSubtotal: []TaxSubtotal{
				{
					TaxableAmount: Amount{CurrencyID: currency, Value: input.TaxSubtotalAmount},
					TaxAmount:     Amount{CurrencyID: currency, Value: input.TaxAmount},
					TaxCategory: TaxCategory{
						ID:        "S", // Standard rate
						Percent:   input.TaxPercent,
						TaxScheme: TaxScheme{ID: "VAT"},
					},
				},
			},
		},

		// GALUTINĖS SUMOS
		LegalMonetaryTotal: MonetaryTotal{
			LineExtensionAmount: Amount{CurrencyID: currency, Value: input.NetAmount},
			TaxExclusiveAmount:  Amount{CurrencyID: currency, Value: input.NetAmount}, // Dažniausiai sutampa su NetAmount, jei nėra ne PVM mokesčių
			TaxInclusiveAmount:  Amount{CurrencyID: currency, Value: input.PayableAmount},
			PayableAmount:       Amount{CurrencyID: currency, Value: input.PayableAmount},
		},
	}

	// EILUTĖS
	for _, line := range input.Lines {
		ublLine := InvoiceLine{
			ID:                  line.ID,
			InvoicedQuantity:    Quantity{UnitCode: "H87", Value: line.Quantity}, // H87 = vienetas
			LineExtensionAmount: Amount{CurrencyID: currency, Value: line.Amount},
			Item: Item{
				Name: line.Description,
			},
			Price: &Price{
				PriceAmount: Amount{CurrencyID: currency, Value: line.Price},
			},
		}
		ubl.InvoiceLine = append(ubl.InvoiceLine, ublLine)
	}

	// 3. Išvedimas
	output, err := xml.MarshalIndent(ubl, "", "  ")
	if err != nil {
		panic(err)
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
