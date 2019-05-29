package helpers

import (
	// "auth/models"
	"testing"
	// "net/http/httptest"
	// "net/http"
	// "time"
)

func TestPasswords(t *testing.T) {
	password := "password"
	hash := HashPassword(password)
	result := CheckPasswordHash(password, hash)
	if !result {
		t.Error("hashing is not equal")
	}

}

func TestLogs(t *testing.T) {
	LogMsg("hello world")
}

// func TestDelCookie(t *testing.T)  {
// 	var w *http.ResponseWriter
// 	cookie := &http.Cookie{
// 		Name:     "sessionid",
// 		Value:    "hjebvjhfdbvikjdf",
// 		Expires:  time.Now(),
// 		HttpOnly: true,
// 	}
// 	DeleteCookie(w, cookie)
// }
