package core

import (
	"fmt"
	"regexp"
	"strings"

	"cmdb-api/database"
	"gorm.io/gorm"
)

var safeSortField = regexp.MustCompile("^[a-zA-Z0-9_]+$")

type CISearchBuilder struct {
	db       *gorm.DB
	query    string
	page     int
	pageSize int
	sort     string
}

func NewCISearchBuilder() *CISearchBuilder {
	return &CISearchBuilder{
		db:       database.DB,
		page:     1,
		pageSize: 25,
	}
}

func (b *CISearchBuilder) WithQuery(q string) *CISearchBuilder {
	b.query = q
	return b
}

func (b *CISearchBuilder) WithPagination(page, pageSize int) *CISearchBuilder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 25
	}
	b.page = page
	b.pageSize = pageSize
	return b
}

func (b *CISearchBuilder) WithSort(sort string) *CISearchBuilder {
	b.sort = sort
	return b
}

func (b *CISearchBuilder) Build() (*gorm.DB, error) {
	db := b.db.Model(&CI{}).Where("deleted_at IS NULL")

	if b.query == "" {
		return db, nil
	}

	conditions := strings.Split(b.query, ",")
	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)
		if cond == "" {
			continue
		}

		// Handle _type:Server
		if strings.HasPrefix(cond, "_type:") {
			typeName := strings.TrimPrefix(cond, "_type:")
			db = db.Where("ci_type_id IN (SELECT id FROM cmdb_core.ci_types WHERE name = ?)", typeName)
			continue
		}

		// Handle attr:value
		parts := strings.SplitN(cond, ":", 2)
		if len(parts) != 2 {
			continue
		}
		attrName := parts[0]
		attrValue := parts[1]

		// Comparison operators
		if strings.HasPrefix(attrValue, ">=") {
			val := strings.TrimPrefix(attrValue, ">=")
			db = db.Where("attr_values->>? >= ?", attrName, val)
		} else if strings.HasPrefix(attrValue, ">") {
			val := strings.TrimPrefix(attrValue, ">")
			db = db.Where("attr_values->>? > ?", attrName, val)
		} else if strings.HasPrefix(attrValue, "<=") {
			val := strings.TrimPrefix(attrValue, "<=")
			db = db.Where("attr_values->>? <= ?", attrName, val)
		} else if strings.HasPrefix(attrValue, "<") {
			val := strings.TrimPrefix(attrValue, "<")
			db = db.Where("attr_values->>? < ?", attrName, val)
		} else if strings.Contains(attrValue, "*") {
			// Wildcard LIKE
			pattern := strings.ReplaceAll(attrValue, "*", "%")
			db = db.Where("attr_values->>? LIKE ?", attrName, pattern)
		} else {
			// Exact match
			db = db.Where("attr_values->>? = ?", attrName, attrValue)
		}
	}

	// Sorting
	if b.sort != "" {
		direction := "ASC"
		field := b.sort
		if strings.HasPrefix(b.sort, "-") {
			direction = "DESC"
			field = strings.TrimPrefix(b.sort, "-")
		}
		if !safeSortField.MatchString(field) {
			return nil, fmt.Errorf("invalid sort field: %s", field)
		}
		db = db.Order(fmt.Sprintf("attr_values->>'%s' %s NULLS LAST", field, direction))
	}

	return db, nil
}

func (b *CISearchBuilder) Execute() ([]CI, int64, error) {
	db, err := b.Build()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	db.Count(&total)

	var cis []CI
	err = db.Offset((b.page - 1) * b.pageSize).Limit(b.pageSize).Find(&cis).Error
	return cis, total, err
}
