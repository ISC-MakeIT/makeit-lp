package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DB_USER")
	cfg.Passwd = os.Getenv("DB_PASSWORD")
	cfg.Net = "tcp"
	cfg.Addr = os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT")
	cfg.DBName = os.Getenv("DB_NAME")
	cfg.ParseTime = true
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") // dev only
		status := map[string]string{"api": "ok", "db": "ok"}
		if err := db.PingContext(r.Context()); err != nil {
			status["db"] = err.Error()
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		json.NewEncoder(w).Encode(status)
	})

	// Cookieの署名・暗号化用の鍵を読み込む。
	readKey := func(name string) []byte {
		key, err := base64.StdEncoding.DecodeString(os.Getenv(name))
		if err != nil || len(key) != 32 {
			log.Fatalf("%s must be a base64-encoded 32-byte key", name)
		}
		return key
	}

	cookieStore := sessions.NewCookieStore(
		readKey("SESSION_AUTH_KEY"),
		readKey("SESSION_ENCRYPT_KEY"),
	)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600, // ログイン操作の有効期間：10分
		HttpOnly: true,
		Secure:   true, // localhostのHTTPで開発するときだけfalse
		SameSite: http.SameSiteLaxMode,
	}
	cookieStore.MaxAge(600)

	googleConfig := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "email", "profile"},
	}

	callbackHandler, err := newGoogleCallbackHandler(
		context.Background(),
		googleConfig,
		cookieStore,
		func(w http.ResponseWriter, r *http.Request, identity GoogleIdentity) {
			// 検証成功後に呼ばれる。
			// ここでidentity.Subjectを使ったユーザーの検索・作成と、
			// アプリ用ログインセッションの発行を実装する。
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Write([]byte("Google認証の検証に成功しました"))
		},
	)
	if err != nil {
		log.Fatalf("Google authentication initialization failed: %v", err)
	}

	http.HandleFunc("/api/auth/google/callback", callbackHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
