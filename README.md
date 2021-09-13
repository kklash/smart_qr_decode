# Smart QR Decode

I recently received a [SMART digital health card](https://smarthealth.cards) to prove my vaccination status. The official guidance was that it would allow me to prove my vaccination status without sharing my personal information, and was resistant to forgery. I was sceptical of such a claim, so I wrote a tool (this repo) to decode the data encoded in the SMART health card's QR code, so that I could inspect it myself. I published it so that you, dear reader, could also examine the information inside _your_ health card.

## Building & Running

This tool is written in [Golang](https://golang.org) and requires a Go compiler to build.

```
# local build inside the working directory:
go build

# if you want it accessible in your PATH:
go install

# run tests
go test
```

## Usage:

First, you will need to decode the data from your smart health card's QR code (outside the scope of this tool, for the moment).

 Most phone app QR code readers can let you copy text from a QR code. Alternatively, try [`zbarimg`](http://zbar.sourceforge.net/index.html) (`apt install zbar-tools` on Linux) to scan static image files.

Once you have the data, which is text that looks like `QR-Code:shc:/5676290952...` or just `shc:/5676290952...`, pass that data to this `smart_qr_decode` as a positional argument:

```
./smart_qr_decode 'QR-Code:shc:/5676295953...'
```

There's [an example image](./qrtest.png) available in this repo if you'd like to try it out from the beginning (scanning), or you can try this:

```
./smart_qr_decode 'QR-Code:shc:/567629095243206034602924374044603122295953265460346029254077280433602870286471674522280928613331456437653141590640220306450459085643550341424541364037063665417137241236380304375622046737407532323925433443326057360104413333527570752435036945652955560456532134341275386403316824406473453857292204440900303559452952006826505600550339666359775936093042093843697708700410035935662630262131684000336442677526685945265961385654232160612010575336603372416830094368087268583034220858653868405566737300676125522604302403684570446967365009635270603428764571693256573264396370103435043129720355532071362041710338407377455561123629414403691039073772267629582653595341365920676509072704756762393223665909237432547227110341752758564033075344643024205259213433223275231145112344442077504133374423657738732465036268045333500372292428680305752055066161205908590727574039273565100365453467620676070466562374456011640503052612603433352534627643734526297670536641110369320974692670435941656559530672526365744144392929372455403655296826416910330331702770525959062807257522393924313335716156446225321038684203745070327705435227392923353242241273761130335734312154002724683467287609315543752975265426124341547269386510063170735444033524323109036562255462334100753044505239593704336853601144044303067643563009455908314400052367685333606845283543332368052172450304695777120961747623331028343573401038623404535367200023435327675723580300766876627154233043006804230025103369507231500670446306070903575832014453337756590712383705676069074064360005643067283632385063256409252323366162400606120354527428123073066511712009226844663207632755572431336840456368371037205910610777322474'
```

Output:

```
Decoded JSON Web Signature vaccine card data

Header:
---------------------------------------------
{"zip":"DEF","alg":"ES256","kid":"3Kfdg-XwP-7gXyywtUfUADwBumDOPKMQx-iELL11W9s"}

Payload:
---------------------------------------------
{"iss":"https://c19.cards/issuer","nbf":1591037940,"vc":{"type":["https://smarthealth.cards#covid19","https://smarthealth.cards#health-card","https://smarthealth.cards#immunization"],"credentialSubject":{"fhirVersion":"4.0.1","fhirBundle":{"resourceType":"Bundle","type":"collection","entry":[{"fullUrl":"resource:0","resource":{"resourceType":"Patient","name":[{"family":"Anyperson","given":["Johnathan","Biggleston III"]}],"birthDate":"1951-01-20"}},{"fullUrl":"resource:1","resource":{"resourceType":"Immunization","meta":{"security":[{"code":"IAL1.2"}]},"status":"completed","vaccineCode":{"coding":[{"system":"http://hl7.org/fhir/sid/cvx","code":"207"}]},"patient":{"reference":"resource:0"},"occurrenceDateTime":"2021-01-01","location":{"reference":"resource:3"},"performer":[{"actor":{"display":"ABC General Hospital"}}],"lotNumber":"Lot #0000001"}},{"fullUrl":"resource:2","resource":{"resourceType":"Immunization","status":"completed","vaccineCode":{"coding":[{"system":"http://hl7.org/fhir/sid/cvx","code":"207"}]},"patient":{"reference":"resource:0"},"occurrenceDateTime":"2021-01-29","performer":[{"actor":{"display":"ABC General Hospital"}}],"lotNumber":"Lot #0000007"}}]}}}}

Signature:
---------------------------------------------
61b3737a1e3d491da98abe14990fb698aa4840c4bf9459ba1430d08e4537dfdd1c6b023d2afde7f2d03a0aa62833894775f10b36a51996a47b44087b8f8ccc13
```

## Sources

I primarily used [this article on Medium by Vishnu Ravi](https://vishnuravi.medium.com/how-do-verifiable-covid-19-vaccination-records-with-smart-health-cards-work-df099370b27a) as a guide, and later referred to the [Official SMART health card documentation](https://spec.smarthealth.cards) to further my understanding, specifically [this section](https://spec.smarthealth.cards/#health-cards-are-encoded-as-compact-serialization-json-web-signatures-jws) was very useful.

## So What's In It?

If your vaccine card QR code is like mine, then it contains:
- Your name
- Your birth date
- The lot numbers, dates, and locations of each of your vaccine doses

Here's an example card payload taken from Dr Ravi's article:

```json
{
  "iss": "https://spec.smarthealth.cards/examples/issuer",
  "nbf": 1622690247.979,
  "vc": {
    "type": [
      "https://smarthealth.cards#health-card",
      "https://smarthealth.cards#immunization",
      "https://smarthealth.cards#covid19"
    ],
    "credentialSubject": {
      "fhirVersion": "4.0.1",
      "fhirBundle": {
        "resourceType": "Bundle",
        "type": "collection",
        "entry": [
          {
            "fullUrl": "resource:0",
            "resource": {
              "resourceType": "Patient",
              "name": [
                {
                  "family": "Anyperson",
                  "given": [
                    "John",
                    "B."
                  ]
                }
              ],
              "birthDate": "1951-01-20"
            }
          },
          {
            "fullUrl": "resource:1",
            "resource": {
              "resourceType": "Immunization",
              "status": "completed",
              "vaccineCode": {
                "coding": [
                  {
                    "system": "http://hl7.org/fhir/sid/cvx",
                    "code": "207"
                  }
                ]
              },
              "patient": {
                "reference": "resource:0"
              },
              "occurrenceDateTime": "2021-01-01",
              "performer": [
                {
                  "actor": {
                    "display": "ABC General Hospital"
                  }
                }
              ],
              "lotNumber": "0000001"
            }
          },
          {
            "fullUrl": "resource:2",
            "resource": {
              "resourceType": "Immunization",
              "status": "completed",
              "vaccineCode": {
                "coding": [
                  {
                    "system": "http://hl7.org/fhir/sid/cvx",
                    "code": "207"
                  }
                ]
              },
              "patient": {
                "reference": "resource:0"
              },
              "occurrenceDateTime": "2021-01-29",
              "performer": [
                {
                  "actor": {
                    "display": "ABC General Hospital"
                  }
                }
              ],
              "lotNumber": "0000007"
            }
          }
        ]
      }
    }
  }
}
```
