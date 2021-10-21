## How run azblob integration tests

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
export STORAGE_AZBLOB_INTEGRATION_TEST=on
export STORAGE_AZBLOB_CREDENTIAL=hmac:access_key:secret_key
export STORAGE_AZBLOB_NAME=bucketname
export STORAGE_AZBLOB_ENDPOINT=https:accountname.blob.core.windows.net
```

Run tests

```shell
make integration_test
```
