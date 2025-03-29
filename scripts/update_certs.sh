#!/bin/bash

set -e

if [[ -z "$1" ]]; then
	echo "Output directory is not provided as first argument" >&2
	exit 1
fi

OUTPUT_DIR=$1

CRT_OUT=$OUTPUT_DIR/root.crt.pem
KEY_OUT=$OUTPUT_DIR/root.key.pem

echo "Generate on `date`" > $OUTPUT_DIR/release.md
echo "`openssl version`" >> $OUTPUT_DIR/release.md

openssl ecparam -out $KEY_OUT \
	-name secp521r1 -genkey -noout

openssl req -new -key $KEY_OUT -out $OUTPUT_DIR/server.csr \
	-config $OUTPUT_DIR/csr.cnf

openssl x509 -req -days 265 -in $OUTPUT_DIR/server.csr -signkey $KEY_OUT \
	-out $CRT_OUT -extensions v3_req -extfile $OUTPUT_DIR/csr.cnf
