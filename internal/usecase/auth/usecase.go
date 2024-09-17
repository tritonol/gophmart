package auth

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

const (
	salt      = "abcd"
	secretJwt = "secret"
)

type Claims struct {
	jwt.RegisteredClaims
	Id user.UserID
}

type authUsecase struct {
	repo AuthRepo
}

type AuthRepo interface {
	Create(ctx context.Context, credentials user.UserCredentials) (user.UserID, error)
	CheckByCredentials(ctx context.Context, credentials user.UserCredentials) (user.UserID, error)
}

func New(repo AuthRepo) *authUsecase {
	return &authUsecase{
		repo: repo,
	}
}

func (uc *authUsecase) Register(ctx context.Context, credetials user.UserCredentials) (string, error) {
	hashedPass, err := sha1Hash(credetials.Password)
	if err != nil {
		return "", err
	}
	credetials.Password = hashedPass

	userId, err := uc.repo.Create(ctx, credetials)
	if err != nil {
		return "", err
	}

	return buildJwt(userId)
}

func (uc *authUsecase) Login(ctx context.Context, credentials user.UserCredentials) (string, error) {
	hashedPass, err := sha1Hash(credentials.Password)
	if err != nil {
		return "", err
	}
	credentials.Password = hashedPass

	userId, err := uc.repo.CheckByCredentials(ctx, credentials)
	if err != nil {
		return "", nil
	}

	return buildJwt(userId)
}

func (uc *authUsecase) ValidateToken(token string) (user.UserID, error) {
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretJwt), nil
		})

	if err != nil {
		return 0, err
	}
	if !parsedToken.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	return claims.Id, nil
}

func sha1Hash(pass string) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(pass)); err != nil {
		return "", err
	}

	hashedBytes := hash.Sum([]byte(salt))
	hashedString := hex.EncodeToString(hashedBytes)

	return hashedString, nil
}

func buildJwt(id user.UserID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		Claims{
			Id: id,
		},
	)

	tokenStr, err := token.SignedString([]byte(secretJwt))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
