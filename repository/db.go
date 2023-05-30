package repository

import (
	"fmt"

	"nikwallet/config"
	"nikwallet/repository/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQL struct {
	DB *gorm.DB
}

func (p *PostgreSQL) Connect(c *config.Config) error {
	dsn := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", c.DbUser, c.DbPassword, c.DbName)
	var err error
	p.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	err = p.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Ledger{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate tables: %w", err)
	}

	return nil
}

func (p *PostgreSQL) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	err = p.DB.Migrator().DropTable(
		&models.User{},
		&models.Wallet{},
		&models.Ledger{},
	)
	if err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}

	return sqlDB.Close()
}
