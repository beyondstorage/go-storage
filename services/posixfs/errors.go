package posixfs

func handleOsError(err error) error {
	if err == nil {
		panic("error must not be nil")
	}

	// TODO: handle PathError here.
	return err
}
