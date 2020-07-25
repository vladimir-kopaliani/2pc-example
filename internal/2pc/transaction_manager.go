package twopc

import (
	"context"
	"sync"
)

const (
	logPrefix = "2pc coordinator:"
)

// Querier represents query of database
type Querier interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Prepare(ctx context.Context) error
}

type Logger interface {
	Debug(...interface{})
	Error(...interface{})
}

type Coordinator struct {
	logger *Logger
}

// NewCoordinatior returns new instance of Coordinator
func NewCoordinatior() *Coordinator {
	return &Coordinator{}
}

// Do makes two process transaction
func (c Coordinator) Do(ctx context.Context, q1, q2 Querier /*queries ...Querier*/) error {
	var err error

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "starting new transaction")
	}

	err = q1.Prepare(ctx)
	if err != nil {
		if c.logger != nil {
			(*c.logger).Error(logPrefix, "error: prepare 1st:", err)
		}

		q1.Rollback(ctx)
		return err
	}

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "1st is prepared")
	}

	err = q2.Prepare(ctx)
	if err != nil {
		if c.logger != nil {
			(*c.logger).Error(logPrefix, "error: prepare 2nd:", err)
		}

		q1.Rollback(ctx)
		q2.Rollback(ctx)
		return err
	}

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "2nd is prepared")
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = q1.Commit(ctx)
		if err != nil {
			q1.Rollback(ctx)
			q2.Rollback(ctx)

			if c.logger != nil {
				(*c.logger).Error(logPrefix, "error: commiting 1st:", err)
			}

			// return err
		}

		if c.logger != nil {
			(*c.logger).Debug(logPrefix, "1nd is commited")
		}
	}()

	go func() {
		defer wg.Done()
		err = q2.Commit(ctx)
		if err != nil {
			q1.Rollback(ctx)
			q2.Rollback(ctx)

			if c.logger != nil {
				(*c.logger).Error(logPrefix, "error: commiting 2nd:", err)
			}

			// return err
		}

		if c.logger != nil {
			(*c.logger).Debug(logPrefix, "1nd is commited")
		}
	}()

	wg.Wait()

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "transaction finished")
	}

	return err
}

func (c *Coordinator) SetLogger(logger Logger) {
	c.logger = &logger
}
