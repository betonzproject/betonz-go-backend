#!/bin/bash

# Source the .env file
. .env

# Parse the DATABASE_URL using awk
parsed_string=$(echo "$DATABASE_URL" | sed 's/\?.*//')

# Set the schema name
schema=$(echo "$DATABASE_URL" | sed -n 's/.*search_path=\([^&]*\).*/\1/p')

# Execute psql command to connect to the database and run SQL files
psql "$parsed_string" -c "SET search_path TO $schema;" -f ./db/seed.sql