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
	qs     []Querier
	logger *Logger
}

// NewCoordinatior returns new instance of Coordinator
func NewCoordinatior() *Coordinator {
	return &Coordinator{
		qs: make([]Querier, 0, 2),
	}
}

// Do makes two process transaction
func (c Coordinator) Do(ctx context.Context) error {
	var err error

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "starting new transaction")
	}

	// commit phase
	wg := sync.WaitGroup{}
	wg.Add(len(c.qs))

	for i := range c.qs {
		go func(i int) {
			defer wg.Done()
			err = c.qs[i].Commit(ctx)
			if err != nil {
				if c.logger != nil {
					(*c.logger).Error(logPrefix, "error: commiting:", err)
				}
			}

			if c.logger != nil {
				(*c.logger).Debug(logPrefix, "commited")
			}
		}(i)
	}

	wg.Wait()

	if c.logger != nil {
		(*c.logger).Debug(logPrefix, "transaction finished")
	}

	return err
}

func (c *Coordinator) SetLogger(logger Logger) {
	c.logger = &logger
}

func (c *Coordinator) Register(ctx context.Context, q Querier) error {
	err := q.Prepare(ctx)
	if err != nil {
		if c.logger != nil {
			(*c.logger).Error(logPrefix, "error: prepare:", err)
		}

		q.Rollback(ctx)

		for i := range c.qs {
			if c.logger != nil {
				(*c.logger).Debug(logPrefix, "rollbacked")
			}

			c.qs[i].Rollback(ctx)
		}

		return err
	}

	c.qs = append(c.qs, q)
	return nil
}
