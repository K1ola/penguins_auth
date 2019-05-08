package main

import (
	"auth/database"
	"auth/models"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/net/context"
)

func TestAuthManager(t *testing.T) {
	manager := NewAuthManager()
	var user models.UserProto
	var ctx context.Context
	user.Login = time.Now().Format("20060102150405") + user.Login
	user.Email = time.Now().Format("20060102150405") + user.Email
	user.Password = "testtest"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.SetMock(db)
	rows := sqlmock.NewRows([]string{"id", "login", "email", "hashpassword", "score", "name", "games"}).
		AddRow(1, "login", "login@mail.ru", "9hevHKd0671tI-j-EsVtQdJaItgeiPZ8bG7g1A==", 10, "name", "5").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	jwt, err := manager.LoginUser(ctx, &user)
	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("some error"))
	_, _ = manager.LoginUser(ctx, &user)
	rows = sqlmock.NewRows([]string{"id", "login", "email", "hashpassword", "score", "name", "games"}).
		AddRow(1, "login", "login@mail.ru", "9hevHK0671tI-j-EsVtQdJaItgeiPZ8bG7g1A==", 10, "name", "5").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	_, _ = manager.LoginUser(ctx, &user)
	fmt.Println("jwt", jwt)
	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("some error"))
	manager.GetUser(ctx, jwt)
	rows = sqlmock.NewRows([]string{"id", "login", "email", "hashpassword", "score", "name", "games"}).
		AddRow(1, "login", "login@mail.ru", "hdfbkfbdj", 10, "name", "5").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	manager.GetUser(ctx, jwt)
	manager.DeleteUser(ctx, jwt)

}

func TestRegister(t *testing.T) {
	manager := NewAuthManager()
	var user models.UserProto
	var ctx context.Context
	user.Login = time.Now().Format("20060102150405") + user.Login
	user.Email = time.Now().Format("20060102150405") + user.Email
	user.Password = "password"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	database.SetMock(db)
	rows := sqlmock.NewRows([]string{"id", "login", "email", "hashpassword", "score", "name", "games"}).
		AddRow(1, "login", "login@mail.ru", "$2a$08$Q1nN3cy96NhOW7jOx31atuzY.QuRXbnWRitfkwZDHbC3dY83bw53i", 10, "name", "5").
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	manager.RegisterUser(ctx, &user)
	rows = sqlmock.NewRows([]string{"id", "login", "score"}).
		AddRow(1, "login", 10).
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("some error"))
	mock.ExpectQuery("INSERT").WillReturnRows(rows)
	manager.RegisterUser(ctx, &user)
	mock.ExpectQuery("SELECT").WillReturnError(fmt.Errorf("some error"))
	mock.ExpectQuery("INSERT").WillReturnError(fmt.Errorf("some error"))
	manager.RegisterUser(ctx, &user)
	manager.ChangeUser(ctx, &user)
	rows = sqlmock.NewRows([]string{"games", "name", "score"}).
		AddRow(10, "name", 5).
		RowError(1, fmt.Errorf("error"))
	mock.ExpectQuery("UPDATE").WillReturnRows(rows)
	manager.ChangeUser(ctx, &user)
}

func TestSetConfig(t *testing.T) {
	setConfig()
	go func() {
		main()
	}()
}
