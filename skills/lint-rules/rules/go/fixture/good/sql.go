package good

import "database/sql"

func query(db *sql.DB) error {
	rows, err := db.Query("select 1")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return rows.Err()
}
