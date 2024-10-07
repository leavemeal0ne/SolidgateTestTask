package handler

import (
	"bytes"
	"encoding/json"
	"github.com/leavemeal0ne/SolidgateTestTask/internal/domen"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_ValidateCardHandler(t *testing.T) {

	testCases := []struct {
		qualifier     string
		card          domen.Card
		expectedError string
	}{
		{
			qualifier: "Valid Card",
			card: domen.Card{
				PAN:             4929164492674746,
				ExpirationMonth: 12,
				ExpirationYear:  2055, // update until 2056! :)
			},
			expectedError: "",
		},
		{
			qualifier: "Failed Luhn Check",
			card: domen.Card{
				PAN:             4916527284180041,
				ExpirationMonth: 12,
				ExpirationYear:  2024,
			},
			expectedError: domen.LuhnError.Error(),
		},
		{
			qualifier: "Expired Card",
			card: domen.Card{
				PAN:             4916527284180046,
				ExpirationMonth: 12,
				ExpirationYear:  2000,
			},
			expectedError: domen.ExpError.Error(),
		},
		{
			qualifier: "Unknown Issuer",
			card: domen.Card{
				PAN:             3530111333300000, //JCB [3528–3589] 16-19
				ExpirationMonth: 12,
				ExpirationYear:  2055,
			},
			expectedError: domen.UnknownIINError.Error(),
		},
	}

	validator, err := domen.InitCardValidator("../../card_classification_data/CardData.json")
	if err != nil {
		t.Fatalf("Unable to init validator: %s", err.Error())
	}

	handler := InitHandler(validator).InitRoutes()

	for _, testCase := range testCases {

		t.Run(testCase.qualifier, func(t *testing.T) {

			body, err := json.Marshal(testCase.card)
			if err != nil {
				t.Fatalf("Unable to marshal card: %s", err.Error())

			}

			req := httptest.NewRequest(http.MethodPost, "/validate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)
			response := domen.Response{}
			err = json.NewDecoder(recorder.Result().Body).Decode(&response)
			if err != nil {
				t.Errorf("Unable to decode response: %s", err.Error())
			}
			switch response.Error == nil {
			case true:
				if testCase.expectedError != "" {
					t.Errorf("\nExpected error: %s, but got none\n Card: %#v", testCase.expectedError, testCase.card)
				}
			case false:
				if response.Error.Message != testCase.expectedError || response.Valid {
					t.Errorf("\nExpected: %s\nActual: %s\n Card: %#v", testCase.expectedError, response.Error.Message, testCase.card)
				}
			}
		})
	}
}
