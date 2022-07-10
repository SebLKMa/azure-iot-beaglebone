package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/go-amqp"
)

// DigiCert Baltimore Root (sha1 fingerprint=d4de20d05e66fc53fe1a50882c78db2852cae474)
// Microsoft RSA TLS CA 01 (sha1 fingerprint=703d7a8f0ebf55aaa59f98eaf4a206004eb2516a)
// Microsoft RSA TLS CA 02 (sha1 fingerprint=b0c2d2d13cdd56cdaa6ab6e2c04440be4a429c75)
// Microsoft Azure TLS Issuing CA 01 (sha1 fingerprint=2f2877c5d778c31e0f29c7e371df5471bd673173)
// Microsoft Azure TLS Issuing CA 02 (sha1 fingerprint=e7eea674ca718e3befd90858e09f8372ad0ae2aa)
// Microsoft Azure TLS Issuing CA 05 (sha1 fingerprint=6c3af02e7f269aa73afd0eff2a88a4a1f04ed1e5)
// Microsoft Azure TLS Issuing CA 06 (sha1 fingerprint=30e01761ab97e59a06b41ef20af6f2de7ef4f7b0)
var caCerts = []byte(`-----BEGIN CERTIFICATE-----
MIIDdzCCAl+gAwIBAgIEAgAAuTANBgkqhkiG9w0BAQUFADBaMQswCQYDVQQGEwJJ
RTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJlclRydXN0MSIwIAYD
VQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTAwMDUxMjE4NDYwMFoX
DTI1MDUxMjIzNTkwMFowWjELMAkGA1UEBhMCSUUxEjAQBgNVBAoTCUJhbHRpbW9y
ZTETMBEGA1UECxMKQ3liZXJUcnVzdDEiMCAGA1UEAxMZQmFsdGltb3JlIEN5YmVy
VHJ1c3QgUm9vdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKMEuyKr
mD1X6CZymrV51Cni4eiVgLGw41uOKymaZN+hXe2wCQVt2yguzmKiYv60iNoS6zjr
IZ3AQSsBUnuId9Mcj8e6uYi1agnnc+gRQKfRzMpijS3ljwumUNKoUMMo6vWrJYeK
mpYcqWe4PwzV9/lSEy/CG9VwcPCPwBLKBsua4dnKM3p31vjsufFoREJIE9LAwqSu
XmD+tqYF/LTdB1kC1FkYmGP1pWPgkAx9XbIGevOF6uvUA65ehD5f/xXtabz5OTZy
dc93Uk3zyZAsuT3lySNTPx8kmCFcB5kpvcY67Oduhjprl3RjM71oGDHweI12v/ye
jl0qhqdNkNwnGjkCAwEAAaNFMEMwHQYDVR0OBBYEFOWdWTCCR1jMrPoIVDaGezq1
BE3wMBIGA1UdEwEB/wQIMAYBAf8CAQMwDgYDVR0PAQH/BAQDAgEGMA0GCSqGSIb3
DQEBBQUAA4IBAQCFDF2O5G9RaEIFoN27TyclhAO992T9Ldcw46QQF+vaKSm2eT92
9hkTI7gQCvlYpNRhcL0EYWoSihfVCr3FvDB81ukMJY2GQE/szKN+OMY3EU/t3Wgx
jkzSswF07r51XgdIGn9w/xZchMB5hbgF/X++ZRGjD8ACtPhSNzkE1akxehi/oCr0
Epn3o0WC4zxe9Z2etciefC7IpJ5OCBRLbf1wbWsaY71k5h+3zvDyny67G7fyUIhz
ksLi4xaNmjICq44Y3ekQEe5+NauQrz4wlHrQMz2nZQ/1/I6eYs9HRCwBXbsdtTLS
R9I4LtD+gdwyah617jzV/OeBHRnDJELqYzmp
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFWjCCBEKgAwIBAgIQDxSWXyAgaZlP1ceseIlB4jANBgkqhkiG9w0BAQsFADBa
MQswCQYDVQQGEwJJRTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJl
clRydXN0MSIwIAYDVQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTIw
MDcyMTIzMDAwMFoXDTI0MTAwODA3MDAwMFowTzELMAkGA1UEBhMCVVMxHjAcBgNV
BAoTFU1pY3Jvc29mdCBDb3Jwb3JhdGlvbjEgMB4GA1UEAxMXTWljcm9zb2Z0IFJT
QSBUTFMgQ0EgMDEwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCqYnfP
mmOyBoTzkDb0mfMUUavqlQo7Rgb9EUEf/lsGWMk4bgj8T0RIzTqk970eouKVuL5R
IMW/snBjXXgMQ8ApzWRJCZbar879BV8rKpHoAW4uGJssnNABf2n17j9TiFy6BWy+
IhVnFILyLNK+W2M3zK9gheiWa2uACKhuvgCca5Vw/OQYErEdG7LBEzFnMzTmJcli
W1iCdXby/vI/OxbfqkKD4zJtm45DJvC9Dh+hpzqvLMiK5uo/+aXSJY+SqhoIEpz+
rErHw+uAlKuHFtEjSeeku8eR3+Z5ND9BSqc6JtLqb0bjOHPm5dSRrgt4nnil75bj
c9j3lWXpBb9PXP9Sp/nPCK+nTQmZwHGjUnqlO9ebAVQD47ZisFonnDAmjrZNVqEX
F3p7laEHrFMxttYuD81BdOzxAbL9Rb/8MeFGQjE2Qx65qgVfhH+RsYuuD9dUw/3w
ZAhq05yO6nk07AM9c+AbNtRoEcdZcLCHfMDcbkXKNs5DJncCqXAN6LhXVERCw/us
G2MmCMLSIx9/kwt8bwhUmitOXc6fpT7SmFvRAtvxg84wUkg4Y/Gx++0j0z6StSeN
0EJz150jaHG6WV4HUqaWTb98Tm90IgXAU4AW2GBOlzFPiU5IY9jt+eXC2Q6yC/Zp
TL1LAcnL3Qa/OgLrHN0wiw1KFGD51WRPQ0Sh7QIDAQABo4IBJTCCASEwHQYDVR0O
BBYEFLV2DDARzseSQk1Mx1wsyKkM6AtkMB8GA1UdIwQYMBaAFOWdWTCCR1jMrPoI
VDaGezq1BE3wMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYI
KwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYI
KwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTA6BgNVHR8EMzAxMC+g
LaArhilodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vT21uaXJvb3QyMDI1LmNybDAq
BgNVHSAEIzAhMAgGBmeBDAECATAIBgZngQwBAgIwCwYJKwYBBAGCNyoBMA0GCSqG
SIb3DQEBCwUAA4IBAQCfK76SZ1vae4qt6P+dTQUO7bYNFUHR5hXcA2D59CJWnEj5
na7aKzyowKvQupW4yMH9fGNxtsh6iJswRqOOfZYC4/giBO/gNsBvwr8uDW7t1nYo
DYGHPpvnpxCM2mYfQFHq576/TmeYu1RZY29C4w8xYBlkAA8mDJfRhMCmehk7cN5F
JtyWRj2cZj/hOoI45TYDBChXpOlLZKIYiG1giY16vhCRi6zmPzEwv+tk156N6cGS
Vm44jTQ/rs1sa0JSYjzUaYngoFdZC4OfxnIkQvUIA4TOFmPzNPEFdjcZsgbeEz4T
cGHTBPK4R28F44qIMCtHRV55VMX53ev6P3hRddJb
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIFWjCCBEKgAwIBAgIQD6dHIsU9iMgPWJ77H51KOjANBgkqhkiG9w0BAQsFADBa
MQswCQYDVQQGEwJJRTESMBAGA1UEChMJQmFsdGltb3JlMRMwEQYDVQQLEwpDeWJl
clRydXN0MSIwIAYDVQQDExlCYWx0aW1vcmUgQ3liZXJUcnVzdCBSb290MB4XDTIw
MDcyMTIzMDAwMFoXDTI0MTAwODA3MDAwMFowTzELMAkGA1UEBhMCVVMxHjAcBgNV
BAoTFU1pY3Jvc29mdCBDb3Jwb3JhdGlvbjEgMB4GA1UEAxMXTWljcm9zb2Z0IFJT
QSBUTFMgQ0EgMDIwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQD0wBlZ
qiokfAYhMdHuEvWBapTj9tFKL+NdsS4pFDi8zJVdKQfR+F039CDXtD9YOnqS7o88
+isKcgOeQNTri472mPnn8N3vPCX0bDOEVk+nkZNIBA3zApvGGg/40Thv78kAlxib
MipsKahdbuoHByOB4ZlYotcBhf/ObUf65kCRfXMRQqOKWkZLkilPPn3zkYM5GHxe
I4MNZ1SoKBEoHa2E/uDwBQVxadY4SRZWFxMd7ARyI4Cz1ik4N2Z6ALD3MfjAgEED
woknyw9TGvr4PubAZdqU511zNLBoavar2OAVTl0Tddj+RAhbnX1/zypqk+ifv+d3
CgiDa8Mbvo1u2Q8nuUBrKVUmR6EjkV/dDrIsUaU643v/Wp/uE7xLDdhC5rplK9si
NlYohMTMKLAkjxVeWBWbQj7REickISpc+yowi3yUrO5lCgNAKrCNYw+wAfAvhFkO
eqPm6kP41IHVXVtGNC/UogcdiKUiR/N59IfYB+o2v54GMW+ubSC3BohLFbho/oZZ
5XyulIZK75pwTHmauCIeE5clU9ivpLwPTx9b0Vno9+ApElrFgdY0/YKZ46GfjOC9
ta4G25VJ1WKsMmWLtzyrfgwbYopquZd724fFdpvsxfIvMG5m3VFkThOqzsOttDcU
fyMTqM2pan4txG58uxNJ0MjR03UCEULRU+qMnwIDAQABo4IBJTCCASEwHQYDVR0O
BBYEFP8vf+EG9DjzLe0ljZjC/g72bPz6MB8GA1UdIwQYMBaAFOWdWTCCR1jMrPoI
VDaGezq1BE3wMA4GA1UdDwEB/wQEAwIBhjAdBgNVHSUEFjAUBggrBgEFBQcDAQYI
KwYBBQUHAwIwEgYDVR0TAQH/BAgwBgEB/wIBADA0BggrBgEFBQcBAQQoMCYwJAYI
KwYBBQUHMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTA6BgNVHR8EMzAxMC+g
LaArhilodHRwOi8vY3JsMy5kaWdpY2VydC5jb20vT21uaXJvb3QyMDI1LmNybDAq
BgNVHSAEIzAhMAgGBmeBDAECATAIBgZngQwBAgIwCwYJKwYBBAGCNyoBMA0GCSqG
SIb3DQEBCwUAA4IBAQCg2d165dQ1tHS0IN83uOi4S5heLhsx+zXIOwtxnvwCWdOJ
3wFLQaFDcgaMtN79UjMIFVIUedDZBsvalKnx+6l2tM/VH4YAyNPx+u1LFR0joPYp
QYLbNYkedkNuhRmEBesPqj4aDz68ZDI6fJ92sj2q18QvJUJ5Qz728AvtFOat+Ajg
K0PFqPYEAviUKr162NB1XZJxf6uyIjUlnG4UEdHfUqdhl0R84mMtrYINksTzQ2sH
YM8fEhqICtTlcRLr/FErUaPUe9648nziSnA0qKH7rUZqP/Ifmbo+WNZSZG1BbgOh
lk+521W+Ncih3HRbvRBE0LWYT8vWKnfjgZKxwHwJ
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQCq+mxcpjxFFB6jvh98dTFzANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwMTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAMedcDrkXufP7pxVm1FHLDNA9IjwHaMoaY8arqqZ4Gff4xyr
RygnavXL7g12MPAx8Q6Dd9hfBzrfWxkF0Br2wIvlvkzW01naNVSkHp+OS3hL3W6n
l/jYvZnVeJXjtsKYcXIf/6WtspcF5awlQ9LZJcjwaH7KoZuK+THpXCMtzD8XNVdm
GW/JI0C/7U/E7evXn9XDio8SYkGSM63aLO5BtLCv092+1d4GGBSQYolRq+7Pd1kR
EkWBPm0ywZ2Vb8GIS5DLrjelEkBnKCyy3B0yQud9dpVsiUeE7F5sY8Me96WVxQcb
OyYdEY/j/9UpDlOG+vA+YgOvBhkKEjiqygVpP8EZoMMijephzg43b5Qi9r5UrvYo
o19oR/8pf4HJNDPF0/FJwFVMW8PmCBLGstin3NE1+NeWTkGt0TzpHjgKyfaDP2tO
4bCk1G7pP2kDFT7SYfc8xbgCkFQ2UCEXsaH/f5YmpLn4YPiNFCeeIida7xnfTvc4
7IxyVccHHq1FzGygOqemrxEETKh8hvDR6eBdrBwmCHVgZrnAqnn93JtGyPLi6+cj
WGVGtMZHwzVvX1HvSFG771sskcEjJxiQNQDQRWHEh3NxvNb7kFlAXnVdRkkvhjpR
GchFhTAzqmwltdWhWDEyCMKC2x/mSZvZtlZGY+g37Y72qHzidwtyW7rBetZJAgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUDyBd16FXlduSzyvQx8J3BM5ygHYwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAlFvNh7QgXVLAZSsNR2XRmIn9iS8OHFCBA
WxKJoi8YYQafpMTkMqeuzoL3HWb1pYEipsDkhiMnrpfeYZEA7Lz7yqEEtfgHcEBs
K9KcStQGGZRfmWU07hPXHnFz+5gTXqzCE2PBMlRgVUYJiA25mJPXfB00gDvGhtYa
+mENwM9Bq1B9YYLyLjRtUz8cyGsdyTIG/bBM/Q9jcV8JGqMU/UjAdh1pFyTnnHEl
Y59Npi7F87ZqYYJEHJM2LGD+le8VsHjgeWX2CJQko7klXvcizuZvUEDTjHaQcs2J
+kPgfyMIOY1DMJ21NxOJ2xPRC/wAh/hzSBRVtoAnyuxtkZ4VjIOh
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQDGrpfM7VmYOGkKAKnqUyFDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwMjCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAOBiO1K6Fk4fHI6t3mJkpg7lxoeUgL8tz9wuI2z0UgY8vFra
3VBo7QznC4K3s9jqKWEyIQY11Le0108bSYa/TK0aioO6itpGiigEG+vH/iqtQXPS
u6D804ri0NFZ1SOP9IzjYuQiK6AWntCqP4WAcZAPtpNrNLPBIyiqmiTDS4dlFg1d
skMuVpT4z0MpgEMmxQnrSZ615rBQ25vnVbBNig04FCsh1V3S8ve5Gzh08oIrL/g5
xq95oRrgEeOBIeiegQpoKrLYyo3R1Tt48HmSJCBYQ52Qc34RgxQdZsLXMUrWuL1J
LAZP6yeo47ySSxKCjhq5/AUWvQBP3N/cP/iJzKKKw23qJ/kkVrE0DSVDiIiXWF0c
9abSGhYl9SPl86IHcIAIzwelJ4SKpHrVbh0/w4YHdFi5QbdAp7O5KxfxBYhQOeHy
is01zkpYn6SqUFGvbK8eZ8y9Aclt8PIUftMG6q5BhdlBZkDDV3n70RlXwYvllzfZ
/nV94l+hYp+GLW7jSmpxZLG/XEz4OXtTtWwLV+IkIOe/EDF79KCazW2SXOIvVInP
oi1PqN4TudNv0GyBF5tRC/aBjUqply1YYfeKwgRVs83z5kuiOicmdGZKH9SqU5bn
Kse7IlyfZLg6yAxYyTNe7A9acJ3/pGmCIkJ/9dfLUFc4hYb3YyIIYGmqm2/3AgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUAKuR/CFiJpeaqHkbYUGQYKliZ/0wHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAzo/KdmWPPTaYLQW7J5DqxEiBT9QyYGUfe
Zd7TR1837H6DSkFa/mGM1kLwi5y9miZKA9k6T9OwTx8CflcvbNO2UkFW0VCldEGH
iyx5421+HpRxMQIRjligePtOtRGXwaNOQ7ySWfJhRhKcPKe2PGFHQI7/3n+T3kXQ
/SLu2lk9Qs5YgSJ3VhxBUznYn1KVKJWPE07M55kuUgCquAV0PksZj7EC4nK6e/UV
bPumlj1nyjlxhvNud4WYmr4ntbBev6cSbK78dpI/3cr7P/WJPYJuL0EsO3MgjS3e
DCX7NXp5ylue3TcpQfRU8BL+yZC1wqX98R4ndw7X4qfGaE7SlF7I
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQDXvt6X2CCZZ6UmMbi90YvTANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwNTCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBAKplDTmQ9afwVPQelDuu+NkxNJ084CNKnrZ21ABewE+UU4GK
DnwygZdK6agNSMs5UochUEDzz9CpdV5tdPzL14O/GeE2gO5/aUFTUMG9c6neyxk5
tq1WdKsPkitPws6V8MWa5d1L/y4RFhZHUsgxxUySlYlGpNcHhhsyr7EvFecZGA1M
fsitAWVp6hiWANkWKINfRcdt3Z2A23hmMH9MRSGBccHiPuzwrVsSmLwvt3WlRDgO
bJkE40tFYvJ6GXAQiaGHCIWSVObgO3zj6xkdbEFMmJ/zr2Wet5KEcUDtUBhA4dUU
oaPVz69u46V56Vscy3lXu1Ylsk84j5lUPLdsAxtultP4OPQoOTpnY8kxWkH6kgO5
gTKE3HRvoVIjU4xJ0JQ746zy/8GdQA36SaNiz4U3u10zFZg2Rkv2dL1Lv58EXL02
r5q5B/nhVH/M1joTvpRvaeEpAJhkIA9NkpvbGEpSdcA0OrtOOeGtrsiOyMBYkjpB
5nw0cJY1QHOr3nIvJ2OnY+OKJbDSrhFqWsk8/1q6Z1WNvONz7te1pAtHerdPi5pC
HeiXCNpv+fadwP0k8czaf2Vs19nYsgWn5uIyLQL8EehdBzCbOKJy9sl86S4Fqe4H
GyAtmqGlaWOsq2A6O/paMi3BSmWTDbgPLCPBbPte/bsuAEF4ajkPEES3GHP9AgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQUx7KcfxzjuFrv6WgaqF2UwSZSamgwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQAe+G+G2RFdWtYxLIKMR5H/aVNFjNP7Jdeu
+oZaKaIu7U3NidykFr994jSxMBMV768ukJ5/hLSKsuj/SLjmAfwRAZ+w0RGqi/kO
vPYUlBr/sKOwr3tVkg9ccZBebnBVG+DLKTp2Ox0+jYBCPxla5FO252qpk7/6wt8S
Zk3diSU12Jm7if/jjkhkGB/e8UdfrKoLytDvqVeiwPA5FPzqKoSqN75byLjsIKJE
dNi07SY45hN/RUnsmIoAf93qlaHR/SJWVRhrWt3JmeoBJ2RDK492zF6TGu1moh4a
E6e00YkwTPWreuwvaLB220vWmtgZPs+DSIb2d9hPBdCJgvcho1c7
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIF8zCCBNugAwIBAgIQAueRcfuAIek/4tmDg0xQwDANBgkqhkiG9w0BAQwFADBh
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH
MjAeFw0yMDA3MjkxMjMwMDBaFw0yNDA2MjcyMzU5NTlaMFkxCzAJBgNVBAYTAlVT
MR4wHAYDVQQKExVNaWNyb3NvZnQgQ29ycG9yYXRpb24xKjAoBgNVBAMTIU1pY3Jv
c29mdCBBenVyZSBUTFMgSXNzdWluZyBDQSAwNjCCAiIwDQYJKoZIhvcNAQEBBQAD
ggIPADCCAgoCggIBALVGARl56bx3KBUSGuPc4H5uoNFkFH4e7pvTCxRi4j/+z+Xb
wjEz+5CipDOqjx9/jWjskL5dk7PaQkzItidsAAnDCW1leZBOIi68Lff1bjTeZgMY
iwdRd3Y39b/lcGpiuP2d23W95YHkMMT8IlWosYIX0f4kYb62rphyfnAjYb/4Od99
ThnhlAxGtfvSbXcBVIKCYfZgqRvV+5lReUnd1aNjRYVzPOoifgSx2fRyy1+pO1Uz
aMMNnIOE71bVYW0A1hr19w7kOb0KkJXoALTDDj1ukUEDqQuBfBxReL5mXiu1O7WG
0vltg0VZ/SZzctBsdBlx1BkmWYBW261KZgBivrql5ELTKKd8qgtHcLQA5fl6JB0Q
gs5XDaWehN86Gps5JW8ArjGtjcWAIP+X8CQaWfaCnuRm6Bk/03PQWhgdi84qwA0s
sRfFJwHUPTNSnE8EiGVk2frt0u8PG1pwSQsFuNJfcYIHEv1vOzP7uEOuDydsmCjh
lxuoK2n5/2aVR3BMTu+p4+gl8alXoBycyLmj3J/PUgqD8SL5fTCUegGsdia/Sa60
N2oV7vQ17wjMN+LXa2rjj/b4ZlZgXVojDmAjDwIRdDUujQu0RVsJqFLMzSIHpp2C
Zp7mIoLrySay2YYBu7SiNwL95X6He2kS8eefBBHjzwW/9FxGqry57i71c2cDAgMB
AAGjggGtMIIBqTAdBgNVHQ4EFgQU1cFnOsKjnfR3UltZEjgp5lVou6UwHwYDVR0j
BBgwFoAUTiJUIBiV5uNu5g/6+rkS7QYXjzkwDgYDVR0PAQH/BAQDAgGGMB0GA1Ud
JQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjASBgNVHRMBAf8ECDAGAQH/AgEAMHYG
CCsGAQUFBwEBBGowaDAkBggrBgEFBQcwAYYYaHR0cDovL29jc3AuZGlnaWNlcnQu
Y29tMEAGCCsGAQUFBzAChjRodHRwOi8vY2FjZXJ0cy5kaWdpY2VydC5jb20vRGln
aUNlcnRHbG9iYWxSb290RzIuY3J0MHsGA1UdHwR0MHIwN6A1oDOGMWh0dHA6Ly9j
cmwzLmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5jcmwwN6A1oDOG
MWh0dHA6Ly9jcmw0LmRpZ2ljZXJ0LmNvbS9EaWdpQ2VydEdsb2JhbFJvb3RHMi5j
cmwwHQYDVR0gBBYwFDAIBgZngQwBAgEwCAYGZ4EMAQICMBAGCSsGAQQBgjcVAQQD
AgEAMA0GCSqGSIb3DQEBDAUAA4IBAQB2oWc93fB8esci/8esixj++N22meiGDjgF
+rA2LUK5IOQOgcUSTGKSqF9lYfAxPjrqPjDCUPHCURv+26ad5P/BYtXtbmtxJWu+
cS5BhMDPPeG3oPZwXRHBJFAkY4O4AF7RIAAUW6EzDflUoDHKv83zOiPfYGcpHc9s
kxAInCedk7QSgXvMARjjOqdakor21DTmNIUotxo8kHv5hwRlGhBJwps6fEVi1Bt0
trpM/3wYxlr473WSPUFZPgP1j519kLpWOJ8z09wxay+Br29irPcBYv0GMXlHqThy
8y4m/HyTQeI2IMvMrQnwqPpY+rLIXyviI2vLoI+4xKE4Rn38ZZ8m
-----END CERTIFICATE-----
`)

