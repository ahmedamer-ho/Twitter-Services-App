#!/bin/bash
# Initialization script for NATS JetStream Streams and Consumers
# Execute this script once the NATS cluster is running in Kubernetes.
# Setup: Port-forward the NATS service to your local machine:
# kubectl port-forward svc/nats 4222:4222

NATS_URL="nats://localhost:4222"

echo "Creating NATS Streams for Microservices Nervous System..."

# 1. UserStream
# Limits based retention, discards oldest messages on size/time limits.
nats stream add UserStream \
  --subjects "user.*" \
  --retention limits \
  --discard old \
  --max-age 7d \
  --max-bytes 1000000000 \
  --storage file \
  --replicas 3 \
  -s $NATS_URL

# 2. TweetStream
# Limits based retention, durable
nats stream add TweetStream \
  --subjects "tweet.*" \
  --retention limits \
  --discard old \
  --max-age 7d \
  --max-bytes 5000000000 \
  --storage file \
  --replicas 3 \
  -s $NATS_URL

# 3. SearchStream
# WorkQueue retention, messages are removed immediately upon ACK.
nats stream add SearchStream \
  --subjects "search.*" \
  --retention workq \
  --storage file \
  --replicas 3 \
  -s $NATS_URL

echo "Creating explicit durable consumers for guaranteed delivery..."

# Timeline Consumer on TweetStream
nats consumer add TweetStream timeline-consumer \
  --ack explicit \
  --wait 30s \
  --deliver all \
  --replay instant \
  --durable timeline-consumer \
  -s $NATS_URL

# Search Indexer Consumer on TweetStream
nats consumer add TweetStream search-indexer-consumer \
  --ack explicit \
  --wait 30s \
  --deliver all \
  --replay instant \
  --durable search-indexer-consumer \
  -s $NATS_URL

echo "Streams and Consumers configured successfully!"
