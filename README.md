# tapresp
A [golang](https://go.dev/) ResponseWriter that lets you tap into the underlying HTTP response 

Package `tapresp` provides a wrapper for the [http.ResponseWriter](https://pkg.go.dev/net/http#ResponseWriter) interface that taps into the response information. It is somewhat similar to the [httptest.ResponseRecorder](https://pkg.go.dev/net/http/httptest#ResponseRecorder) in that it records its mutations and provides access to HTTP response properties. Where they differ is this is meant to wrap an existing `http.ResponseWriter` and "tap" into it.

## Install

```bash
go get -u github.com/justsimmons/tapresp
```

## Examples

Functions exactly the same as a standard `http.ResponseWriter`:
```golang
func YourHandler(w http.ResponseWriter, r *http.Request) {
    // Wrap the old ResponseWriter
    trw := tapresp.New(w)
    
    // Do the stuff....
    
    trw.WriteHeader(http.StatusOK)
    trw.Header().Set("Content-Type", "application/json")
    trw.Write([]byte("Success!))
}
```

Provides access to the response body:
```golang
func yourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       // Wrap the old ResponseWriter
       trw := tapresp.New(w)
       
       // Call next handler in the chain
       next.ServeHTTP(w, r)
       
       // Note this is a copy, any modification will not effect the original
       trw.Body() // []byte
    })
}
```

Provides access to the status code:
```golang
func yourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       // Wrap the old ResponseWriter
       trw := tapresp.New(w)
       
       // Call next handler in the chain
       next.ServeHTTP(w, r)
       
       statusCode := trw.StatusCode()
       
       if statusCode != http.StatusOK {
          // Do something special....
       }
    })
}
```


Can be used to log response values in middleware:
```golang
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       // Wrap the old ResponseWriter
       trw := tapresp.New(w)
       
       // Call next handler in the chain
       next.ServeHTTP(w, r)
       
       // Log desired response info
       log.Println("Status: ", trw.StatusCode())
       log.Println("Headers: ", trw.Header())
       log.Printf("Body: %s\n", trw.Body())
    })
}
```

Provides an `http.Response` struct that contains all of the response information in the `http.ResponseWriter`. This is a somwhat uncommon use case, but recently when trying to comply with the [HTTP Message Signing Standard](https://datatracker.ietf.org/doc/draft-ietf-httpbis-message-signatures/), I ran into a need to have access to the underlying `http.Response` not just the header values. 

```golang
// Signs HTTP Response according to IETF standard: https://datatracker.ietf.org/doc/draft-ietf-httpbis-message-signatures/
func Sign(resp *http.Response) http.Header {
  signedHeaders := make(http.Header, 2)
  
  // Do the signing logic (Adds "Signature" and "Signature-Input" headers)
  
  return signedHeaders
}

func YourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       // Wrap the old ResponseWriter
       trw := tapresp.New(w)
       
       // Call next handler in the chain
       next.ServeHTTP(w, r)
       
       // Provides an http.Response struct
       if err := Sign(trw.Response()); err != nil {
           trw.WriteHeader(http.StatusBadRequest)
           trw.Write([]byte("unable to sign HTTP response"))
           return 
       }
    })
}
```

## License

See the LICENSE file for details.
