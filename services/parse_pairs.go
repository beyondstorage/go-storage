package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	. "github.com/beyondstorage/go-storage/v4/types"
)

var (
	// ErrPairTypeNotParsable means the pair's type is not parseable.
	ErrPairTypeNotParsable = NewErrorCode("type not parseable")
	// ErrPairNotRegistered means the pair is not registered.
	ErrPairNotRegistered = NewErrorCode("pair not registered")
	// ErrPairValueInvalid means the pair's value is invalid.
	ErrPairValueInvalid = NewErrorCode("pair value invalid")
)

func parseString(config string) (ty string, ps []Pair, err error) {
	return "", nil, nil
}

func parse(m map[string]string, k string, v string) (pair Pair, err error) {
	vType, ok := m[k]
	if !ok {
		err = fmt.Errorf("%w: %v", ErrPairNotRegistered, k)
		return Pair{}, err
	}

	pair.Key = k

	switch vType {
	case "string":
		pair.Value, err = v, nil
	case "bool":
		pair.Value, err = strconv.ParseBool(v)
	case "int":
		var i int64
		i, err = strconv.ParseInt(v, 0, 0)
		pair.Value = int(i)
	case "int64":
		pair.Value, err = strconv.ParseInt(v, 0, 64)
	default:
		return Pair{}, fmt.Errorf("%w: %v, %v", ErrPairTypeNotParsable, k, vType)
	}

	if err != nil {
		pair = Pair{}
		err = fmt.Errorf("%w: %v, %v, %v: %v", ErrPairValueInvalid, k, vType, v, err)
	}
	return
}

func parseListMode(s string) (ListMode, error) {
	v64, err := strconv.ParseUint(s, 0, 8)
	if err == nil || errors.Is(err, strconv.ErrRange) {
		return ListMode(v64), nil
	}
	if strings.TrimSpace(s) == "" {
		return ListMode(0), nil
	}
	modes := strings.Split(s, "|")
	var l ListMode
	for _, mode := range modes {
		mode = strings.TrimSpace(mode)
		switch mode {
		case ListModeBlock.String():
			l |= ListModeBlock
		case ListModeDir.String():
			l |= ListModeDir
		case ListModePart.String():
			l |= ListModePart
		case ListModePrefix.String():
			l |= ListModePrefix
		default:
			return ListMode(0), errors.New("invalid mode " + mode)
		}
	}
	return l, nil
}
