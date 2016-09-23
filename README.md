# DockerDev HTTP Proxy

Based on excellent work here: https://github.com/codekitchen/dinghy-http-proxy

The difference of this proxy with the original made by @codekitchen is that this one uses a "shared" network to connect all the exposed containers using aliases. This way is more easy to make connections between the containers using the same host names.

## Run

Before running create a file on your Mac, located at `/etc/resolver/dev` with the content:

```
nameserver 127.0.0.1
port 19322
```

Then, make sure the `shared` Docker network exists:

```bash
docker network create shared
```

Now start the proxy container:

```bash
docker run -d --restart=always \
  -v /var/run/docker.sock:/tmp/docker.sock:ro \
  -p 80:80 -p 443:443 -p 19322:19322/udp \
  --log-opt max-size=10m --log-opt max-file=3 \
  --network shared --name dockerdev \
  juanwaj/dockerdev
```


## Usage

Once the proxy is running, new containers started using Docker Compose will be joined to the `shared` network and proxied by the `dockerdev` container.

For example, if you have a container `web` in a project named `foo`, once the container is running, point your browser to `web.foo.dev`.


## Development

If you want to make changes or just don't want to use the prebuilt image, after cloning
this repository, first build the `monitor` service running:

```bash
docker-compose run monitor
```

This will leave a `monitor` executable on the working directory.

Now start the proxy:

```bash
docker-compose up -d proxy
```

That's all!
