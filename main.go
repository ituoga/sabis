package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

// =========================================================
// 1. JSON INPUT STRUKTŪRA
// =========================================================

type InputInvoice struct {
	ID                string      `json:"id"`
	IssueDate         string      `json:"issue_date"`
	DueDate           string      `json:"due_date"`
	Currency          string      `json:"currency"`
	Note              string      `json:"note"`
	OrderID           string      `json:"order_id"`
	ContractID        string      `json:"contract_id"`
	Project           ProjectInfo `json:"project"`
	Supplier          Company     `json:"supplier"`
	Customer          Company     `json:"customer"`
	Lines             []LineItem  `json:"lines"`
	TaxAmount         float64     `json:"tax_amount"`
	TaxSubtotalAmount float64     `json:"tax_subtotal_amount"`
	TaxPercent        float64     `json:"tax_percent"`
	NetAmount         float64     `json:"net_amount"`
	PayableAmount     float64     `json:"payable_amount"`
}

type ProjectInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
	TaxPercent  float64 `json:"tax_percent"`
}

// =========================================================
// 2. UBL XML STRUKTŪRA
// =========================================================

type Invoice struct {
	XMLName xml.Name `xml:"Invoice"`
	Xmlns   string   `xml:"xmlns,attr"`
	Cac     string   `xml:"xmlns:cac,attr"`
	Cbc     string   `xml:"xmlns:cbc,attr"`

	UBLVersionID         string `xml:"cbc:UBLVersionID"`
	CustomizationID      string `xml:"cbc:CustomizationID"`
	ProfileID            string `xml:"cbc:ProfileID"`
	ID                   string `xml:"cbc:ID"`
	IssueDate            string `xml:"cbc:IssueDate"`
	DueDate              string `xml:"cbc:DueDate"`
	InvoiceTypeCode      string `xml:"cbc:InvoiceTypeCode"`
	Note                 string `xml:"cbc:Note,omitempty"`
	DocumentCurrencyCode string `xml:"cbc:DocumentCurrencyCode"`

	// Nuorodos
	OrderReference            *DocReference       `xml:"cac:OrderReference,omitempty"`
	ContractDocumentReference *DocReference       `xml:"cac:ContractDocumentReference,omitempty"`
	ProcurementProject        *ProcurementProject `xml:"cac:ProcurementProject,omitempty"`

	// Šalys
	AccountingSupplierParty SupplierParty `xml:"cac:AccountingSupplierParty"`
	AccountingCustomerParty CustomerParty `xml:"cac:AccountingCustomerParty"`

	// Mokesčiai ir Sumos
	TaxTotal           TaxTotal           `xml:"cac:TaxTotal"`
	LegalMonetaryTotal LegalMonetaryTotal `xml:"cac:LegalMonetaryTotal"`

	// Eilutės
	InvoiceLine []InvoiceLine `xml:"cac:InvoiceLine"`
}

// --- Bendriniai tipai ---

type AmountType struct {
	Text       string `xml:",chardata"`
	CurrencyID string `xml:"currencyID,attr"`
}

type QuantityType struct {
	Text     string `xml:",chardata"`
	UnitCode string `xml:"unitCode,attr"`
}

type DocReference struct {
	ID string `xml:"cbc:ID"`
}

type ProcurementProject struct {
	ID   string `xml:"cbc:ID"`
	Name string `xml:"cbc:Name,omitempty"`
}

type EndpointIDType struct {
	Text     string `xml:",chardata"`
	SchemeID string `xml:"schemeID,attr"`
}

// --- Šalių struktūros ---

type SupplierParty struct {
	Party PartyType `xml:"cac:Party"`
}

type CustomerParty struct {
	Party PartyType `xml:"cac:Party"`
}

type PartyType struct {
	EndpointID       *EndpointIDType   `xml:"cbc:EndpointID,omitempty"`
	PartyName        *PartyNameType    `xml:"cac:PartyName,omitempty"`
	PostalAddress    *AddressType      `xml:"cac:PostalAddress,omitempty"`
	PartyTaxScheme   *PartyTaxScheme   `xml:"cac:PartyTaxScheme,omitempty"`
	PartyLegalEntity *PartyLegalEntity `xml:"cac:PartyLegalEntity,omitempty"`
}

type PartyNameType struct {
	Name string `xml:"cbc:Name"`
}

