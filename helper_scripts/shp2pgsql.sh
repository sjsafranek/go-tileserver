#!/bin/bash

set -x

SHAPE_FILE="$1"
DB_TABLE="$2"

shp2pgsql -I "$SHAPE_FILE" "$DB_TABLE" | psql
