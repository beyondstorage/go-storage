package services

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beyondstorage/go-storage/v5/pairs"
	"github.com/beyondstorage/go-storage/v5/types"
)

type (
	// NewServicerFunc is a function that can initiate a new servicer.
	NewServicerFunc func(ps ...types.Pair) (types.Servicer, error)
	// NewStoragerFunc is a function that can initiate a new storager.
	NewStoragerFunc func(ps ...types.Pair) (types.Storager, error)
)

var (
	servicerFnMap map[string]NewServicerFunc
	storagerFnMap map[string]NewStoragerFunc
)

// RegisterServicer will register a servicer.
func RegisterServicer(ty string, fn NewServicerFunc) {
	servicerFnMap[ty] = fn
}

// RegisterStorager will register a storager.
func RegisterStorager(ty string, fn NewStoragerFunc) {
	storagerFnMap[ty] = fn
}

func init() {
	servicerFnMap = make(map[string]NewServicerFunc)
	storagerFnMap = make(map[string]NewStoragerFunc)
}

var (
	servicePairMaps map[string]map[string]string
)

// RegisterSchema will register a service's pair map.
//
// Users SHOULD NOT call this function.
func RegisterSchema(ty string, m map[string]string) {
	servicePairMaps[ty] = m
}

var (
	// ErrConnectionStringInvalid means the connection string is invalid.
	ErrConnectionStringInvalid = NewErrorCode("connection string is invalid")
)

// <type>://[<name>][<work_dir>][?key1=value1&...&keyN=valueN]
func parseConnectionString(ConnStr string) (ty string, ps []types.Pair, err error) {
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
		// both <name> and <work_dir> missing
	} else {
		slash := strings.Index(path, "/")
		if slash == -1 {
			name := path
			ps = append(ps, pairs.WithName(name))
		} else if slash == 0 {
			workDir := path
			ps = append(ps, pairs.WithWorkDir(workDir))
		} else {
			name := path[:slash]
			workDir := path[slash:]
			ps = append(ps, pairs.WithName(name), pairs.WithWorkDir(workDir))
		}
	}

	for _, v := range strings.Split(rest, "&") {
		opt := strings.SplitN(v, "=", 2)
		if len(opt) != 2 {
			if v != "" {
				opt = append(opt, "true")
			} else {
				// &&, ignore
				continue
			}
		}
		pair, parseErr := parse(m, opt[0], opt[1])
		if parseErr != nil {
			ps = nil
			err = fmt.Errorf("%w: %v", ErrConnectionStringInvalid, parseErr)
			return
		}
		ps = append(ps, pair)
	}
	return
}

func parse(m map[string]string, k string, v string) (pair types.Pair, err error) {
	vType, ok := m[k]
	if !ok {
		err = fmt.Errorf("pair not registered: %v", k)
		return types.Pair{}, err
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
	case "uint64":
		pair.Value, err = strconv.ParseUint(v, 0, 64)
	case "[]byte":
		pair.Value, err = base64.RawStdEncoding.DecodeString(v)
	case "time.Duration":
		var i int64
		i, err = strconv.ParseInt(v, 0, 64)
		pair.Value = time.Duration(i)
	default:
		return types.Pair{}, fmt.Errorf("type not parseable: %v, %v", k, vType)
	}

	if err != nil {
		pair = types.Pair{}
		err = fmt.Errorf("pair value invalid: %v, %v, %v: %v", k, vType, v, err)
	}
	return
}

func init() {
	servicePairMaps = make(map[string]map[string]string)
}
