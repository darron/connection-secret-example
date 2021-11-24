#!/bin/sh

while [ 1 ]; do
    curl -s rtest2.default.svc.cluster.local/redis > /dev/null
done