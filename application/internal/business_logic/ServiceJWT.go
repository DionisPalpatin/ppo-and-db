package bl

import (
	"fmt"
	"github.com/DionisPalpatin/ppo-and-db/tree/master/application/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var secretKey = []byte("your_secret_key")

func generateToken(userID int, userRole string) (*models.Token, error) {
	expirationTime := time.Now().Add(730 * time.Hour) // Месяц

	claims := &models.CustomClaims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "ProgrammerNotebook",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the token: %v", err)
	}

	return &models.Token{
		AccessToken: signedToken,
		TokenType:   "Bearer",
		ExpiresIn:   int64(expirationTime.Sub(time.Now()).Seconds()),
	}, nil
}

func checkToken(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return secretKey, nil
}

func ValidateAndParseToken(signedToken string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &models.CustomClaims{}, checkToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the token: %v", err)
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
