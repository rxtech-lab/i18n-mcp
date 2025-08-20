#!/bin/bash

# Exit on any error
set -e

# Check if required variables are set
if [ -z "${APPLICATION_SIGNING_CERTIFICATE_NAME}" ]; then
  echo "Warning: APPLICATION_SIGNING_CERTIFICATE_NAME is not set, skipping signing"
  exit 0
fi

# Import binary configuration
source "$(dirname "$0")/binaries.sh"

# Sign each binary
for binary in "${BINARIES[@]}"; do
  BINARY_PATH="bin/${binary}"
  
  if [ ! -f "${BINARY_PATH}" ]; then
    echo "Error: Binary ${BINARY_PATH} not found"
    exit 1
  fi
  
  echo "Signing ${BINARY_PATH}..."
  codesign --force --deep --sign "${APPLICATION_SIGNING_CERTIFICATE_NAME}" \
    --options runtime \
    --timestamp \
    "${BINARY_PATH}"
  
  # Verify the signature
  echo "Verifying signature for ${BINARY_PATH}..."
  codesign --verify --verbose "${BINARY_PATH}"
done

echo "All binaries signed successfully"