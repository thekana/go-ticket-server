package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func init() {
	Migrations = append(Migrations, &Migration{
		Number: 5,
		Name:   "Create reservations table",
		Forwards: func(db *gorm.DB) error {
			const sql = `
				CREATE TABLE reservations ( 
					reservation_id         BIGSERIAL PRIMARY KEY NOT NULL,
					user_id                BIGINT,
					event_id               BIGINT,
					quota                  BIGINT NOT NULL,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);
				ALTER TABLE reservations
					ADD CONSTRAINT FK_reserved_by
						FOREIGN KEY (user_id)
							REFERENCES users (id)                
				ON DELETE SET NULL
				;
				ALTER TABLE reservations
					ADD CONSTRAINT FK_reserved_for
						FOREIGN KEY (event_id)
							REFERENCES events (id)
				ON DELETE CASCADE
				;
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create reservations table")
		},
	})
}
