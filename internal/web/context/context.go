package context

import (
	"context"

	"github.com/glasskube/glasskube/internal/clicontext"

	"github.com/glasskube/glasskube/internal/web/types"
)

const coreListersContextKey clicontext.ContextKey = 100

func ContextWithCoreListers(parent context.Context, coreListers *types.CoreListers) context.Context {
	return context.WithValue(parent, coreListersContextKey, coreListers)
}

func CoreListersFromContext(ctx context.Context) *types.CoreListers {
	value := ctx.Value(coreListersContextKey)
	if value != nil {
		if coreListers, ok := value.(*types.CoreListers); ok {
			return coreListers
		}
	}
	return nil
}
