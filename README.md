# MangoMail

MangoMail is a CLI tool for sending bulk emails using HTML templates and CSV data with [Resend](https://resend.com). Features built-in rate limiting to comply with API restrictions.

## Installation

### Option 1: Docker

```bash
docker pull ghcr.io/themagicmango/mangomail:latest
```

### Option 2: Build from Source

```bash
git clone https://github.com/TheMagicMango/mangomail.git
cd mangomail
go build -o mangomail cmd/main.go
```

## Quick Start

### Using Local Installation

```bash
# 1. Store your Resend API key
mkdir -p ~/.mangomail/secrets
echo "re_your_api_key_here" > ~/.mangomail/secrets/resend_api_key

# 2. Send your campaign
mangomail my-campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Hello {{name}}!" \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

### Using Docker

```bash
docker run --rm \
  -v $(pwd):/app -w /app \
  ghcr.io/themagicmango/mangomail:latest \
  my-campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Hello {{name}}!" \
  --resend-api-key-file /app/.secrets/resend_api_key
```

## Configuration

### API Key (Required)

You can provide your Resend API key in three ways:

#### Option 1: File (Recommended)

```bash
# Create secrets directory
mkdir -p ~/.mangomail/secrets

# Store API key
echo "re_your_api_key_here" > ~/.mangomail/secrets/resend_api_key

# Use with --resend-api-key-file flag
mangomail campaign --resend-api-key-file ~/.mangomail/secrets/resend_api_key ...

# Or set environment variable
export MANGOMAIL_RESEND_API_KEY_FILE="$HOME/.mangomail/secrets/resend_api_key"
```

#### Option 2: Environment Variable

```bash
export MANGOMAIL_RESEND_API_KEY="re_your_api_key_here"
```

#### Option 3: Flag (Not recommended for production)

```bash
mangomail campaign --resend-api-key "re_your_api_key_here" ...
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MANGOMAIL_RESEND_API_KEY` | Resend API key | - |
| `MANGOMAIL_RESEND_API_KEY_FILE` | Path to file containing API key | - |
| `MANGOMAIL_RATE_LIMIT` | Max emails per second | `2` |
| `MANGOMAIL_LOG_LEVEL` | Log level (debug, info, warn, error) | `info` |

## Usage

### Basic Command

```bash
mangomail <campaign-name> \
  --html <path-to-html> \
  --sample <path-to-csv> \
  --from <sender-email> \
  --subject <subject-line>
```

### CSV File Format

Your CSV must have an `email` column. Additional columns become template placeholders.

**Example** (`contacts.csv`):

```csv
email,name,company
alice@example.com,Alice Johnson,TechCorp
bob@example.com,Bob Smith,StartupXYZ
```

### HTML Templates

Use `{{placeholder}}` syntax to insert CSV values.

**Example** (`template.html`):

```html
<!DOCTYPE html>
<html>
<body>
  <h1>Hello {{name}}!</h1>
  <p>Thank you for being part of {{company}}.</p>
</body>
</html>
```

**Result for alice@example.com:**
```html
<h1>Hello Alice Johnson!</h1>
<p>Thank you for being part of TechCorp.</p>
```

### Template Rules

- ✅ Use `{{column_name}}` syntax
- ✅ Placeholders are case-sensitive
- ✅ Whitespace allowed: `{{ name }}` or `{{name}}`
- ✅ Works in HTML body and `--subject` flag
- ⚠️ CSV must have `email` column

### Complete Example

```bash
mangomail welcome-campaign \
  --html templates/welcome.html \
  --sample data/contacts.csv \
  --from "hello@mycompany.com" \
  --subject "Welcome {{name}}!" \
  --reply-to "support@mycompany.com" \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

### With Attachments

```bash
mangomail newsletter \
  --html templates/newsletter.html \
  --sample data/subscribers.csv \
  --from "news@example.com" \
  --subject "{{name}}, check out our latest updates" \
  --attachments "https://example.com/report.pdf,https://example.com/guide.pdf" \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

## Hosting Assets with AWS S3

For email images and attachments, upload files to S3 using the [AWS CLI](https://docs.aws.amazon.com/cli/v1/userguide/cli-services-s3-commands.html):

```bash
# Upload files
aws s3 cp logo.png s3://your-bucket/assets/ --acl public-read

# Use in template
<img src="https://your-bucket.s3.amazonaws.com/assets/logo.png">

# Or as attachment
mangomail campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Newsletter" \
  --attachments "https://your-bucket.s3.amazonaws.com/files/report.pdf"
```

### Rate Limiting

MangoMail automatically batches emails to respect Resend's rate limits (default: 2 emails/second).

**Adjust rate limit:**

```bash
# Send 5 emails per second (if your plan allows)
mangomail campaign \
  --html template.html \
  --sample contacts.csv \
  --from "hello@example.com" \
  --subject "Hello!" \
  --rate-limit 5 \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

**How it works:**
- Sends emails in batches of `rate-limit` size
- Waits 1 second between batches
- Logs progress: `"Waiting before next batch sent=2 total=10"`

## Campaign Reports

After each run, MangoMail generates a report in `.mangomail/<timestamp>.md`:

```markdown
# Email Campaign Report: welcome-campaign

## Campaign Details
- Campaign Name: welcome-campaign
- HTML Template: template.html
- CSV File: contacts.csv
- From: hello@example.com
- Subject: Welcome {{name}}!

## Statistics
- Total Recipients: 100

## Execution Summary
- Started at: 2025-10-06T14:30:00Z
- Completed at: 2025-10-06T14:31:45Z
- Duration: 1m45s
```

## Command Reference

### Flags

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--html` | string | ✅ Yes | Path to HTML template file |
| `--sample` | string | ✅ Yes | Path to CSV file |
| `--from` | string | ✅ Yes | Sender email address |
| `--subject` | string | ✅ Yes | Email subject (supports `{{placeholders}}`) |
| `--reply-to` | string | No | Reply-to email address |
| `--attachments` | string | No | Comma-separated attachment URLs |
| `--resend-api-key` | string | No | Resend API key (use env var instead) |
| `--resend-api-key-file` | string | No | Path to API key file |
| `--rate-limit` | uint64 | No | Max emails per second (default: 2) |
| `--log-level` | string | No | Log level: debug, info, warn, error |

### Example Workflows

#### Production Setup

```bash
# 1. One-time setup
mkdir -p ~/.mangomail/secrets
echo "re_your_production_key" > ~/.mangomail/secrets/resend_api_key
chmod 600 ~/.mangomail/secrets/resend_api_key

# 2. Add to ~/.bashrc or ~/.zshrc
echo 'export MANGOMAIL_RESEND_API_KEY_FILE="$HOME/.mangomail/secrets/resend_api_key"' >> ~/.bashrc

# 3. Run campaigns
mangomail campaign --html template.html --sample data.csv --from "hello@company.com" --subject "Hi {{name}}!"
```

#### Test Campaign (3 emails)

```bash
mangomail test-campaign \
  --html template.html \
  --sample test-contacts.csv \
  --from "test@example.com" \
  --subject "Test for {{name}}" \
  --resend-api-key-file ~/.mangomail/secrets/resend_api_key
```

**Console output:**
```
2025/10/06 14:30:00 INFO Email sent successfully to=[alice@example.com] subject="Test for Alice" id=abc123
2025/10/06 14:30:00 INFO Email sent successfully to=[bob@example.com] subject="Test for Bob" id=def456
2025/10/06 14:30:00 INFO Waiting before next batch sent=2 total=3 delay=1s
2025/10/06 14:30:01 INFO Email sent successfully to=[carol@example.com] subject="Test for Carol" id=ghi789
2025/10/06 14:30:01 INFO Campaign completed campaign=test-campaign duration=1.2s
2025/10/06 14:30:01 INFO Report generated path=.mangomail/20251006-143001.md
```

## Troubleshooting

### Rate Limit Errors

**Error:** `Too many requests. You can only make 2 requests per second`

**Solution:** The default rate limit is already set to 2. If you still see this error, your CSV might have issues or batching failed. Check:
- Ensure you're using the latest version
- Check `.mangomail/<timestamp>.md` report for details

### Missing Placeholders

**Issue:** Placeholders like `{{name}}` appear literally in sent emails

**Solution:**
- Verify CSV has matching column headers (case-sensitive)
- Check for typos: CSV has `Name` but template uses `{{name}}`

### API Key Not Found

**Error:** `MANGOMAIL_RESEND_API_KEY is required`

**Solution:**
```bash
# Check if file exists
cat ~/.mangomail/secrets/resend_api_key

# Use flag directly
mangomail campaign --resend-api-key-file ~/.mangomail/secrets/resend_api_key ...
```

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or PR on [GitHub](https://github.com/TheMagicMango/mangomail).
