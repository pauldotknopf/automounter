#!/usr/bin/env bash

. ./.smb-params.sh

curl --silent \
    --request POST \
    --data '{"server":"'$SERVER'", "share":"'$SHARE'", "folder":"'$FOLDER'", "secure":'$SECURE', "username":"'$USERNAME'", "password":"'$PASSWORD'"}' \
     http://localhost:3000/smb/test | jq

