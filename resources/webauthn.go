package resources

import "github.com/go-webauthn/webauthn/webauthn"

var (
	webAuthn *webauthn.WebAuthn
)

func WebAuthn() *webauthn.WebAuthn {
	return webAuthn
}

func initWebAuthIn() (err error) {
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPID:          "passvault.fun",
		RPDisplayName: "PassVault",
		RPOrigins: []string{
			"https://www.passvault.fun",
		},
		Debug: true,
	})
	return
}
