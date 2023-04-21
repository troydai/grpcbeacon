#!/bin/sh

grpcurl -plaintext -proto etc/protos/beacon.proto localhost:7000 grpcbeacon.Beacon/Signal
