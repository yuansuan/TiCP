package credential

// Credential represents a credential used to access the API
type Credential struct {
	appKey    string
	appSecret string
}

func (cred *Credential) GetAppKey() string {
	return cred.appKey
}

func (cred *Credential) GetAppSecret() string {
	return cred.appSecret
}

// NewCredential creates a credential
func NewCredential(appKey, appSecret string) *Credential {
	return &Credential{
		appKey:    appKey,
		appSecret: appSecret,
	}
}
