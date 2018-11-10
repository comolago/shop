package domain

type AuthRequest struct {
  ClientID     string `json:"clientId"`
  ClientSecret string `json:"clientSecret"`
}

type AuthResponse struct {
  Token string `json:"token,omitempty"`
  Err   string `json:"err,omitempty"`
}

// AuthService provides authentication service
type AuthHandler interface {
  Auth(string, string) (string, error)
  GetSecret() []byte
  SetSecret(secret []byte)
}
