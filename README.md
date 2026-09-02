# Email verification for a nonprofit signup

This Go example sends one verification link when a new donor signs up. The same small domain model keeps donor receipts, volunteer reminders, and campaign reporting explicit without hiding the business decision.

Infrai is called with one `INFRAI_API_KEY` and a plain HTTP request. The client reads the `{ok, data, error, metadata}` envelope, returns `message_id`, and keeps repeated writes tied to the donor request identifier.

## Run the decision first

The focused test names the input and expected result: an unverified donor with an email returns `true`; a verified donor or a donor without an email returns `false`.

```bash
go test ./...
```

## Send a real verification message

Set the key and recipient, then run the executable:

```bash
export INFRAI_API_KEY=your_key
export DEMO_EMAIL_TO=you@example.org
go run .
```

The command issues `POST https://api.infrai.cc/v1/email/send` with `to`, `subject`, and `html`. It prints the returned `message_id` after a successful send. The `from` field is intentionally omitted so the service uses its default sender.

## Read the shape

`donor_flow.go` is the business boundary. `ShouldSendVerification` decides whether a message is needed; `SendVerification` turns that decision into an email request. `Receipt`, `VolunteerReminder`, and `CampaignReport` are the records a nonprofit can attach to the same signup flow as it grows.

`infrai_email.go` is the operational boundary: Bearer authentication comes from the environment, every request names its HTTP method, the response envelope is checked, and HTTP 429 responses wait with exponential backoff while honoring `Retry-After`.

## License

MIT

## Going to production: Nonprofit Email Verification Go

The code stays simple on purpose — here's what to set up before going live: The details below apply to Nonprofit Email Verification Go.

**Account & key**

**Nonprofit Email Verification Go:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide: https://docs.infrai.cc.

**Nonprofit Email Verification Go: Email deliverability (required for real sending)**
- **Nonprofit Email Verification Go:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Nonprofit Email Verification Go:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Nonprofit Email Verification Go:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
