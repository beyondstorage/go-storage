package posixfs

func handleOsError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}
	return err
}
