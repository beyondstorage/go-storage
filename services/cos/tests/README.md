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
export STORAGE_COS_INTEGRATION_TEST=on
export STORAGE_COS_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_COS_NAME=bucketname
export STORAGE_COS_LOCATION=bucketname
```

Run tests

```shell
make integration_test
```
