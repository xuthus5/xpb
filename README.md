# xpb

> a pastebin cli

## RUN

docker run:

```shell
docker run -d --network=host \
      -p 21330:21330 \
      --name xpb \ 
      -e "mongo_uri=mongodb://admin:123456@192.168.3.6:27017" \
      -e "mongo_dbname=pastebin" \
      xuthus5/xpb
```
