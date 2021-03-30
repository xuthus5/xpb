#!/bin/bash

configfile=/src/config.yaml

if [[ ! -x "$configfile" ]]; then
echo "# 项目名称
project_name: pastebin
# debug
debug: false
# 端口" > $configfile

if [ $inner_port ]; then
  echo "port: ${inner_port}" >> $configfile
else
  echo "port: :21330" >> $configfile
fi

echo "# mongodb配置
mongodb:" >> $configfile

if [ $mongo_uri ]; then
  echo "  uri: ${mongo_uri}" >> $configfile
else
  echo "  uri: mongodb://127.0.0.1:27017" >> $configfile
fi

if [ $mongo_name ]; then
  echo "  dbname: ${mongo_name}" >> $configfile
else
  echo "  dbname: pastebin" >> $configfile
fi
fi

echo "done"

cat /src/config.yaml

/src/xpb serve
