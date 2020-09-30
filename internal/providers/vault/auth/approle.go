package auth

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"os"
	"strconv"
)

// DefaultAppRoleLoginPath path to AppRole login in Vault.
const DefaultAppRoleLoginPath = "auth/approle/login"

// EnvVaultAppRoleRoleID environment variable name of role ID used in AppRole auth method.
const EnvVaultAppRoleRoleID = "VAULT_APPROLE_ROLE_ID"

// EnvVaultAppRoleSecretID environment variable name of (possibly wrapped) secret ID used in AppRole auth method.
const EnvVaultAppRoleSecretID = "VAULT_APPROLE_SECRET_ID"

// EnvVaultAppRoleUnwrap environment variable name of unwrapping of secret ID in AppRole auth method.
const EnvVaultAppRoleUnwrap = "VAULT_APPROLE_UNWRAP"

// AppRoleConfig holds AppRole auth method configuration.
type AppRoleConfig struct {
	LoginPath string
	RoleID string
	SecretID string
	Unwrap bool
}

// AppRoleAuthMethodFactory performs AppRole auth method for given client.
func AppRoleAuthMethodFactory(client *vault.Client) (*vault.Secret, error) {
	config := getAppRoleConfigFromEnv()

	var secret *vault.Secret
	var err error
	if secret, err = getTokenFromAppRoleLogin(client, config); err != nil {
		return nil, err
	}

	return secret, nil
}

// getAppRoleConfigFromEnv returns AppRole auth method configuration from environment variables.
func getAppRoleConfigFromEnv() *AppRoleConfig {
	config := &AppRoleConfig{}
	config.RoleID = os.Getenv(EnvVaultAppRoleRoleID)
	config.SecretID = os.Getenv(EnvVaultAppRoleSecretID)
	// Ignore error, default unwrap = false
	config.Unwrap, _ = strconv.ParseBool(os.Getenv(EnvVaultAppRoleUnwrap))
	config.LoginPath = getLoginPathFromEnv(DefaultAppRoleLoginPath)
	return config
}

// getTokenFromAppRoleLogin performs AppRole login for client using given AppRole auth method configuration.
func getTokenFromAppRoleLogin(client *vault.Client, config *AppRoleConfig) (*vault.Secret, error) {
	// Conditionally unwrap secret ID (referred to as pull mode in Vault)
	if config.Unwrap {
		secret, err := client.Logical().Unwrap(config.SecretID)
		if err != nil {
			return nil, fmt.Errorf("HashiCorp Vault provider unwrapping failed: %s", err)
		}

		// Replace with unwrapped secret ID in config
		config.SecretID = secret.Data["secret_id"].(string)
	}

	data := map[string]interface{}{
		"role_id":   config.RoleID,
		"secret_id": config.SecretID,
	}

	return client.Logical().Write(config.LoginPath, data)
}