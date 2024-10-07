package domen

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	LuhnError       = errors.New("the credit card number you entered failed the Luhn Check")
	IINError        = errors.New("the length of the card does not match the manufacturer's standards")
	UnknownIINError = errors.New("we couldn't find an institution that matched your credit card number")
	ExpError        = errors.New("your credit card is expired")
	UnexpectedError = errors.New("unexpected error")
)

type CardValidator struct {
	CardIssuerInfo map[int][]CardIssuer
}

// read and transform json data into map[int][]CardIssuer
func InitCardValidator(initData string) (*CardValidator, error) {

	data, err := os.ReadFile(initData)
	if err != nil {
		return nil, err
	}

	var jsonDataArray []map[string][]CardIssuer

	err = json.Unmarshal(data, &jsonDataArray)
	if err != nil {
		return nil, err
	}

	issuersMap := make(map[int][]CardIssuer)

	for _, item := range jsonDataArray {
		for key, value := range item {

			intKey, err := strconv.Atoi(key)
			if err != nil {
				return nil, err
			}
			issuersMap[intKey] = value
		}
	}

	//compile rgx instance from string presentation
	for key, issuerList := range issuersMap {
		for i := range issuerList {
			rgx, err := regexp.Compile(issuerList[i].RgxStr)
			//validation check for regex
			if err != nil {
				return nil, fmt.Errorf("failed to compile regex for issuer %s: %w", issuerList[i], err)
			}
			issuersMap[key][i].Rgx = rgx
		}
	}

	err = validateSegments(&issuersMap)
	if err != nil {
		return nil, err
	}

	return &CardValidator{
		CardIssuerInfo: issuersMap,
	}, nil
}

// the json data check function about the permissible PAN lengths for different issuers of bank cards
func validateSegments(issuersMap *map[int][]CardIssuer) error {
	mtErrText := "mismatched type, unable convert to integer"
	dfErrText := "wrong data format for segments"

	for _, cardIssuers := range *issuersMap {
		for _, issuer := range cardIssuers {
			//valid data format it`s single integer or 2 integers separated by `-`
			for _, segment := range issuer.CardLength {
				if strings.Contains(segment, "-") {
					parts := strings.Split(segment, "-")
					if len(parts) != 2 {
						return fmt.Errorf("%s\nin range: %s\nIssuer: %s", mtErrText, segment, issuer)
					} else {
						start, startErr := strconv.Atoi(parts[0])
						end, endErr := strconv.Atoi(parts[1])
						//check that used correct data types and values
						if startErr != nil || endErr != nil || end <= start {
							return fmt.Errorf("%v\nin range: %s\nmin val: %v\nmax val: %v\nIssuer: %s",
								dfErrText, segment, start, end, issuer)
						}
					}
				} else {
					if _, err := strconv.Atoi(segment); err != nil {
						return fmt.Errorf("%v\nvalue: %s\nIssuer: %s", mtErrText, segment, issuer)
					}
				}
			}
		}
	}
	return nil
}

func (v *CardValidator) Validate(card Card) error {
	// first stage - validate PAN with Luhn algorithm
	luhnValidation := v.ValidateSumLuhn(card.PAN)
	if !luhnValidation {
		return LuhnError
	}
	//second stage - check whether the card is not expired
	ExpDateValidation := v.ValidateExpDate(card)
	if !ExpDateValidation {
		return ExpError
	}

	//last stage - validate issuer (regex for BIN, PAN within the permissible values )
	err := v.ValidateIssuer(card)

	if errors.Is(err, IINError) || errors.Is(err, UnknownIINError) {
		return err
	} else if err != nil {
		log.Println("Unexpected error during card validation: ", err)
		return UnexpectedError
	}

	return nil
}

func (v *CardValidator) ValidateIssuer(card Card) error {
	firstSignInCard := strconv.Itoa(card.PAN)[0]
	firstDigitInCard, err := strconv.Atoi(string(firstSignInCard))

	if err != nil {
		return err
	}
	//iterate over
	for _, issuer := range v.CardIssuerInfo[firstDigitInCard] {
		if issuer.Rgx.MatchString(strconv.Itoa(card.PAN)) {
			if isInRange(len(strconv.Itoa(card.PAN)), issuer.CardLength) {
				return nil
			} else {
				return IINError
			}
		}
	}
	return UnknownIINError
}

func isInRange(value int, ranges []string) bool {
	for _, segment := range ranges {
		if strings.Contains(segment, "-") {
			parts := strings.Split(segment, "-")
			start, _ := strconv.Atoi(parts[0])
			end, _ := strconv.Atoi(parts[1])
			if value >= start && value <= end {
				return true
			}
		} else {
			if num, err := strconv.Atoi(segment); err == nil {
				if value == num {
					return true
				}
			}
		}
	}
	return false
}

func (v *CardValidator) ValidateExpDate(card Card) bool {
	currentTime := time.Now()
	currentMonth := int(currentTime.Month())
	currentYear := currentTime.Year()

	if card.ExpirationYear >= currentYear && card.ExpirationMonth >= currentMonth {
		return true
	} else {
		return false
	}
}

func (v *CardValidator) ValidateSumLuhn(cardNumber int) bool {

	data := strconv.Itoa(cardNumber)

	sum := 0

	for i := len(data) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(data[i]))
		if (len(data)-i)%2 == 0 {
			digit = doubleAndSumDigits(digit)
		}
		sum += digit
	}
	return sum%10 == 0
}

func doubleAndSumDigits(digit int) int {
	doubled := digit * 2
	if doubled > 9 {
		return doubled - 9
	}
	return doubled
}