func errorf(format string, v ...interface{}) error {
	return fmt.Errorf("iotservice: "+format, v...)
}

// Message is a common message format for all device-facing protocols.
// This message format is used for both device-to-cloud and cloud-to-device messages.
// See: https://docs.microsoft.com/en-us/azure/iot-hub/iot-hub-devguide-messages-construct
type Message struct {
	// MessageID is a user-settable identifier for the message used for request-reply patterns.
	MessageID string `json:"MessageId,omitempty"`

	// To is a destination specified in cloud-to-device messages.
	To string `json:"To,omitempty"`

	// ExpiryTime is time of message expiration.
	ExpiryTime *time.Time `json:"ExpiryTimeUtc,omitempty"`

	// EnqueuedTime is time the Cloud-to-Device message was received by IoT Hub.
	EnqueuedTime *time.Time `json:"EnqueuedTime,omitempty"`

	// CorrelationID is a string property in a response message that typically
	// contains the MessageId of the request, in request-reply patterns.
	CorrelationID string `json:"CorrelationId,omitempty"`

	// UserID is an ID used to specify the origin of messages.
	UserID string `json:"UserId,omitempty"`

	// ConnectionDeviceID is an ID set by IoT Hub on device-to-cloud messages.
	// It contains the deviceId of the device that sent the message.
	ConnectionDeviceID string `json:"ConnectionDeviceId,omitempty"`

	// ConnectionDeviceGenerationID is an ID set by IoT Hub on device-to-cloud messages.
	// It contains the generationId (as per Device identity properties)
	// of the device that sent the message.
	ConnectionDeviceGenerationID string `json:"ConnectionDeviceGenerationId,omitempty"`

	// ConnectionAuthMethod is an authentication method set by IoT Hub on
	// device-to-cloud messages. This property contains information about
	// the authentication method used to authenticate the device sending the message.
	ConnectionAuthMethod *ConnectionAuthMethod `json:"ConnectionAuthMethod,omitempty"`

	// MessageSource determines a device-to-cloud message transport.
	MessageSource string `json:"MessageSource,omitempty"`

	// Payload is message data.
	Payload []byte `json:"Payload,omitempty"`

	// Properties are custom message properties (property bags).
	Properties map[string]string `json:"Properties,omitempty"`

	// TransportOptions transport specific options.
	TransportOptions map[string]interface{} `json:"-"`
}

