Session Store
=============

Just for fun, lege Sessions auf dem Stack ab. Nicht geeignet für Big Data, 1000 request per seconds und auf für sonst nichts :D.

```go
type Session struct {
    Token string
    UserID string
}

store := smem.NewStore()
expire := time.Now().Add(1 * time.Hour)
token, err := store.NewSession("user-ID", expire)
session, err := store.ReadSession(token)
err := store.RemoveExpire
```

Gin Beispiel für Authentication Middleware
```go

func handler(c *gin.Context) {
    c.String(200, "Useless")
}

func auth(store, fn gin.HandlerFunc) gin.HandlerFunc {
    // token from cookie
    return func(c *gin.Context) {
          session, ok := store.ReadSession(token) 
          if !ok {
            c.String(401, "No, no no!")
          }

          fn(c)
    }
}

func main() {
    store := smem.NewStore()
    h := auth(store, handler)
    r := gin.New()
    r.GET("/", h)

    r.Run()
}

```
