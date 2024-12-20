// internal/infrastructure/db/db.go
package db

import (
	"fmt"
	"sync"

	"github.com/Casagrande-Lucas/dnd/internal/domain/race/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB defines the contract for database operations and access.
type DB interface {
	GetDB() *gorm.DB
}

// dbImpl is the concrete implementation of DB.
type dbImpl struct {
	conn *gorm.DB
}

// GetDB returns the underlying *gorm.DB instance.
func (d *dbImpl) GetDB() *gorm.DB {
	return d.conn
}

// FactoryDB is responsible for creating and managing DB instances.
type FactoryDB struct {
	connections map[string]DB
	mu          sync.Mutex
}

var (
	factoryInstance *FactoryDB
	once            sync.Once
)

// GetDBFactory returns a singleton instance of FactoryDB.
func GetDBFactory() *FactoryDB {
	once.Do(func() {
		factoryInstance = &FactoryDB{
			connections: make(map[string]DB),
		}
	})
	return factoryInstance
}

// CreatePostgresConnection creates a new Postgres connection and registers it with the factory.
func (f *FactoryDB) CreatePostgresConnection(name, host, port, user, password, dbname string) (DB, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if conn, exists := f.connections[name]; exists {
		return conn, nil
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(
		&entities.Race{},
		&entities.Age{},
		&entities.Trait{},
		&entities.Subrace{},
		&entities.Language{},
		&entities.Proficiency{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	connection := &dbImpl{conn: db}
	f.connections[name] = connection
	return connection, nil
}

// GetConnection retrieves an existing connection by name.
func (f *FactoryDB) GetConnection(name string) (DB, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	conn, exists := f.connections[name]
	return conn, exists
}
