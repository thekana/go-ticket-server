package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func init() {
	Migrations = append(Migrations, &Migration{
		Number: 4,
		Name:   "Create events table",
		Forwards: func(db *gorm.DB) error {
			const sql = `
			CREATE TABLE events ( 
				id            BIGSERIAL PRIMARY KEY NOT NULL,
				name          TEXT NOT NULL,
				quota         BIGINT NOT NULL,
				owner         BIGINT NOT NULL,
				created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
			  );
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create events table")
		},
	})
}
