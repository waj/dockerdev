# DockerDev HTTP Proxy

Based on excellent work here: https://github.com/codekitchen/dinghy-http-proxy

The difference of this proxy with the original made by @codekitchen is that this one uses a "shared" network to connect all the exposed containers using aliases. This way is more easy to make connections between the containers using the same host names.

## How to build

First, build the `monitor` service running:

```bash
docker-compose run monitor
```

This will leave a `monitor` executable on the working directory.

Now build the main image:

```bash
docker-compose build proxy
```

This makes the `juanwaj/dockerdev` Docker image.

## Run

Before running create a file on your Mac, located at `/etc/resolvers/dev` with the content:

```
nameserver 127.0.0.1
port 19322
```

Then, make sure the `shared` Docker network exists:

```bash
docker network create shared
```

### Run from the working copy

If you cloned this repository, the proxy can be started executing:

```bash
docker-compose up -d proxy
```

### Run pulling the image from Docker Hub

To run the proxy using the prebuilt image, run in your command line:

```bash
docker run -d --restart=always \
  -v /var/run/docker.sock:/tmp/docker.sock:ro \
  -p 80:80 -p 443:443 -p 19322:19322/udp \
  --log-opt max-size=10m --log-opt max-file=3 \
  --network shared --name dockerdev \
  juanwaj/dockerdev
```
