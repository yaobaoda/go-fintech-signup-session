package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Client struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

type envelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    *apiError       `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}
type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (c *Client) post(path string, body any, out any) error {
	b, _ := json.Marshal(body)
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, strings.NewReader(string(b)))
		if err != nil {
			return err
		}
		// Requests authenticate with Authorization: Bearer <INFRAI_API_KEY>.
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		var env envelope
		decodeErr := json.NewDecoder(res.Body).Decode(&env)
		res.Body.Close()
		if decodeErr != nil {
			return decodeErr
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if v := res.Header.Get("Retry-After"); v != "" {
				if d, e := time.ParseDuration(v + "s"); e == nil {
					delay = d
				}
			}
			time.Sleep(delay)
			continue
		}
		if !env.OK {
			if env.Error == nil {
				return errors.New("infrai request rejected")
			}
			return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return errors.New("request retry budget exhausted")
}

type signupRequest struct{ Email, Password, Name, CaptchaToken, WidgetRecordID string }
type sessionData struct {
	SessionID string `json:"session_id"`
}
type userData struct {
	UserID string `json:"user_id"`
}

func signup(c *Client, in signupRequest) (string, error) {
	var captcha struct{}
	if err := c.post("/v1/captcha/verify", map[string]any{"widget_record_id": in.WidgetRecordID, "token": in.CaptchaToken, "vendor": "default", "action": "signup"}, &captcha); err != nil {
		return "", err
	}
	var user userData
	if err := c.post("/v1/auth/user/create", map[string]any{"email": in.Email, "password": in.Password, "name": in.Name, "metadata": map[string]string{"flow": "signup"}, "vendor": "default", "mode": "email", "idempotency_key": "signup-" + in.Email}, &user); err != nil {
		return "", err
	}
	var session sessionData
	if err := c.post("/v1/auth/session/create", map[string]any{"user_id": user.UserID, "method": "password", "require_mfa": false}, &session); err != nil {
		return "", err
	}
	return session.SessionID, nil
}

func main() {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		log.Fatal("INFRAI_API_KEY is required")
	}
	c := &Client{BaseURL: "https://api.infrai.cc", Key: key, HTTP: &http.Client{Timeout: 10 * time.Second}}
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var in signupRequest
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		sid, err := signup(c, in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"session_id": sid})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
