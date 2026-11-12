package main

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

//cookieの有効時間１０分にする
const (
    googleLoginCookieName = "google_login"
    googleLoginTTLSeconds = 10 * 60
)

func main(){
signingKey := loadCookieKey("COOKIE_SIGNING_KEY")//cookieが書き換えられてないか検証
encryptionKey := loadCookieKey("COOKIE_ENCRYPTION_KEY")	 //cookieの中身を暗号化

//cookieを保存するためのストアを作成
cookieStore := sessions.NewCookieStore(signingKey, encryptionKey)
//cookie属性を設定
cookieStore.Options = &sessions.Options{
    Path:     "/auth/google",
    HttpOnly: true,
    Secure:   os.Getenv("APP_ENV") != "local",
    SameSite: http.SameSiteLaxMode,
}
cookieStore.MaxAge(googleLoginTTLSeconds)

//google認証に必要な情報をまとめる
oauthConfig := &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),	
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"openid", "email", "profile"},
	Endpoint:     google.Endpoint,
}
//URLと対応付け
mux := http.NewServeMux()
mux.HandleFunc(
	"/auth/google",
	startGoogleLogin(oauthConfig, cookieStore),
)

server := &http.Server{
	Addr:         ":8080",
	Handler:      mux,
	ReadHeaderTimeout:  5 * time.Second,
}
log.Println("listening on :8080")
log.Fatal(server.ListenAndServe())
}

//リクエストを処理する関数を作ってreturn
func startGoogleLogin(
	config *oauth2.Config,cookieStore *sessions.CookieStore,)http.HandlerFunc { //リクエストを処理する関数を作って返す
		return func(response http.ResponseWriter, request *http.Request) {
			//キャッシュを無効化するヘッダーを設定
			//Google へ移動するときに、自分のサイトの URL を Referer ヘッダーで送らないするように
			response.Header().Set("Cache-Control", "no-store")
			response.Header().Set("Referrer-Policy", "no-referrer")
			//GETメソッド以外のリクエストを拒否
			if request.Method != http.MethodGet{
				response.Header().Set("Allow", http.MethodGet)
				http.Error(response, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			
		}
		//stateとnonceを生成し、生成したものと異なるものの場合エラーを吐く。
		state, err := generateRandomToken()
		if err != nil {
			http.Error(response,"login failed",500)
			return
		}
		nonce, err := generateRandomToken()
		if err != nil {
			http.Error(response,"login failed",500)
			return
		}
		verifier:= oauth2.GenerateVerifier()
		
		//作った情報をcookieに保存する
		session := sessions.NewSession(
			cookieStore,
			googleLoginCookieName,
		) 	//cookieを毎回新しいのにするよ
		options := *cookieStore.Options
		// 共通設定をコピーして、今回用の設定を作る
		session.Options = &options
		// 今回のセッションには、そのコピーを使わせる
		session.Values["state"] = state
		session.Values["nonce"] = nonce
		session.Values["pkce_verifier"] = verifier

		if err := session.Save(request, response); err != nil {
   		 http.Error(response, "login initialization failed", 500)
   		 return
		}
			//Googleの認証ページにリダイレクトする
		authorizationURL := config.AuthCodeURL(
   			 state,
   			 oauth2.SetAuthURLParam("nonce", nonce),
 		     oauth2.S256ChallengeOption(verifier),
			)
		http.Redirect(
			response,
			request,
			authorizationURL,
			http.StatusFound,
				)
			}
		func generateRandomToken()(string,error){
			randomBytes := make([]byte, 32)
			if _ , err := rand.Read(randomBytes); err != nil {
				return "", err
			}
			return base64.URLEncoding.EncodeToString(randomBytes), nil
		}
			func requireEnv(name string) string{
				value := os.Getenv(name)
				if value == "" {
					log.Fatalf("environment variable %s is not set", name)
				}
				return value
			}
			func loadCookieKey(name string) []byte {
   				 key, err := base64.StdEncoding.DecodeString(requireEnv(name))
  				  if err != nil || len(key) != 32 {
      				  log.Fatalf("%s must be a base64-encoded 32-byte key", name)
 				   }
  				  return key
			}

		}
