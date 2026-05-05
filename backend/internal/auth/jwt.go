package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService gestiona la creación y validación de tokens
type JWTService struct {
	secretKey []byte
	issuer    string
}

// CustomClaims define la estructura de los datos contenidos en el token
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTService crea una nueva instancia del servicio
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{
		secretKey: []byte(secretKey),
		issuer:    "seismic-monitor-api",
	}
}

// GenerateToken crea un nuevo token para un usuario específico
func (s *JWTService) GenerateToken(userID string) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 horas de validez
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken comprueba si un token es válido y devuelve sus claims
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token inválido")
}
