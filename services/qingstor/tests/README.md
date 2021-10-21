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
export STORAGE_QINGSTOR_INTEGRATION_TEST=on
export STORAGE_QINGSTOR_CREDENTIAL=hamc:access_key:secret_key
export STORAGE_QINGSTOR_ENDPOINT=https:qingstor.com:443
export STORAGE_QINGSTOR_NAME=bucketname
```

Run tests

```shell
make integration_test
```
