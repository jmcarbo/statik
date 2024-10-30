# Statik

Serve static files from a directory. HTML files can have Go template syntax embedded in them.

## Available variables

- `{{.Headers}}`: The request headers.
- `{{.Method}}`: The request method.
- `{{.Path}}`: The request path.
- `{{.Query}}`: The request query string.
- `{{.Token}}`: Authentication bearer token or blank if it does not exist.


## Usage

```
docker run -ti --rm -p 3000:3000 -v ./public:/public ghcr.io/jmcarbo/statik:latest -dir /public
```

