.PHONY: intermediate

# https://jamielinux.com/docs/openssl-certificate-authority/sign-server-and-client-certificates.html
# https://stackoverflow.com/questions/13732826/convert-pem-to-crt-and-key

# 0. Setup dirs
init:
	@mkdir -p certs crl newcerts private && chmod 700 private
	@mkdir -p intermediate
	@cd intermediate && mkdir -p certs crl csr newcerts private && chmod 700 private && touch index.txt && echo 1000 > serial
	@cd intermediate && echo 1000 > crlnumber
	@touch index.txt
	@echo 1000 > serial

# 1. Create the root key
ca:
	@openssl genrsa -aes256 -out private/ca.key.pem 4096
	@chmod 600 private/ca.key.pem

# 2. Create the root certificate
root:
	@openssl req -config openssl.cnf \
          -key private/ca.key.pem \
          -new -x509 -days 7300 -sha256 -extensions v3_ca \
          -out certs/ca.cert.pem
	@chmod 444 certs/ca.cert.pem

# 3. Create the intermediate pair
intermediate:
	@openssl genrsa -aes256 -out intermediate/private/intermediate.key.pem 4096
	@chmod 400 intermediate/private/intermediate.key.pem
	@openssl req -config intermediate/openssl.cnf -new -sha256 \
           -key intermediate/private/intermediate.key.pem \
           -out intermediate/csr/intermediate.csr.pem
	@openssl ca -config openssl.cnf -extensions v3_intermediate_ca \
           -days 3650 -notext -md sha256 \
           -in intermediate/csr/intermediate.csr.pem \
           -out intermediate/certs/intermediate.cert.pem
	@chmod 444 intermediate/certs/intermediate.cert.pem
	@openssl x509 -noout -text \
           -in intermediate/certs/intermediate.cert.pem
	@openssl verify -CAfile certs/ca.cert.pem \
           intermediate/certs/intermediate.cert.pem
	@cat intermediate/certs/intermediate.cert.pem \
             certs/ca.cert.pem > intermediate/certs/ca-chain.cert.pem
	@chmod 444 intermediate/certs/ca-chain.cert.pem

# 4. Server Key
server_pass:
	@openssl genrsa -aes256 -out intermediate/private/www.example.com.key.pem 2048
	@ chmod 400 intermediate/private/www.example.com.key.pem

key:
	@./gen_key.sh

@server_post:
	@openssl req -config intermediate/openssl.cnf \
            -key intermediate/private/www.example.com.key.pem \
            -new -sha256 -out intermediate/csr/www.example.com.csr.pem

clean:
	@rm -f private/* intermediate/private/*
	@rm -f certs/* intermediate/certs/*
