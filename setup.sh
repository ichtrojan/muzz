#!/bin/bash

# Prompt the user for database details
echo "Enter database host (default: 127.0.0.1):"
read -r DB_HOST
DB_HOST=${DB_HOST:-127.0.0.1}

echo "Enter database port (default: 3306):"
read -r DB_PORT
DB_PORT=${DB_PORT:-3306}

echo "Enter database name (default: muzz):"
read -r DB_NAME
DB_NAME=${DB_NAME:-muzz}

echo "Enter database user (default: root):"
read -r DB_USER
DB_USER=${DB_USER:-root}

echo "Enter database password:"
read -r -s DB_PASS

# Create the .env file
cat <<EOF > .env
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASS=$DB_PASS
EOF

echo ".env file created successfully."

# Connect to MySQL and create the database
mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME;"

if [ $? -eq 0 ]; then
  echo "Database '$DB_NAME' created successfully."

  # Check if the migration.sql file exists and run it
  SQL_DUMP_FILE="migration.sql"
  if [ -f "$SQL_DUMP_FILE" ]; then
    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$SQL_DUMP_FILE"
    if [ $? -eq 0 ]; then
      echo "SQL dump file '$SQL_DUMP_FILE' imported successfully into database '$DB_NAME'."
    else
      echo "Failed to import SQL dump file '$SQL_DUMP_FILE'."
    fi
  else
    echo "SQL dump file '$SQL_DUMP_FILE' does not exist."
  fi

else
  echo "Failed to create database '$DB_NAME'."
fi
