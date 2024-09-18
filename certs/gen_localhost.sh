#!/bin/bash

# generate private key (localhost.key) and public key (localhost.csr)
openssl req -new -subj "/C=US/ST=Utah/CN=localhost" -newkey rsa:2048 -nodes -keyout localhost.key -out localhost.csr

# sign certificate
openssl x509 -req -days 3650 -in localhost.csr -signkey localhost.key -out localhost.crt