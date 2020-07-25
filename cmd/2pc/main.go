package main

import (
	"context"
	"time"

	twopc "github.com/vladimir-kopaliani/2pc-example/internal/2pc"
	"github.com/vladimir-kopaliani/2pc-example/internal/logger"
	"github.com/vladimir-kopaliani/2pc-example/internal/model"

	mongorepo "github.com/vladimir-kopaliani/2pc-example/internal/mongo"
	postgresrepo "github.com/vladimir-kopaliani/2pc-example/internal/postgres"

	uuid "github.com/satori/go.uuid"
)

func main() {
	ctx /*, cancel*/, _ := context.WithCancel(context.Background())

	// // handle interupt signal
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// go func() {
	// 	select {
	// 	case <-signalChan:
	// 		log.Println("Got Interrupt signal. Shutting down...")
	// 		cancel()
	// 	}
	// }()

	// ---

	l := logger.Logger{}

	repo1, err := mongorepo.New(ctx, &mongorepo.Configuration{
		URI: "mongodb://mongo:27017/test_db",
	})
	if err != nil {
		panic(err)
	}

	repo2, err := postgresrepo.New(ctx, &postgresrepo.Configuration{
		Host:     "localhost",
		User:     "postgres",
		Password: "postgres",
		DBName:   "test_db",
	})
	if err != nil {
		panic(err)
	}

	smth := model.Something{
		UserID:    uuid.NewV4().String(),
		CreatedAt: time.Now(),
	}

	q1, err := repo1.DoSomething(ctx, smth)
	if err != nil {
		panic(err)
	}

	q2, err := repo2.DoSomething(ctx, smth)
	if err != nil {
		panic(err)
	}

	coord := twopc.NewCoordinatior()
	coord.SetLogger(l)

	err = coord.Do(ctx, &q1, &q2)
	if err != nil {
		panic(err)
	}

	// ---

	// <-ctx.Done()
}
