package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func init() {
	Migrations = append(Migrations, &Migration{
		Number: 3,
		Name:   "Create user_roles table",
		Forwards: func(db *gorm.DB) error {
			const sql = `
			CREATE TABLE user_roles (
 				id SERIAL PRIMARY KEY,
				user_id         BIGINT NOT NULL,
				role_id         INTEGER NOT NULL
			  )
			  ;
			  ALTER TABLE user_roles
				  ADD CONSTRAINT FK_REFERENCE_1
					  FOREIGN KEY (user_id)
						  REFERENCES users (id)
			  ON DELETE CASCADE
			   ;
			  ALTER TABLE user_roles
				  ADD CONSTRAINT FK_REFERENCE_2
					  FOREIGN KEY (role_id)
						  REFERENCES roles (id)
			  ON DELETE CASCADE
			   ;
			`
			err := db.Exec(sql).Error
			return errors.Wrap(err, "unable to create user_roles table")
		},
	})
}
