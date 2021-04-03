#!/usr/bin/env sh

config_file=config.yaml

if [ ! -f "$config_file" ]; then
touch "$config_file"
echo "# 项目名称
project_name: pastebin
# debug
debug: false
# 端口" >> $config_file

if [ -n "$inner_port" ]; then
  echo "port: ${inner_port}" >> $config_file
else
  echo "port: :21330" >> $config_file
fi

echo "# mongodb配置
mongodb:" >> $config_file

if [ -n "$mongo_uri" ]; then
  echo "  uri: ${mongo_uri}" >> $config_file
else
  echo "  uri: mongodb://127.0.0.1:27017" >> $config_file
fi

if [ -n "$mongo_name" ]; then
  echo "  dbname: ${mongo_name}" >> $config_file
else
  echo "  dbname: pastebin" >> $config_file
fi
fi

./pbx serve
