```

# use go mod vendor, copy dependencies to local directory

# docker build
# Create Dockerfile in the directory
docker build -t wncbb/httpserver:v0.0.1 .

# docker image rename
docker tag  httpserver:v0.0.1 wncbb/httpserver:v0.0.1

# docker image delete
docker image rm httpserver:v0.0.1

# docker image push
docker image push wncbb/httpserver:v0.0.1
```

```
# start container
docker run -d -p7878:7878 --rm wncbb/httpserver:v0.0.2

# curl test
curl -v -H 'demoHeader:thisIsDemoHeader' http://127.0.0.1:7878
```
