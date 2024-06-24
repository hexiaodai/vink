#!/bin/bash

set -eu

buf generate --timeout 10m -v \
  --path common/ \
  --path management/

for d in common/ management/; do
    for f in $(find $d -name "*.proto"); do
        protoc --validate_out="paths=source_relative,lang=go:." $f
    done
done

for d in sdks/ts/management; do
  for f in $(find $d -type f -name "*.ts"); do
    if [ "$(uname)" = "Darwin" ]; then
      sed -i "" -r 's#(^type Base.*)#/* vink modified */ export \1#g' $f
    else
      sed -i -r 's#(^type Base.*)#/* vink modified */ export \1#g' $f
    fi
  done
done
