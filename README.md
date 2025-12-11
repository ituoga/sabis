# SABIS tinkamo XML konvertavimas iš JSON formato

version alpha 1.1


xsd's https://docs.oasis-open.org/ubl/os-UBL-2.1/xsd/
validator https://ecosio.com/en/peppol-e-invoice-xml-document-validator/
learning https://sabis-mok.nbfc.lt/

# SABIS Invoice Converter (JSON to UBL XML)

Tai lengvasvoris **Go** įrankis, skirtas konvertuoti sąskaitas faktūras iš paprasto **JSON** formato į sudėtingą **UBL 2.1 XML** standartą, tinkamą Lietuvos **SABIS** (Sąskaitų administravimo bendroji informacinė sistema) ir kitiems **Peppol** tinklo dalyviams.

## Kodėl šis įrankis?

UBL (Universal Business Language) XML struktūra yra itin gili ir sudėtinga, o XSD schemos dažnai sukelia problemų generuojant kodą automatiškai (dėl *circular dependencies* ir *namespace clashing*).

Daugelis programuotojų nori tiesiog suformuoti JSON iš savo duomenų bazės ir gauti validų XML failą siuntimui. Šis įrankis būtent tai ir daro.

**Pagrindinės savybės:**
* **Teisingi Namespaces:** Generuoja „švarius“ `cbc:` ir `cac:` prefiksus, kurių reikalauja griežti validatoriai (pvz., SABIS).
* **Suderinamumas:** Struktūra atitinka EN 16931 (Peppol BIS Billing 3.0) reikalavimus.
* **Paprastumas:** Paslepia XML sudėtingumą po paprasta Go struktūra.

## Naudojimas

### 1. Duomenų paruošimas
Sukurkite failą `input.json` tame pačiame kataloge. Struktūros pavyzdys:

```json
{
  "id": "SASK-2023-001",
  "issue_date": "2023-12-11",
  "due_date": "2024-01-11",
  "currency": "EUR",
  "supplier": {
    "name": "Pardavėjas UAB",
    "company_id": "300000000",
    "vat_id": "LT100000000000",
    "street": "Gedimino pr. 1",
    "city": "Vilnius",
    "country": "LT"
  },
  "customer": {
    "name": "Pirkėjas VšĮ",
    "company_id": "100000000",
    "vat_id": "LT200000000000",
    "street": "Konstitucijos pr. 20",
    "city": "Vilnius",
    "country": "LT"
  },
  "lines": [
    {
      "id": "1",
      "description": "Programavimo paslaugos",
      "quantity": 10,
      "price": 50.00,
      "amount": 500.00
    }
  ],
  "tax_amount": 105.00,
  "tax_subtotal_amount": 500.00,
  "tax_percent": 21.00,
  "net_amount": 500.00,
  "payable_amount": 605.00
}
```

### 2. Paleidimas

`go run . > output.xml`

### 3. Rezultatas
Gautas `output.xml` failas bus validus UBL 2.1 XML formatas, paruoštas siuntimui per SABIS ar Peppol tinklą
