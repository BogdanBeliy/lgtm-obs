# lgtm-obs

[![Latest tag](https://img.shields.io/github/v/tag/BogdanBeliy/lgtm-obs?sort=semver)](https://github.com/BogdanBeliy/lgtm-obs/tags)

OpenTelemetry observability library for Go: traces, metrics, logs → OTLP → LGTM stack.

## Install

```bash
go get github.com/BogdanBeliy/lgtm-obs@v0.1.0
```

## Usage

```go
obs, err := observability.InitObservability(ctx, &observability.Configs{
    ServiceName: "my-service",
    Protocol:    "http", // or grpc
    Path:        "localhost:4318",
    Insecure:    true,
})
defer obs.Shutdown(ctx)
```

## Packages

| Import | Purpose |
|--------|---------|
| `github.com/BogdanBeliy/lgtm-obs/observability` | Init, traces, metrics, logs, slog |
| `github.com/BogdanBeliy/lgtm-obs/observability/gin` | Gin middleware |
| `github.com/BogdanBeliy/lgtm-obs/observability/fiber` | Fiber middleware |

## Examples

```bash
go run ./examples/http
go run ./examples/gin
go run ./examples/fiber
```

Requires a running OTLP collector (e.g. local LGTM on `:4318` HTTP / `:4317` gRPC).

## Quality checks

```bash
make tools     # install pinned golangci-lint and govulncheck versions
make fmt       # format Go code
make check     # formatting, linters, security, vulnerabilities, race tests
```

Pull requests to `main` must use a release prefix:

- `patch/*` increments `vX.Y.Z`
- `minor/*` increments `vX.Y.0`
- `major/*` increments `vX.0.0`

The next version is shown in the PR checks. The tag is created after merge.