// ConnectionAuthMethod is an authentication method of device-to-cloud communication.
type ConnectionAuthMethod struct {
	Scope  string `json:"scope"`
	Type   string `json:"type"`
	Issuer string `json:"issuer"`
}

// ====================================================================================
// SharedAccessKey is SAS token generator.
type SharedAccessKey struct {
	HostName            string
	SharedAccessKeyName string
	SharedAccessKey     string
}

// Token generates a shared access signature for the named resource and lifetime.
func (c *SharedAccessKey) Token(
	resource string, lifetime time.Duration,
) (*SharedAccessSignature, error) {
	return NewSharedAccessSignature(
		resource, c.SharedAccessKeyName, c.SharedAccessKey, time.Now().Add(lifetime),
	)
}

// NewSharedAccessSignature initialized a new shared access signature
// and generates signature fields based on the given input.
func NewSharedAccessSignature(
	resource, policy, key string, expiry time.Time,
) (*SharedAccessSignature, error) {
	sig, err := mksig(resource, key, expiry)
	if err != nil {
		return nil, err
	}
	return &SharedAccessSignature{
		Sr:  resource,
		Sig: sig,
		Se:  expiry,
		Skn: policy,
	}, nil
}

func mksig(sr, key string, se time.Time) (string, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}
	h := hmac.New(sha256.New, b)
	if _, err := fmt.Fprintf(h, "%s\n%d", url.QueryEscape(sr), se.Unix()); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

