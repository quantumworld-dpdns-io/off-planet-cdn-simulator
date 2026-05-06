#!/usr/bin/env bash
set -e

docker compose up -d
supabase start
echo "Dev services started"
