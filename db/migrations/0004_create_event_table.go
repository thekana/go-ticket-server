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
					id            	BIGSERIAL PRIMARY KEY NOT NULL,
					name          	TEXT NOT NULL,
					owner         	BIGINT NOT NULL references users(id),
					quota         	BIGINT NOT NULL CHECK(quota >= remaining_quota),
					remaining_quota BIGINT NOT NULL CHECK(remaining_quota >= 0),
					created_at    	TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at    	TIMESTAMPTZ NOT NULL DEFAULT NOW()
				  );
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create events table")
		},
	})
}
