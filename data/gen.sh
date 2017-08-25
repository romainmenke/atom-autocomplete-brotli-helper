#!/bin/bash

sed 's/^\s*./  "/g' data/raw_brotli.txt | sed 's/.$/",/' | sed '$ s/.$//' > data/tmp

echo "[" > data/brotli.json
cat data/tmp >> data/brotli.json
echo "]" >> data/brotli.json

rm data/tmp
