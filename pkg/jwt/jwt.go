package jwt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ribice/chisk/model"

	jwt "github.com/dgrijalva/jwt-go"
)

// Key to use when setting the request ID.
type ctxJWTKey int

const (
	// TokenKey is the key that holds jwt token
	TokenKey ctxJWTKey = iota
	// UserIDKey is the key that holds user's ID
	UserIDKey
	// UserDisplayNameKey is the key that holds user's DisplayName
	UserDisplayNameKey
	// UserEmailKey is the key that holds user's Email
	UserEmailKey
	// UserRoleKey is the key that holds user's Role
	UserRoleKey
)

var (
	errorParsingToken = errors.New("error parsing JWT token")
)

// New instantiates new JWT service
func New(key string, duration int, algo string, sess SessionStorer) *JWT {
	signingMethod := jwt.GetSigningMethod(algo)
	if signingMethod == nil {
		panic("invalid signing method")
	}
	return &JWT{
		key:      []byte(key),
		duration: time.Duration(duration) * time.Minute,
		algo:     signingMethod,
		sess:     sess,
	}
}

// JWT contains data necessery for jwt auth
type JWT struct {
	// Secret key used for signing.
	key []byte

	// Duration for which the jwt token is valid.
	duration time.Duration

	// JWT signing algorithm
	algo jwt.SigningMethod

	// Session storer interface
	sess SessionStorer
}

// SessionStorer represents session store interface
type SessionStorer interface {
	Get(string) (*chisk.AuthUser, error)
}

// GenerateToken generates new jwt token
func (j *JWT) GenerateToken() (string, error) {
	t := time.Now()
	return jwt.NewWithClaims(j.algo, jwt.MapClaims{
		"iat": t.Unix(),
		"exp": t.Add(j.duration).Unix(),
	}).SignedString(j.key)
}

// ParseToken parses JWT token
func (j *JWT) ParseToken(token string) error {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if token.Method != j.algo {
			return nil, errorParsingToken
		}
		return j.key, nil
	})

	if err != nil || !t.Valid {
		return err
	}

	return nil
}

type errMsg struct {
	Message string `json:"message"`
}

// MWFunc is a middleware func for JWT Authorization
func (j *JWT) MWFunc(next http.Handler) http.Handler {
	var (
		missingAuthorizationHeader, _ = json.Marshal(errMsg{Message: "Missing Authorization header"})
		missingBearerKeyword, _       = json.Marshal(errMsg{Message: "Missing Bearer keyword"})
		cannotParseToken, _           = json.Marshal(errMsg{Message: errorParsingToken.Error()})
		cannotRetreiveSession, _      = json.Marshal(errMsg{Message: "Error retreiving session"})
	)

	fn := func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header.Get("Authorization")
		if ah == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(missingAuthorizationHeader)
			return
		}

		spl := strings.Split(ah, " ")
		if spl[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(missingBearerKeyword)
			return
		}

		token := spl[1]
		if err := j.ParseToken(token); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(cannotParseToken)
			return
		}

		user, err := j.sess.Get(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(cannotRetreiveSession)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, TokenKey, token)
		ctx = context.WithValue(ctx, UserIDKey, user.ID)
		ctx = context.WithValue(ctx, UserDisplayNameKey, user.DisplayName)
		ctx = context.WithValue(ctx, UserRoleKey, user.Role)
		ctx = context.WithValue(ctx, UserEmailKey, user.Email)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
