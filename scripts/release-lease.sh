#!/usr/bin/env bash

LEASE_ID="$1"

curl --silent \
    --request POST \
    --data '{"leaseId":"'$LEASE_ID'"}' \
     http://localhost:3000/leases/release | jq

