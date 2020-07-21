package postgresrepo

import (
	"context"
	"fmt"
	"github.com/vladimir-kopaliani/2pc-example/internal/model"

	// postgres driver
	_ "github.com/lib/pq"
)

// DoSomething ...
func (r *Repository) DoSomething(ctx context.Context, smth model.Something) (Query, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Println("error: Postgres: Begin Transaction:", err)
		return Query{}, err
	}

	q := Query{
		tx: tx,
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO test_transactions (user_id, created_at) VALUES ($1, $2);`, smth.UserID, smth.CreatedAt.UTC())
	if err != nil {
		fmt.Println("error: Postgres: Insert:", err)
		return q, err
	}

	return q, nil
}
