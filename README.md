# Journie

Journie is an AI-powered Journaling Genie on Telegram. Powered by Google's Gemini GenAi

## Clone the project

```
$ git clone https://github.com/watahak/journie.git

```

## Development

Create `.env` from `.env.example`.

- `TELEGRAM_TOKEN`: Token of telegram bot to interface with [How to create a new bot](https://core.telegram.org/bots/tutorial)
- `GEMINI_API_KEY`: Key from Gemini API [Creating Gemini Key](https://aistudio.google.com/app/apikey)
- `FIREBASE_CREDENTIALS`: Firebase Credentials in JSON string [Firebase Credentials Instructions](https://firebase.google.com/docs/admin/setup)
- `FIREBASE_PROJECT_ID`: Firebase Project ID (retrieve from firebase console)

```
$ go mod tidy
$ go run ./cmd/journie/main.go
```

## Testing

We are using Go's built in unit testing. Refer to docs on [how to add tests](https://go.dev/doc/tutorial/add-a-test)

To run tests in project

```
$ go test ./pkg/... -v
```
