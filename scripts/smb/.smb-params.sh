#!/usr/bin/env bash

SERVER="$1"
SHARE="$2"
SECURE="$3"
USERNAME="$4"
PASSWORD="$5"

if [ "$SERVER" == "" ]; then
    SERVER="localhost"
fi

if [ "$SHARE" == "" ]; then
    SHARE="custom"
fi

if [ "$SECURE" == "" ]; then
    SECURE="true"
fi

if [ "$USERNAME" == "" ]; then
    USERNAME="pknopf"
fi

if [ "$PASSWORD" == "" ]; then
    PASSWORD="password"
fi

echo "server: $SERVER"
echo "share: $SHARE"
echo "secure: $SECURE"
echo "username: $USERNAME"
echo "password: $PASSWORD"
