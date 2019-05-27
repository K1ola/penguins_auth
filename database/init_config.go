package database

import (
	"auth/helpers"
	h "auth/helpers"

	"github.com/jackc/pgx"
	"github.com/spf13/viper"
)

// var connection *sq.DB = nil

// var connectionConfig pgx.ConnConfig
// var connectionPoolConfig = pgx.ConnPoolConfig{
// 	MaxConnections: 8,
// }

// var ImagesAddress string

// func initConfig() error {
// 	viper.AddConfigPath("./configs")
// 	viper.SetConfigName("config")
// 	if err := viper.ReadInConfig(); err != nil {
// 		helpers.LogMsg("Can't find db config: ", err)
// 		return err
// 	}
// 	connectionConfig = pgx.ConnConfig{
// 		Host:     viper.GetString("db.host"),
// 		Port:     uint16(viper.GetInt("db.port")),
// 		Database: viper.GetString("db.database"),
// 		User:     viper.GetString("db.user"),
// 		Password: viper.GetString("db.password"),
// 	}
// 	psqlURI := "postgresql://" + connectionConfig.User
// 	if len(connectionConfig.Password) > 0 {
// 		psqlURI += ":" + connectionConfig.Password
// 	}
// 	psqlURI += "@" + connectionConfig.Host + ":" + strconv.Itoa(int(connectionConfig.Port)) + "/" + connectionConfig.Database + "?sslmode=disable"
// 	fmt.Println(psqlURI)
// 	var err error
// 	connection, err = sq.Open("postgres", psqlURI)
// 	if err != nil {
// 		helpers.LogMsg("Can't connect to db: ", err)
// 		return err
// 	}
// 	viper.SetConfigName("fileserver")
// 	if err := viper.ReadInConfig(); err != nil {
// 		helpers.LogMsg("Can't find images address: ", err)
// 		return err
// 	}
// 	ImagesAddress = viper.GetString("address")
// 	return nil
// }

// func SetMock(databaseMock *sq.DB) {
// 	connection = databaseMock
// }

// func Connect() error {
// 	if connection != nil {
// 		return nil
// 	}
// 	err := initConfig()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func Disconnect() {
// 	if connection != nil {
// 		connection.Close()
// 		connection = nil
// 	}
// }

var connection *pgx.ConnPool = nil

var connectionConfig pgx.ConnConfig
var connectionPoolConfig = pgx.ConnPoolConfig{
	MaxConnections: 20,
}

var ImagesAddress string

//TODO check connect
func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		helpers.LogMsg("Can't find db config: ", err)
		return err
	}
	connectionConfig = pgx.ConnConfig{
		Host:     viper.GetString("db.host"),
		Port:     uint16(viper.GetInt("db.port")),
		Database: viper.GetString("db.database"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
	}

	connectionPoolConfig.ConnConfig = connectionConfig
	viper.SetConfigName("fileserver")
	if err := viper.ReadInConfig(); err != nil {
		helpers.LogMsg("Can't find images address: ", err)
		return err
	}
	ImagesAddress = viper.GetString("address")
	return nil
}

func Connect() error {
	if connection != nil {
		return nil
	}
	err := initConfig()
	if err != nil {
		return err
	}
	connection, err = pgx.NewConnPool(connectionPoolConfig)
	if err != nil {
		h.LogMsg("Connect DB error: " + err.Error())
		return err
	}
	return nil
}

func Disconnect() {
	if connection != nil {
		connection.Close()
		connection = nil
	}
}
