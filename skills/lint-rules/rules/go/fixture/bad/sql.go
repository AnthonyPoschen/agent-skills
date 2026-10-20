package bad

import "database/sql"

func missingClose(db *sql.DB) error {
	rows, err := db.Query("select 1")
	if err != nil {
		return err
	}
	return rows.Err()
}

func missingErr(db *sql.DB) error {
	rows, err := db.Query("select 2")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return nil
}
