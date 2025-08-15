package main

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

var (
	db      *sql.DB
	privkey jwk.Key
	pubkey  jwk.Key
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	privkey, err = setupPrivate()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := setupDB(); err != nil {
		slog.Error(err.Error())
		return
	}

	setupAuth()

	router := gin.Default()

	router.GET("/jwks", func(ctx *gin.Context) {
		jwks, err := JWK()
		if err != nil {
			slog.Error(err.Error())
			return
		}

		data, err := json.MarshalIndent(jwks, "", "  ")
		if err != nil {
			slog.Error(err.Error())
			return
		}

		ctx.Data(http.StatusOK, "application/json", data)
	})

	router.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "html/auth.html")
	})

	router.GET("/auth/:provider/callback", func(ctx *gin.Context) {
		user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
		if err != nil {
			slog.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		id := uuid.NewString()
		if _, err := db.Exec("insert into users(id,name,email,avatar) values(?,?,?,?) on duplicate key update email=email",
			id, user.Name,
			user.Email,
			user.AvatarURL); err != nil {
			slog.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, err := createToken(id, user.Name, user.Email, user.AvatarURL)
		if err != nil {
			slog.Error(err.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		url := fmt.Sprintf("http://localhost:7000/auth/callback?token=%s", token)
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	})

	router.GET("/auth/:provider", func(ctx *gin.Context) {
		q := ctx.Request.URL.Query()
		q.Add("provider", ctx.Param("provider"))
		ctx.Request.URL.RawQuery = q.Encode()
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	})

	router.Run("localhost:3000")
}

func setupPrivate() (privkey jwk.Key, err error) {
	buff := os.Getenv("PRIVATE_KEY")

	rsaKey, err := decode([]byte(buff))
	if err != nil {
		return
	}

	privkey, err = jwk.Import(rsaKey)
	if err != nil {
		return
	}

	_ = privkey.Set(jwk.KeyIDKey, "main-2025")
	_ = privkey.Set(jwk.AlgorithmKey, jwa.RS256)
	_ = privkey.Set(jwk.KeyUsageKey, "sig")

	return
}

func setupAuth() {
	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback", "profile"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:3000/auth/github/callback"),
	)
}

func setupDB() (err error) {
	db, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	return nil
}

func JWK() (jwks jwk.Set, err error) {
	pub, err := jwk.PublicKeyOf(privkey)
	if err != nil {
		return nil, err
	}

	key := pub
	_ = key.Set(jwk.KeyIDKey, "main-2025")
	_ = key.Set(jwk.AlgorithmKey, "RS256")
	_ = key.Set(jwk.KeyUsageKey, "sig")

	jwks = jwk.NewSet()
	jwks.AddKey(key)

	return jwks, nil
}

func createToken(id string, username string, email string, avatar string) (signed []byte, err error) {
	token, err := jwt.NewBuilder().
		Issuer("algernon").
		IssuedAt(time.Now()).
		Expiration(time.Now().Add(24*time.Hour)).
		Audience([]string{"begadangz"}).
		Claim("user_id", id).
		Claim("username", username).
		Claim("email", email).
		Claim("avatar", avatar).Claim("kid", "main-2025").
		Build()
	if err != nil {
		return nil, err
	}

	signed, err = jwt.Sign(token, jwt.WithKey(jwa.RS256(), privkey))
	if err != nil {
		return nil, err
	}

	return signed, nil
}

func decode(pemData []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var key any
	var err error

	switch block.Type {
	case "PRIVATE KEY":
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported key type: %s", block.Type)
	}

	if err != nil {
		return nil, err
	}

	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA private key")
	}

	return rsaKey, nil
}
