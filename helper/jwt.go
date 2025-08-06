package helper

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nanoLeinz/librarium/model/dto"
	log "github.com/sirupsen/logrus"
)

type JWTClaims struct {
	MemberID string
	Role     string
	Email    string
	jwt.RegisteredClaims
}

func GenerateJWTToken(member *dto.MemberResponse) (string, error) {

	expiryStr := os.Getenv("EXPIRYINMINUTE")
	expiry, err := strconv.Atoi(expiryStr)
	if err != nil {

		log.WithError(err).Error("Invalid EXPIRYINMINUTE env var, falling back to 60 minutes")
		expiry = 60
	}
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

	log.WithFields(log.Fields{
		"memberID":  claims.MemberID,
		"role":      claims.Role,
		"expiresAt": claims.ExpiresAt,
	}).Info("Generating new JWT")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.WithError(err).Error("Error when signing token")
		return "", fmt.Errorf("error signing token: %w", err)
	}

	log.WithField("memberID", claims.MemberID).Info("Token created successfully")
	return signedToken, nil
}

func ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("SECRETJWT")

	log.Info("Validating incoming JWT")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.WithError(err).Warn("Error parsing or validating token")
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Warn("Token is invalid or claims could not be parsed")
		return nil, fmt.Errorf("invalid token")
	}

	log.WithFields(log.Fields{
		"memberID": claims.MemberID,
		"issuer":   claims.Issuer,
	}).Info("Token validated successfully")

	return claims, nil
}
