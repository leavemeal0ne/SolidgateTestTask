package domen

import "testing"

func TestCardValidator_ValidateSumLuhn(t *testing.T) {
	validator := CardValidator{}
	testcases := []struct {
		cardNumber     int
		expectedResult bool
	}{
		{
			cardNumber:     4121502025433724,
			expectedResult: true,
		},
		{
			cardNumber:     6011488915652052,
			expectedResult: true,
		},
		{
			cardNumber:     347349026752694,
			expectedResult: true,
		},
		{
			cardNumber:     4121502025433720,
			expectedResult: false,
		},
		{
			cardNumber:     6011488915652051,
			expectedResult: false,
		},
		{
			cardNumber:     347349026752695,
			expectedResult: false,
		},
	}

	for _, testcase := range testcases {
		if validator.ValidateSumLuhn(testcase.cardNumber) != testcase.expectedResult {
			t.Errorf("Expected %t bot got %t {%d}",
				testcase.expectedResult, validator.ValidateSumLuhn(testcase.cardNumber), testcase.cardNumber)
		}
	}

}
