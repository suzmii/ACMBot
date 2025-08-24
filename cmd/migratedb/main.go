package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/database"
	"github.com/suzmii/ACMBot/internal/database/dbmodel"
	"gorm.io/gorm"
)

func main() {
	db := database.Init()

	// 特殊迁移：修改races表的id字段类型并转换现有数据
	if err := migrateRacesIDField(db); err != nil {
		logrus.Fatal("Failed to migrate races ID field:", err)
	}

	// 常规自动迁移
	err := db.AutoMigrate(database.AllModels...)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("Database migration completed successfully")
}

func migrateRacesIDField(db *gorm.DB) error {
	// 检查races表是否存在
	if !db.Migrator().HasTable("races") {
		logrus.Info("races table does not exist, skipping ID field migration")
		return nil
	}

	// 检查当前id字段类型
	columnType, err := db.Migrator().ColumnTypes("races")
	if err != nil {
		return err
	}

	var currentIDType string
	for _, col := range columnType {
		if col.Name() == "id" {
			currentIDType = col.DatabaseTypeName()
			break
		}
	}

	// 如果id字段已经是varchar类型，跳过迁移
	if currentIDType == "VARCHAR" || currentIDType == "varchar" {
		logrus.Info("races.id field is already VARCHAR type, skipping migration")
		return nil
	}

	logrus.Info("Migrating races.id field from bigint to varchar and converting existing data...")

	// 1. 首先获取所有现有数据
	var existingRaces []dbmodel.Races
	if err := db.Find(&existingRaces).Error; err != nil {
		return fmt.Errorf("failed to fetch existing races: %w", err)
	}

	logrus.Infof("Found %d existing races to migrate", len(existingRaces))

	// 2. 创建临时表来存储转换后的数据
	tempTableSQL := `
	CREATE TABLE races_temp (
		id VARCHAR(64) NOT NULL PRIMARY KEY,
		created_at DATETIME(3),
		updated_at DATETIME(3),
		deleted_at DATETIME(3),
		resource BIGINT,
		title LONGTEXT,
		start_at DATETIME(3),
		end_at DATETIME(3),
		link LONGTEXT,
		INDEX idx_races_temp_deleted_at (deleted_at),
		INDEX idx_races_temp_start_at (start_at),
		INDEX idx_races_temp_end_at (end_at)
	)`

	if err := db.Exec(tempTableSQL).Error; err != nil {
		return fmt.Errorf("failed to create temp table: %w", err)
	}

	// 3. 转换数据并插入临时表
	for _, race := range existingRaces {
		// 生成新的ID
		newID := race.GenerateID()

		insertSQL := `
		INSERT INTO races_temp (id, created_at, updated_at, deleted_at, resource, title, start_at, end_at, link)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		if err := db.Exec(insertSQL, newID, race.CreatedAt, race.UpdatedAt, race.DeletedAt,
			race.Resource, race.Title, race.StartAt, race.EndAt, race.Link).Error; err != nil {
			// 清理临时表
			db.Exec("DROP TABLE races_temp")
			return fmt.Errorf("failed to insert race with new ID: %w", err)
		}
	}

	// 4. 删除原表并重命名临时表
	if err := db.Exec("DROP TABLE races").Error; err != nil {
		// 清理临时表
		db.Exec("DROP TABLE races_temp")
		return fmt.Errorf("failed to drop original races table: %w", err)
	}

	if err := db.Exec("RENAME TABLE races_temp TO races").Error; err != nil {
		return fmt.Errorf("failed to rename temp table: %w", err)
	}

	logrus.Info("races.id field migration completed successfully")
	return nil
}
