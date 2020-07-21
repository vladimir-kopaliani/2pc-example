package postgresrepo

import (
	"context"
	"database/sql"
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

// Query ...
type Query struct {
	tx   *sql.Tx
	uuid string
}

// Commit ...
func (q *Query) Commit(ctx context.Context) error {
	// var err error
	// if q != nil && q.tx != nil {
	// 	err = q.tx.Commit()
	// }
	//
	// return err

	_, err := q.tx.QueryContext(ctx, "COMMIT PREPARED '"+q.uuid+"';")
	if err != nil {
		fmt.Println("error: Postgres: Commit Transaction:", q.uuid, err)
		return err
	}

	fmt.Println("Postgres: Commit Transaction:", q.uuid)

	return nil
}

// Rollback ...
func (q *Query) Rollback(ctx context.Context) error {
	// var err error
	// if q != nil && q.tx != nil {
	// 	err = q.tx.Rollback()
	// }
	//
	// return err

	_, err := q.tx.QueryContext(ctx, "ROLLBACK PREPARED '"+q.uuid+"';")
	if err != nil {
		fmt.Println("error: Postgres: Rollback Transaction:", q.uuid, err)
		return err
	}

	fmt.Println("Postgres: Rollback Transaction:", q.uuid)

	return nil
}

// Prepare ...
func (q *Query) Prepare(ctx context.Context) error {
	q.uuid = uuid.NewV4().String()
	fmt.Println("Postgres: Prepare Transaction:", q.uuid)
	_, err := q.tx.QueryContext(ctx, "PREPARE TRANSACTION '"+q.uuid+"';")
	return err
}

// // RollbackPrepared ...
// func (q *Query) RollbackPrepared(ctx context.Context) error {
// 	_, err := q.tx.QueryContext(ctx, "ROLLBACK PREPARED $1", q.uuid)
// 	return err
// }

// // CommitPrepared ...
// func (q *Query) CommitPrepared(ctx context.Context) error {
// 	_, err := q.tx.QueryContext(ctx, "COMMIT PREPARED $1", q.uuid)
// 	return err
// }
