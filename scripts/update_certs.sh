#!/bin/bash

set -e

if [[ -z "$1" ]]; then
	echo "Output directory is not provided as first argument" >&2
	exit 1
fi

if [[ -z "$2" ]]; then
	echo "CN subject is not provided as second argument" >&2
	exit 1
fi

OUTPUT_DIR=$1
CN_SUBJECT=$2

CRT_OUT=$OUTPUT_DIR/root.crt.pem
KEY_OUT=$OUTPUT_DIR/root.key.pem

echo "Generate on `date`" > $OUTPUT_DIR/release.md
echo "`openssl version`" >> $OUTPUT_DIR/release.md

openssl ecparam -out $KEY_OUT \
	-name secp521r1 -genkey -noout

openssl req -x509 -sha256 -new -nodes -key $KEY_OUT -out $CRT_OUT \
	-days 30 -subj $CN_SUBJECT