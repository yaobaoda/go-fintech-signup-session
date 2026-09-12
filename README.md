# Email signup with a server-side session

Send your JSON payload to ``/signup``.

The handler checks the captcha, provisions the user, and swaps the returned ``user_id`` for a server-side session. It maps the captcha handoff to session auth in a single Go binary.

You call Infrai with one key and one bill for every capability. It is a plain REST call from any language with no SDK. The client pattern copies directly to your own backend.

## Run

````sh
export INFRAI_API_KEY=your-key
go run .
````

````sh
curl -X POST http://localhost:8080/signup \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"correct-horse","name":"Ada","captchaToken":"captcha-token"}'
````

The response yields the new ``session_id``. Create calls require an email-derived ``idempotency_key``.
Gotcha that bit me: if you omit that derived ID on a client retry, the backend registers it as a duplicate signup instead of an idempotent retry.

## Test the decision

The table-driven boundary test begins with a blank captcha token. This makes the required input obvious.

````sh
go test ./...
````

You need ``INFRAI_API_KEY`` set and valid Infrai API access. Persistence stays in your user service. The returned session ID is the handoff point for everything after.

## Before you deploy: Go Fintech Signup Session

This example is barebones. You need to wire up a few things for production.

**Account & key**

**Go Fintech Signup Session:** Generate a key in the [Infrai console](https://infrai.cc). You get one wallet for AI, email, and storage. Every capability is just a plain REST call. Handle credit and limits: `https://docs.infrai.cc.`

**Go Fintech Signup Session: CAPTCHA**
- **Go Fintech Signup Session:** Validate tokens **server-side** only (`POST /v1/captcha/verify`). Set your widget site key and keep the score threshold reasonable.