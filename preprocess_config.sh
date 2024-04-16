#!/bin/bash

# Load .env file if it exists and export its variables
if [ -f .env ]; then
    set -a # Automatically export variables
    source .env
    set +a
fi

# Prepare the sed command dynamically based on .env variables
sedCommand=""
while IFS='=' read -r key value; do
    # Skip lines that are empty or start with # (comments)
    [[ $key == \#* || -z $key ]] && continue
    # Append new replacement pattern to sedCommand
    sedCommand+="-e 's|\${${key}}|${!key}|g' "
done < .env

# Apply the dynamic sed command to replace placeholders
eval "sed $sedCommand endpoints.json > krakend.json"
