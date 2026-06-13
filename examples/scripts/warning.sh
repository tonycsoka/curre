#!/bin/bash
# Produces stderr output and a non-zero exit code

echo "Starting warning check..."
echo "Step 1: Checking disk space"

# Write to stderr
>&2 echo "WARNING: Disk usage above 80%"
>&2 echo "WARNING: /tmp has only 2GB remaining"

sleep 1

echo "Step 2: Checking memory"
>&2 echo "WARNING: Memory usage at 85%"

sleep 1

echo "Step 3: Checking logs"
>&2 echo "WARNING: Found 3 errors in application logs"

# Simulate a failure
>&2 echo "ERROR: Warning check failed — too many issues found"
exit 1
