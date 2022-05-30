package morpheus

import (
	"database/sql"

	"github.com/spoonboy-io/dozer/internal"
)

// GetProcessTypes is used to initialise a map of code/names for processTypes
// since ideally we'll reference by code in the YAML config but need to use name when
// looking for matches in the process data or we need to complicate the SQL query
func GetProcessTypes(db *sql.DB, pt map[string]string) error {
	rows, err := db.Query("SELECT id, code, name, image_code FROM process_type;")
	if err != nil {
		return err
	}

	for rows.Next() {
		var processType internal.ProcessType
		err := rows.Scan(&processType.Id, &processType.Code, &processType.Name, &processType.ImageCode)
		if err != nil {
			return err
		}

		pt[processType.Code.String] = processType.Name.String
	}

	// we will keep the data in the internal namespace
	internal.ProcessTypes = pt

	return nil
}
