#!/bin/bash

USER=$1
PASSWORD=$2
HOST=$3
PORT=$4
DB=$5
SSLMODE=$6

pg_copy(){
    local table=$1
    local file=$2

    echo "Seeding $table..."
    psql "postgresql://$USER:$PASSWORD@$HOST:$PORT/$DB?sslmode=$SSLMODE" \
        -c "\copy $table FROM './db/seeds/$file' DELIMITER ',' CSV HEADER;"
}

pg_copy users users.csv
pg_copy user_details user_details.csv
