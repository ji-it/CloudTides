package config

import (
	"os"
	// "database/sql"
	"fmt"
	"tides-server/pkg/models"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	config *Config
	db     *gorm.DB
	err    error
)

func init() {
	initConfig()
}

func initConfig() {
	config = &Config{}
	config.Port = defaultPort
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort != "" {
		config.Port = serverPort
	}
	startDB()
}

// GetConfig returns a pointer to the current config.
func GetConfig() *Config {
	return config
}

func GetDB() *gorm.DB {
	return db
}

func startDB() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.Policy{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.VM{})
	db.AutoMigrate(&models.ResourceUsage{})
	db.AutoMigrate(&models.VMUsage{})
	fmt.Println("DB connection success")

}
