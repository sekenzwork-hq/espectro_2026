package database

import "gorm.io/gorm"

type TransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return TransactionManager{db: db}
}

func (t TransactionManager) Run(fn func() error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}
