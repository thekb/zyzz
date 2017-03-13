#!/usr/bin/env bash

echo "generating messages from flatbuffer files"

echo "generating go bindings"
rm -rf message/*.go
flatc --grpc --go message/*.fbs
echo "generating javascript bindings"
rm -rf message/*.js
flatc --js -o message/ message/*.fbs