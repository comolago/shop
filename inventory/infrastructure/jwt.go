// infrastructure implementation details
package infrastructure

import (
   "errors"
   "time"
   jwt "github.com/dgrijalva/jwt-go"
   "github.com/go-kit/kit/endpoint"
   gokitjwt "github.com/go-kit/kit/auth/jwt"
   "github.com/comolago/shop/inventory/domain"
   "fmt"
)

type CustomClaims struct {
   ClientID string `json:"clientId"`
   jwt.StandardClaims
}

type AuthService struct {
   Key     []byte
   AuthDb domain.DbHandler
}

const expiration = 120

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

func (as AuthService)SetSecret(secret []byte){
   as.Key = secret
}

func (as AuthService)GetSecret() []byte{
   return as.Key
}

func (as AuthService) Auth(clientID string, clientSecret string) (string, error) {
   //if as.Clients[clientID] == clientSecret {
   //var err *error
   var id int
   id , _ = as.AuthDb.AuthenticateUser(clientID,clientSecret)
   if id > 0 {
      fmt.Println(id)
      //signed, err := generateToken(as.Key, clientID)
      signed, err := generateToken(as.Key, clientID)
      if err != nil {
         return "", errors.New(err.Error())
      }
      return signed, nil
   }
   return "", ErrAuth
}

func MakeSecureEndpoint(endpoint endpoint.Endpoint,auth domain.AuthHandler) endpoint.Endpoint{
   key := auth.GetSecret()
   keys := func(token *jwt.Token) (interface{}, error) {
      return key, nil
   }
   return gokitjwt.NewParser(keys, jwt.SigningMethodHS256, func() jwt.Claims { return &CustomClaims{} })(endpoint)
}

// ErrAuth is returned when credentials are incorrect
var ErrAuth = errors.New("Incorrect credentials")
