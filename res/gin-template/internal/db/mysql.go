package db

import (
	"github.com/IfanTsai/go-lib/mysql"
	"github.com/pkg/errors"
)

var (
	db *mysql.DB

	models = []interface{}{&User{}}
)

func Init(filename string) error {
	var err error
	db, err = mysql.NewDBInConfig(filename)
	if err != nil {
		return err
	}

	return migrate()
}

func migrate() error {
	for _, model := range models {
		if err := db.Write.AutoMigrate(model); err != nil {
			return errors.Wrapf(err, "failed to auto migrate %v to master db", model)
		}

		if err := db.Read.AutoMigrate(model); err != nil {
			return errors.Wrapf(err, "failed to auto migrate %v to slave db", model)
		}
	}

	return nil
}
