---
title: "authelia crypto certificate rsa generate"
description: "Reference for the authelia crypto certificate rsa generate command."
lead: ""
date: 2022-06-27T18:27:57+10:00
draft: false
images: []
menu:
  reference:
    parent: "cli-authelia"
weight: 905
toc: true
---

## authelia crypto certificate rsa generate

Generate an RSA private key and certificate

### Synopsis

Generate an RSA private key and certificate.

This subcommand allows generating an RSA private key and certificate.

```
authelia crypto certificate rsa generate [flags]
```

### Examples

```
authelia crypto certificate rsa generate --help
```

### Options

```
  -b, --bits int                      number of RSA bits for the certificate (default 2048)
      --ca                            create the certificate as a certificate authority certificate
  -n, --common-name string            certificate common name
      --country strings               certificate country
  -d, --directory string              directory where the generated keys, certificates, etc will be stored
      --duration duration             duration of time the certificate is valid for (default 8760h0m0s)
      --extended-usage strings        specify the extended usage types of the certificate
      --file.ca-certificate string    certificate authority certificate to use when signing this certificate (default "ca.public.crt")
      --file.ca-private-key string    certificate authority private key to use to signing this certificate (default "ca.private.pem")
      --file.certificate string       name of the file to export the certificate data to (default "public.crt")
      --file.private-key string       name of the file to export the private key data to (default "private.pem")
  -h, --help                          help for generate
  -l, --locality strings              certificate locality
      --not-before string             earliest date and time the certificate is considered valid formatted as Jan 2 15:04:05 2006 (default is now)
  -o, --organization strings          certificate organization (default [Authelia])
      --organizational-unit strings   certificate organizational unit
      --path.ca string                source directory of the certificate authority files, if not provided the certificate will be self-signed
  -p, --postcode strings              certificate postcode
      --province strings              certificate province
      --sans strings                  subject alternative names
      --signature string              signature algorithm for the certificate (default "SHA256")
  -s, --street-address strings        certificate street address
```

### Options inherited from parent commands

```
  -c, --config strings                        configuration files or directories to load, for more information run 'authelia -h authelia config' (default [configuration.yml])
      --config.experimental.filters strings   list of filters to apply to all configuration files, for more information run 'authelia -h authelia filters'
```

### SEE ALSO

* [authelia crypto certificate rsa](authelia_crypto_certificate_rsa.md)	 - Perform RSA certificate cryptographic operations

