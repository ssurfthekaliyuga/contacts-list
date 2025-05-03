package auth

import (
	"contacts-list/internal/domain/errs"
	"contacts-list/pkg/sl"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"strings"
)

type key struct{}

type auth struct {
	secret       []byte
	jwtUserIDKey string
	logUserIDKey string
}

func New(secret []byte) fiber.Handler {
	mw := auth{
		secret: secret,
	}

	return mw.middleware
}

func Extract(c *fiber.Ctx) uuid.UUID {
	value := c.UserContext().Value(key{})
	if value == nil {
		return uuid.Nil
	}

	id, ok := value.(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return id
}

func (a *auth) middleware(c *fiber.Ctx) error {
	tokenStr, err := a.token(c)
	if err != nil {
		return err
	}

	var claims jwt.MapClaims
	token, err := jwt.ParseWithClaims(tokenStr, claims, a.parse) //todo может быть лучше созлать парсер а потом его использовать
	if err != nil || !token.Valid {
		return errs.NewUnauthorized("invalid token") //todo
	}

	sub, ok := claims[a.jwtUserIDKey].(string)
	if !ok {
		msg := fmt.Sprintf("invalid token missing [%s] field", a.jwtUserIDKey)
		return errs.NewUnauthorized(msg)
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return errs.NewUnauthorized("invalid uuid in token in [%s] field")
	}

	ctx := c.UserContext()
	ctx = context.WithValue(ctx, key{}, id)
	ctx = sl.ContextWithAttrs(ctx, slog.String(a.logUserIDKey, id.String()))

	c.SetUserContext(ctx)

	return c.Next()
}

func (a *auth) token(c *fiber.Ctx) (string, error) {
	header := c.Get("Authorization")
	if header == "" {
		return "", nil
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errs.NewUnauthorized("authorization header format must be Bearer")
	}

	return parts[1], nil
}

func (a *auth) parse(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		msg := fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
		return nil, errs.NewUnauthorized(msg)
	}

	return a.secret, nil
}
