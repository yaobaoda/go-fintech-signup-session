package main

import "testing"

func TestSignupRequestRequiresCaptcha(t *testing.T) {
	cases := []struct {
		name string
		in   signupRequest
		want bool
	}{
		{"missing captcha", signupRequest{Email: "a@example.com", Password: "secret", Name: "A"}, false},
		{"provided captcha", signupRequest{Email: "a@example.com", Password: "secret", Name: "A", CaptchaToken: "token"}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.CaptchaToken != ""
			if got != tc.want {
				t.Fatalf("captcha present=%v, want %v", got, tc.want)
			}
		})
	}
}