type AddressType struct {
	StreetName string      `xml:"cbc:StreetName"`
	CityName   string      `xml:"cbc:CityName"`
	Country    CountryType `xml:"cac:Country"`
}

type CountryType struct {
	IdentificationCode string `xml:"cbc:IdentificationCode"`
}

type PartyTaxScheme struct {
	CompanyID string        `xml:"cbc:CompanyID"`
	TaxScheme TaxSchemeType `xml:"cac:TaxScheme"`
}

type TaxSchemeType struct {
	ID string `xml:"cbc:ID"`
}

type PartyLegalEntity struct {
	RegistrationName string `xml:"cbc:RegistrationName"`
	CompanyID        string `xml:"cbc:CompanyID"`
}

// --- Mokesčių struktūros ---

type TaxTotal struct {
	TaxAmount   AmountType    `xml:"cbc:TaxAmount"`
	TaxSubtotal []TaxSubtotal `xml:"cac:TaxSubtotal,omitempty"`
}

type TaxSubtotal struct {
	TaxableAmount AmountType      `xml:"cbc:TaxableAmount"`
	TaxAmount     AmountType      `xml:"cbc:TaxAmount"`
	TaxCategory   TaxCategoryType `xml:"cac:TaxCategory"`
}

type TaxCategoryType struct {
	ID        string        `xml:"cbc:ID"`
	Percent   string        `xml:"cbc:Percent"`
	TaxScheme TaxSchemeType `xml:"cac:TaxScheme"`
}

// --- Sumų struktūros ---

type LegalMonetaryTotal struct {
	LineExtensionAmount AmountType `xml:"cbc:LineExtensionAmount"`
	TaxExclusiveAmount  AmountType `xml:"cbc:TaxExclusiveAmount"`
	TaxInclusiveAmount  AmountType `xml:"cbc:TaxInclusiveAmount"`
	PayableAmount       AmountType `xml:"cbc:PayableAmount"`
}

// --- Eilučių struktūros ---

type InvoiceLine struct {
	ID                  string       `xml:"cbc:ID"`
	InvoicedQuantity    QuantityType `xml:"cbc:InvoicedQuantity"`
	LineExtensionAmount AmountType   `xml:"cbc:LineExtensionAmount"`
	Item                ItemType     `xml:"cac:Item"`
	Price               *PriceType   `xml:"cac:Price,omitempty"`
}

type ItemType struct {
	Name string `xml:"cbc:Name"`
	// NAUJAS LAUKAS: Privalomas klasifikuotas mokesčio tarifas prekei
	ClassifiedTaxCategory *TaxCategoryType `xml:"cac:ClassifiedTaxCategory,omitempty"`
}

type PriceType struct {
	PriceAmount AmountType `xml:"cbc:PriceAmount"`
}

// =========================================================
// 3. MAIN FUNKCIJA
// =========================================================

