#!/bin/sh
set -e

host="$1"
shift
cmd="$@"

until pg_isready -h "$host" -p 5432; do
  >&2 echo "Waiting for PostgreSQL at $host:5432..."
  sleep 1
done

>&2 echo "PostgreSQL is up - executing command"
exec $cmd
