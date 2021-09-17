/*
Package tests need to ensure the generator work as intended so that we can generate code correctly.
So we ignored `generated.go`, and generated every time to test generator and check `service.toml`.
If the test failed, the generator SHOULD NOT be used in specific service.
*/

package tests

//go:generate go run -tags tools github.com/beyondstorage/go-storage/v4/cmd/definitions --debug service.toml
