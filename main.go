package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// --- UBL XML Namespaces ---
const (
	NsInvoice = "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2"
	NsCac     = "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2"
	NsCbc     = "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2"
)

// --- Structs: JSON -> XML mapping ---

// InputInvoice - Supaprastinta struktūra JSON nuskaitymui
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

// --- Structs: XML UBL 2.1 Structure ---

type UBLInvoice struct {
	XMLName xml.Name `xml:"urn:oasis:names:specification:ubl:schema:xsd:Invoice-2 Invoice"`

	// Namespace Prefixes
	XmlnsCac string `xml:"xmlns:cac,attr"`
	XmlnsCbc string `xml:"xmlns:cbc,attr"`

	UBLVersionID     string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 UBLVersionID"`
	CustomizationID  string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CustomizationID"`
	ID               string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	IssueDate        string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IssueDate"`
	DueDate          string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 DueDate"`
	InvoiceTypeCode  string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InvoiceTypeCode"`
	DocumentCurrency string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 DocumentCurrencyCode"`

	AccountingSupplierParty SupplierParty `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AccountingSupplierParty"`
	AccountingCustomerParty CustomerParty `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 AccountingCustomerParty"`

	TaxTotal           TaxTotal           `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxTotal"`
	LegalMonetaryTotal LegalMonetaryTotal `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 LegalMonetaryTotal"`

	InvoiceLine []InvoiceLine `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 InvoiceLine"`
}

// --- Helper types for UBL ---

type SupplierParty struct {
	Party Party `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Party"`
}

type CustomerParty struct {
	Party Party `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Party"`
}

type Party struct {
	PartyName        PartyName        `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyName"`
	PostalAddress    PostalAddress    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PostalAddress"`
	PartyTaxScheme   PartyTaxScheme   `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyTaxScheme"`
	PartyLegalEntity PartyLegalEntity `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 PartyLegalEntity"`
}

type PartyName struct {
	Name string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
}

type PostalAddress struct {
	StreetName string      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 StreetName"`
	CityName   string      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CityName"`
	Country    CountryType `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Country"`
}

type CountryType struct {
	IdentificationCode string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 IdentificationCode"`
}

type PartyTaxScheme struct {
	CompanyID string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CompanyID"`
	TaxScheme TaxScheme `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxScheme"`
}

type PartyLegalEntity struct {
	RegistrationName string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 RegistrationName"`
	CompanyID        string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 CompanyID"`
}

type TaxScheme struct {
	ID string `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
}

type TaxTotal struct {
	TaxAmount   Amount        `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxAmount"`
	TaxSubtotal []TaxSubtotal `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxSubtotal"`
}

type TaxSubtotal struct {
	TaxableAmount Amount      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxableAmount"`
	TaxAmount     Amount      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxAmount"`
	TaxCategory   TaxCategory `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxCategory"`
}

type TaxCategory struct {
	ID        string    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	Percent   float64   `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Percent"`
	TaxScheme TaxScheme `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 TaxScheme"`
}

type LegalMonetaryTotal struct {
	LineExtensionAmount Amount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineExtensionAmount"`
	TaxExclusiveAmount  Amount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxExclusiveAmount"`
	TaxInclusiveAmount  Amount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 TaxInclusiveAmount"`
	PayableAmount       Amount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PayableAmount"`
}

type InvoiceLine struct {
	ID                  string   `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 ID"`
	InvoicedQuantity    Quantity `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 InvoicedQuantity"`
	LineExtensionAmount Amount   `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 LineExtensionAmount"`
	Item                Item     `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Item"`
	Price               Price    `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 Price"`
}

type Item struct {
	Name                  string      `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 Name"`
	ClassifiedTaxCategory TaxCategory `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2 ClassifiedTaxCategory"`
}

type Price struct {
	PriceAmount Amount `xml:"urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2 PriceAmount"`
}

