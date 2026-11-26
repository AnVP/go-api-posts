package jwt_test

import (
	"go-api-posts/pkg/jwt"
	"testing"
)

func TestJWT_GenerateToken(t *testing.T) {
	const email = "a@a.ru"
	jwtService := jwt.NewJWT("/2+XnmJGz1j3ehIVI/5P9kll_CghrE3DcS7rnT+gar5w=")
	token, err := jwtService.GenerateToken(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}
	if data.Email != email {
		t.Fatal("Email is not valid")
	}
}
