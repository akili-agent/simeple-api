# simple-api

A small, production-ready HTTP JSON API written in Go on top of the
[Okapi](https://github.com/jkaninda/okapi) web framework. It ships with
OpenAPI 3.1 documentation (Scalar UI), health/readiness probes, an in-memory
example resource, tests, a multi-stage Dockerfile and CI.

## Layout

```
cmd/server/          entrypoint: env config, Okapi setup, graceful shutdown
internal/routes/     routes declared as data ([]okapi.RouteDefinition) under /api/v1
internal/handlers/   thin handlers: probes + items example resource
internal/models/     request/response types (also the OpenAPI schemas)
internal/store/      concurrency-safe in-memory item store
```

## Endpoints

| Method | Path                  | Description                     |
|--------|-----------------------|---------------------------------|
| GET    | `/api/v1/healthz`     | Liveness probe                  |
| GET    | `/api/v1/readyz`      | Readiness probe                 |
| GET    | `/api/v1/items`       | List items                      |
| GET    | `/api/v1/items/{id}`  | Get one item                    |
| POST   | `/api/v1/items`       | Create an item (`{"name":...}`) |
| GET    | `/docs`               | Scalar API reference UI         |
| GET    | `/openapi.json`       | OpenAPI 3.1 specification       |

## Configuration

Everything is configured through environment variables:

| Variable    | Default | Description                                    |
|-------------|---------|------------------------------------------------|
| `PORT`      | `8080`  | TCP port the HTTP server listens on            |
| `LOG_LEVEL` | `info`  | Log level: `debug`, `info`, `warning`, `error` |

## Run locally

Requires Go 1.26+.

```sh
make run                 # go run ./cmd/server
# or with configuration
PORT=9000 LOG_LEVEL=debug make run
```

Then:

```sh
curl localhost:8080/api/v1/healthz
curl -X POST localhost:8080/api/v1/items -d '{"name":"My first item"}'
curl localhost:8080/api/v1/items/1
open http://localhost:8080/docs
```

## Test and build

```sh
make test   # go test -race ./...
make vet    # go vet ./...
make build  # static binary into ./bin/server
```

## Docker

```sh
make docker                       # build the image (multi-stage, non-root)
docker run --rm -p 8080:8080 simple-api:latest
```

The runtime image is Alpine-based, runs as a non-root user and contains a
static CGO-disabled binary.

## CI

The Gitea workflow (`.gitea/workflows/ci.yml`) runs `go vet ./...` and
`go test -race ./...` on pushes to `main` and on pull requests.

An identical GitHub Actions workflow is provided as
`ci/github-actions.yml.txt`. Move it into place to activate it:

```sh
mkdir -p .github/workflows && mv ci/github-actions.yml.txt .github/workflows/ci.yml
```

(It ships as a `.txt` file because pushing workflow files to GitHub requires a
token with the `workflow` scope.)
