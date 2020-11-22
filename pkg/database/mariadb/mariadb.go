package mariadb

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"

	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/configs"
)

var db *gorm.DB

// Setup the database instance
func Setup() {
	viper.AutomaticEnv()
	viper.SetConfigName("config-database")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(viper.GetString("PROJECT_PATH") + "/pkg/configs/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error when reading config file, %s", err)
	}

	var configuration config.Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	db, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			configuration.Mariadb.User,
			viper.GetString("DATABASE_PASSWORD"),
			configuration.Mariadb.Host,
			configuration.Mariadb.Name),
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("Mariadb Setup error: %v", err)
	}

	if err = db.AutoMigrate(&Account{}, &Product{}, &Cart{}, &Order{}, &OrderItem{}); err != nil {
		log.Fatal(err.Error())
	}
}

// CloseMariadb operation will be defer
func CloseMariadb() {
	mysqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Cannot get gorm DB when setting defer closeDB")
	}

	defer mysqlDB.Close()
}
