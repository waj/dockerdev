FROM jwilder/nginx-proxy:latest

RUN apt-get update \
 && apt-get install -y -q --no-install-recommends \
    dnsmasq \
 && apt-get clean \
 && rm -r /var/lib/apt/lists/*

# override nginx configs
COPY *.conf /etc/nginx/conf.d/

# override nginx-proxy templating
COPY nginx.tmpl Procfile monitor /app/

# COPY htdocs /var/www/default/htdocs/

ENV DOMAIN_TLD dev
ENV DNS_IP 127.0.0.1
ENV HOSTMACHINE_IP 127.0.0.1
