package database

import (
	"testing"
	"cmdb-api/config"
	"github.com/stretchr/testify/assert"
)

func TestInitPostgres(t *testing.T) {
	cfg := &config.Config{
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUser:     "cmdb_test",
		DBPassword: "cmdb_test",
		DBName:     "cmdb_test",
	}
	// This test requires a running PostgreSQL instance.
	// For CI/CD, use a test container or mock.
	// Skip if no database is available.
	t.Skip("Skipping: requires PostgreSQL test database")

	assert.NotPanics(t, func() { InitPostgres(cfg) })
	assert.NotNil(t, DB)
}