func main() {
	file, err := os.Open("input.json")
	if err != nil {
		fmt.Println("KLAIDA: Nėra input.json failo")
		return
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var input InputInvoice
	json.Unmarshal(byteValue, &input)

	curr := input.Currency

	inv := Invoice{
		Xmlns:                "urn:oasis:names:specification:ubl:schema:xsd:Invoice-2",
		Cac:                  "urn:oasis:names:specification:ubl:schema:xsd:CommonAggregateComponents-2",
		Cbc:                  "urn:oasis:names:specification:ubl:schema:xsd:CommonBasicComponents-2",
		UBLVersionID:         "2.1",
		CustomizationID:      "urn:cen.eu:en16931:2017#compliant#urn:fdc:peppol.eu:2017:poacc:billing:3.0",
		ProfileID:            "urn:fdc:peppol.eu:2017:poacc:billing:01:1.0",
		ID:                   input.ID,
		IssueDate:            input.IssueDate,
		DueDate:              input.DueDate,
		InvoiceTypeCode:      "380",
		Note:                 input.Note,
		DocumentCurrencyCode: curr,

		// PARDAVĖJAS
		AccountingSupplierParty: SupplierParty{
			Party: PartyType{
				EndpointID: &EndpointIDType{
					Text:     input.Supplier.CompanyID,
					SchemeID: "0200",
				},
				PartyName: &PartyNameType{Name: input.Supplier.Name},
				PostalAddress: &AddressType{
					StreetName: input.Supplier.Street,
					CityName:   input.Supplier.City,
					Country:    CountryType{IdentificationCode: input.Supplier.Country},
				},
				PartyTaxScheme: &PartyTaxScheme{
					CompanyID: input.Supplier.VatID,
					TaxScheme: TaxSchemeType{ID: "VAT"},
				},
				PartyLegalEntity: &PartyLegalEntity{
					RegistrationName: input.Supplier.Name,
					CompanyID:        input.Supplier.CompanyID,
				},
			},
		},

		// PIRKĖJAS
		AccountingCustomerParty: CustomerParty{
			Party: PartyType{
				EndpointID: &EndpointIDType{
					Text:     input.Customer.CompanyID,
					SchemeID: "0200",
				},
				PartyName: &PartyNameType{Name: input.Customer.Name},
				PostalAddress: &AddressType{
					StreetName: input.Customer.Street,
					CityName:   input.Customer.City,
					Country:    CountryType{IdentificationCode: input.Customer.Country},
				},
				PartyTaxScheme: &PartyTaxScheme{
					CompanyID: input.Customer.VatID,
					TaxScheme: TaxSchemeType{ID: "VAT"},
				},
				PartyLegalEntity: &PartyLegalEntity{
					RegistrationName: input.Customer.Name,
					CompanyID:        input.Customer.CompanyID,
				},
			},
		},

		// Mokesčiai
		TaxTotal: TaxTotal{
			TaxAmount: AmountType{Text: fmt.Sprintf("%.2f", input.TaxAmount), CurrencyID: curr},
			TaxSubtotal: []TaxSubtotal{
				{
					TaxableAmount: AmountType{Text: fmt.Sprintf("%.2f", input.TaxSubtotalAmount), CurrencyID: curr},
					TaxAmount:     AmountType{Text: fmt.Sprintf("%.2f", input.TaxAmount), CurrencyID: curr},
					TaxCategory: TaxCategoryType{
						ID:        "S",
						Percent:   fmt.Sprintf("%.2f", input.TaxPercent),
						TaxScheme: TaxSchemeType{ID: "VAT"},
					},
				},
			},
		},

		// Sumos
		LegalMonetaryTotal: LegalMonetaryTotal{
			LineExtensionAmount: AmountType{Text: fmt.Sprintf("%.2f", input.NetAmount), CurrencyID: curr},
			TaxExclusiveAmount:  AmountType{Text: fmt.Sprintf("%.2f", input.NetAmount), CurrencyID: curr},
			TaxInclusiveAmount:  AmountType{Text: fmt.Sprintf("%.2f", input.PayableAmount), CurrencyID: curr},
			PayableAmount:       AmountType{Text: fmt.Sprintf("%.2f", input.PayableAmount), CurrencyID: curr},
		},
	}

	if input.OrderID != "" {
		inv.OrderReference = &DocReference{ID: input.OrderID}
	}
	if input.ContractID != "" {
		inv.ContractDocumentReference = &DocReference{ID: input.ContractID}
	}
	if input.Project.ID != "" {
		inv.ProcurementProject = &ProcurementProject{ID: input.Project.ID, Name: input.Project.Name}
	} else if input.ContractID == "" {
		inv.ProcurementProject = &ProcurementProject{ID: "YES", Name: "IsNotForPublication"}
	}

	// Eilutės
	for _, l := range input.Lines {
		line := InvoiceLine{
			ID:                  l.ID,
			InvoicedQuantity:    QuantityType{Text: fmt.Sprintf("%.2f", l.Quantity), UnitCode: "H87"},
			LineExtensionAmount: AmountType{Text: fmt.Sprintf("%.2f", l.Amount), CurrencyID: curr},
			Item: ItemType{
				Name: l.Description,
				// PILDOMAS NAUJAS LAUKAS (Pagal Jūsų XML pavyzdį):
				ClassifiedTaxCategory: &TaxCategoryType{
					ID:      "S",                               // Standartinis tarifas
					Percent: fmt.Sprintf("%.2f", l.TaxPercent), // Pvz "21.00"
					TaxScheme: TaxSchemeType{
						ID: "VAT",
					},
				},
			},
			Price: &PriceType{
				PriceAmount: AmountType{Text: fmt.Sprintf("%.2f", l.Price), CurrencyID: curr},
			},
		}
		inv.InvoiceLine = append(inv.InvoiceLine, line)
	}

	output, _ := xml.MarshalIndent(inv, "", "  ")
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
