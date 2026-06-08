package discovery

import "cmdb-api/database"

func Migrate() error {
	schemas := []interface{}{
		&DiscoveryRule{}, &DiscoveryTask{}, &DiscoveryResult{}, &Agent{},
	}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	return nil
}
