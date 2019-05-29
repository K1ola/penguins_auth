package main

import (
	db "auth/database"
	"auth/helpers"
	"auth/models"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func setConfig() string {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	var port string
	if err := viper.ReadInConfig(); err == nil {
		port = ":" + viper.GetString("port")
		SECRET = []byte(viper.GetString("secret"))
	}
	return port
}

func main() {
	port := setConfig()
	lis, _ := net.Listen("tcp", port)
	// if err != nil {
	// 	helpers.LogMsg("Can`t listen port ", err)
	// }

	server := grpc.NewServer()

	models.RegisterAuthCheckerServer(server, NewAuthManager())

	_ = db.Connect()
	// if err != nil {
	// 	helpers.LogMsg("Connection error: ", err)
	// 	return
	// }
	defer db.Disconnect()

	helpers.LogMsg("AuthServer started at", port)
	server.Serve(lis)
}
