package player

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUnauthorized = errors.New("not authorized")
	ErrInternal     = errors.New("authorization internal error")
)

type PlayerClaim struct {
	PlayerId string
	GameId   string
}

const secret = "3489453kjhkdayf98di54jk34hksadjfnjkas4h378yfkj3"

func GenerateToken(claim PlayerClaim) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"playerId": claim.PlayerId,
		"gameId":   claim.GameId,
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		slog.Error("Error Generating Token")
	}
	return tokenString
}

func DecodeToken(stringToken string) (PlayerClaim, error) {
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	claim := PlayerClaim{}
	switch {
	case token != nil && token.Valid:
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if playerId, ok := claims["playerId"].(string); ok {
				claim.PlayerId = playerId
			}
			if gameId, ok := claims["gameId"].(string); ok {
				claim.GameId = gameId
			}
			return claim, nil
		}
	case errors.Is(err, jwt.ErrTokenMalformed):
		return claim, ErrUnauthorized
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return claim, ErrUnauthorized
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return claim, ErrUnauthorized
	default:
		slog.Error(err.Error())
	}
	return claim, ErrInternal
}

func AuthorizeRequest(r *http.Request) (PlayerClaim, error) {
	Token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	slog.Debug(Token)
	return DecodeToken(Token)
}
