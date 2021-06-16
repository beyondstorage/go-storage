package services

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/beyondstorage/go-storage/v4/pairs"
	. "github.com/beyondstorage/go-storage/v4/types"
)

var (
	// ErrConnectionStringInvalid means the connection string is invalid.
	ErrConnectionStringInvalid = NewErrorCode("connection string is invalid")
)

// <type>://[<name>][<work_dir>][?key1=value1&...&keyN=valueN]
func parseString(ConnStr string) (ty string, ps []Pair, err error) {
	colon := strings.Index(ConnStr, ":")
	if colon == -1 {
		err = fmt.Errorf("%w: %s, %s", ErrConnectionStringInvalid, "service type missing", ConnStr)
		return
	}
	ty = ConnStr[:colon]
	rest := ConnStr[colon+1:]
	m, ok := servicePairMaps[ty]
	if !ok {
		err = ErrServiceNotRegistered
		return
	}

	if !strings.HasPrefix(rest, "//") {
		err = fmt.Errorf("%w: %s", ErrConnectionStringInvalid, ConnStr)
		return
	}
	rest = rest[2:]

	// [<name>][<work_dir>][?key1=value1&...&keyN=valueN]
	// <name> does not contain '/'
	// <work_dir> begins with '/'
	question := strings.Index(rest, "?")
	var path string
	if question == -1 {
		path = rest
		rest = ""
	} else {
		path = rest[:question]
		rest = rest[question+1:]
	}

	if len(path) == 0 {
		err = fmt.Errorf("%w: %s, %s", ErrConnectionStringInvalid, "both <name> and <work_dir> missing", ConnStr)
		return
	} else {
		slash := strings.Index(path, "/")
		if slash == -1 {
			name := path
			ps = append(ps, pairs.WithName(name))
		} else if slash == 0 {
			work_dir := path
			ps = append(ps, pairs.WithWorkDir(work_dir))
		} else {
			name := path[:slash]
			work_dir := path[slash:]
			ps = append(ps, pairs.WithName(name), pairs.WithWorkDir(work_dir))
		}
	}

	for _, v := range strings.Split(rest, "&") {
		opt := strings.SplitN(v, "=", 2)
		if len(opt) != 2 {
			// && or &key&, ignore
			continue
		}
		pair, err1 := parse(m, opt[0], opt[1])
		if err1 != nil {
			ps = nil
			err = fmt.Errorf("%w: %v", ErrConnectionStringInvalid, err1)
			return
		}
		ps = append(ps, pair)
	}
	return
}

func parse(m map[string]string, k string, v string) (pair Pair, err error) {
	vType, ok := m[k]
	if !ok {
		err = fmt.Errorf("pair not registered: %v", k)
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
	case "[]byte":
		pair.Value, err = base64.RawStdEncoding.DecodeString(v)
	default:
		return Pair{}, fmt.Errorf("type not parseable: %v, %v", k, vType)
	}

	if err != nil {
		pair = Pair{}
		err = fmt.Errorf("pair value invalid: %v, %v, %v: %v", k, vType, v, err)
	}
	return
}
