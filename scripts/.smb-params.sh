#!/usr/bin/env bash

SERVER="$1"
SHARE="$2"
FOLDER="$3"
SECURE="$4"
USERNAME="$5"
PASSWORD="$6"

if [ "$SERVER" == "" ]; then
    SERVER="localhost"
fi

if [ "$SHARE" == "" ]; then
    SHARE="custom"
fi

if [ "$FOLDER" == "" ]; then
    FOLDER="/"
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
echo "folder: $FOLDER"
echo "secure: $SECURE"
echo "username: $USERNAME"
echo "password: $PASSWORD"
