package webhooks

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPost(t *testing.T) {
	var gotSig string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get("X-Signature")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	err := Post(context.Background(), srv.Client(), srv.URL, "secret", map[string]string{"hello": "world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotSig == "" {
		t.Fatalf("expected signature header")
	}
}

func TestVerifySignature(t *testing.T) {
	secret := "test-secret"
	payload := []byte(`{"key":"value"}`)
	
	// Generate a valid signature
	validSig := sign(secret, payload)
	
	// Test valid signature
	if !VerifySignature(secret, payload, validSig) {
		t.Fatal("expected valid signature to verify")
	}
	
	// Test invalid signature
	if VerifySignature(secret, payload, "invalid-signature") {
		t.Fatal("expected invalid signature to fail verification")
	}
	
	// Test wrong secret
	wrongSecretSig := sign("wrong-secret", payload)
	if VerifySignature(secret, payload, wrongSecretSig) {
		t.Fatal("expected signature with wrong secret to fail verification")
	}
	
	// Test empty signature
	if VerifySignature(secret, payload, "") {
		t.Fatal("expected empty signature to fail verification")
	}
}
