#!/bin/sh
# scripts/wait-for-db.sh
# Script para aguardar o banco estar pronto

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "🔄 Aguardando banco de dados em $host:$port..."

until nc -z "$host" "$port"; do
  echo "🔄 Banco ainda não está pronto - aguardando..."
  sleep 2
done

echo "✅ Banco de dados está pronto!"
exec $cmd