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
export STORAGE_S3_INTEGRATION_TEST=on
export STORAGE_S3_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_S3_NAME=bucketname
export STORAGE_S3_LOCATION=location
```

Run tests

```shell
make integration_test
```
