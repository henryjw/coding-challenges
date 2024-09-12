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
var UserAlreadyExistsError = errors.New("user already exists")

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
		// Should use a more reliable method to check for error (e.g., error code), but it's fine
		// for this toy project
		if err.Error() == "sql: no rows in result set" {
			return "", InvalidLoginError
		} else {
			return "", err
		}
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = receiver.db.Exec("INSERT INTO Users (username, password) VALUES (?, ?)", user.Username, hashedPassword)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return UserAlreadyExistsError
		}

		return err
	}

	return nil
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
