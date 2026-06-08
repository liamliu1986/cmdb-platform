package ipam

import "cmdb-api/database"

func Migrate() error {
	schemas := []any{&Subnet{}, &IPAddress{}, &IPAMHistory{}, &UserIPAddress{}}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	database.DB.Exec("CREATE INDEX IF NOT EXISTS idx_ip_addresses_ci_id_status ON cmdb_ipam.ip_addresses(ci_id, status)")
	return nil
}
