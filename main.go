package main

import (
	db "main/database"
	"main/helpers"
	"main/models"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func setConfig() string {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("auth")
	var port string
	if err := viper.ReadInConfig(); err != nil {
		port = ":8083"
		SECRET = []byte("")
	} else {
		port = ":" + viper.GetString("port")
		SECRET = []byte(viper.GetString("secret"))
	}
	return port
}

func main() {
	port := setConfig()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		helpers.LogMsg("Can`t listen port ", err)
	}

	server := grpc.NewServer()

	models.RegisterAuthCheckerServer(server, NewAuthManager())

	err = db.Connect()
	if err != nil {
		helpers.LogMsg("Connection error: ", err)
		return
	}
	defer db.Disconnect()

	helpers.LogMsg("AuthServer started at", port)
	server.Serve(lis)
}
