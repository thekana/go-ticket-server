package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var createAdminTableMigration = &Migration{
	Number: 2,
	Name:   "Create admin table",
	Forwards: func(db *gorm.DB) error {
		const sql = `
			CREATE TABLE admins(
 				id BIGSERIAL PRIMARY KEY NOT NULL,
				username VARCHAR(255) UNIQUE NOT NULL,
				active BOOLEAN NOT NULL DEFAULT TRUE,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
		`

		err := db.Exec(sql).Error
		return errors.Wrap(err, "unable to create users table")
	},
}

func init() {
	Migrations = append(Migrations, createAdminTableMigration)
}
