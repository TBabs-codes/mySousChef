package auth

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

//Tests HashPassword and CheckPasswordHash functions when using correct and incorrect passwords
func TestHashPassword(t *testing.T){
	password := "ReallyGoodPassword"

	//Hash password
	hash, err := HashPassword(password)
	if err != nil || len(hash) == 0 {
		t.Errorf("Hashing failed: %v", err)
	}

	//Valid password validation
	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("Valid password failed validation check: %v", err)
	}

	//Invalid password validation
	err = CheckPasswordHash("WrongPassword", hash)
	if err == nil {
		t.Errorf("Expected error for incorrect password: %v", err)
	}

	//Hash uniqueness
	hash2, err := HashPassword(password)
	if err != nil {
		t.Errorf("HashPassword failed to hash: %v", err)
	}
	if hash == hash2 {
		t.Errorf("HashPassword didn't produce unique hash")
	}
}

// Test JSON webtoken generation and validation functions (MakeJWT and ValidateJWT)
func TestJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret-key123"
	token_expiresIn := time.Second*1

	token, err := MakeJWT(userID, secret, token_expiresIn)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}

	userID_test, err := ValidateJWT(token, secret)
	if err != nil {
		t.Errorf("Error validating JWT: %v", err)
	}
	if userID_test != userID {
		t.Errorf("Validation resulted in wrong userID.")
	}

	time.Sleep(time.Second*2)

	userID_test2, err := ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("No error validating JWT after expiration: %v", err)
	}

	if userID_test2 == userID {
		t.Error("UserID still validated after token expiration")
	}
}

//func GetBearerToken(headers http.Header) (string, error)
func TestGetAPIKey(t *testing.T) {
	type output struct {
		key string
		err error
	}

	type test struct {
		header          http.Header
		expected_result output
	}

	tests := []test{
		{header: http.Header{"Authorization": []string{"Bearer 12345678-abcd-90ef-gh12-ijklmnopqrst"}}, expected_result: output{key: "12345678-abcd-90ef-gh12-ijklmnopqrst", err: nil}},
		{header: http.Header{"Authorization": []string{"Bear 12345678-abcd-90ef-gh12-ijklmnopqrst"}}, expected_result: output{key: "", err: errors.New("malformed authorization header")}},
		{header: http.Header{}, expected_result: output{key: "", err: errors.New("no authorization header included")}},
	}

	for i, tc := range tests {
		key, err := GetBearerToken(tc.header)
		if !reflect.DeepEqual(tc.expected_result.key, key) {
			t.Fatalf("test %d failed: expected key: %s, actual key: %s", i+1, tc.expected_result.key, key)
		} else if err != nil && tc.expected_result.err.Error() != err.Error() {
			t.Fatalf("test %d failed: expected err: %v, actual err: %v", i+1, tc.expected_result.err, err)
		}
	}
}
