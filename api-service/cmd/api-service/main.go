package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"api-service/api"
	"api-service/config"
	"api-service/router"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/database/mariadb"
	"github/littlepaulhi/highly-concurrent-e-commerce-lightweight-system/pkg/logger"
)

var (
	eg            errgroup.Group
	configuration config.Configuration
)

func init() {
	mariadb.Setup()
	viper.AutomaticEnv()
	viper.SetConfigName("config-server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PROJECT_PATH/api-service/config/")

	if err := viper.ReadInConfig(); err != nil {
		logger.InitLog.Errorf("Error when reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.InitLog.Errorf("Unable to decode into struct, %v", err)
	}
}

func main() {
	gin.SetMode(configuration.Server.RunMode)

	ginRouter := router.Initialize()
	readTimeout := configuration.Server.ReadTimeout
	writeTimeout := configuration.Server.WriteTimeout
	endPoint := fmt.Sprintf("%s:%d", configuration.Server.Addr, configuration.Server.Port)

	api.InitAllCache(&configuration.Cache)

	apiServer := &http.Server{
		Addr:         endPoint,
		Handler:      ginRouter,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	logger.APILog.Printf("[Info] Start http server, listening on port %s", endPoint)

	eg.Go(func() error {
		return apiServer.ListenAndServe()
	})

	if err := eg.Wait(); err != nil {
		logger.InitLog.Warnln("EG Wait error: ", err)
	}
}