// SharedAccessSignature is a shared access signature instance.
type SharedAccessSignature struct {
	Sr  string
	Sig string
	Se  time.Time
	Skn string
}

// String converts the signature to a token string.
func (sas *SharedAccessSignature) String() string {
	s := "SharedAccessSignature " +
		"sr=" + url.QueryEscape(sas.Sr) +
		"&sig=" + url.QueryEscape(sas.Sig) +
		"&se=" + url.QueryEscape(strconv.FormatInt(sas.Se.Unix(), 10))
	if sas.Skn != "" {
		s += "&skn=" + url.QueryEscape(sas.Skn)
	}
	return s
}

// ====================================================================================

// rootCAs root CA certificates pool for connecting to the cloud.
func rootCAs() *x509.CertPool {
	p := x509.NewCertPool()
	if ok := p.AppendCertsFromPEM(caCerts); !ok {
		panic("tls: unable to append certificates")
	}
	return p
}

// EventHandler handles incoming cloud-to-device events.
//type EventHandler func(e *Event) error
type EventHandler func(e *amqp.Message) error

// Event is a device-to-cloud message.
type Event struct {
	*Message
}

const userAgent = "iothub-golang-sdk/dev"

// putTokenContinuously writes token first time in blocking mode and returns
// maintaining token updates in the background until the client is closed.
func putTokenContinuously(ctx context.Context, conn *amqp.Client) error {
	const (
		tokenUpdateInterval = time.Hour

		// we need to update tokens before they expire to prevent disconnects
		// from azure, without interrupting the message flow
		tokenUpdateSpan = 10 * time.Minute
	)

	sess, err := conn.NewSession()
	if err != nil {
		return err
	}

	if err := putToken(ctx, sess, tokenUpdateInterval); err != nil {
		_ = sess.Close(context.Background())
		return err
	}

	var done chan struct{}
	go func() {
		defer sess.Close(context.Background())
		ticker := time.NewTimer(tokenUpdateInterval - tokenUpdateSpan)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := putToken(context.Background(), sess, tokenUpdateInterval); err != nil {
					log.Printf("put token error: %v\n", err)
					return
				}
				ticker.Reset(tokenUpdateInterval - tokenUpdateSpan)
				log.Println("token updated")
			case <-done:
				return
			}
		}
	}()
	return nil
}

