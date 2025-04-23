package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	todo "rest_API"
	"rest_API/pkg/handler"
	"rest_API/pkg/repository"
	"rest_API/pkg/service"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	err := InitConfig()
	if err != nil {
		logrus.Fatalf("Eror, when begining db - %s", err)
	}

	err = godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error can not load env %s", err.Error())
	}

	db, err := repository.NewPosgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("Error in install contact with DataBase: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	err = srv.Start("7777", handlers.InitRoutes())
	if err != nil {
		logrus.Fatalf("Error, when server were started : %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
