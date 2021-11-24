#!/bin/sh

while [ 1 ]; do
    curl -s rtest.default.svc.cluster.local/redis > /dev/null
done