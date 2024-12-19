package main

import (
	"event-manager/config"
	"event-manager/delivery/httpserver"
	"event-manager/delivery/httpserver/eventhandler"
	"event-manager/delivery/httpserver/userhandler"
	"event-manager/repository/mysql"
	"event-manager/repository/mysql/mysqlevent"
	"event-manager/repository/mysql/mysqluser"
	"event-manager/service/authservice"
	"event-manager/service/eventservice"
	"event-manager/service/userservice"
)

func main() {
	cnf := getConfig()
	
	authSvc, userSvc, eventSvc := services(cnf)
	uh, eh := handlers(userSvc, authSvc, eventSvc)
	hserver := httpserver.New(cnf, uh, eh)

	hserver.Serve()
}

func handlers(us userservice.UserService, as authservice.AuthService, es eventservice.EventService) (userhandler.UserHandler, eventhandler.EventHandler) {
	uh := userhandler.New(us, as)
	eh := eventhandler.New(es)
	
	return uh, eh
}

func services(cnf config.Config) (authservice.AuthService, userservice.UserService, eventservice.EventService) {
	mysql := mysql.New()
	userRepo := mysqluser.New(mysql)
	eventRepo := mysqlevent.New(mysql)

	authSvc := authservice.New(cnf.AuthConfig)
	userSvc := userservice.New(userRepo, authSvc)
	eventSvc := eventservice.New(eventRepo)
	
	return authSvc, userSvc, eventSvc
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
