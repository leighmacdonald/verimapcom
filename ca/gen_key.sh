#!/bin/bash
if test -z "$SSLHOST"
then
  echo "Please set SSLHOST to the common name/dns"
  exit 1
fi

echo "GENERATING KEY FOR $SSLHOST"
#openssl genrsa -aes256 -out intermediate/private/"$SSLHOST".key.pem 2048
openssl genrsa -out intermediate/private/"$SSLHOST".key.pem 2048
chmod 400 intermediate/private/"$SSLHOST".key.pem

# chmod 400 intermediate/private/www.example.com.key.pem
openssl req -config intermediate/openssl.cnf \
      -key intermediate/private/"$SSLHOST".key.pem \
      -new -sha256 -out intermediate/csr/"$SSLHOST".csr.pem

openssl ca -config intermediate/openssl.cnf \
      -extensions server_cert -days 375 -notext -md sha256 \
      -in intermediate/csr/"$SSLHOST".csr.pem \
      -out intermediate/certs/"$SSLHOST".cert.pem
chmod 444 intermediate/certs/"$SSLHOST".cert.pem

openssl pkey -in intermediate/private/"$SSLHOST".key.pem -out intermediate/private/"$SSLHOST".key
openssl x509 -outform der -in intermediate/certs/"$SSLHOST".cert.pem -out intermediate/certs/"$SSLHOST".crt
# Verify the certificate
openssl x509 -noout -text \
      -in intermediate/certs/"$SSLHOST".cert.pem
openssl verify -CAfile intermediate/certs/ca-chain.cert.pem \
      intermediate/certs/"$SSLHOST".cert.pem

# Deploy
#    ca-chain.cert.pem
#    www.example.com.key.pem
#    www.example.com.cert.pem
