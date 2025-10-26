#!/bin/bash

mkdir -p keys

openssl genrsa -out keys/jwt_private.pem 2048
openssl rsa -in keys/jwt_private.pem -pubout -out keys/jwt_public.pem

echo "JWT keys generated successfully in ./keys/ directory"
