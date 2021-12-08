#!/bin/bash
for i in {1..600}
do
   curl -X POST -H "Content-Type: application/json" -d '{"id": "28"}' http://producer.$CLUSTER_BASE_URL
done
