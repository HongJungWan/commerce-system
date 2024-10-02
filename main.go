package main

import (
	"flag"
	"fmt"
	"github.com/HongJungWan/commerce-system/docs"
	"github.com/HongJungWan/commerce-system/internal/helper"
	configs "github.com/HongJungWan/commerce-system/internal/infrastructure/configs"
	"github.com/HongJungWan/commerce-system/internal/infrastructure/router"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
	"os"
)

var (
	conf = configs.Config{}
	file string
)

// @title Commerce-System Service API
// @version 1.0
// @description A commerce-system service API in Go using Gin framework
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@commerce-system.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3031
// @BasePath /api
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if !parseConfig() {
		helper.ShowHelp()
		os.Exit(-1)
	}
	initializeSwaggerHost(&conf)
	db := configs.ConnectionDB(&conf)
	startServer(db)
}

func loadConfig() bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	viper.SetConfigFile(file)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(conf.DBHost)
	err = viper.GetViper().Unmarshal(&conf)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func parseConfig() bool {
	flag.StringVar(&file, "c", "config.toml", "config file")
	flag.Parse()
	if !loadConfig() {
		return false
	}

	return true
}

func startServer(db *gorm.DB) {
	routers := router.NewRouter(conf, db)

	server := &http.Server{
		Addr:    "0.0.0.0:3031",
		Handler: routers,
	}
	err := server.ListenAndServe()
	if err != nil {
		helper.ErrorPanic(err)
	}
}

func initializeSwaggerHost(conf *configs.Config) {
	docs.SwaggerInfo.Host = conf.Host
	docs.SwaggerInfo.Schemes = conf.Scheme
	docs.SwaggerInfo.Version = conf.Version
	docs.SwaggerInfo.BasePath = conf.BasePath
	docs.SwaggerInfo.Title = conf.Title

	fmt.Printf(
		"설정된 Swagger 정보:\nHost: %s\nSchemes: %v\nVersion: %s\nBasePath: %s\nTitle: %s\n",
		conf.Host,
		conf.Scheme,
		conf.Version,
		conf.BasePath,
		conf.Title,
	)
}
