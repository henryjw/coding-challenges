package auth

import (
	"database/sql"
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var jwtSecret = []byte("this is a learning project. it's ok to hardcode this here")

var InvalidLoginError = errors.New("invalid username or password")
var JTWCreateError = errors.New("error generating jwt")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct for JWT token
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthService struct {
	db *sql.DB
}

func New(db *sql.DB) *AuthService {
	return &AuthService{db}
}

// Login Returns JWT on successful login
func (receiver AuthService) Login(user User) (string, error) {
	var passwordHash string
	err := receiver.db.QueryRow("SELECT password FROM Users WHERE username = ?", user.Username).Scan(&passwordHash)

	if err != nil {
		return "", InvalidLoginError
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(user.Password))

	if err != nil {
		return "", InvalidLoginError
	}

	tokenExpiration := time.Now().Add(5 * time.Minute)
	token, err := generateJwt(user.Username, tokenExpiration)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (receiver AuthService) SignUp(user User) error {
	return errors.New("not yet implemented")
}

func generateJwt(username string, expires time.Time) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)

	if err != nil {
		log.Printf("Unexpected error generating JTW: %v\n", err)
		return "", JTWCreateError
	}

	return token, nil
}
