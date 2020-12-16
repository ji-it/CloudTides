package config

import (
	"fmt"
	"os"
	"tides-server/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	config *Config
	db     *gorm.DB
	err    error
)

func InitConfig() {
	config = &Config{}
	config.Port = defaultPort
	serverPort := os.Getenv("envkey_SERVER_PORT")
	if serverPort != "" {
		config.Port = serverPort
	}
	StartDB()
}

// GetConfig returns a pointer to the current config.
func GetConfig() *Config {
	return config
}

func GetDB() *gorm.DB {
	return db
}

func StartDB() {
	var dbinfo string

	dbHost := os.Getenv("envkey_DB_HOST")
	dbPort := os.Getenv("envkey_DB_PORT")
	dbUser := os.Getenv("envkey_DB_USER")
	dbPassword := os.Getenv("envkey_DB_PASSWORD")
	dbName := os.Getenv("envkey_DB_NAME")
	dbinfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err = gorm.Open(postgres.Open(dbinfo), &gorm.Config{})
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Project{})
	db.AutoMigrate(&models.Template{})
	db.AutoMigrate(&models.Policy{})
	db.AutoMigrate(&models.VcdPolicy{})
	db.AutoMigrate(&models.Resource{})
	db.AutoMigrate(&models.Vsphere{})
	db.AutoMigrate(&models.Vcd{})
	db.AutoMigrate(&models.VM{})
	db.AutoMigrate(&models.ResourceUsage{})
	db.AutoMigrate(&models.ResourcePastUsage{})
	db.AutoMigrate(&models.VMUsage{})
	fmt.Println("DB connection success")
	CreateAdmin()
	TemplateSetup()
}

func CreateAdmin() {
	db := GetDB()
	var adm models.User
	adminUser := os.Getenv("envkey_ADMIN_USER")
	adminPwd := os.Getenv("envkey_ADMIN_PASSWORD")
	if db.Where("username = ?", adminUser).First(&adm).RowsAffected == 0 {
		admin := models.User{
			Username: adminUser,
			Password: adminPwd,
			Priority: models.UserPriorityHigh,
		}
		db.Create(&admin)
	}
}

func TemplateSetup() {
	db := GetDB()
	var tem models.Template
	if db.Where("name = ?", "tides-boinc-attached").First(&tem).RowsAffected == 0 {
		newTem := models.Template{
			GuestOS:          "Ubuntu-18.04",
			MemorySize:       8,
			Name:             "tides-boinc-attached",
			ProvisionedSpace: 16,
			VmName:           "tides-gromacs",
		}
		db.Create(&newTem)
	}
}
