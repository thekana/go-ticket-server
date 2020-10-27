package migrations

// Actually not used
func init() {
	//Migrations = append(Migrations, &Migration{
	//	Number: 6,
	//	Name:   "Create triggers and functions",
	//	Forwards: func(db *gorm.DB) error {
	//		const sql = `
	//		CREATE OR REPLACE FUNCTION make_res() RETURNS TRIGGER AS $make_res$
	//			DECLARE
	//			q integer;
	//			s integer;
	//			BEGIN
	//				select quota, sold into q, s from events
	//				where events.id = NEW.event_id;
	//				IF FOUND THEN
	//					RAISE NOTICE 'RemainingQuota % Sold % Want %', q,s, NEW.quota;
	//					IF (q - s >= NEW.quota) THEN
	//						UPDATE events SET sold = sold + NEW.quota where events.id = NEW.event_ID;
	//					ELSE
	//						RAISE EXCEPTION 'NOT ENOUGH QUOTA';
	//						RETURN NULL;
	//					END IF;
	//					RETURN NEW;
	//				END IF;
	//				RAISE EXCEPTION 'EVENT NOT FOUND';
	//				RETURN NULL;
	//			END;
	//		$make_res$ LANGUAGE plpgsql;
	//
	//		CREATE TRIGGER make_reservation BEFORE INSERT ON reservations
	//			FOR EACH ROW EXECUTE FUNCTION make_res();
	//
	//		DROP TRIGGER IF EXISTS make_reservation on reservations;
	//
	//		CREATE OR REPLACE FUNCTION cancel_res() RETURNS TRIGGER AS $cancel_res$
	//			BEGIN
	//				UPDATE events SET sold = sold - OLD.quota WHERE events.id = OLD.event_ID;
	//				RETURN OLD;
	//			END;
	//		$cancel_res$ LANGUAGE plpgsql;
	//
	//		CREATE TRIGGER cancel_reservation AFTER DELETE ON reservations
	//			FOR EACH ROW EXECUTE FUNCTION cancel_res();
	//
	//		DROP TRIGGER IF EXISTS cancel_reservation on reservations;
	//		`
	//		err := db.Exec(sql).Error
	//		return errors.Wrap(err, "unable to create functions")
	//	},
	//})
}
