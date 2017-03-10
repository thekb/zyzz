#!/usr/bin/env bash

echo "generating messages from flatbuffer files"

echo "generating go bindings"
flatc --grpc --go message/*.fbs
echo "generating javascript bindings"
flatc --js -o message/ message/*.fbs