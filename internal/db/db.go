package db

import (
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ScanHistory struct {
	ID             uint      `gorm:"primaryKey"`
	CreatedAt      time.Time
	TotalContainers int
	StoppedContainers int
	TotalImages     int
	DanglingImages  int
	TotalVolumes    int
	OrphanedVolumes int
	TotalNetworks         int
	UnusedNetworks        int
	HealthScore           int
	RecoverableSpaceBytes int64
}

var DB *gorm.DB

func InitDB() error {
	var err error
	
	// Usamos sqlite pure-go para evitar dependencias CGO y asegurar compatibilidad cross-platform
	DB, err = gorm.Open(sqlite.Open("docker_doctor.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("error al abrir base de datos: %w", err)
	}

	err = DB.AutoMigrate(&ScanHistory{})
	if err != nil {
		return fmt.Errorf("error migrando esquema de base de datos: %w", err)
	}

	return nil
}

func SaveScan(scan ScanHistory) error {
	if DB == nil {
		return fmt.Errorf("base de datos no inicializada")
	}
	return DB.Create(&scan).Error
}

func GetLatestScans(limit int) ([]ScanHistory, error) {
	var scans []ScanHistory
	if DB == nil {
		return scans, fmt.Errorf("base de datos no inicializada")
	}
	err := DB.Order("created_at desc").Limit(limit).Find(&scans).Error
	return scans, err
}

func GetLastScan() (ScanHistory, error) {
	var scan ScanHistory
	if DB == nil {
		return scan, fmt.Errorf("base de datos no inicializada")
	}
	err := DB.Order("created_at desc").First(&scan).Error
	return scan, err
}
