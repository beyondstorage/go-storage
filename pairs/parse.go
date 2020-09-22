package pairs

import (
	"strconv"
)

func parseInt(in string) (int, error) {
	i, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), err
}

func parseInt64(in string) (int64, error) {
	return strconv.ParseInt(in, 10, 64)
}
