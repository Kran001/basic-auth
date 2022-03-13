package validators

import (
	"errors"
	"regexp"

	crypter "github.com/NGRsoftlab/ngr-crypter"

	"github.com/Kran001/basic-auth/pkg/logging"
)

//pattern.
const (
	PswPattern   = `^[\w@]{8,}$`
	LoginPattern = `^[a-zA-Z]{1}[\w@\.]{5,}$`

	PasswordPhrase = `testpassword phrase`
)

// ValidateByPattern - Validating pattern chars
func ValidateByPattern(word, pattern string) bool {
	matched, err := regexp.MatchString(pattern, word)
	if err != nil {
		logging.Logger.Errorf("Error matching strings. Reason: %s. Returning invalid data. ", err.Error())

		return false
	}

	return matched
}

// ValidateAndCryptPsw - Validating password chars
func ValidateAndCryptPsw(psw, key string) (string, error) {
	matched := ValidateByPattern(psw, PswPattern)
	if !matched {
		logging.Logger.Error("Error match string")

		return "", errors.New("no matched string")
	}

	return crypter.Encrypt([]byte(key), psw)
}
