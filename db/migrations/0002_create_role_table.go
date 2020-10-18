package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func init() {
	Migrations = append(Migrations, &Migration{
		Number: 2,
		Name:   "Create roles table",
		Forwards: func(db *gorm.DB) error {
			const sql = `
				CREATE TABLE roles(
					id SERIAL PRIMARY KEY NOT NULL,
					role VARCHAR(10) UNIQUE NOT NULL
				);
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create roles table")
		},
	})
}
