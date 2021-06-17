package config

import (
	"fmt"
	"os"
	"tides-server/pkg/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	config *Config
	db     *gorm.DB
	err    error
)

const (
	URLSuffix = "cloudtides.vthink.cloud"
)

func init() {
	initConfig()
}

func initConfig() {
	godotenv.Load("../.env")
	config = &Config{}
	config.ServerIP = os.Getenv("SERVER_IP")
	config.ServerPort = os.Getenv("SERVER_PORT")
	config.PostgresHost = os.Getenv("POSTGRES_HOST")
	config.PostgresPort = os.Getenv("POSTGRES_PORT")
	config.PostgresUser = os.Getenv("POSTGRES_USER")
	config.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	config.PostgresDB = os.Getenv("POSTGRES_DB")
	config.SecretKey = os.Getenv("SECRET_KEY")
	config.AdminUser = os.Getenv("ADMIN_USER")
	config.AdminPassword = os.Getenv("ADMIN_PASSWORD")
	StartDB()
}

// GetConfig returns a pointer to the current config.
func GetConfig() *Config {
	return config
}

// GetDB returns a pointer to the database
func GetDB() *gorm.DB {
	return db
}

// StartDB initiates database with user configuration, also migrates db schema
func StartDB() {
	var dbinfo string

	dbinfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost, config.PostgresPort, config.PostgresUser, config.PostgresPassword, config.PostgresDB)
	db, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.VMachine{})
	db.AutoMigrate(&models.VMTemp{})
	db.AutoMigrate(&models.Policy{})
	db.AutoMigrate(&models.VcdPolicy{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.Vsphere{})
	db.AutoMigrate(&models.Vcd{})
	db.AutoMigrate(&models.VM{})
	db.AutoMigrate(&models.ResourceUsage{})
	db.AutoMigrate(&models.ResourcePastUsage{})
	db.AutoMigrate(&models.VMUsage{})
	db.AutoMigrate(&models.Vendor{})
	db.AutoMigrate(&models.Vapp{})
	db.AutoMigrate(&models.Port{})
	fmt.Println("DB connection success")
	CreateAdmin()
	TemplateSetup()
}

// CreateAdmin sets up an admin account
func CreateAdmin() {
	db := GetDB()
	var adm models.User
	if db.Where("username = ?", config.AdminUser).First(&adm).RowsAffected == 0 {
		admin := models.User{
			Username: config.AdminUser,
			Password: config.AdminPassword,
			Priority: models.UserPriorityHigh,
		}
		db.Create(&admin)
	}
}

// TemplateSetup sets up a VM template instance
func TemplateSetup() {
	db := GetDB()
	var tem models.Template
	if db.Where("name = ?", "tides-boinc-attached").First(&tem).RowsAffected == 0 {
		newTem := models.Template{
			GuestOS:          "Ubuntu-18.04",
			MemorySize:       8,
			Name:             "tides-boinc-attached",
			ProvisionedSpace: 16,
			VMName:           "tides-gromacs",
		}
		db.Create(&newTem)
	}
}
