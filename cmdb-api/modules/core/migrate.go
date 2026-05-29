package core

import "cmdb-api/database"

func Migrate() error {
	schemas := []interface{}{
		&CIType{}, &Attribute{}, &CITypeAttribute{},
		&RelationType{}, &CITypeRelation{},
		&CI{}, &CIRelation{}, &OperationLog{},
	}
	for _, schema := range schemas {
		if err := database.DB.AutoMigrate(schema); err != nil {
			return err
		}
	}
	return nil
}
