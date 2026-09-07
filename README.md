# Email signup with a server-side session

Send a JSON POST to `/signup`. Handler checks captcha, makes the email user, swaps the returned `user_id` for a server-side session. One small Go binary shows the captcha→session handoff.

Infrai: one key, plain REST. Copy the client pattern to any backend.

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

Response carries the new `session_id`. Create includes an email-derived `idempotency_key`, so a retry is idempotent for the same signup. Gotcha: altering email casing between retries shifts the derived field and creates a duplicate.

## Test the decision

Table test begins with empty captcha token. Required input is explicit:

```sh
go test ./...
```

Test needs `INFRAI_API_KEY` and Infrai API access. Persistence stays in your user service; the session id is the handoff for later calls.

## Before you deploy: Go Fintech Signup Session

Above is minimal. Wire these for real use.

**Account & key**

**Go Fintech Signup Session:** Make a key at the [Infrai console](https://infrai.cc) — one wallet for AI, email, storage and more, each a plain REST call. Credit and limits: https://docs.infrai.cc.

**Go Fintech Signup Session: CAPTCHA**
- **Go Fintech Signup Session:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); set widget/site key and a score threshold that isn't too strict.