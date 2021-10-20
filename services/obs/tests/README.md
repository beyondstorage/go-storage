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
export STORAGE_OBS_INTEGRATION_TEST=on
export STORAGE_OBS_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_OBS_NAME=bucketname
export STORAGE_OBS_ENDPOINT=https:obs.<region>.myhuaweicloud.com
```

Run tests

```shell
make integration_test
```