func putToken(
	ctx context.Context, sess *amqp.Session, lifetime time.Duration,
) error {
	send, err := sess.NewSender(
		amqp.LinkTargetAddress("$cbs"),
	)
	if err != nil {
		return err
	}
	defer send.Close(context.Background())

	recv, err := sess.NewReceiver(
		amqp.LinkSourceAddress("$cbs"),
	)
	if err != nil {
		return err
	}
	defer recv.Close(context.Background())

	// Seb: https://docs.microsoft.com/en-us/rest/api/eventhub/generate-sas-token
	//      See the NodeJS version
	sak := &SharedAccessKey{
		HostName:            "seb-hub.azure-devices.net",
		SharedAccessKeyName: "service",
		SharedAccessKey:     "sTuhHHftCxGl35I0XtLV4B4GwQg8Di0t5nPjXYg6CqA=",
	}
	sas, err := sak.Token(sak.HostName, lifetime)
	if err != nil {
		log.Printf("putToken: %v\n", err)
		return err
	}
	log.Printf("putToken generated: %s\n", sas.String())

	if err = send.Send(ctx, &amqp.Message{
		Value: sas.String(),
		Properties: &amqp.MessageProperties{
			To:      "$cbs",
			ReplyTo: "cbs",
		},
		ApplicationProperties: map[string]interface{}{
			"operation": "put-token",
			"type":      "servicebus.windows.net:sastoken",
			"name":      sak.HostName,
		},
	}); err != nil {
		log.Printf("putToken send error: %v\n", err)
		return err
	}

	msg, err := recv.Receive(ctx)
	if err != nil {
		log.Printf("putToken recv error: %v\n", err)
		return err
	}
	if err = msg.Accept(ctx); err != nil {
		log.Printf("putToken Accept error: %v\n", err)
		return err
	}
	return checkMessageResponse(msg)
}

