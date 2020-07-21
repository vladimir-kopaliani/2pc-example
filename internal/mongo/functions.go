package mongorepo

import (
	"context"
	"fmt"
	"github.com/vladimir-kopaliani/2pc-example/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
)

// DoSomething ...
func (r *Repository) DoSomething(ctx context.Context, smth model.Something) (Query, error) {
	q := Query{}

	session, err := r.client.StartSession()
	if err != nil {
		fmt.Println("Mongo: Start session: ", err)
		return q, err
	}
	q.session = session

	err = session.StartTransaction()
	if err != nil {
		fmt.Println("Mongo: Start transaction: ", err)
		return q, err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := r.transactionCollection.InsertOne(ctx, smth)
		if err != nil {
			fmt.Println("Mongo: Insert: ", err)
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println("Mongo: With Session: ", err)
		return q, err
	}

	return q, nil
}
