package main

import (
	"context"
	"crypto/subtle"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type GoogleIdentity struct {
	Subject       string `json:"google_sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func newGoogleCallbackHandler(
	ctx context.Context,
	config *oauth2.Config,
	cookieStore *sessions.CookieStore,
	onVerified func(http.ResponseWriter, *http.Request, GoogleIdentity),
) (http.HandlerFunc, error) {
	if config == nil {
		return nil, errors.New("Google OAuth configuration is required")
	}

	if config.ClientID == "" {
		return nil, errors.New("Google OAuth ClientID is required")
	}

	if config.ClientSecret == "" {
		return nil, errors.New("Google OAuth ClientSecret is required")
	}

	if config.RedirectURL == "" {
		return nil, errors.New("Google OAuth RedirectURL is required")
	}

	if cookieStore == nil {
		return nil, errors.New("Google login cookie store is required")
	}

	if onVerified == nil {
		return nil, errors.New("Google login success callback is required")
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	providerContext := oidc.ClientContext(ctx, httpClient)

	provider, err := oidc.NewProvider(
		providerContext,
		"https://accounts.google.com",
	)
	if err != nil {
		return nil, errors.New("google OIDC discovery failed")
	}

	verifier := provider.VerifierContext(providerContext, &oidc.Config{
		ClientID:             config.ClientID,
		SupportedSigningAlgs: []string{"RS256"},
	})

	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Cache-Control", "no-store")
		response.Header().Set("Referrer-Policy", "no-referrer")

		// OAuth errors can contain credentials, so raw errors and query strings
		// are deliberately excluded from application logs.
		reject := func(status int, stage string) {
			slog.WarnContext(request.Context(), "google callback rejected", "stage", stage)
			http.Error(response, "Google authentication failed. Please start again.", status)
		}

		if request.Method != http.MethodGet {
			response.Header().Set("Allow", http.MethodGet)
			reject(http.StatusMethodNotAllowed, "method")
			return
		}

		parameters, err := url.ParseQuery(request.URL.RawQuery)
		if err != nil {
			reject(http.StatusBadRequest, "query")
			return
		}

		session, err := cookieStore.Get(request, "google_login")
		if err != nil || session.IsNew {
			reject(http.StatusBadRequest, "login_cookie")
			return
		}

		expectedState, hasState := session.Values["state"].(string)
		expectedNonce, hasNonce := session.Values["nonce"].(string)
		pkceVerifier, hasVerifier := session.Values["pkce_verifier"].(string)
		if !hasState || !hasNonce || !hasVerifier || expectedState == "" ||
			expectedNonce == "" || pkceVerifier == "" {
			reject(http.StatusBadRequest, "login_cookie_values")
			return
		}

		states := parameters["state"]
		if len(states) != 1 ||
			subtle.ConstantTimeCompare([]byte(states[0]), []byte(expectedState)) != 1 {
			reject(http.StatusBadRequest, "state")
			return
		}

		// A callback with an unrelated state must not clear an active attempt.
		// Only delete the cookie after the browser-bound state has matched.
		options := *session.Options
		options.MaxAge = -1
		session.Options = &options
		if err := session.Save(request, response); err != nil {
			reject(http.StatusInternalServerError, "cookie_cleanup")
			return
		}

		if oauthErrors, exists := parameters["error"]; exists {
			if len(oauthErrors) != 1 || oauthErrors[0] == "" || parameters.Has("code") {
				reject(http.StatusBadRequest, "oauth_error_parameters")
				return
			}
			if oauthErrors[0] == "access_denied" {
				reject(http.StatusForbidden, "cancelled")
				return
			}
			reject(http.StatusBadRequest, "oauth_error")
			return
		}

		codes := parameters["code"]
		if len(codes) != 1 || codes[0] == "" {
			reject(http.StatusBadRequest, "code")
			return
		}

		exchangeContext := oidc.ClientContext(request.Context(), httpClient)
		verificationContext, cancel := context.WithTimeout(exchangeContext, 10*time.Second)
		defer cancel()

		token, err := config.Exchange(
			verificationContext,
			codes[0],
			oauth2.VerifierOption(pkceVerifier),
		)
		if err != nil {
			reject(http.StatusBadGateway, "code_exchange")
			return
		}

		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok || rawIDToken == "" {
			reject(http.StatusBadGateway, "missing_id_token")
			return
		}

		idToken, err := verifier.Verify(verificationContext, rawIDToken)
		if err != nil {
			reject(http.StatusUnauthorized, "id_token")
			return
		}

		// go-oidc verifies signature, issuer, audience and expiry, but the
		// application must compare the nonce itself.
		if subtle.ConstantTimeCompare([]byte(idToken.Nonce), []byte(expectedNonce)) != 1 {
			reject(http.StatusUnauthorized, "nonce")
			return
		}

		if idToken.Subject == "" {
			reject(http.StatusUnauthorized, "subject")
			return
		}

		var claims struct {
			Email           string `json:"email"`
			EmailVerified   bool   `json:"email_verified"`
			Name            string `json:"name"`
			Picture         string `json:"picture"`
			AuthorizedParty string `json:"azp"`
		}
		if err := idToken.Claims(&claims); err != nil {
			reject(http.StatusUnauthorized, "claims")
			return
		}

		if (claims.AuthorizedParty != "" && claims.AuthorizedParty != config.ClientID) ||
			(len(idToken.Audience) > 1 && claims.AuthorizedParty == "") {
			reject(http.StatusUnauthorized, "authorized_party")
			return
		}

		if idToken.AccessTokenHash != "" {
			if err := idToken.VerifyAccessToken(token.AccessToken); err != nil {
				reject(http.StatusUnauthorized, "access_token_hash")
				return
			}
		}

		onVerified(response, request, GoogleIdentity{
			Subject:       idToken.Subject,
			Email:         claims.Email,
			EmailVerified: claims.EmailVerified,
			Name:          claims.Name,
			Picture:       claims.Picture,
		})
	}, nil
}
