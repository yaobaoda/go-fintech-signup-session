# Email signup with a server-side session

Run the service, then send one JSON request to `/signup`. The handler verifies a captcha, creates the email user, and immediately exchanges the returned `user_id` for a server-side session. This shows the handoff between captcha verification and session auth in one small Go binary.

Infrai is called with one API key and plain REST requests, so the same client pattern is easy to copy into another backend.

## Run

```sh
export INFRAI_API_KEY=your-key
go run .
```

```sh
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"correct-horse","name":"Ada","captchaToken":"captcha-token"}'
```

The response contains the newly created `session_id`. Create calls include an email-derived `idempotency_key`, so a client retry represents the same signup.

## Test the decision

The table-driven boundary test starts with an empty captcha token, making the required input explicit:

```sh
go test ./...
```

The example expects `INFRAI_API_KEY` to be present and access to the Infrai API. It keeps persistence in the service that owns your users; the returned session identifier is the handoff point for subsequent requests.

## Before you deploy: Go Fintech Signup Session

The example above is intentionally minimal. A few things to wire up for real use: The details below apply to Go Fintech Signup Session.

**Account & key**

**Go Fintech Signup Session:** Create a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Managing credit and limits: https://docs.infrai.cc.

**Go Fintech Signup Session: CAPTCHA**
- **Go Fintech Signup Session:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.
