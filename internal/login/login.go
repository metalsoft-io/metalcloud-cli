// Package login signs the CLI in through the browser: the server starts a request, the user approves
// it in the UI, and the CLI collects a session token. It talks plain HTTP, not the SDK.
package login

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var (
	ErrNotSupported = errors.New("This MetalCloud doesn't support signing in from the CLI; use an API token.")
	ErrDenied       = errors.New("The sign-in was denied.")
	ErrExpired      = errors.New("The sign-in request has expired; run `metalcloud-cli login` again.")
)

// slowDownStep is added to the polling interval each time the server answers 429.
var slowDownStep = 5 * time.Second

// openBrowser opens a URL in the user's browser. It is a variable so tests can replace it.
var openBrowser = func(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

type Client struct {
	Endpoint string
	Insecure bool
	HTTP     *http.Client
}

type Started struct {
	DeviceCode              string
	UserCode                string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresIn               time.Duration
	Interval                time.Duration
}

type Approved struct {
	Token            string
	ExpiresTimestamp string
	UserID           int64
	UserEmail        string
}

func (c Client) httpClient() *http.Client {
	if c.HTTP != nil {
		return c.HTTP
	}
	client := &http.Client{Timeout: 30 * time.Second}
	if c.Insecure {
		client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	return client
}

// post sends a JSON POST and returns the status and body. token, when not empty, goes in as Bearer.
func (c Client) post(ctx context.Context, path string, body interface{}, token string) (int, []byte, error) {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.Endpoint, "/")+path, reader)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "metalcloud-cli")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.httpClient().Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}
	return resp.StatusCode, data, nil
}

func statusError(what string, status int, body []byte) error {
	var msg struct {
		Message interface{} `json:"message"`
	}
	if json.Unmarshal(body, &msg) == nil && msg.Message != nil {
		return fmt.Errorf("%s failed: %d - %v", what, status, msg.Message)
	}
	return fmt.Errorf("%s failed: %d %s", what, status, http.StatusText(status))
}

// Start asks the server for a new sign-in request.
func Start(ctx context.Context, c Client, clientName, version string) (*Started, error) {
	status, body, err := c.post(ctx, "/api/v2/cli-login", map[string]string{
		"clientName":    clientName,
		"clientVersion": version,
	}, "")
	if err != nil {
		return nil, err
	}
	// An older server answers the unauthenticated call with 401 (its auth guard runs before the
	// fallback route) or 404.
	if status == http.StatusUnauthorized || status == http.StatusNotFound {
		return nil, ErrNotSupported
	}
	if status != http.StatusOK && status != http.StatusCreated {
		return nil, statusError("starting the sign-in", status, body)
	}

	var r struct {
		DeviceCode              string `json:"deviceCode"`
		UserCode                string `json:"userCode"`
		VerificationURI         string `json:"verificationUri"`
		VerificationURIComplete string `json:"verificationUriComplete"`
		ExpiresIn               int    `json:"expiresIn"`
		Interval                int    `json:"interval"`
	}
	if err := json.Unmarshal(body, &r); err != nil || r.DeviceCode == "" {
		return nil, fmt.Errorf("unexpected answer from the server when starting the sign-in")
	}

	return &Started{
		DeviceCode:              r.DeviceCode,
		UserCode:                r.UserCode,
		VerificationURI:         r.VerificationURI,
		VerificationURIComplete: r.VerificationURIComplete,
		ExpiresIn:               time.Duration(r.ExpiresIn) * time.Second,
		Interval:                time.Duration(r.Interval) * time.Second,
	}, nil
}

// Wait polls until the user approves or denies the request, it expires, or ctx ends.
func Wait(ctx context.Context, c Client, s *Started) (*Approved, error) {
	interval := s.Interval

	for {
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}

		status, body, err := c.post(ctx, "/api/v2/cli-login/token", map[string]string{"deviceCode": s.DeviceCode}, "")
		if err != nil {
			return nil, err
		}

		switch status {
		case http.StatusOK:
			var r struct {
				Status           string `json:"status"`
				Token            string `json:"token"`
				ExpiresTimestamp string `json:"expiresTimestamp"`
				User             struct {
					ID    int64  `json:"id"`
					Email string `json:"email"`
				} `json:"user"`
			}
			if err := json.Unmarshal(body, &r); err != nil {
				return nil, fmt.Errorf("unexpected answer from the server while waiting for the sign-in")
			}
			switch {
			case r.Status == "pending":
			case r.Status == "approved" && r.Token != "":
				return &Approved{Token: r.Token, ExpiresTimestamp: r.ExpiresTimestamp, UserID: r.User.ID, UserEmail: r.User.Email}, nil
			default:
				return nil, fmt.Errorf("unexpected sign-in status %q", r.Status)
			}
		case http.StatusTooManyRequests:
			interval += slowDownStep
		case http.StatusForbidden:
			return nil, ErrDenied
		case http.StatusGone:
			return nil, ErrExpired
		default:
			return nil, statusError("waiting for the sign-in", status, body)
		}
	}
}

// Logout ends the session on the server.
func Logout(ctx context.Context, c Client, token string) error {
	status, body, err := c.post(ctx, "/api/v2/logout", nil, token)
	if err != nil {
		return err
	}
	if status < 200 || status > 299 {
		return statusError("logout", status, body)
	}
	return nil
}

// formatUserCode shows the code as XXXX-XXXX.
func formatUserCode(code string) string {
	plain := strings.ToUpper(strings.ReplaceAll(code, "-", ""))
	if len(plain) == 8 {
		return plain[:4] + "-" + plain[4:]
	}
	return code
}

// formatExpiry shows an ISO timestamp in UTC; anything unparsable is shown as received.
func formatExpiry(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.UTC().Format("2006-01-02 15:04 UTC")
}

// Login runs the whole sign-in: start, show the URL and code, open the browser unless noBrowser,
// wait for the approval, and save the session for the endpoint.
func Login(ctx context.Context, c Client, clientName, version string, noBrowser bool, out io.Writer) error {
	s, err := Start(ctx, c, clientName, version)
	if err != nil {
		return err
	}

	url := s.VerificationURIComplete
	if url == "" {
		url = s.VerificationURI
	}

	fmt.Fprintf(out, "To sign in, open this URL in a browser on any device:\n  %s\n", url)
	fmt.Fprintf(out, "Check that the page shows this code: %s\n", formatUserCode(s.UserCode))

	if !noBrowser {
		if err := openBrowser(url); err != nil {
			fmt.Fprintf(out, "Could not open the browser (%v); open the URL yourself.\n", err)
		} else {
			fmt.Fprintln(out, "Opened it in your browser.")
		}
	}

	fmt.Fprintf(out, "Waiting for approval (the code is valid for %s)...\n", s.ExpiresIn)

	approved, err := Wait(ctx, c, s)
	if err != nil {
		return err
	}

	err = SaveSession(c.Endpoint, Session{
		Token:            approved.Token,
		ExpiresTimestamp: approved.ExpiresTimestamp,
		UserEmail:        approved.UserEmail,
		CreatedTimestamp: time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return fmt.Errorf("signed in, but the session could not be saved: %w", err)
	}

	fmt.Fprintf(out, "Signed in as %s. The session is valid until %s.\n", approved.UserEmail, formatExpiry(approved.ExpiresTimestamp))
	return nil
}
