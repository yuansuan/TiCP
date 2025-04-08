package v1

import (
	"testing"
)

func TestGenerateCredentials(t *testing.T) {
	accessKey, secretKey, err := GenerateCredentials()
	if err != nil {
		t.Errorf("Error generating credentials: %v", err)
	}

	if len(accessKey) < accessKeyMinLen || len(accessKey) > accessKeyMaxLen {
		t.Errorf("Access key length out of range: %v", accessKey)
	}

	if len(secretKey) < secretKeyMinLen || len(secretKey) > secretKeyMaxLen {
		t.Errorf("Secret key length out of range: %v", secretKey)
	}
}

func TestGetNewCredentialsWithMetadata(t *testing.T) {
	m := make(map[string]interface{})
	m["key1"] = "value1"
	m["key2"] = "value2"

	accessKey, secretKey, sessionToken, err := GetNewCredentialsWithMetadata(m)
	if err != nil {
		t.Errorf("Error generating credentials with metadata: %v", err)
	}

	if len(accessKey) < accessKeyMinLen || len(accessKey) > accessKeyMaxLen {
		t.Errorf("Access key length out of range: %v", accessKey)
	}

	if len(secretKey) < secretKeyMinLen || len(secretKey) > secretKeyMaxLen {
		t.Errorf("Secret key length out of range: %v", secretKey)
	}

	if len(sessionToken) == 0 {
		t.Errorf("Session token not generated")
	}

	// if _, err := JWTSignWithAccessKey(accessKey, sessionToken, globalSecretKey); err != nil {
	// 	t.Errorf("Error verifying JWT token: %v", err)
	// }

	if m["accessKey"] != accessKey {
		t.Errorf("Access key not added to metadata")
	}
}

func TestJWTSignWithAccessKey(t *testing.T) {
	m := make(map[string]interface{})
	m["key1"] = "value1"
	m["key2"] = "value2"

	accessKey, _, err := GenerateCredentials()
	if err != nil {
		t.Errorf("Error generating credentials: %v", err)
	}

	token, err := JWTSignWithAccessKey(accessKey, m, globalSecretKey)
	if err != nil {
		t.Errorf("Error signing JWT token: %v", err)
	}

	if len(token) == 0 {
		t.Errorf("JWT token not generated")
	}

	// if _, err := JWTVerifyWithAccessKey(accessKey, token, globalSecretKey); err != nil {
	// 	t.Errorf("Error verifying JWT token: %v", err)
	// }
}
