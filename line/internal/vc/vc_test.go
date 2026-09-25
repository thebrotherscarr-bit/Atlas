package vc

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const probe = "# probe\n\n```json\n" +
	`{"id": "probe", "can_approve": false, "covenant": "abc123", "office": "O", "reports_to": "boss"}` +
	"\n```\n"

// ONE SOURCE, THE MANIFEST (operator, 2026-09-25). The issuer is minted in
// the record's own covenant, so subject, reporting line and issuer share one
// namespace by construction; a record with no covenant gets no issuer and no
// credential; a hand that names an issuer is obeyed.
func TestTheIssuerIsMintedInTheRecordsOwnCovenant(t *testing.T) {
	block, prose, err := ParseUS(probe)
	if err != nil {
		t.Fatal(err)
	}
	issuer, err := IssuerFor(block)
	if err != nil || issuer != "did:atlas:abc123:operator" {
		t.Fatalf("issuer %q, %v", issuer, err)
	}
	v, err := ToVC(block, prose, issuer)
	if err != nil {
		t.Fatal(err)
	}
	if v.Issuer != issuer || v.CredentialSubject.ID != "did:atlas:abc123:probe" ||
		v.CredentialSubject.ReportsTo != "did:atlas:abc123:boss" {
		t.Fatalf("one namespace for all three, got issuer %s subject %+v", v.Issuer, v.CredentialSubject)
	}
	if _, err := IssuerFor(map[string]any{"id": "probe", "can_approve": false}); err == nil ||
		!strings.Contains(err.Error(), "no covenant") {
		t.Fatalf("a record with no covenant must be refused by name: %v", err)
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "probe.us")
	if err := os.WriteFile(path, []byte(probe), 0o644); err != nil {
		t.Fatal(err)
	}
	minted, err := FromFile(path, "")
	if err != nil || minted.Issuer != "did:atlas:abc123:operator" {
		t.Fatalf("FromFile with no issuer named must mint from the record: %+v %v", minted, err)
	}
	named, err := FromFile(path, "did:atlas:zzz:auditor")
	if err != nil || named.Issuer != "did:atlas:zzz:auditor" {
		t.Fatalf("a named issuer is a hand's choice and is used as given: %+v %v", named, err)
	}
	bare := filepath.Join(dir, "bare.us")
	os.WriteFile(bare, []byte("# bare\n\n```json\n{\"id\": \"bare\", \"can_approve\": false, \"office\": \"O\"}\n```\n"), 0o644)
	if _, err := FromFile(bare, ""); err == nil || !strings.Contains(err.Error(), "no covenant") {
		t.Fatalf("no covenant, no issuer, no credential: %v", err)
	}
}
