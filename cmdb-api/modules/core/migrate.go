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

func InitBuiltinCITypes() error {
	var count int64
	if err := database.DB.Model(&CIType{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	builtinTypes := []struct {
		Name       string
		Alias      string
		UniqueAttr string
		Attributes []struct {
			Name       string
			Alias      string
			ValueType  string
			IsUnique   bool
			IsIndex    bool
			IsRequired bool
		}
	}{
		{
			Name: "Region", Alias: "地域", UniqueAttr: "name",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "name", Alias: "名称", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "description", Alias: "描述", ValueType: "text"},
			},
		},
		{
			Name: "IDC", Alias: "数据中心", UniqueAttr: "name",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "name", Alias: "名称", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "address", Alias: "地址", ValueType: "text"},
				{Name: "contact", Alias: "联系人", ValueType: "text"},
			},
		},
		{
			Name: "ServerRoom", Alias: "机房", UniqueAttr: "name",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "name", Alias: "名称", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "floor", Alias: "楼层", ValueType: "text"},
			},
		},
		{
			Name: "Rack", Alias: "机柜", UniqueAttr: "name",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "name", Alias: "名称", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "total_u", Alias: "总U位", ValueType: "integer", IsRequired: true},
			},
		},
		{
			Name: "Subnet", Alias: "子网", UniqueAttr: "cidr",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "cidr", Alias: "CIDR", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "vlan_id", Alias: "VLAN ID", ValueType: "text"},
			},
		},
		{
			Name: "IPAddress", Alias: "IP地址", UniqueAttr: "ip",
			Attributes: []struct {
				Name       string
				Alias      string
				ValueType  string
				IsUnique   bool
				IsIndex    bool
				IsRequired bool
			}{
				{Name: "ip", Alias: "IP地址", ValueType: "text", IsUnique: true, IsRequired: true},
				{Name: "status", Alias: "状态", ValueType: "choice", IsRequired: true},
			},
		},
	}

	for _, bt := range builtinTypes {
		var uniqueAttrID uint
		for i, attrDef := range bt.Attributes {
			attr := Attribute{
				Name:      attrDef.Name,
				Alias:     attrDef.Alias,
				ValueType: attrDef.ValueType,
				IsUnique:  attrDef.IsUnique,
				IsIndex:   attrDef.IsIndex,
			}
			if err := database.DB.Create(&attr).Error; err != nil {
				return err
			}
			if attrDef.Name == bt.UniqueAttr {
				uniqueAttrID = attr.ID
			}

			// Will create CITypeAttribute after CIType is created
			_ = i
		}

		ciType := CIType{
			Name:         bt.Name,
			Alias:        bt.Alias,
			UniqueAttrID: uniqueAttrID,
			IsBuiltin:    true,
		}
		if err := database.DB.Create(&ciType).Error; err != nil {
			return err
		}

		for _, attrDef := range bt.Attributes {
			var attr Attribute
			if err := database.DB.Where("name = ?", attrDef.Name).First(&attr).Error; err != nil {
				return err
			}
			link := CITypeAttribute{
				CITypeID:    ciType.ID,
				AttributeID: attr.ID,
				IsRequired:  attrDef.IsRequired,
			}
			if err := database.DB.Create(&link).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
