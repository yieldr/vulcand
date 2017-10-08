# Vulcand

Yieldr's fork of Mailgun's [vulcand](https://github.com/vulcand/vulcand) built using [`vbundle`](http://vulcand.github.io/middlewares.html#vbundle) to add custom middleware.

# Usage

## Kubernetes

Typically this project is deployed alongside [romulus](https://github.com/albertrdixon/romulus) as an [Ingress Controller](https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-controllers).

Please refer to the [romulus wiki](https://github.com/albertrdixon/romulus/wiki/Annotations) for more details.

## Standalone

Build the binaries using `make`.

	make build OS=linux

Then start docker compose

	docker-compose up -d

When all services are up and running, you'll need to configure vulcan to route traffic to the upstream servers. A dummy upstream is included for convenience that simply displays the nginx default web page.

### Frontends, backends  and servers

The following configuration will create a new upstream backend (`nginx`), a server (`nginx-srv`) and a new frontend (`nginx`).

	vctl backend upsert -id nginx
	vctl server upsert -b nginx -id nginx-srv -url http://nginx:80
	vctl frontend upsert -id nginx -b nginx -route 'PathRegexp("/.*")'

### Middleware

This command will add `oauth2` middleware to the `nginx` frontend.

	vctl oauth2 upsert -f nginx -id nginx-oauth \
		-domain $AUTH0_DOMAIN \
		-clientId $AUTH0_CLIENT_ID \
		-clientSecret $AUTH0_CLIENT_SECRET \
		-redirectUrl http://localhost:8181/callback