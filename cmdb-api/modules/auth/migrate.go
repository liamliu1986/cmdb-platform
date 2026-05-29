package auth

import "cmdb-api/database"

func Migrate() error {
	schemas := []interface{}{
		&User{}, &Role{}, &RoleRelation{},
		&ResourceType{}, &Resource{}, &ResourceGroup{}, &ResourceGroupItem{},
		&Permission{}, &RolePermission{}, &UserRole{},
	}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	return nil
}
