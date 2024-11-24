package handler_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

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

func TestGetProducts(t *testing.T) {
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

	t.Run("GET /products", func(t *testing.T) {
		req, err := createHttpRequest(&httpRequestOpts{
			method: http.MethodGet,
			path:   "/products",
		})
		assert.Nil(t, err)
		res := httptest.NewRecorder()
		h.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)
	})
}