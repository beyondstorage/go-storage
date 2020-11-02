package types

import (
	"context"
)

type Interceptor func(ctx context.Context, method string) func(error)

func ChainInterceptor(interceptors ...Interceptor) Interceptor {
	n := len(interceptors)
	fns := make([]func(error), n)

	return func(ctx context.Context, method string) func(err error) {
		for k, v := range interceptors {
			fn := v(ctx, method)
			fns[n-k-1] = fn
		}
		return func(err error) {
			for _, v := range fns {
				v(err)
			}
		}
	}
}
