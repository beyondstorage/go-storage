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
export STORAGE_FTP_INTEGRATION_TEST=on
export STORAGE_FTP_CREDENTIAL=basic:user:password
export STORAGE_FTP_ENDPOINT=tcp:127.0.0.1:2121
```

Run tests

```shell
make integration_test
```