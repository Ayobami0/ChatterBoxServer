package service

import (
	"time"

	"github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/repository/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserAccountClaim struct {
	jwt.RegisteredClaims
	Email      string  `json:"email"`
	UserId     string  `json:"userId"`
}

func NewAuthService(r user.UserRepository, secret string) AuthService {
	return AuthService{r, []byte(secret)}
}

type AuthService struct {
	r      user.UserRepository
	secret []byte
}

func (a *AuthService) GenenerateJWTToken(id, username, email string, profileUrl *string) (string, error) {
	newClaim := &UserAccountClaim{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
		email,
    id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)

	if t, err := token.SignedString(a.secret); err != nil {
		return "", err
	} else {
		return t, nil
	}
}

func (a *AuthService) JWTConfig() echo.MiddlewareFunc {
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(UserAccountClaim)
		},
		SigningKey: a.secret,
	}

	return echojwt.WithConfig(config)
}

func (a *AuthService) VerifyUser(password, usernameOrEmail string) (*model.User, error) {
	user := a.r.GetUser(usernameOrEmail)
	if user == nil {
		return nil, errors.ErrUserNotExist(usernameOrEmail)
	}

	if !compareHash(password, user.Password) {
		return nil, errors.ErrPasswordIncorrect{}
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func (a *AuthService) RegisterUser(u model.UserCreate) (*model.User, error) {
	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	pwHash, err := hashPassword(u.Password)

	if err != nil {
		return nil, errors.ErrUnexpected("an unexpected error occured")
	}

	nUSer := model.User{
		UserBase: u.UserBase,
		ID:       id.String(),
		Password: pwHash,
	}
	err = a.r.CreateUser(nUSer)

	if err != nil {
		return nil, errors.ErrUserExist{
			Username: u.Username,
			Email:    u.Email,
		}
	}

	return &nUSer, nil
}
