package twopc

import (
	"context"
	"sync"
)

// Querier represents query of database
type Querier interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Prepare(ctx context.Context) error
}

// Do makes two process transaction
func Do(ctx context.Context, q1, q2 Querier /*queries ...Querier*/) error {
	var err error

	err = q1.Prepare(ctx)
	if err != nil {
		q1.Rollback(ctx)
		return err
	}

	err = q2.Prepare(ctx)
	if err != nil {
		q1.Rollback(ctx)
		q2.Rollback(ctx)
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = q1.Commit(ctx)
		if err != nil {
			q1.Rollback(ctx)
			q2.Rollback(ctx)
			// return err
		}
	}()

	go func() {
		defer wg.Done()
		err = q2.Commit(ctx)
		if err != nil {
			q1.Rollback(ctx)
			q2.Rollback(ctx)
			// return err
		}
	}()

	wg.Wait()
	return err
}
