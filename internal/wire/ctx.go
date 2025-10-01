package wire

import "context"

func ProvideContext() context.Context {
    return context.Background()
}