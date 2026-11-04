package testapp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
)

func (a *TestApp) Request(
	t *testing.T,
	request any,
	path string,
	cookies []*http.Cookie,
) *http.Response {
	t.Helper()

	body, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshaling request: %v", err)
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodPost,
		a.address+path,
		bytes.NewReader(body),
	)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	if cookies != nil {
		for _, cookie := range cookies {
			req.AddCookie(cookie)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := a.http.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	return resp
}

func (a *TestApp) Get(
	t *testing.T,
	path string,
	query url.Values,
	cookies []*http.Cookie,
) *http.Response {
	t.Helper()

	u, err := url.Parse(a.address + path)
	if err != nil {
		t.Fatalf("parsing URL: %v", err)
	}

	if query != nil {
		u.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(
		t.Context(),
		http.MethodGet,
		u.String(),
		nil,
	)
	if err != nil {
		t.Fatalf("creating request: %v", err)
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := a.http.Do(req)
	if err != nil {
		t.Fatalf("sending request: %v", err)
	}

	return resp
}

func (a *TestApp) GetSessionId(t *testing.T, response *http.Response) *uuid.UUID {
	t.Helper()

	for _, cookie := range response.Cookies() {
		if cookie.Name == a.Config.Application.SessionCookieName {
			parsed, err := uuid.Parse(cookie.Value)
			if err != nil {
				t.Fatalf("parsing session cookie: %v", err)
			}

			return &parsed
		}
	}

	return nil
}
