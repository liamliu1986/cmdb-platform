package dcim

import "cmdb-api/database"

func Migrate() error {
	schemas := []interface{}{&IDC{}, &ServerRoom{}, &Rack{}, &RackLayout{}, &LocationCoord{}}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	return nil
}
