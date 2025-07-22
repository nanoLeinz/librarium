package helper

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoLeinz/librarium/model/dto"
)

type JWTClaims struct {
	MemberID string
	Role     string
	Email    string
	jwt.RegisteredClaims
}

func GenerateJWTToken(member *dto.MemberResponse) (string, error) {

	expiry, _ := strconv.Atoi(os.Getenv("EXPIRYINMINUTE"))
	secretKey := os.Getenv("SECRETJWT")

	claims := JWTClaims{
		MemberID: member.ID.String(),
		Email:    member.Email,
		Role:     member.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "librarium",
		},
	}

	log.Printf("Claims Created : %+v\n", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", fmt.Errorf("error signing token : %w", err)
	}

	return signedToken, nil

}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("SECRETJWT")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, fmt.Errorf("error parsing token : %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("error parsing claims : %w", err)
	}

	return claims, nil
}
