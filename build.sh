#!/bin/bash
CUR_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
FRONTEDN=$CUR_DIR/frontend
STATIC_DIR=$CUR_DIR/framework/http/middleware/frontend


cd $CUR_DIR/frontend && npm run build
mkdir -p $STATIC_DIR
cp -a $FRONTEDN/dist $STATIC_DIR

