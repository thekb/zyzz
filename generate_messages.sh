#!/usr/bin/env bash

echo "generating messages from flatbuffer files"

echo "generating go and python bindings"
flatc --go --python message/*.fbs
echo "generating javascript bindings"
flatc --js -o message/ message/*.fbs