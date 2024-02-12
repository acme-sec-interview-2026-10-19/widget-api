# widget-api

Internal HTTP API for Acme Widgets. It serves a public health-style greeting and a
small set of authenticated admin endpoints used by the ops dashboard.

## Endpoints

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/hello` | none | Returns `hello, world`. |
| `GET` | `/admin/status` | Bearer JWT | Service name, region, and status as JSON. |

Admin tokens are HS256 JWTs with `aud` set to the service name (`widget-api`). The ops
team issues them; ask in #ops if you need one.

## Configuration

Settings live in `config/settings.yaml`. `${VAR}` references in that file are expanded
from the environment at startup. Set `WIDGET_API_CONFIG` to load a different file.

| Variable | Required | Description |
| --- | --- | --- |
| `JWT_SIGNING_KEY` | yes | HMAC key used to verify admin tokens. |
| `AWS_ACCESS_KEY_ID` | no | AWS credentials for inventory sync. |
| `AWS_SECRET_ACCESS_KEY` | no | AWS credentials for inventory sync. |
| `WIDGET_API_CONFIG` | no | Path to the settings file. Defaults to `config/settings.yaml`. |

## Running locally

```bash
export JWT_SIGNING_KEY=local-dev-key
go run .
```

```bash
curl localhost:8080/hello
```

## Tests

```bash
go test ./...
```

## Docker

```bash
docker build -t widget-api .
docker run -p 8080:8080 -e JWT_SIGNING_KEY=local-dev-key widget-api
```

## CI

GitHub Actions builds and tests every push to `main` and every pull request, then
deploys `main`.
