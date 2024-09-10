## Usage

### Run via docker:

To run app and it's dependencies in docker, use next command:
```bash
task -d scripts docker_prod -v
```

### Run via source files:

To run application via source files, use next commands:
```shell
go run ./cmd/hmtmsso/main.go
```

## gRPC:

To setup protobuf use next command:
```shell
task -d scripts setup_proto -v
```


### Base files generation:
```shell
task -d scripts grpc_generate -v
```

## Linters

```shell
golangci-lint run -v --fix
```

## Tests

```shell
go test -v ./test...
```
