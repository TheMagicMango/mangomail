# MangoMail

CLI and HTTP API for sending emails using [Resend](https://resend.com).

## Installation

```bash
# Docker
docker pull ghcr.io/themagicmango/mangomail:latest

# From source
go build -o mangomail ./cmd
```

## Quick Start

### 1. Get your API Keys

**Resend API Key:**
1. Sign up at [resend.com](https://resend.com)
2. Go to **API Keys** → **Create API Key**
3. Copy the key (starts with `re_`)

**Authentication Key (for HTTP API):**
```bash
openssl rand -base64 32
```

### 2. Store your keys

```bash
mkdir -p ~/.mangomail/secrets
echo "re_your_api_key_here" > ~/.mangomail/secrets/resend_api_key
echo "your-generated-key" > ~/.mangomail/secrets/api_key
chmod 600 ~/.mangomail/secrets/*
```

### 3. CLI: Bulk Campaigns

**Using local installation:**
```bash
mangomail send my-campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Hello {{name}}!" \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

**Using Docker:**
```bash
docker run --rm \
  -v $(pwd):/app -w /app \
  -v ~/.mangomail/secrets:/secrets:ro \
  ghcr.io/themagicmango/mangomail:latest \
  send my-campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Hello {{name}}!" \
  --resend-api-key-file /secrets/resend_api_key
```

### 4. HTTP API: Server Mode

**Using local installation:**
```bash
./mangomail serve --resend-api-key-file ~/.mangomail/secrets/resend_api_key \
  --api-key-file ~/.mangomail/secrets/api_key
```

**Using Docker:**
```bash
docker run -d \
  -v ~/.mangomail/secrets:/secrets:ro \
  -e MANGOMAIL_RESEND_API_KEY_FILE=/secrets/resend_api_key \
  -e MANGOMAIL_API_KEY_FILE=/secrets/api_key \
  -p 8080:8080 \
  -p 8081:8081 \
  ghcr.io/themagicmango/mangomail:latest serve
```

**Send email:**
```bash
curl -X POST http://localhost:8081/api/v1/send \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $(cat ~/.mangomail/secrets/api_key)" \
  -d '{"from": "hello@example.com", "to": "user@example.com", "subject": "Hello!"}'
```

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MANGOMAIL_RESEND_API_KEY` | Resend API key | required |
| `MANGOMAIL_API_KEY` | HTTP API authentication key | required (serve) |
| `MANGOMAIL_HTTP_ADDRESS` | API server address | `:8081` |
| `MANGOMAIL_TELEMETRY_ADDRESS` | Health check server address | `:8080` |
| `MANGOMAIL_RESEND_RATE_LIMIT` | Max emails per second | `2` |
| `MANGOMAIL_LOG_LEVEL` | Log level | `info` |