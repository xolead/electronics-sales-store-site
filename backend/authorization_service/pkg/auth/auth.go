package auth

import (
	"authoriz-service/pkg/dbwork"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

func init() {
	log.Logger = log.With().Str("package", "auth").Logger()
}

type claims struct {
	GUID  string `json:"GUID"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func CreateAccessToken(db *dbwork.DataBase, GUID string, admin bool) (string, error) {
	expires, err := strconv.Atoi(os.Getenv("expires_jwt"))
	if err != nil {
		log.Error().Msgf("Ошибка expires_access: %v", err)
		expires = 1
	}

	claims := &claims{
		GUID:  GUID,
		Admin: admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expires) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = db.EnableSession(ctx, GUID)
	if err != nil {
		return "", fmt.Errorf("auth/CreateAccessToken enableSession: %v", err)
	}
	return token.SignedString([]byte(os.Getenv("jwt_secret")))
}

func CheckAccessToken(access string) (string, bool, error) {
	claim := &claims{}

	token, err := jwt.ParseWithClaims(
		access,
		claim,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwt_secret")), nil
		},
	)
	if err != nil || !token.Valid {
		return uuid.Nil.String(), false, fmt.Errorf("auth/CheckAccessToken parse token: %v", err)
	}

	return claim.GUID, claim.Admin, nil
}

func CreateRefreshToken(ctx context.Context, db *dbwork.DataBase, GUID string) (string, error) {
	for range 3 {
		var refresh [32]byte

		_, err := rand.Read(refresh[:])
		if err != nil {
			return "", fmt.Errorf("auth/CreateRefreshToken rand.Read: %v", err)
		}

		strRefresh := base64.StdEncoding.EncodeToString(refresh[:])
		hashRefresh, err := bcrypt.GenerateFromPassword([]byte(strRefresh), bcrypt.DefaultCost)
		if err != nil {
			return "", fmt.Errorf("auth/CreateRefreshToken hasgRefresh: %v", err)
		}
		if err = db.CheckCollisionRefresh(ctx, string(hashRefresh)); err != nil {
			if err == dbwork.DuplicateRefresh {
				continue
			}
			return "", fmt.Errorf("auth/CreateRefreshToken CheckCollision: %v", err)
		}

		err = db.CreateRefresh(ctx, string(hashRefresh), GUID)
		if err != nil {
			return "", fmt.Errorf("auth/CreateRefreshToken CreateRefresh: %v", err)
		}

		return strRefresh, nil
	}
	log.Error().Msg("Не удалось создать refresh токен")
	return "", fmt.Errorf("Не удалось создать refresh token")
}
