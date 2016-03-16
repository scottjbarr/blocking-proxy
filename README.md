# Blocking Proxy

A little proxy that 404's a list of paths.

Any path not on the block list is forwarded to http://localhost:3000

Listens on http://localhost:8080

Usage

   blocking-proxy  /some/url /foo.json
