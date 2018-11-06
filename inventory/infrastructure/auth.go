// infrastructure implementation details
package infrastructure

import (
        "errors"
        "time"
        jwt "github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
        ClientID string `json:"clientId"`
        jwt.StandardClaims
}

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
}

type AuthService struct {
        Key     []byte
        Clients map[string]string
}

const expiration = 120

// define a type for the middleware helpers
type AuthMiddleware func(AuthHandler) AuthHandler

func generateToken(signingKey []byte, clientID string) (string, error) {
        claims := CustomClaims{
                clientID,
                jwt.StandardClaims{
                        ExpiresAt: time.Now().Add(time.Second * expiration).Unix(),
                        IssuedAt:  jwt.TimeFunc().Unix(),
                },
        }
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        return token.SignedString(signingKey)
}

func (as AuthService) Auth(clientID string, clientSecret string) (string, error) {
        if as.Clients[clientID] == clientSecret {
                signed, err := generateToken(as.Key, clientID)
                if err != nil {
                        return "", errors.New(err.Error())
                }
                return signed, nil
        }
        return "", ErrAuth
}

// ErrAuth is returned when credentials are incorrect
var ErrAuth = errors.New("Incorrect credentials")
