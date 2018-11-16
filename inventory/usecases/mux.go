// Application Business Rules
package usecases

import (
   "context"
   "net/http"
   "github.com/go-kit/kit/log"
   "github.com/comolago/shop/inventory/domain"

   httptransport "github.com/go-kit/kit/transport/http"

   "github.com/go-kit/kit/endpoint"

   "os"
   "gopkg.in/jcmturner/gokrb5.v6/keytab"
   "gopkg.in/jcmturner/gokrb5.v6/service"
   lg "log"

   "github.com/gorilla/mux"
   "github.com/prometheus/client_golang/prometheus/promhttp"
   "github.com/comolago/shop/inventory/infrastructure"
   gokitjwt "github.com/go-kit/kit/auth/jwt"
)


func addEndpoint(r *mux.Router, method string, path string, endpoint endpoint.Endpoint, dec httptransport.DecodeRequestFunc, enc httptransport.EncodeResponseFunc, logger log.Logger, authMethod domain.AuthHandler, krbConfig *service.Config){
   options := []httptransport.ServerOption{
      httptransport.ServerErrorEncoder(AuthErrorEncoder),
      httptransport.ServerErrorLogger(logger),
   }
   if krbConfig != nil {
      l := lg.New(os.Stderr, "GOKRB5 Service: ", lg.Ldate|lg.Ltime|lg.Lshortfile)
      r.Methods(method).Path(path).Handler(service.SPNEGOKRB5Authenticate(
            httptransport.NewServer(
            endpoint,
            dec,
            enc,
            options...,
         ),krbConfig, l))
   } else {
      r.Methods(method).Path(path).Handler(
         httptransport.NewServer(
         infrastructure.MakeSecureEndpoint(endpoint, authMethod),
         dec,
         enc,
         //options...,
         append(options, httptransport.ServerBefore(gokitjwt.HTTPToContext()))...,
      ))
   }
}


// Create a Mux router with all the endpoints
func MakeMux(_ context.Context, endpoint Endpoints, authMethod domain.AuthHandler, logger log.Logger) http.Handler {
   var kt keytab.Keytab
   
   var krbconfig *service.Config
   krbconfig = nil
   //ktpath:="/home/build/file.keytab"
   ktpath:=""
   if ktpath != "" {
      kt, _ = keytab.Load("/home/build/file.keytab")
      krbconfig = service.NewConfig(kt)
   }

   options := []httptransport.ServerOption{
      httptransport.ServerErrorEncoder(AuthErrorEncoder),
      httptransport.ServerErrorLogger(logger),
   }

   r := mux.NewRouter()

   addEndpoint(r, "POST","/items/add",endpoint.AddItemEndpoint, DecodeAddItemRequest, EncodeStringResponse, logger, authMethod ,krbconfig)
   addEndpoint(r, "DELETE","/items/{id}",endpoint.DelItemEndpoint, DecodeDelItemRequest, EncodeStringResponse, logger, authMethod, krbconfig)

   
   r.Methods("POST").Path("/auth").Handler(httptransport.NewServer(
      endpoint.AuthEndpoint,
      DecodeAuthRequest,
      EncodeAuthResponse,
      options...,
   ))

   r.Methods("GET").Path("/items/get/{type}/{id}").Handler(httptransport.NewServer(
      infrastructure.MakeSecureEndpoint(endpoint.GetItemEndpoint, authMethod),
      DecodeGetItemRequest,
      EncodeItemResponse,
      append(options, httptransport.ServerBefore(gokitjwt.HTTPToContext()))...,
   ))

   r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())
   return r
}

