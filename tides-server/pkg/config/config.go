package config

import (
	"os"
	// "database/sql"
	"fmt"
	"tides-server/pkg/models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
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
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err = gorm.Open("postgres", dbinfo)
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.Policy{}).AddForeignKey("user_ref", "users(id)", "CASCADE", "CASCADE").AddForeignKey("project_ref", "projects(id)", "SET NULL", "CASCADE").AddForeignKey("template_ref", "templates(id)", "SET NULL", "CASCADE")
	db.AutoMigrate(&models.Resource{}).AddForeignKey("user_ref", "users(id)", "CASCADE", "CASCADE").AddForeignKey("policy_ref", "policies(id)", "SET NULL", "CASCADE")
	db.AutoMigrate(&models.VM{}).AddForeignKey("resource_ref", "resources(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&models.ResourceUsage{}).AddForeignKey("resource_ref", "resources(id)", "CASCADE", "CASCADE")
	db.AutoMigrate(&models.VMUsage{}).AddForeignKey("vm_ref", "vms(id)", "CASCADE", "CASCADE")
	fmt.Println("DB connection success")

}
