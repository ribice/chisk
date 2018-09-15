package jwt_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ribice/chisk/model"

	"github.com/ribice/chisk/mock"
	"github.com/ribice/chisk/pkg/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name        string
		algo        string
		shouldPanic bool
		wantErr     bool
		wantToken   bool
	}{
		{
			name:        "Invalid algo",
			algo:        "HS128",
			shouldPanic: true,
		},
		{
			name:      "Success",
			algo:      "HS256",
			wantToken: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(func() {
					jwt.New("testkey", 10, tt.algo, nil)
				})
				return
			}
			j := jwt.New("testkey", 10, tt.algo, nil)
			token, err := j.GenerateToken()
			assert.Equal(tt.wantToken, token != "")
			assert.Equal(tt.wantErr, err != nil)
		})
	}
}

func TestParseToken(t *testing.T) {
	cases := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Algorithm missmatch",
			token:   "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTkyNjIzOTAyMn0.iVZCNigNdVvjPr8nD8He-Qb9t1sN-lTyi9tmAJEvDAp92Yqv1KPUpuYDGWxANHhG",
			wantErr: true,
		},
		{
			name:  "Success",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTkyNjIzOTAyMn0.ZTBhczNrhnpQqOFDOlQahRzvj8XuoIn0eyvYY8uhC8c",
		},
	}
	j := jwt.New("testingsecret", 20, "HS256", nil)
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := j.ParseToken(tt.token)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func testHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
	}
	return http.HandlerFunc(fn)
}

func TestMWFunc(t *testing.T) {
	assert := assert.New(t)
	cases := []struct {
		name        string
		headers     map[string]string
		wantMessage string
		wantStatus  int
		sess        *mock.Session
		sessUser    *chisk.AuthUser
	}{
		{
			name:        "Missing authorization header",
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Missing Authorization header",
		},
		{
			name:       "Missing bearer keyword",
			wantStatus: http.StatusUnauthorized,
			headers: map[string]string{
				"Authorization": "tokengoeshere",
			},
			wantMessage: "Missing Bearer keyword",
		},
		{
			name:       "Invalid token",
			wantStatus: http.StatusUnauthorized,
			headers: map[string]string{
				"Authorization": "Bearer invalidtoken",
			},
			wantMessage: "error parsing JWT token",
		},
		{
			name:       "Invalid token",
			wantStatus: http.StatusUnauthorized,
			headers: map[string]string{
				"Authorization": "Bearer invalidtoken",
			},
			wantMessage: "error parsing JWT token",
		},
		{
			name:       "Fail on session extraction",
			wantStatus: http.StatusUnauthorized,
			headers: map[string]string{
				"Authorization": mock.ValidJWTToken,
			},
			sess: &mock.Session{
				GetFn: func(string) (*chisk.AuthUser, error) {
					return nil, mock.ErrGeneric
				},
			},
			wantMessage: "Error retreiving session",
		},
		{
			name:       "Success",
			wantStatus: http.StatusOK,
			headers: map[string]string{
				"Authorization": mock.ValidJWTToken,
			},
			sess: &mock.Session{
				GetFn: func(string) (*chisk.AuthUser, error) {
					return &chisk.AuthUser{
						ID:          "uid",
						DisplayName: "johndoe",
						Email:       "johndoe@mail.com",
						Role:        1,
					}, nil
				},
			},
			sessUser: &chisk.AuthUser{
				ID:          "uid",
				DisplayName: "johndoe",
				Email:       "johndoe@mail.com",
				Role:        1,
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			j := jwt.New("testingsecret", 10, "HS256", tt.sess)
			ts := httptest.NewServer(j.MWFunc(testHandler()))
			defer ts.Close()

			req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
			assert.NoError(err)

			ctx := context.Background()
			req = req.WithContext(ctx)

			for k, v := range tt.headers {
				req.Header.Add(k, v)
			}

			res, err := ts.Client().Do(req)
			assert.NoError(err)

			assert.Equal(tt.wantStatus, res.StatusCode)

			if tt.wantStatus == http.StatusOK {
				// TODO: Test values in context
				// assert.Equal(strings.Split(tt.headers["Authorization"], " ")[1], ctx.Value(jwt.TokenKey).(string))
				// assert.Equal(tt.sessUser.ID, ctx.Value(jwt.UserIDKey).(string))
				// assert.Equal(tt.sessUser.Role, ctx.Value(jwt.UserRoleKey).(int))
				// assert.Equal(tt.sessUser.DisplayName, ctx.Value(jwt.UserDisplayNameKey).(string))
				// assert.Equal(tt.sessUser.Email, ctx.Value(jwt.UserEmailKey).(string))
			} else {
				defer res.Body.Close()

				b, err := ioutil.ReadAll(res.Body)
				msg := &errMsg{}
				assert.NoError(json.Unmarshal(b, msg))

				assert.NoError(err)
				assert.Equal(tt.wantMessage, msg.Message)
			}
		})
	}
}

type errMsg struct {
	Message string `json:"message"`
}
