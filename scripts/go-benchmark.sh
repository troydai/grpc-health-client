#!/bin/bash

set -e

# Start server in background and save its PID
bin/server &
SERVER_PID=$!

# Cleanup function to kill the server
cleanup() {
    echo "Cleaning up server process..."
    kill $SERVER_PID
}

echo "Server started with PID $SERVER_PID"

# Start timing
start_time=$(date +%s.%N)

# Run client 1000 times
for i in {1..1000}; do
    # Print progress every 50 cycles
    if [ $((i % 50)) -eq 0 ]; then
        echo "Progress: $i/1000 cycles completed"
    fi

    bin/client
done

# Calculate total duration
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
echo "Total duration: $duration seconds for 1000 runs"

# Register cleanup function to run on script exit
trap cleanup EXIT