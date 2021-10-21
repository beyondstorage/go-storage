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
export STORAGE_DROPBOX_INTEGRATION_TEST=on
export STORAGE_DROPBOX_CREDENTIAL=apikey:apikey
```

Run tests

```shell
make integration_test
```
