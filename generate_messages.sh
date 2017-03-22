#!/usr/bin/env bash

echo "generating messages from flatbuffer files"

echo "generating go bindings"
rm -rf message/*.go
flatc --go message/*.fbs
echo "generating javascript bindings"
rm -rf message/*.js
flatc --js -o message/ message/*.fbs
echo "generating java bindling"
flatc --java message/*.fbs