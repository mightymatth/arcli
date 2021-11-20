#!/bin/sh

mkdir -p certs

openssl req -x509 -out certs/server.crt -keyout certs/server.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=redmine.local' -extensions EXT -config <( \
   printf "[dn]\nCN=redmine.local\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:redmine.local\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
