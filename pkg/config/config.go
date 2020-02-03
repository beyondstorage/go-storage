package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Xuanwo/storage/pkg/credential"
	"github.com/Xuanwo/storage/pkg/endpoint"
	"github.com/Xuanwo/storage/types"
	"github.com/Xuanwo/storage/types/pairs"
)

var (
	// ErrInvalidConfig will be returned when config is invalid.
	ErrInvalidConfig = errors.New("invalid config")
)

// Parse will parse config string and return service type and namespace.
//
// TODO: Options didn't support for now.
func Parse(cfg string) (t string, opt []*types.Pair, err error) {
	errorMessage := "parse config [%s]: <%w>"

	// Parse type from: "<type>://<config>"
	s := strings.Split(cfg, "://")
	if len(s) != 2 || s[0] == "" || s[1] == "" {
		err = fmt.Errorf(errorMessage, cfg, ErrInvalidConfig)
		return
	}
	t = s[0]

	// Split <credential>@<endpoint>/<name>?<options> into tow parts.
	s = strings.SplitN(s[1], "/", 2)
	// Handle credential and endpoint
	// credential and endpoint could be missing for some service.
	if s[0] != "" {
		// Split <credential>@<endpoint> into tow parts.
		ce := strings.Split(s[0], "@")
		if len(ce) == 0 || len(ce) > 2 {
			return "", nil, fmt.Errorf(errorMessage, cfg, ErrInvalidConfig)
		}

		// We always have credential part
		cred, err := credential.Parse(ce[0])
		if err != nil {
			return "", nil, fmt.Errorf(errorMessage, cfg, err)
		}
		opt = append(opt, pairs.WithCredential(cred))

		if len(ce) == 2 {
			end, err := endpoint.Parse(ce[1])
			if err != nil {
				return "", nil, fmt.Errorf(errorMessage, cfg, err)
			}
			opt = append(opt, pairs.WithEndpoint(end))
		}
	}
	// Handle name and options.
	s = strings.SplitN(s[1], "?", 2)
	if len(s) == 1 {
		// We don't have options, return directly.
		return
	}
	ops := strings.Split(s[1], "&")
	for _, v := range ops {
		x := strings.SplitN(v, "=", 2)
		// TODO: type support
		opt = append(opt, &types.Pair{
			Key:   x[0],
			Value: x[1],
		})
	}
	return
}
