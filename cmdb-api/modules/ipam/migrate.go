package ipam

import "cmdb-api/database"

func Migrate() error {
	schemas := []interface{}{&Subnet{}, &IPAddress{}, &IPAMHistory{}}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	return nil
}
