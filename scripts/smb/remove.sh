#!/usr/bin/env bash

MEDIA_ID="$1"

curl --silent \
    --request POST \
    --data '{"mediaId":"'$MEDIA_ID'"}' \
     http://localhost:3000/smb/remove | jq

