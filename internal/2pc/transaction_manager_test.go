package twopc

import (
	"context"
	"errors"
	"testing"
)

const (
	initial    = "initial"
	prepeared  = "prepeared"
	commited   = "commited"
	rollbacked = "rollbacked"
)

const (
	prepearedErr = "prepearedErr"
)

type mock struct {
	value  string
	status string
	_value string
	err    *string
}

func (m *mock) Commit(ctx context.Context) error {
	if m.status != prepeared {
		return errors.New("not prepeared")
	}

	m.status = commited
	m.value = m._value
	m._value = ""

	return nil
}

func (m *mock) Rollback(ctx context.Context) error {
	m.status = rollbacked
	m._value = ""
	return nil
}

func (m *mock) Prepare(ctx context.Context) error {
	if m.err != nil {
		return errors.New("error")
	}

	if m.status != initial {
		return errors.New("not initial")
	}

	m.status = prepeared

	return nil
}

func (m *mock) Change() error {
	m.status = initial
	m._value = "new"

	return nil
}

func (m *mock) setPreparedError() {
	pr := prepearedErr
	m.err = &pr
}

func TestDo(t *testing.T) {
	test1(t)
	test2(t)
	test3(t)
	test4(t)
}

func test1(t *testing.T) {
	m1 := mock{
		status: initial,
	}
	m2 := mock{
		status: initial,
	}

	m1.Change()
	m2.Change()

	err := Do(context.TODO(), &m1, &m2)
	if err != nil {
		t.Error(err)
	}

	if m1.status != commited {
		t.Errorf("expected %q, got %q", commited, m1.status)
	}
	if m2.status != commited {
		t.Errorf("expected %q, got %q", commited, m2.status)
	}
	if m1.value != "new" {
		t.Errorf("expected %q, got %q", "new", m1.value)
	}
	if m2.value != "new" {
		t.Errorf("expected %q, got %q", "new", m2.value)
	}
}

func test2(t *testing.T) {
	m1 := mock{
		status: initial,
	}
	m2 := mock{
		status: initial,
	}

	m1.Change()
	m2.Change()

	m1.setPreparedError()

	err := Do(context.TODO(), &m1, &m2)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	if m1.status != rollbacked {
		t.Errorf("expected %q, got %q", rollbacked, m1.status)
	}
	if m2.status != initial {
		t.Errorf("expected %q, got %q", initial, m2.status)
	}
	if m1.value != "" {
		t.Errorf("expected %q, got %q", "", m1.value)
	}
	if m2.value != "" {
		t.Errorf("expected %q, got %q", "", m2.value)
	}
}

func test3(t *testing.T) {
	m1 := mock{
		status: initial,
	}
	m2 := mock{
		status: initial,
	}

	m1.Change()
	m2.Change()

	m2.setPreparedError()

	err := Do(context.TODO(), &m1, &m2)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	if m1.status != rollbacked {
		t.Errorf("expected %q, got %q", rollbacked, m1.status)
	}
	if m2.status != rollbacked {
		t.Errorf("expected %q, got %q", rollbacked, m2.status)
	}
	if m1.value != "" {
		t.Errorf("expected %q, got %q", "", m1.value)
	}
	if m2.value != "" {
		t.Errorf("expected %q, got %q", "", m2.value)
	}
}

func test4(t *testing.T) {
	m1 := mock{
		status: initial,
	}
	m2 := mock{
		status: initial,
	}

	m1.Change()
	m2.Change()

	m1.setPreparedError()
	m2.setPreparedError()

	err := Do(context.TODO(), &m1, &m2)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}

	if m1.status != rollbacked {
		t.Errorf("expected %q, got %q", rollbacked, m1.status)
	}
	if m2.status != initial {
		t.Errorf("expected %q, got %q", initial, m2.status)
	}
	if m1.value != "" {
		t.Errorf("expected %q, got %q", "", m1.value)
	}
	if m2.value != "" {
		t.Errorf("expected %q, got %q", "", m2.value)
	}
}
