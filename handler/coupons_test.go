package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/blobstore"
	"github.com/rohitxdev/go-api-starter/config"
	"github.com/rohitxdev/go-api-starter/database"
	"github.com/rohitxdev/go-api-starter/email"
	"github.com/rohitxdev/go-api-starter/handler"
	"github.com/rohitxdev/go-api-starter/kvstore"
	"github.com/rohitxdev/go-api-starter/logger"
	"github.com/rohitxdev/go-api-starter/repo"
	"github.com/stretchr/testify/assert"
)

func TestCoupons(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	//Set up logger
	logr := logger.New(os.Stderr, cfg.IsDev)
	//Connect to postgres database
	db, err := database.NewPostgreSQL(cfg.DatabaseURL)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic("Failed to close database: " + err.Error())
		}
	}()

	//Connect to KV store
	kv, err := kvstore.New("kv", time.Minute*5)
	if err != nil {
		panic("Failed to connect to KV store: " + err.Error())
	}

	defer func() {
		kv.Close()
	}()

	// Create repo
	r, err := repo.New(db)
	if err != nil {
		panic("Failed to create repo: " + err.Error())
	}
	defer r.Close()

	bs, err := blobstore.New(cfg.S3Endpoint, cfg.S3DefaultRegion, cfg.AWSAccessKeyID, cfg.AWSAccessKeySecret)
	if err != nil {
		panic("Failed to connect to S3 client: " + err.Error())
	}
	e, err := email.New(&email.SMTPCredentials{})
	assert.Nil(t, err)

	h, err := handler.New(&handler.Services{
		BlobStore: bs,
		Config:    cfg,
		KVStore:   kv,
		Logger:    logr,
		Repo:      r,
		Email:     e,
	})
	assert.Nil(t, err)

	cookie, err := createTestSessionCookie(h, cfg.JWTSecret)
	assert.Nil(t, err)

	t.Run("GET /coupons", func(t *testing.T) {
		type args struct {
			isAuthenticated bool
			reqOpts         *httpRequestOpts
		}
		tests := []struct {
			name       string
			args       args
			wantStatus int
		}{
			{
				name: "Unauthorized",
				args: args{
					reqOpts: &httpRequestOpts{
						method: http.MethodGet,
						path:   "/coupons",
					},
				},
				wantStatus: http.StatusUnauthorized,
			},
			{
				name: "Authorized",
				args: args{
					reqOpts: &httpRequestOpts{
						method: http.MethodGet,
						path:   "/coupons",
						headers: map[string]string{
							"Cookie": cookie,
						},
					},
					isAuthenticated: true,
				},
				wantStatus: http.StatusOK,
			},
			{
				name: "Authorized for admin",
				args: args{
					reqOpts: &httpRequestOpts{
						method: http.MethodGet,
						path:   "/coupons/all",
						headers: map[string]string{
							"Cookie":       cookie,
							"Content-Type": "application/json",
						},
					},
					isAuthenticated: true,
				},
				wantStatus: http.StatusOK,
			},
			{
				name: "Create coupon",
				args: args{
					reqOpts: &httpRequestOpts{
						method: http.MethodPost,
						path:   "/coupons",
						headers: map[string]string{
							"Cookie":       cookie,
							"Content-Type": "application/json",
						},
						body: echo.Map{
							"discountPercent": 10,
							"userId":          1,
						},
					},
					isAuthenticated: true,
				},
				wantStatus: http.StatusOK,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := createHttpRequest(tt.args.reqOpts)
				if tt.args.isAuthenticated {
					req.Header.Set("Cookie", cookie)
				}
				assert.Nil(t, err)
				res := httptest.NewRecorder()
				h.ServeHTTP(res, req)
				assert.Equal(t, tt.wantStatus, res.Code)
			})
		}
	})
}
