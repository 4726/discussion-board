package common

import (
	"github.com/cenkalti/backoff/v3"
	"github.com/jinzhu/gorm"
)

func OpenDB(dialect string, args ...interface{}) (*gorm.DB, error) {
	var db *gorm.DB
	op := func() error {
		var err error
		db, err = gorm.Open(dialect, args...)
		return err
	}

	err := backoff.Retry(op, backoff.NewExponentialBackOff())

	return db, err
}
