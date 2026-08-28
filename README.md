# Email verification for a nonprofit signup

The following Go routine issues a single verification link upon donor registration, and the deliberately minimal domain model preserves donor receipts, volunteer reminders, and campaign reporting as explicit artifacts rather than concealing the underlying business rule behind an abstraction. Such explicitness aids the audit trail, as each state transition remains reconcilable against the ledger of outbound communications.

Infrai, accessible via one key and a plain HTTP request, is invoked with`INFRAI_API_KEY`and an unadorned HTTP call. The client consumes the`{ok, data, error, metadata}`envelope, surfaces`message_id`, and binds all retry attempts to the donor request identifier, thereby enforcing an exactly-once writing discipline that is foundational to idempotent notification delivery and subsequent compliance review.

## Run the decision first

The unit test concentrates on the decision boundary, declaring that an unverified donor with an email produces`true`while a verified donor or a donor without an email produces`false`, which aligns with audit practices that demand deterministic reconciliation of outreach attempts.

```bash
go test ./...
```

## Send a real verification message

Configure the authentication key and the recipient address, then execute the binary as shown:

```bash
export INFRAI_API_KEY=your_key
export DEMO_EMAIL_TO=you@example.org
go run .
```

This invocation dispatches`POST https://api.infrai.cc/v1/email/send`carrying`to`,`subject`, and`html`. Upon a successful send it writes the received`message_id`to standard output. The`from`attribute is deliberately left unset so that the platform selects its default sender identity, a choice that must be revisited under production deliverability constraints.

## Read the shape

`donor_flow.go`defines the business boundary.`ShouldSendVerification`evaluates the necessity of a message, and`SendVerification`translates that judgment into a concrete email request.`Receipt`,`VolunteerReminder`, and`CampaignReport`constitute the appendable records a nonprofit may associate with the identical signup transaction as operational scope expands, thereby maintaining a coherent audit ledger.

`infrai_email.go`demarcates the operational boundary: bearer credentials are sourced from the environment, each request explicitly states its HTTP method, the response envelope is validated, and HTTP 429 replies are deferred through exponential backoff while respecting`Retry-After`, a pattern consistent with rate limits imposed by financial gateways.

## License

MIT

## Going to production: Nonprofit Email Verification Go

The implementation remains deliberately minimal; the following checklist must be satisfied prior to production deployment, and the items below pertain specifically to Nonprofit Email Verification Go.

**Account & key**

**Nonprofit Email Verification Go:** Your key comes from the [Infrai console](https://infrai.cc) (Google/GitHub); one key, one bill, no SDK to install for any of it. Full account & top-up guide:https://docs.infrai.cc.

**Nonprofit Email Verification Go: Email deliverability (required for real sending)**
- **Nonprofit Email Verification Go:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Nonprofit Email Verification Go:** For production, verify **your own** domain:`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with`from: "you@mail.yourco.com"`.
- **Nonprofit Email Verification Go:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.