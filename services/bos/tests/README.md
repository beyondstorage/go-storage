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
export STORAGE_BOS_INTEGRATION_TEST=on
export STORAGE_BOS_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_BOS_NAME=bucketname
export STORAGE_BOS_ENDPOINT=https:<region>.bcebos.com
```

Run tests

```shell
make integration_test
```
