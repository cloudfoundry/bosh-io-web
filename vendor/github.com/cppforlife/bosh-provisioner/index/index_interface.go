package index

import (
	"errors"
)

var (
	ErrNotFound = errors.New("Record is not found")

	// ErrChanged is returned from SaveLocked
	// if record with the same key was inserted by another client
	// between FindLocked-SaveLocked cycle; e.g.:
	//   client1: FindLocked -> ErrNotFound
	//   client2: FindLocked -> ErrNotFound
	//   client1: SaveLocked -> no err
	//   client2: SaveLocked -> ErrChanged
	ErrChanged = errors.New("Record changed")
)

type Index interface {
	ListKeys(interface{}) error
	List(interface{}) error

	Find(interface{}, interface{}) error
	Save(interface{}, interface{}) error
	Remove(interface{}) error

	// FindLocked returns and locks record such that another client cannot read or write it
	// before associated locked record is released is called.
	// (In databases typically this is done with row level locks and transactions.)
	FindLocked(interface{}, interface{}) (LockedRecord, error)
}

type LockedRecord interface {
	Release() error

	Save(interface{}) error
	Remove() error
}
