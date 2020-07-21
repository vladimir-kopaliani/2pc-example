package mongorepo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// Query ...
type Query struct {
	session mongo.Session
}

// Commit ...
func (q *Query) Commit(ctx context.Context) error {
	defer q.session.EndSession(ctx)

	err := q.session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Rollback ...
func (q *Query) Rollback(ctx context.Context) error {
	defer q.session.EndSession(ctx)

	err := q.session.AbortTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

// Prepare ...
func (q *Query) Prepare(ctx context.Context) error {
	return nil
}

// // RollbackPrepared ...
// func (q *Query) RollbackPrepared(ctx context.Context) error {
// 	defer q.session.EndSession(ctx)
//
// 	err := q.session.AbortTransaction(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// // CommitPrepared ...
// func (q *Query) CommitPrepared(ctx context.Context) error {
// 	defer q.session.EndSession(ctx)
//
// 	err := q.session.CommitTransaction(ctx)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
