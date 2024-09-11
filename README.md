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

To setup protobuf, use next command:
```shell
task -d scripts setup_proto -v
```


To generate files from .proto, use next command:
```shell
task -d scripts grpc_generate -v
```

## Linters

```shell
task -d scripts linters -v
```

## Tests

```shell
task -d scripts tests -v
```
