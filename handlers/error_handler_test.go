package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadRequestErrorHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	BadRequestErrorHandler(recorder, request)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestNotFoundErrorHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	NotFoundErrorHandler(recorder, request)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}

func TestInternalServerErrorHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	err := error(nil)

	InternalServerErrorHandler(recorder, request, err)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Equal(t, JsonContentType, recorder.Header().Get(ContentType))
}
