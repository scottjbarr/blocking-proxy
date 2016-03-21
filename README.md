# Blocking Proxy

A little proxy that 404's a list of resources.

Any path not on the block list is forwarded to the backend.

Listens on http://localhost:8080

The exampe below blocks

- all HTTP verbs to `/foo`
- `PUT /moo.json`

Usage

    blocking-proxy -bind :8080 \
                   -backend localhost:3000
                   -block "*:/foo"
                   -block PUT:/moo.json
