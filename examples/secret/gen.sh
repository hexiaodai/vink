#!/bin/bash

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=operator-ca" -days 3650 -out ca.crt

openssl genrsa -out tls.key 2048
openssl req -new -key tls.key -out server.csr -config server-csr.conf

openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out tls.crt -days 365 -extensions v3_req -extfile server-csr.conf

cat ca.crt | base64 | tr -d '\n'
