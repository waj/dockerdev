你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# DockerDev HTTP Proxy

Based on excellent work here: https://github.com/codekitchen/dinghy-http-proxy

The difference of this proxy with the original made by @codekitchen is that this one uses a "shared" network to connect all the exposed containers using aliases. This way is more easy to make connections between the containers using the same host names.

## Run

First, make sure the `shared` Docker network exists:

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
  -e DOMAIN_TLD=lvh.me \
  juanwaj/dockerdev
```


## Usage

Once the proxy is running, **new** containers started using Docker Compose will be joined to the `shared` network and proxied by the `dockerdev` container.

For example, if you have a container `web` in a project named `foo`, once the container is running, point your browser to `web.foo.lvh.me`.

### Important note

If the container was created **before** dockerdev was installed, **you need to recreate the container**. This means, if you are using docker compose, you need to run `docker-compose down` and `docker-compose up` to allow the new container to be added to the `shared` network.

### Using HTTPS/SSL

You can also offload SSL to the dockerdev container. That way, you can expose your applications through HTTPS without modifying their source code.

Let's say you have a service called `app` in project called `myproject`, and you're starting up `dockerdev` with `DOMAIN_TLD=lvh.me`. With that setup, using `dockerdev` you'll be accessing `app` by issuing requests to `app.myproject.lvh.me`.

If you want `https://app.myproject.lvh.me` to work, follow these steps:

1. Generate a crt/key pair. You can generate both at the same time using this OpenSSL command (replace names to match your service/project/tld combination):

```
openssl req -x509 -newkey rsa:2048 -keyout app.myproject.lvh.me.key \
-out app.myproject.lvh.me.crt -days 365 -nodes \
-subj "/C=US/ST=Oregon/L=Portland/O=Company Name/OU=Org/CN=app.myproject.lvh.me" \
-config <(cat /etc/ssl/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:app.myproject.lvh.me")) \
-reqexts SAN -extensions SAN
```

2. Put these files in a per-host directory. I use `~/.dockerdev/certs`, but any other will do. 

3. Start `dockerdev` mounting that directory in `/etc/nginx/certs`. See example below:

```
docker run -d --restart=always \
  -v /var/run/docker.sock:/tmp/docker.sock:ro \
  -v ~/.dockerdev/certs:/etc/nginx/certs \
  -p 80:80 -p 443:443 -p 19322:19322/udp \
  --log-opt max-size=10m --log-opt max-file=3 \
  --network shared --name dockerdev \
  -e DOMAIN_TLD=lvh.me \
  juanwaj/dockerdev
```

4. Restart containers so they reregister, which lets the underlying NGINX container properly configure vhosts now that there are certs.

5. That's it! In a browser, navigate to https://app.myproject.lvh.me and it should be working. Note: browsers will complain that they don't trust the certificates you're providing. In order to avoid that, you'll need to tell your OS to trust those certificates. In macOS you can do that by following these steps:

- Open Keychain Access app.
- Drag and drop .crt files to the app.
- Double click each file, open the `Trust` section, select `Always trust`, close.

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
