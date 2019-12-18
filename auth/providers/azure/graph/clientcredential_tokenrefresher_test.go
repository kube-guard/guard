package graph

import (
	"fmt"
	"net/http"
	"testing"
)

func TestClientCredentialTokenRefresher(t *testing.T) {
	const (
		inputAccessToken    = "inputAccessToken"
		oboAccessToken      = "oboAccessToken"
		clientID            = "fakeID"
		clientSecret        = "fakeSecret"
		scope               = "https://graph.microsoft.com/.default"
		oboResponse         = `{"token_type":"Bearer","expires_in":3599,"access_token":"%s"}`
		expectedContentType = "application/x-www-form-urlencoded"
		expectedGrantType   = "client_credentials"
		expectedTokneType   = "Bearer"
	)

	t.Run("Upon Success Response", func(t *testing.T) {
		s := startTestServer(t, func(rw http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodPost {
				t.Errorf("expected http method %s, actual %s", http.MethodPost, req.Method)
			}
			if req.Header.Get("Content-Type") != expectedContentType {
				t.Errorf("expected content type: %s, actual: %s", expectedContentType, req.Header.Get("Content-Type"))
			}
			if req.FormValue("client_id") != clientID {
				t.Errorf("expected client_id: %s, actual: %s", clientID, req.FormValue("client_id"))
			}
			if req.FormValue("client_secret") != clientSecret {
				t.Errorf("expected client_secret: %s, actual: %s", clientSecret, req.FormValue("client_secret"))
			}
			if req.FormValue("scope") != scope {
				t.Errorf("expected scope: %s, actual: %s", scope, req.FormValue("scope"))
			}
			if req.FormValue("grant_type") != expectedGrantType {
				t.Errorf("expected grant_type: %s, actual: %s", expectedGrantType, req.FormValue("grant_type"))
			}
			_, _ = rw.Write([]byte(fmt.Sprintf(oboResponse, oboAccessToken)))
		})

		defer stopTestServer(t, s)

		r := NewClientCredentialTokenRefresher(clientID, clientSecret, s.URL, scope)
		resp, err := r.Refresh(inputAccessToken)
		if err != nil {
			t.Fatalf("refresh should not return error: %s", err)
		}

		if resp.Token != oboAccessToken {
			t.Errorf("returned obo token '%s' doesn't match expected '%s'", resp.Token, oboAccessToken)
		}
		if resp.TokenType != expectedTokneType {
			t.Errorf("expected token type: Bearer, actual: %s", resp.TokenType)
		}
	})

	t.Run("Upon Error Response", func(t *testing.T) {
		s := startTestServer(t, func(rw http.ResponseWriter, req *http.Request) {
			if req.Method != http.MethodPost {
				t.Errorf("expected http method %s, actual %s", http.MethodPost, req.Method)
			}
			if req.Header.Get("Content-Type") != expectedContentType {
				t.Errorf("expected content type: %s, actual: %s", expectedContentType, req.Header.Get("Content-Type"))
			}
			if req.FormValue("client_id") != clientID {
				t.Errorf("expected client_id: %s, actual: %s", clientID, req.FormValue("client_id"))
			}
			if req.FormValue("client_secret") != clientSecret {
				t.Errorf("expected client_secret: %s, actual: %s", clientSecret, req.FormValue("client_secret"))
			}
			if req.FormValue("scope") != scope {
				t.Errorf("expected scope: %s, actual: %s", scope, req.FormValue("scope"))
			}
			if req.FormValue("grant_type") != expectedGrantType {
				t.Errorf("expected grant_type: %s, actual: %s", expectedGrantType, req.FormValue("grant_type"))
			}

			rw.WriteHeader(http.StatusBadRequest)
			_, _ = rw.Write([]byte(`{"error":{"code":"Authorization_RequestDenied","message":"Insufficient privileges to complete the operation.","innerError":{"request-id":"6e73da70-96f3-4415-8c6a-a940cb1ba0e2","date":"2019-12-17T21:57:17"}}}`))
		})

		defer stopTestServer(t, s)

		r := NewClientCredentialTokenRefresher(clientID, clientSecret, s.URL, scope)
		resp, err := r.Refresh(inputAccessToken)
		if err == nil {
			t.Error("refresh should return error")
		}

		if resp.Token != "" {
			t.Errorf("returned obo token '%s' should be empty", resp.Token)
		}
		if resp.TokenType != "" {
			t.Errorf("expected token type: %s should be empty", resp.TokenType)
		}
	})
}
