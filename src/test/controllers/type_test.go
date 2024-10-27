package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fariasBP/acapela-api/src/controllers"
	"github.com/jaswdr/faker"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTypeProduct(t *testing.T) {
	// inicializando
	if err := godotenv.Load("./../../../.env"); err != nil {
		t.Errorf("Error loading env vars filed: %s", err)
	}
	fake := faker.New()
	// test
	var typeBody = `{"name": "` + fake.App().Name() + `", "description": "` + fake.Lorem().Sentence(30) + `"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(typeBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controllers.CreateTypeProduct(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		// assert.Equal(t, typeBody, rec.Body.String())
		t.Log(rec.Body.String())
	}
}
