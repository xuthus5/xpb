# xpb

> a pastebin cli

## RUN

docker run:

```shell
docker run -d --network=host \
      -p 21330:21330 \
      --name xpb \ 
      -e "mongo_uri=mongodb://admin:admin@127.0.0.1:27017" \
      -e "mongo_dbname=pastebin" \
      xuthus5/xpb
```
