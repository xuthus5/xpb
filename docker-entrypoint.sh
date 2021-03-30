#!/usr/bin/env bash

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

if [ $mongo_uri ]; then
  echo "refresh_time_interval: ${time_interval}" >> $configfile
else
  echo "refresh_time_interval: 15" >> $configfile
fi

echo "# mongodb配置
mongodb:" >> $configfile

if [ $mongo_dbname ]; then
  echo "  uri: ${mongo_dbname}" >> $configfile
else
  echo "  uri: mongodb://127.0.0.1:27017" >> $configfile
fi
fi

./xpb serve