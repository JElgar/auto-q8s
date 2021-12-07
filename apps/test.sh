#!/bin/bash
for i in {1..1000}
do
   curl -X POST -H "Content-Type: application/json" -d '{"id": "28"}' localhost:3000
done
