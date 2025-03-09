#!/bin/sh

read_input() {
    if [ -n "$1" ]; then
        echo "${*}"
    elif [ -p /dev/stdin ]; then
        read -r input
        echo "$input"
    fi
}


SCRIPT_DIR=$(dirname -- "$( readlink -f -- "$0"; )");
if [ -n "$1" ]; then
  if [ -f "$1" ]; then
    DEVICE_CERT="$1"
    shift
  else
    echo "$1 does not seem to exist; please check the location of your device certificate file"
  fi
else
  echo "Please provide a device certificate file as first argument to this script"
  exit 1
fi

MESSAGE=$(read_input $*)
if [ -z "$MESSAGE" ]; then
  echo "Please provide a message to send"
  exit 1
fi

curl -X POST -d"$MESSAGE" --cacert $SCRIPT_DIR/ca.pem --cert $DEVICE_CERT --key $DEVICE_CERT  https://receivers.iot-exchange.io/uplink 