type Amount struct {
	CurrencyID string  `xml:"currencyID,attr"`
	Value      float64 `xml:",chardata"`
}

type Quantity struct {
	UnitCode string  `xml:"unitCode,attr"`
	Value    float64 `xml:",chardata"`
}

// --- Main Converter Logic ---

func main() {
	// 1. Nuskaitome JSON failą
	file, err := os.Open("input.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var input InputInvoice
	json.Unmarshal(byteValue, &input)

	// 2. Mapiname JSON duomenis į UBL XML struktūrą
	// Čia vyksta pagrindinė konversija

	currency := input.Currency

	ubl := UBLInvoice{
		XmlnsCac:         NsCac,
		XmlnsCbc:         NsCbc,
		UBLVersionID:     "2.1",
		CustomizationID:  "urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0", // Peppol
		ID:               input.ID,
		IssueDate:        input.IssueDate,
		DueDate:          input.DueDate,
		InvoiceTypeCode:  "380",
		DocumentCurrency: currency,

		// Pardavėjas
		AccountingSupplierParty: SupplierParty{
			Party: Party{
				PartyName: PartyName{Name: input.Supplier.Name},
				PostalAddress: PostalAddress{
					StreetName: input.Supplier.Street,
					CityName:   input.Supplier.City,
					Country:    CountryType{IdentificationCode: input.Supplier.Country},
				},
				PartyTaxScheme: PartyTaxScheme{
					CompanyID: input.Supplier.VatID,
					TaxScheme: TaxScheme{ID: "VAT"},
				},
				PartyLegalEntity: PartyLegalEntity{
					RegistrationName: input.Supplier.Name,
					CompanyID:        input.Supplier.CompanyID,
				},
			},
		},

		// Pirkėjas
		AccountingCustomerParty: CustomerParty{
			Party: Party{
				PartyName: PartyName{Name: input.Customer.Name},
				PostalAddress: PostalAddress{
					StreetName: input.Customer.Street,
					CityName:   input.Customer.City,
					Country:    CountryType{IdentificationCode: input.Customer.Country},
				},
				PartyTaxScheme: PartyTaxScheme{
					CompanyID: input.Customer.VatID,
					TaxScheme: TaxScheme{ID: "VAT"},
				},
				PartyLegalEntity: PartyLegalEntity{
					RegistrationName: input.Customer.Name,
					CompanyID:        input.Customer.CompanyID,
				},
			},
		},

		// Mokesčių totalai
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

		// Galutinės sumos
		LegalMonetaryTotal: LegalMonetaryTotal{
			LineExtensionAmount: Amount{CurrencyID: currency, Value: input.NetAmount},
			TaxExclusiveAmount:  Amount{CurrencyID: currency, Value: input.NetAmount},
			TaxInclusiveAmount:  Amount{CurrencyID: currency, Value: input.PayableAmount},
			PayableAmount:       Amount{CurrencyID: currency, Value: input.PayableAmount},
		},
	}

	// Eilučių konversija
	for _, line := range input.Lines {
		ublLine := InvoiceLine{
			ID:                  line.ID,
			InvoicedQuantity:    Quantity{UnitCode: "H87", Value: line.Quantity}, // H87 = vienetai
			LineExtensionAmount: Amount{CurrencyID: currency, Value: line.Amount},
			Item: Item{
				Name: line.Description,
				ClassifiedTaxCategory: TaxCategory{
					ID:        "S",
					Percent:   input.TaxPercent,
					TaxScheme: TaxScheme{ID: "VAT"},
				},
			},
			Price: Price{
				PriceAmount: Amount{CurrencyID: currency, Value: line.Price},
			},
		}
		ubl.InvoiceLine = append(ubl.InvoiceLine, ublLine)
	}

	// 3. Išvedame XML
	output, _ := xml.MarshalIndent(ubl, "", "  ")
	fmt.Println(xml.Header + string(output))
}
