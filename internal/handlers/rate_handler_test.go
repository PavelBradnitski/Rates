package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PavelBradnitski/Rates/internal/mocks"
	"github.com/PavelBradnitski/Rates/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_GetAllRates(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		mockSvcBehavior func(mockSvc *mocks.MockRateService)
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "Success",
			mockSvcBehavior: func(mockSvc *mocks.MockRateService) {
				mockSvc.On("GetAllRates", mock.Anything).Return([]models.Rate{
					{ID: 315,
						CurID:           440,
						Date:            "2025-01-10",
						CurAbbreviation: "AUD",
						CurScale:        1,
						CurOfficialRate: 2.1555},
					{ID: 316,
						CurID:           510,
						Date:            "2025-01-10",
						CurAbbreviation: "AMD",
						CurScale:        1000,
						CurOfficialRate: 8.7102},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"Cur_Abbreviation":"AUD", "Cur_ID":440, "Cur_OfficialRate":2.1555, "Cur_Scale":1, "Date":"2025-01-10", "ID":315},{"Cur_Abbreviation":"AMD", "Cur_ID":510, "Cur_OfficialRate":8.7102, "Cur_Scale":1000, "Date":"2025-01-10", "ID":316}]`,
		},
		{
			name: "Service Returns Error",
			mockSvcBehavior: func(mockSvc *mocks.MockRateService) {
				mockSvc.On("GetAllRates", mock.Anything).Return([]models.Rate{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to retrieve rates"}`,
		},
		{
			name: "No Rates Found",
			mockSvcBehavior: func(mockSvc *mocks.MockRateService) {
				mockSvc.On("GetAllRates", mock.Anything).Return([]models.Rate{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Rates not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mocks.MockRateService{}
			tt.mockSvcBehavior(mockSvc)

			handler := RateHandler{service: mockSvc}

			router := gin.Default()
			handler.RegisterRoutes(router)

			req, _ := http.NewRequest("GET", "/rate/", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.JSONEq(t, tt.expectedBody, resp.Body.String())

			mockSvc.AssertExpectations(t)
		})
	}
}
