package main

import (
	"github.com/dgrijalva/jwt-go"
	db "main/database"
	"main/helpers"
	"main/models"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	//"io/ioutil"
	//"net/http"
	"sync"
	"time"
)

var SECRET = []byte("myawesomesecret")

type AuthManager struct {
	mu       sync.RWMutex
	token 	 *models.JWT
	user     *models.User
}

func NewAuthManager() *AuthManager {
	return &AuthManager{
		mu:       sync.RWMutex{},
		token:    new(models.JWT),
		user:     new(models.User),
	}
}


//TODO check error returns

func (am *AuthManager) LoginUser(ctx context.Context, user *models.User) (*models.JWT, error) {
	found, _ := db.GetUserByEmail(user.Email)

	if found == nil {
		//w.WriteHeader(http.StatusNotFound)
		return nil, status.Errorf(codes.NotFound, "No such user")
	}

	if !helpers.CheckPasswordHash(user.Password, found.HashPassword) {
		//w.WriteHeader(http.StatusForbidden)
		return nil, status.Errorf(codes.PermissionDenied, "Incorrect password")
	}
	ttl := time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    found.ID,
		"userEmail": user.Email,
		"userLogin": user.Login,
		"exp":       time.Now().UTC().Add(ttl).Unix(),
	})

	str, err := token.SignedString(SECRET)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Incorrect secret")
	}
	am.token.Token = str
	return am.token, nil
}

func (am *AuthManager) RegisterUser(ctx context.Context, user *models.User) (*models.JWT, error) {
	foundByEmail, _ := db.GetUserByEmail(user.Email)
	foundByLogin, _ := db.GetUserByLogin(user.Login)

	if foundByEmail != nil || foundByLogin != nil{
		return nil, status.Errorf(codes.AlreadyExists, "Such user already exists")
	}

	user.HashPassword = helpers.HashPassword(user.Password)

	err := db.CreateUser(user)
	if err != nil {
		helpers.LogMsg(err)
		return nil, status.Errorf(codes.Internal, "Server error")
	}

	ttl := time.Hour

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    user.ID,
		"userEmail": user.Email,
		"userLogin": user.Login,
		"exp":       time.Now().UTC().Add(ttl).Unix(),
	})

	str, err := token.SignedString(SECRET)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "Incorrect secret")
	}

	am.token.Token = str
	return am.token, nil
}

func (am *AuthManager) GetUser(ctx context.Context, token *models.JWT) (*models.User, error) {
	t, _ := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})

	claims, _ := t.Claims.(jwt.MapClaims)

	temp := claims["userID"]
	mytemp := uint(temp.(float64))

	user, _ := db.GetUserByID(mytemp)
	if user == nil {
		return nil, status.Errorf(codes.Unknown, "Unauthorized")
	}
	return user, nil
}

func (am *AuthManager) ChangeUser(ctx context.Context, user *models.User) (*models.Nothing, error) {
	_, err := db.UpdateUserByID(user, uint(user.ID))
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "Such user already exists")
	}
	return &models.Nothing{}, nil
}

//TODO DeleteUser() is needed?
func (am *AuthManager) DeleteUser(ctx context.Context, token *models.JWT) (*models.Nothing, error) {
	return &models.Nothing{}, nil
}

