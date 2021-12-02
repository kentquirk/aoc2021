#! /bin/bash

if [ -z $1 ]; then
  echo "Day number required"
  exit 1
fi

cp -r _template day$1
cd day$1
sed s/XXX/day$1/ <../_template/go.mod >go.mod
vc .
