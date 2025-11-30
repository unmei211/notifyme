package messaging

import "context"

type IFetcher interface {
	Fetch(ctx context.Context) error
	Fallback()
}

type IFetcherManager interface {
	Launch(ctx context.Context)
}
