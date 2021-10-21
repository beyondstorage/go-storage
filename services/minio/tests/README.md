## How run integration tests

### Run tests locally

Copy example files and update corresponding values.

```shell
cp Makefile.env.exmaple Makefile.env
```

Run tests

```shell
make integration_test
```

### Run tests in CI

Set following environment variables:

```shell
export STORAGE_MINIO_INTEGRATION_TEST=on
export STORAGE_MINIO_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_MINIO_NAME=bucketname
export STORAGE_MINIO_ENDPOINT=http:host:port
```

Run tests

```shell
make integration_test
```