// newSession connects to IoT Hub's AMQP broker,
// it's needed for sending C2S events and subscribing to events feedback.
//
// It establishes connection only once, subsequent calls return immediately.
func newSession(ctx context.Context, c *amqp.Client, tc *tls.Config) (*amqp.Session, error) {
	//if c != nil {
	//	return c.NewSession() // already connected
	//}
	sakHost := "seb-hub.azure-devices.net"
	conn, err := amqp.Dial("amqps://"+sakHost, //c.sak.HostName,
		amqp.ConnTLSConfig(tc),
		amqp.ConnProperty("com.microsoft:client-version", userAgent),
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	log.Printf("newSession connected to %s", sakHost)

	if err = putTokenContinuously(ctx, conn); err != nil {
		return nil, err
	}

	sess, err := conn.NewSession()
	if err != nil {
		return nil, err
	}
	//.conn = conn
	return sess, nil
}

func connectToEventHub(c *amqp.Client, ctx context.Context, tc *tls.Config) (*amqp.Client, error) {
	sess, err := newSession(ctx, c, tc)
	if err != nil {
		return nil, err
	}

	// iothub broker should redirect us to an eventhub compatible instance
	// straight after subscribing to events stream, for that we need to connect twice
	defer sess.Close(context.Background())

	_, err = sess.NewReceiver(
		amqp.LinkSourceAddress("messages/events/"),
	)
	if err == nil {
		return nil, errorf("expected redirect error")
	}
	rerr, ok := err.(*amqp.Error)

	/*
		2021/04/05 14:49:17 connectToEventHub *Error{Condition: amqp:link:redirect,
			Description: , Info: map[address:amqps://ihsuprodsgres013dednamespace.servicebus.windows.net:5671/iothub-ehub-seb-hub-8717893-e68ae183bb/ hostname:ihsuprodsgres013dednamespace.servicebus.windows.net network-host:ihsuprodsgres013dednamespace.servicebus.windows.net port:5671]}
	*/
	if !ok || rerr.Condition != amqp.ErrorLinkRedirect {
		log.Print("connectToEventHub Error:", err)
		return nil, err
	}
	log.Printf("connectToEventHub error is expected: %v", err.(*amqp.Error))

	// "amqps://{host}:5671/{consumerGroup}/"
	group := rerr.Info["address"].(string)
	group = group[strings.Index(group, ":5671/")+6 : len(group)-1]

	host := rerr.Info["hostname"].(string)
	//c.logger.Debugf("redirected to %s:%s eventhub", host, group)

	// Seb
	log.Printf("connectToEventHub redirected to %s:%s iot eventhub", host, group)

	tlsCloned := tc.Clone()
	tlsCloned.ServerName = "ihsuprodsgres013dednamespace.servicebus.windows.net"
	log.Printf("connectToEventHub tls ServerName: %s\n", tlsCloned.ServerName)
	/*
		amqpclient, err := amqp.Dial("amqps://"+"ihsuprodsgres013dednamespace.servicebus.windows.net",
			amqp.ConnTLSConfig(tlsCloned),
			amqp.ConnSASLPlain("service", "sTuhHHftCxGl35I0XtLV4B4GwQg8Di0t5nPjXYg6CqA="),
			amqp.ConnOption(amqp.ConnProperty("com.microsoft:client-version", userAgent)),
		)
	*/
	//log.Printf("connectToEventHub Dial sak: %+v\n", c.sak)
	log.Printf("connectToEventHub group: %s\n", group)
	amqpclient, err := dial("ihsuprodsgres013dednamespace.servicebus.windows.net", group,
		WithTLSConfig(tlsCloned),
		WithSASLPlain("service", "sTuhHHftCxGl35I0XtLV4B4GwQg8Di0t5nPjXYg6CqA="),
		WithConnOption(amqp.ConnProperty("com.microsoft:client-version", userAgent)),
	)
	if err != nil {
		log.Print("connectToEventHub amqp Dial Error:", err)
		return nil, err
	}

	// Seb
	log.Printf("eventHub host: %s\n", host)
	log.Printf("eventHub group: %s\n", group)
	//log.Printf("eventHub SharedAccessKeyName: %s\n", c.sak.SharedAccessKeyName)
	//log.Printf("eventHub SharedAccessKey: %s\n", c.sak.SharedAccessKey)

	return amqpclient, nil
}

// Option is a client configuration option.
type Option func(c *Client)

// Client is an EventHub client.
type Client struct {
	name string
	conn *amqp.Client
	opts []amqp.ConnOption
}

// WithTLSConfig sets connection TLS configuration.
func WithTLSConfig(tc *tls.Config) Option {
	return WithConnOption(amqp.ConnTLSConfig(tc))
}

// WithSASLPlain configures connection username and password.
func WithSASLPlain(username, password string) Option {
	return WithConnOption(amqp.ConnSASLPlain(username, password))
}

// WithConnOption sets a low-level connection option.
func WithConnOption(opt amqp.ConnOption) Option {
	return func(c *Client) {
		c.opts = append(c.opts, opt)
	}
}

// dial connects to the named EventHub and returns a client instance.
func dial(host, name string, opts ...Option) (*amqp.Client, error) {
	c := &Client{name: name}
	log.Printf("Dial options: %v\n", c)
	for _, opt := range opts {
		opt(c)
	}
	// Seb
	log.Printf("Dial Client: %+v\n", c)
	for _, opt := range c.opts {
		log.Printf("Dial options: %+v\n", opt)
	}

	var err error
	c.conn, err = amqp.Dial("amqps://"+host, c.opts...)
	if err != nil {
		log.Print("dial Error:", err)
		return nil, err
	}
	return c.conn, nil
}

// fromAMQPMessage converts a amqp.Message into common.Message.
//
// Exported to use with a custom stream when devices telemetry is
// routed for example to an EventhHub instance.
func fromAMQPMessage(msg *amqp.Message) *Message {
	m := &Message{
		Payload:    msg.GetData(),
		Properties: make(map[string]string, len(msg.ApplicationProperties)+5),
	}
	if msg.Properties != nil {
		m.UserID = string(msg.Properties.UserID)
		if msg.Properties.MessageID != nil {
			m.MessageID = msg.Properties.MessageID.(string)
		}
		if msg.Properties.CorrelationID != nil {
			m.CorrelationID = msg.Properties.CorrelationID.(string)
		}
		m.To = msg.Properties.To
		m.ExpiryTime = &msg.Properties.AbsoluteExpiryTime
	}
	for k, v := range msg.Annotations {
		switch k {
		case "iothub-enqueuedtime":
			t, _ := v.(time.Time)
			m.EnqueuedTime = &t
		case "iothub-connection-device-id":
			m.ConnectionDeviceID = v.(string)
		case "iothub-connection-auth-generation-id":
			m.ConnectionDeviceGenerationID = v.(string)
		case "iothub-connection-auth-method":
			var am ConnectionAuthMethod
			if err := json.Unmarshal([]byte(v.(string)), &am); err != nil {
				m.Properties[k.(string)] = fmt.Sprint(v)
				continue
			}
			m.ConnectionAuthMethod = &am
		case "iothub-message-source":
			m.MessageSource = v.(string)
		default:
			m.Properties[k.(string)] = fmt.Sprint(v)
		}
	}
	for k, v := range msg.ApplicationProperties {
		if v, ok := v.(string); ok {
			m.Properties[k] = v
		} else {
			m.Properties[k] = ""
		}
	}
	return m
}

func genID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// checkMessageResponse checks for 200 response code otherwise returns an error.
func checkMessageResponse(msg *amqp.Message) error {
	rc, ok := msg.ApplicationProperties["status-code"].(int32)
	if !ok {
		return errors.New("unable to typecast status-code")
	}
	if rc == 200 {
		return nil
	}
	rd, _ := msg.ApplicationProperties["status-description"].(string)
	return fmt.Errorf("code = %d, description = %q", rc, rd)
}

// getPartitionIDs returns partition ids of the hub.
func getPartitionIDs(ctx context.Context, sess *amqp.Session) ([]string, error) {
	replyTo := genID()
	recv, err := sess.NewReceiver(
		amqp.LinkSourceAddress("$management"),
		amqp.LinkTargetAddress(replyTo),
	)
	if err != nil {
		return nil, err
	}
	defer recv.Close(context.Background())

	send, err := sess.NewSender(
		amqp.LinkTargetAddress("$management"),
		amqp.LinkSourceAddress(replyTo),
	)
	if err != nil {
		return nil, err
	}
	defer send.Close(context.Background())

	mid := genID()
	if err := send.Send(ctx, &amqp.Message{
		Properties: &amqp.MessageProperties{
			MessageID: mid,
			ReplyTo:   replyTo,
		},
		ApplicationProperties: map[string]interface{}{
			"operation": "READ",
			"name":      "iothub-ehub-seb-hub-8717893-e68ae183bb", //c.name
			"type":      "com.microsoft:eventhub",
		},
	}); err != nil {
		return nil, err
	}

	msg, err := recv.Receive(ctx)
	if err != nil {
		return nil, err
	}
	if err = checkMessageResponse(msg); err != nil {
		return nil, err
	}
	if msg.Properties.CorrelationID != mid {
		return nil, errors.New("message-id mismatch")
	}
	if err := msg.Accept(ctx); err != nil {
		return nil, err
	}

	val, ok := msg.Value.(map[string]interface{})
	if !ok {
		return nil, errors.New("unable to typecast value")
	}
	ids, ok := val["partition_ids"].([]string)
	if !ok {
		return nil, errors.New("unable to typecast partition_ids")
	}
	return ids, nil
}

type sub struct {
	group string
	opts  []amqp.LinkOption
}

// SubscribeOption is a Subscribe option.
type SubscribeOption func(r *sub)

// WithSubscribeSince requests events that occurred after the given time.
func WithSubscribeSince(t time.Time) SubscribeOption {
	return WithSubscribeLinkOption(amqp.LinkSelectorFilter(
		fmt.Sprintf("amqp.annotation.x-opt-enqueuedtimeutc > '%d'",
			t.UnixNano()/int64(time.Millisecond)),
	))
}

// WithSubscribeLinkOption is a low-level subscription configuration option.
func WithSubscribeLinkOption(opt amqp.LinkOption) SubscribeOption {
	return func(s *sub) {
		s.opts = append(s.opts, opt)
	}
}

// subscribe subscribes to all hub's partitions and registers the given
// handler and blocks until it encounters an error or the context is cancelled.
//
// It's client's responsibility to accept/reject/release events.
func subscribe(
	c *amqp.Client,
	ctx context.Context,
	fn func(evtmsg *amqp.Message) error,
	opts ...SubscribeOption,
) error {
	var s sub
	for _, opt := range opts {
		opt(&s)
	}
	if s.group == "" {
		s.group = "$Default"
	}

	// initialize new session for each subscribe session
	sess, err := c.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close(context.Background())

	ids, err := getPartitionIDs(ctx, sess)
	if err != nil {
		return err
	}

	// stop all goroutines at return
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	msgc := make(chan *amqp.Message)
	errc := make(chan error)

	for _, id := range ids {
		addr := fmt.Sprintf("/%s/ConsumerGroups/%s/Partitions/%s", "iothub-ehub-seb-hub-8717893-e68ae183bb", s.group, id)
		// Seb
		//addr := fmt.Sprintf("/%s/ConsumerGroups/%s/Partitions/%s", "iothub-ehub-seb-hub-8717893-e68ae183bb", "$Default", id)
		recv, err := sess.NewReceiver(
			append([]amqp.LinkOption{amqp.LinkSourceAddress(addr)}, s.opts...)...,
		)
		if err != nil {
			return err
		}

		go func(recv *amqp.Receiver) {
			defer recv.Close(context.Background())
			for {
				msg, err := recv.Receive(ctx)
				if err != nil {
					select {
					case errc <- err:
					case <-ctx.Done():
					}
					return
				}
				select {
				case msgc <- msg:
				case <-ctx.Done():
				}
			}
		}(recv)
	}

	for {
		select {
		case msg := <-msgc:
			if err := fn(msg); err != nil {
				return err
			}
		case err := <-errc:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// subscribeEvents subscribes to D2C events.
// Event handler is blocking, handle asynchronous processing on your own.
func subscribeEvents(c *amqp.Client, ctx context.Context, fn EventHandler, tc *tls.Config) error {
	// a new connection is established for every invocation,
	// this made on purpose because normally an app calls the method once
	eh, err := connectToEventHub(c, ctx, tc)
	if err != nil {
		log.Printf("subscribeEvents error: %v\n", err)
		return err
	}
	defer eh.Close()

	return subscribe(eh, ctx, func(msg *amqp.Message) error {
		//if err := fn(&Event{fromAMQPMessage(Message)}); err != nil {
		if err := fn(msg); err != nil {
			log.Printf("subscribeEvents subscribe error: %v\n", err)
			return err
		}
		return msg.Accept(ctx)
	},
		WithSubscribeSince(time.Now()),
	)

}

func main() {
	// Create client
	// From azure seb-hub|Settings|Shared access policies|service|Connection string-primary key

	hubname := "seb-hub.azure-devices.net"
	mytls := &tls.Config{RootCAs: rootCAs()}

	client, err := amqp.Dial("amqps://"+hubname,
		amqp.ConnTLSConfig(mytls),
		amqp.ConnProperty("com.microsoft:client-version", userAgent),
	)

	if err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}
	defer client.Close()

	log.Printf("connected to %s\n", hubname)

	ctx := context.Background()

	err = subscribeEvents(client, ctx, func(msg *amqp.Message) error {
		//fmt.Printf("%q sends %q\n", msg.ConnectionDeviceID, msg.Payload)
		fmt.Printf("Message received: %+v\n", msg)
		return nil
	}, mytls)

	/*
		subscribeEvents(client, ctx, func(msg *Event) error {
			fmt.Printf("%q sends %q\n", msg.ConnectionDeviceID, msg.Payload)
			return nil
		})
	*/
}
