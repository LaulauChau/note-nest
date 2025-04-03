#!/bin/bash
set -e # Exit immediately if a command exits with a non-zero status.

echo "Running custom migration script..."

# Loop through all .up.sql files in the directory and execute them
for f in /docker-entrypoint-initdb.d/*.up.sql; do
  echo "Applying migration $f"
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done

echo "Custom migration script finished."