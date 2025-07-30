#!/bin/bash
set -e

# 生成 CA 使用 SHA256 签名
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -sha256 -key ca.key -subj "/CN=operator-ca" -days 3650 -out ca.crt

# 生成服务私钥和 CSR
openssl genrsa -out tls.key 2048
openssl req -new -sha256 -key tls.key -out server.csr -config server-csr.conf

# 使用 CA 签发证书（含 SAN）
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
  -out tls.crt -days 365 -sha256 -extensions v3_req -extfile server-csr.conf

# 打印 base64 的 CA 证书（供 caBundle 用）
echo "CA cert (base64):"
cat ca.crt | base64 | tr -d '\n'
echo

# #!/bin/bash

# openssl genrsa -out ca.key 2048
# openssl req -x509 -new -nodes -key ca.key -subj "/CN=operator-ca" -days 3650 -out ca.crt

# openssl genrsa -out tls.key 2048
# openssl req -new -key tls.key -out server.csr -config server-csr.conf

# openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
#   -out tls.crt -days 365 -extensions v3_req -extfile server-csr.conf

# cat ca.crt | base64 | tr -d '\n'
