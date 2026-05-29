package proxy

import (
	"kiro-go/config"
	"testing"
)

func TestRebindAccountUsageCtxKeepsPendingRecordOnFinalAccount(t *testing.T) {
	usageCtxByAccountMu.Lock()
	usageCtxByAccount = map[string]*usageCtx{}
	usageCtxByAccountMu.Unlock()
	t.Cleanup(func() {
		usageCtxByAccountMu.Lock()
		usageCtxByAccount = map[string]*usageCtx{}
		usageCtxByAccountMu.Unlock()
	})

	uc := &usageCtx{
		Endpoint:    "/v1/chat/completions",
		RequestBody: `{"model":"kr/test"}`,
		Streaming:   true,
	}
	setAccountUsageCtx("initial", uc)
	rebindAccountUsageCtx("initial", "final")

	if got := getAccountUsageCtx("initial"); got != nil {
		t.Fatalf("expected initial account context to be cleared")
	}
	if got := getAccountUsageCtx("final"); got != uc {
		t.Fatalf("expected final account context to be rebound")
	}

	recordUsageWithCtx(&config.Account{ID: "final", Email: "final@example.com"}, "kr/test", 11, 7, 0.5, true, 200)

	if uc.Pending == nil {
		t.Fatalf("expected usage record to stay pending in rebound context")
	}
	if uc.Pending.AccountID != "final" {
		t.Fatalf("expected final account ID, got %q", uc.Pending.AccountID)
	}
	if uc.Pending.Endpoint != "/v1/chat/completions" {
		t.Fatalf("expected endpoint to be preserved, got %q", uc.Pending.Endpoint)
	}
	if uc.Pending.RequestBody != `{"model":"kr/test"}` {
		t.Fatalf("expected request body to be preserved, got %q", uc.Pending.RequestBody)
	}
	if !uc.Pending.Streaming {
		t.Fatalf("expected streaming flag to be preserved")
	}

	clearUsageCtx(uc)
	if got := getAccountUsageCtx("final"); got != nil {
		t.Fatalf("expected rebound context to be cleared")
	}
}
