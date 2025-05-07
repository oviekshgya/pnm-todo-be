package pkg

import (
	"fmt"
	"gorm.io/gorm"
)

func WithTransaction(db *gorm.DB, fn func(tz *gorm.DB) (interface{}, error)) (interface{}, error) {

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if tx.Error != nil {
			_ = tx.Rollback()
		} else {
			cerr := tx.Commit().Error
			if cerr != nil {
				tx.Error = fmt.Errorf("error committing transaction: %v", cerr)
			}
		}
	}()
	res, err := fn(tx)
	if err != nil {
		tx.Error = err
		return nil, err
	}

	return res, nil
}
