package main

import (
	"event-manager/config"
	"event-manager/delivery/httpserver"
	"event-manager/repository/mysql"
	"event-manager/repository/mysql/mysqlevent"
	"event-manager/repository/mysql/mysqluser"
	"event-manager/service/authservice"
	"event-manager/service/eventservice"
	"event-manager/service/userservice"
)

func main() {
	cnf := getConfig()
	mysql := mysql.New()
	userRepo := mysqluser.New(mysql)
	eventRepo := mysqlevent.New(mysql)

	authSvc := authservice.New(cnf.AuthConfig)
	userSvc := userservice.New(userRepo, authSvc)
	eventSvc := eventservice.New(eventRepo)

	hserver := httpserver.New(cnf, userSvc, authSvc, eventSvc)

	hserver.Serve()
}

func getConfig() config.Config {
	return config.Config{
		HttpServer: config.HttpServer{Port: 8080},
		AuthConfig: authservice.Config{
			SignKey:               config.JwtSignKey,
			AccessSubject:         config.AccessTokenSubject,
			RefreshSubject:        config.RefreshTokenSubject,
			AccessExpirationTime:  config.AccessTokenExpireDuration,
			RefreshExpirationTime: config.RefreshTokenExpireDuration,
		},
	}
}
