#!/usr/bin/env bash

curl --silent \
    --request GET \
     http://localhost:3000/leases | jq

