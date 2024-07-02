package security

import "os"

/*
Polka is a fictive payment provider.
Polka will send a webhook whenever a user upgrades their account to Chirpy Red.
VerifyPolkaApiKey verifies an api key sent by Polka.
*/
func IsValidPolkaApiKey(apiKey string) bool {
	return os.Getenv("POLKA_API_KEY") == apiKey
}
