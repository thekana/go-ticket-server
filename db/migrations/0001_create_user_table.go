package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func init() {
	Migrations = append(Migrations, &Migration{
		Number: 1,
		Name:   "Create users table",
		Forwards: func(db *gorm.DB) error {
			const sql = `
				CREATE TABLE users(
					id BIGSERIAL PRIMARY KEY NOT NULL,
					username VARCHAR(255) UNIQUE NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create users table")
		},
	})
}
