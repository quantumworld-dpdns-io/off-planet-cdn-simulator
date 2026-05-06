#!/usr/bin/env bash
set -e

supabase stop
docker compose down
echo "Dev services stopped"
