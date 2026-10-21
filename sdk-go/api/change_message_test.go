package api

import (
	"testing"

	"github.com/go-faster/jx"
)

// A change page with expand=message carries inbound and outbound rows in one
// field; the direction discriminator picks the type of each.
func TestChangePageDecodesExpandedMessagesByDirection(t *testing.T) {
	const page = `{
	  "changes": [
	    {"kind": "email.visible", "email_id": "7c1e4a52-9a1d-4c55-9a3e-3c1f0b6c2a11", "sent_email_id": null,
	     "thread_id": null, "status": null, "agent_address": null, "reason": null,
	     "changed_at": "2026-10-20T12:00:00.000001Z",
	     "message": {"direction": "inbound", "id": "7c1e4a52-9a1d-4c55-9a3e-3c1f0b6c2a11", "status": "completed",
	       "sender": "bounce@mailer.example", "recipient": "agent@example.com", "domain": "example.com",
	       "created_at": "2026-10-20T12:00:00Z", "received_at": "2026-10-20T12:00:00Z", "webhook_attempt_count": 0,
	       "reply_count": 0, "last_replied_at": null, "awaiting": "you", "automated": false, "automated_reasons": [],
	       "from_verified": true, "from_header_email": "alex@example.org"}},
	    {"kind": "sent_email.updated", "email_id": null, "sent_email_id": "9d02b7a1-5f3c-4e8e-8a77-2b6f1d0c9e44",
	     "thread_id": null, "status": "delivered", "agent_address": null, "reason": null,
	     "changed_at": "2026-10-20T12:00:01.000001Z",
	     "message": {"direction": "outbound", "id": "9d02b7a1-5f3c-4e8e-8a77-2b6f1d0c9e44", "status": "delivered",
	       "status_changed_at": "2026-10-20T12:00:01Z", "created_at": "2026-10-20T11:59:00Z",
	       "updated_at": "2026-10-20T12:00:01Z", "content_hash": "sha256:abc", "from_header": "agent@example.com",
	       "from_address": "agent@example.com", "to_header": "alex@example.org", "to_address": "alex@example.org",
	       "subject": "Re: Plan", "body_size_bytes": 5}},
	    {"kind": "thread.read_state_changed", "email_id": null, "sent_email_id": null,
	     "thread_id": "a41f0c2e-6b1d-4f7a-9c33-8e2d5b7a1f00", "status": null, "agent_address": null, "reason": null,
	     "changed_at": "2026-10-20T12:00:02.000001Z", "message": null}
	  ],
	  "next_cursor": "djF8", "has_more": false, "baseline": false
	}`

	var decoded ChangePage
	if err := decoded.Decode(jx.DecodeStr(page)); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(decoded.Changes) != 3 {
		t.Fatalf("changes = %d, want 3", len(decoded.Changes))
	}

	inbound, ok := decoded.Changes[0].Message.Get()
	if !ok || inbound.Type != ChangeInboundMessageChangeMessage {
		t.Fatalf("first message type = %v (set %v), want inbound", inbound.Type, ok)
	}
	if got := inbound.ChangeInboundMessage.Sender; got != "bounce@mailer.example" {
		t.Errorf("inbound sender = %q", got)
	}

	outbound, ok := decoded.Changes[1].Message.Get()
	if !ok || outbound.Type != ChangeOutboundMessageChangeMessage {
		t.Fatalf("second message type = %v (set %v), want outbound", outbound.Type, ok)
	}
	if got := outbound.ChangeOutboundMessage.ToAddress; got != "alex@example.org" {
		t.Errorf("outbound to_address = %q", got)
	}

	if !decoded.Changes[2].Message.IsSet() || !decoded.Changes[2].Message.IsNull() {
		t.Errorf("third message should be present and null")
	}
}
