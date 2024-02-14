package internal

import (
	"context"
	"github.com/go-pogo/serv/middleware"
)

type serverNameKey struct{}

func ServerNameMiddleware(name string) middleware.Middleware {
	return middleware.WithContextValue(serverNameKey{}, name)
}

func ServerName(ctx context.Context) string {
	if v := ctx.Value(serverNameKey{}); v != nil {
		return v.(string)
	}
	return ""
}
