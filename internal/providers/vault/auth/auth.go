package auth

import (
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"os"
)

// EnvVarVaultAuthMethod environment variable name of Vault auth method to use.
const EnvVarVaultAuthMethod = "VAULT_AUTH_METHOD"

// EnvVarVaultLoginPath path to auth method's login in Vault.
const EnvVarVaultLoginPath = "VAULT_LOGIN_PATH"

// VaultAuthMethodFactories mapping of supported Vault auth methods by their IDs. Must have either an auth method
// factory or `nil` to signify no action is required, e.g. for Token auth method.
// NOTE: the "Token" auth methods is added for consistency. Token auth method requires no additional authentication
// steps, unlike other auth methods in Vault.
var VaultAuthMethodFactories = map[string]func(*vault.Client) (*vault.Secret, error) {
	"AppRole": AppRoleAuthMethodFactory,
	"Token": nil,
}

// Authenticate a Vault client using any of the supported auth methods.
// NOTE: All auth methods require an initialized Vault API client.
func Authenticate(client *vault.Client) error {
	authMethod := getAuthMethodFromEnv()
	applyAuthMethod, ok := VaultAuthMethodFactories[authMethod]
	if !ok {
		return fmt.Errorf("HashiCorp Vault provider unsupported auth method: %s", authMethod)
	}

	// No (additional) authentication steps required, hence done
	if applyAuthMethod == nil {
		return nil
	}

	secret, err := applyAuthMethod(client)
	if err != nil {
	  return fmt.Errorf("HashiCorp Vault provider auth method failed: %s", err)
	}
	if secret == nil || secret.Auth == nil {
		return fmt.Errorf("HashiCorp Vault provider login failed (no secret or auth info)")
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}

// getAuthMethodFromEnv returns the auth method from environment variable. If not set, it defaults to "Token" for
// backwards compatibility.
func getAuthMethodFromEnv() string {
	authMethod, ok := os.LookupEnv(EnvVarVaultAuthMethod)
	if !ok {
		return "Token"
	}
	return authMethod
}

// getLoginPathFromEnv returns login path from environment variable or returns given default value otherwise.
func getLoginPathFromEnv(defaultLoginPath string) string {
	// Optional path to login, otherwise default path of Vault is used
	path := os.Getenv(EnvVarVaultLoginPath)
	if path != "" {
		return path
	}
	return defaultLoginPath
}
