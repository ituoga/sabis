package main

import "encoding/xml"

// Invoice was generated 2025-12-11 21:27:10 by mind on mb.local with zek 0.1.28.
type ___Invoice_ struct {
	XMLName              xml.Name `xml:"Invoice"`
	Text                 string   `xml:",chardata"`
	Xmlns                string   `xml:"xmlns,attr"`
	Cac                  string   `xml:"cac,attr"`
	Cec                  string   `xml:"cec,attr"`
	Cbc                  string   `xml:"cbc,attr"`
	CustomizationID      string   `xml:"CustomizationID"`
	ProfileID            string   `xml:"ProfileID"`
	ID                   string   `xml:"ID"`
	IssueDate            string   `xml:"IssueDate"`
	DueDate              string   `xml:"DueDate"`
	InvoiceTypeCode      string   `xml:"InvoiceTypeCode"`
	DocumentCurrencyCode string   `xml:"DocumentCurrencyCode"`
	OrderReference       struct {
		Text string `xml:",chardata"`
		ID   string `xml:"ID"`
	} `xml:"OrderReference"`
	ContractDocumentReference struct {
		Text string `xml:",chardata"`
		ID   string `xml:"ID"`
	} `xml:"ContractDocumentReference"`
	AccountingSupplierParty struct {
		Text  string `xml:",chardata"`
		Party struct {
			Text       string `xml:",chardata"`
			EndpointID struct {
				Text     string `xml:",chardata"`
				SchemeID string `xml:"schemeID,attr"`
			} `xml:"EndpointID"`
			PartyName struct {
				Text string `xml:",chardata"`
				Name string `xml:"Name"`
			} `xml:"PartyName"`
			PostalAddress struct {
				Text       string `xml:",chardata"`
				StreetName string `xml:"StreetName"`
				Country    struct {
					Text               string `xml:",chardata"`
					IdentificationCode string `xml:"IdentificationCode"`
				} `xml:"Country"`
			} `xml:"PostalAddress"`
			PartyTaxScheme struct {
				Text      string `xml:",chardata"`
				CompanyID string `xml:"CompanyID"`
				TaxScheme struct {
					Text string `xml:",chardata"`
					ID   string `xml:"ID"`
				} `xml:"TaxScheme"`
			} `xml:"PartyTaxScheme"`
			PartyLegalEntity struct {
				Text             string `xml:",chardata"`
				RegistrationName string `xml:"RegistrationName"`
				CompanyID        string `xml:"CompanyID"`
			} `xml:"PartyLegalEntity"`
			Contact struct {
				Text           string `xml:",chardata"`
				Name           string `xml:"Name"`
				Telephone      string `xml:"Telephone"`
				ElectronicMail string `xml:"ElectronicMail"`
			} `xml:"Contact"`
		} `xml:"Party"`
	} `xml:"AccountingSupplierParty"`
	AccountingCustomerParty struct {
		Text  string `xml:",chardata"`
		Party struct {
			Text       string `xml:",chardata"`
			EndpointID struct {
				Text     string `xml:",chardata"`
				SchemeID string `xml:"schemeID,attr"`
			} `xml:"EndpointID"`
			PartyName struct {
				Text string `xml:",chardata"`
				Name string `xml:"Name"`
			} `xml:"PartyName"`
			PostalAddress struct {
				Text             string `xml:",chardata"`
				StreetName       string `xml:"StreetName"`
				CityName         string `xml:"CityName"`
				PostalZone       string `xml:"PostalZone"`
				CountrySubentity string `xml:"CountrySubentity"`
				Country          struct {
					Text               string `xml:",chardata"`
					IdentificationCode string `xml:"IdentificationCode"`
				} `xml:"Country"`
			} `xml:"PostalAddress"`
			PartyLegalEntity struct {
				Text             string `xml:",chardata"`
				RegistrationName string `xml:"RegistrationName"`
				CompanyID        string `xml:"CompanyID"`
			} `xml:"PartyLegalEntity"`
			Contact struct {
				Text           string `xml:",chardata"`
				Telephone      string `xml:"Telephone"`
				ElectronicMail string `xml:"ElectronicMail"`
			} `xml:"Contact"`
		} `xml:"Party"`
	} `xml:"AccountingCustomerParty"`
	PaymentMeans struct {
		Text                  string `xml:",chardata"`
		PaymentMeansCode      string `xml:"PaymentMeansCode"`
		PayeeFinancialAccount struct {
			Text string `xml:",chardata"`
			ID   string `xml:"ID"`
		} `xml:"PayeeFinancialAccount"`
	} `xml:"PaymentMeans"`
	TaxTotal struct {
		Text      string `xml:",chardata"`
		TaxAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"TaxAmount"`
		TaxSubtotal struct {
			Text          string `xml:",chardata"`
			TaxableAmount struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"TaxableAmount"`
			TaxAmount struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"TaxAmount"`
			TaxCategory struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"ID"`
				Percent   string `xml:"Percent"`
				TaxScheme struct {
					Text string `xml:",chardata"`
					ID   string `xml:"ID"`
				} `xml:"TaxScheme"`
			} `xml:"TaxCategory"`
		} `xml:"TaxSubtotal"`
	} `xml:"TaxTotal"`
	LegalMonetaryTotal struct {
		Text                string `xml:",chardata"`
		LineExtensionAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"LineExtensionAmount"`
		TaxExclusiveAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"TaxExclusiveAmount"`
		TaxInclusiveAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"TaxInclusiveAmount"`
		AllowanceTotalAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"AllowanceTotalAmount"`
		ChargeTotalAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"ChargeTotalAmount"`
		PayableRoundingAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"PayableRoundingAmount"`
		PayableAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"PayableAmount"`
	} `xml:"LegalMonetaryTotal"`
	InvoiceLine struct {
		Text             string `xml:",chardata"`
		ID               string `xml:"ID"`
		InvoicedQuantity struct {
			Text     string `xml:",chardata"`
			UnitCode string `xml:"unitCode,attr"`
		} `xml:"InvoicedQuantity"`
		LineExtensionAmount struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"LineExtensionAmount"`
		Item struct {
			Text                  string `xml:",chardata"`
			Name                  string `xml:"Name"`
			ClassifiedTaxCategory struct {
				Text      string `xml:",chardata"`
				ID        string `xml:"ID"`
				Percent   string `xml:"Percent"`
				TaxScheme struct {
					Text string `xml:",chardata"`
					ID   string `xml:"ID"`
				} `xml:"TaxScheme"`
			} `xml:"ClassifiedTaxCategory"`
		} `xml:"Item"`
		Price struct {
			Text        string `xml:",chardata"`
			PriceAmount struct {
				Text       string `xml:",chardata"`
				CurrencyID string `xml:"currencyID,attr"`
			} `xml:"PriceAmount"`
		} `xml:"Price"`
	} `xml:"InvoiceLine"`
}
