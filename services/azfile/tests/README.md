## How run azfile integration tests

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
export STORAGE_AZFILE_INTEGRATION_TEST=on
export STORAGE_AZFILE_CREDENTIAL=hmac:account_name:account_key
export STORAGE_AZFILE_NAME=sharename
export STORAGE_AZFILE_ENDPOINT=https:accountname.file.core.windows.net
```

Run tests

```shell
make integration_test
```
