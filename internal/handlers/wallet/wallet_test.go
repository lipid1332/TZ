package wallet_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet/internal/handlers/wallet"
	"wallet/internal/handlers/wallet/mocks"
	"wallet/pkg/customValidator"
	"wallet/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name            string
		requestBody     string
		walletRepoError error
		expectedStatus  int
		expectedBody    string
	}{
		{
			name:            "successful deposit",
			requestBody:     `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"DEPOSIT","Amount":1}`,
			walletRepoError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody:    "Balance changed",
		},
		{
			name:           "empty request body",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `request body is empty`,
		},
		{
			name:           "invalid uuid",
			requestBody:    `{"id":"156f5065","operationType":"DEPOSIT","Amount":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field id must be a valid uuid",
		},
		{
			name:           "invalid operationType",
			requestBody:    `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"","Amount":1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field operationType is required",
		},
		{
			name:           "invalid Amount",
			requestBody:    `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"DEPOSIT","Amount":-1}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "field amount is invalid",
		},
		{
			name:            "wallet not found",
			requestBody:     `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"DEPOSIT","Amount":1}`,
			walletRepoError: storage.ErrWalletNotFound,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "this wallet id not found",
		},
		{
			name:            "negative amount",
			requestBody:     `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"WITHDRAW","Amount":1}`,
			walletRepoError: storage.ErrNegativeAmount,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "amount cant be negative",
		},
		{
			name:            "wallet interface error",
			requestBody:     `{"id":"156f5065-4a38-4abe-bf06-2dea11850405","operationType":"WITHDRAW","Amount":1}`,
			walletRepoError: errors.New("database error"),
			expectedStatus:  http.StatusInternalServerError,
			expectedBody:    "something wrong please retry",
		},
	}

	validator := customValidator.NewCustomValidator()
	logger := zap.NewNop()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			walletRepo := mocks.NewWalletRepo(t)

			if tt.walletRepoError != nil || tt.expectedStatus == http.StatusOK {
				walletRepo.
					On("ChangeAmount",
						mock.AnythingOfType("string"),
						mock.AnythingOfType("int"),
						mock.AnythingOfType("string")).
					Return(tt.walletRepoError)
			}

			handler := wallet.ChangeAmountWallet(logger, validator, walletRepo)

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBufferString(tt.requestBody))

			handler(ctx)

			body := w.Body.String()

			var decodedBody string
			require.NoError(t, json.Unmarshal([]byte(body), &decodedBody))
			require.Equal(t, tt.expectedBody, decodedBody)

			require.Equal(t, tt.expectedStatus, w.Code)

			walletRepo.AssertExpectations(t)
		})
	}
}
