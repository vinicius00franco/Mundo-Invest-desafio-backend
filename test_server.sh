#!/bin/bash
export $(cat .env | xargs)
./server > server.log 2>&1 &
