package main

import (
	db "main/database"
	"main/helpers"
	"main/models"

	//"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8083")
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

	helpers.LogMsg("AuthServer started at 8083")
	server.Serve(lis)
}