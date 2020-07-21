package twopc

import (
	"context"
)

// Querier represents query of database
type Querier interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Prepare(ctx context.Context) error
	// RollbackPrepared(ctx context.Context) error
	// CommitPrepared(ctx context.Context) error
}

// Do makes two process transaction
func Do(ctx context.Context, q1, q2 Querier /*queries ...Querier*/) error {
	var err error

	err = q1.Prepare(ctx)
	if err != nil {
		// fmt.Println("error: 1: prepare:", err)
		q1.Rollback(ctx)
		return err
	}

	err = q2.Prepare(ctx)
	if err != nil {
		// fmt.Println("error: 2: prepare:", err)
		q1.Rollback(ctx)
		q2.Rollback(ctx)
		return err
	}

	err = q1.Commit(ctx)
	if err != nil {
		// fmt.Println("error: 1: commit:", err)
		q1.Rollback(ctx)
		q2.Rollback(ctx)
		return err
	}

	err = q2.Commit(ctx)
	if err != nil {
		// fmt.Println("error: 2: commit:", err)
		q1.Rollback(ctx)
		q2.Rollback(ctx)
		return err
	}

	return nil
}
