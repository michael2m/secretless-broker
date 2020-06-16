package mysql

import (
	"github.com/cyberark/secretless-broker/internal/plugin/connectors/tcp/ssl"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectedFields(t *testing.T) {
	credentials := map[string][]byte{
		"host":     []byte("myhost"),
		"port":     []byte("1234"),
		"username": []byte("myusername"),
		"password": []byte("mypassword"),
		"sslmode":  []byte("disable"),
	}

	expectedConnDetails := ConnectionDetails{
		Host: "myhost",
		Options: ssl.NewSSLOptions(credentials),
		Password: "mypassword",
		Port:     1234,
		Username: "myusername",
	}

	actualConnDetails, err := NewConnectionDetails(credentials)
	assert.Nil(t, err)

	if err == nil {
		assert.EqualValues(t, expectedConnDetails, *actualConnDetails)
	}
}

func TestDefaultPort(t *testing.T) {
	credentials := map[string][]byte{
		"host":     []byte("myhost"),
		"username": []byte("myusername"),
		"password": []byte("mypassword"),
	}

	expectedConnDetails := ConnectionDetails{
		Host:     "myhost",
		Port:     DefaultMySQLPort,
		Username: "myusername",
		Password: "mypassword",
		Options: ssl.NewSSLOptions(credentials),
	}

	actualConnDetails, err := NewConnectionDetails(credentials)
	assert.Nil(t, err)

	if err == nil {
		assert.EqualValues(t, expectedConnDetails, *actualConnDetails)
	}
}

func TestAddress(t *testing.T) {
	credentials := map[string][]byte{
		"host": []byte("myhost2"),
		"port": []byte("12345"),
	}

	actualConnDetails, err := NewConnectionDetails(credentials)
	assert.Nil(t, err)

	if err == nil {
		assert.EqualValues(t, "myhost2:12345", (*actualConnDetails).Address())
	}
}
