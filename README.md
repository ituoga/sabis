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
  "id": "SASK-2023-999",
  "issue_date": "2023-12-11",
  "due_date": "2024-01-11",
  "currency": "EUR",
  "order_id": "1",
  "contract_id": "ŽODINĖ",
  "supplier": {
    "name": "Mano Įmonė UAB",
    "company_id": "300000001",
    "vat_id": "LT100000000011",
    "street": "Gedimino pr. 1",
    "city": "Vilnius",
    "country": "LT"
  },
  "customer": {
    "name": "Klientas UAB",
    "company_id": "300000002",
    "vat_id": "LT100000000022",
    "street": "Savanorių pr. 100",
    "city": "Kaunas",
    "country": "LT"
  },
  "lines": [
    {
      "id": "1",
      "description": "Programavimo paslaugos",
      "quantity": 10,
      "price": 50.0,
      "amount": 500.0,
      "tax_percent": 21.0
    },
    {
      "id": "2",
      "description": "Serverio nuoma",
      "quantity": 1,
      "price": 100.0,
      "amount": 100.0,
      "tax_percent": 21.0
    }
  ],
  "tax_amount": 126.0,
  "tax_subtotal_amount": 600.0,
  "tax_percent": 21.0,
  "net_amount": 600.0,
  "payable_amount": 726.0
}

```

### 2. Paleidimas

`go run . ./input.json > output.xml`

### 3. Rezultatas
Gautas `output.xml` failas bus validus UBL 2.1 XML formatas, paruoštas siuntimui per SABIS ar Peppol tinklą
