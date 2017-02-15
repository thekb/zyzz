#!/usr/bin/env bash
echo "generating bin data for migrations..."
go-bindata -pkg db -o db/bindata_migrations.go db/migrations/