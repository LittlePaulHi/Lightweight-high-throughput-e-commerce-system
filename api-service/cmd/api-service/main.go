package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	config "api-service/configs"
	"api-service/router"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
)

var (
	eg errgroup.Group
)

func init() {
	mariadb.Setup()
	viper.AutomaticEnv()
	viper.SetConfigName("config-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/api-service/configs/")
}

func main() {
	var configuration config.Configuration

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error when reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	gin.SetMode(configuration.Server.RunMode)

	ginRouter := router.Initialize()
	readTimeout := configuration.Server.ReadTimeout
	writeTimeout := configuration.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", configuration.Server.Port)

	apiServer := &http.Server{
		Addr:         endPoint,
		Handler:      ginRouter,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	log.Printf("[Info] Start http server, listening on port %s", endPoint)

	eg.Go(func() error {
		return apiServer.ListenAndServe()
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
