# Two Process Commit Example

Implementation of 2pc protocol. It coordinates transaction among several databases if one of participant fail, other will be rollbacked, otherwise commited.

## How to Use

You have to implement `Querier` interface to be able to use `Coordinator`

```go
coord := twopc.NewCoordinatior()

// q1 and q2 are types which implement `Querier`

// Register will prepare transaction, in case of failure other transactions will be rollbacked
err := coordinator.Register(context.TODO(), &q1)
if err != nil {
  panic(err)
}

err = coordinator.Register(context.TODO(), &q2)
if err != nil {
  panic(err)
}

// Do will commit transactions
err = coordinator.Do(context.TODO())
if err != nil {
  panic(err)
}
```

## Examples of implementation `Querier` interface

### MongoDB

This is example for MongoDB at least 4.2 version with support transaction features.

```go
type Query struct {
	session mongo.Session
}

func (q *Query) Commit(ctx context.Context) error {
	defer q.session.EndSession(ctx)

	err := q.session.CommitTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (q *Query) Rollback(ctx context.Context) error {
	defer q.session.EndSession(ctx)

	err := q.session.AbortTransaction(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (q *Query) Prepare(ctx context.Context) error {
	return nil
}
```

Usage:

```go
func InsertSomething(ctx context.Context, smth interface{}) (Query, error) {
	q := Query{}

	session, err := client.StartSession()
	if err != nil {
		return q, err
	}

	err = session.StartTransaction()
	if err != nil {
		return q, err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := transactionCollection.InsertOne(ctx, smth)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return q, err
	}

	return q, nil
}
```

### Postrges

Before you have to run command below (with desired amount) to turn on prepared transactions:

```sql
ALTER SYSTEM SET max_prepared_transactions = 64;
```

```go
type Query struct {
	tx   *sql.Tx
	uuid string
}

func (q *Query) Commit(ctx context.Context) error {
	_, err := q.tx.QueryContext(ctx, "COMMIT PREPARED '"+q.uuid+"';")
	if err != nil {
		return err
	}

	return nil
}

func (q *Query) Rollback(ctx context.Context) error {
	_, err := q.tx.QueryContext(ctx, "ROLLBACK PREPARED '"+q.uuid+"';")
	if err != nil {
		return err
	}

	return nil
}

// uuid here is "github.com/satori/go.uuid"

func (q *Query) Prepare(ctx context.Context) error {
	q.uuid = uuid.NewV4().String()

	_, err := q.tx.QueryContext(ctx, "PREPARE TRANSACTION '"+q.uuid+"';")
	return err
}
```
usage:

```go
func (r *Repository) InsertSomething(ctx context.Context, smth interface{}) (Query, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Query{}, err
	}

	q := Query{
		tx: tx,
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO test_transactions (smth) VALUES ($1);`, smth)
	if err != nil {
		return q, err
	}

	return q, nil
}
```
