```
# docker build
# Create Dockerfile in the directory
docker build -t wncbb/geektimehttpserverdemo:v0.0.1 .

# docker image rename
docker tag  geektimehttpserverdemo:v0.0.1 wncbb/geektimehttpserverdemo:v0.0.1

# docker image delete
docker image rm geektimehttpserverdemo:v0.0.1

# docker image push
docker image push wncbb/geektimehttpserverdemo:v0.0.1
```
