package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/internal/config"
	"wallet/internal/handlers/wallet"
	"wallet/pkg/customValidator"
	"wallet/storage/postgres"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestPingRoute(t *testing.T) {

	config := config.New("../config_test.env")

	storage := postgres.New(config)

	cv := customValidator.NewCustomValidator()

	tests := []struct {
		name            string
		requestBody     string
		walletRepoError error
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:           "successful deposit",
			requestBody:    `{"id":"156f5065-4a38-4abe-bf06-2dea11850408","operationType":"DEPOSIT","Amount":1}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "Balance changed",
		},
		{
			name:           "empty request body",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `request body is empty`,
		},
		{
			name:           "wallet not found",
			requestBody:    `{"id":"156f5065-4a38-4abe-bf06-2dea55850405","operationType":"DEPOSIT","Amount":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "this wallet id not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			handler := wallet.ChangeAmountWallet(&zap.Logger{}, cv, storage)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(tt.requestBody))

			handler(ctx)

			body := w.Body.String()

			var decodedBody string
			require.NoError(t, json.Unmarshal([]byte(body), &decodedBody))
			require.Equal(t, tt.expectedBody, decodedBody)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
