package auth

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func findFixture(t *testing.T, name string) []byte {
	t.Helper()
	for _, c := range []string{
		filepath.Join("tests", "fixtures", name),
		filepath.Join("..", "tests", "fixtures", name),
		filepath.Join("..", "..", "tests", "fixtures", name),
		filepath.Join("..", "..", "..", "tests", "fixtures", name),
	} {
		if b, err := os.ReadFile(c); err == nil {
			return b
		}
	}
	t.Fatalf("fixture %s not found from here", name)
	return nil
}

func TestAuthContract(t *testing.T) {
	raw := findFixture(t, "auth_vectors.json")
	var doc struct {
		KeyRe  string `json:"key_re"`
		KidRe  string `json:"kid_re"`
		Iters  int    `json:"iters"`
		Domain string `json:"domain"`
		KDF    struct {
			Key  string `json:"key"`
			Salt string `json:"salt"`
			Hash string `json:"hash"`
		} `json:"kdf_example"`
		Scope []struct {
			Tenants []string `json:"tenants"`
			Request string   `json:"request"`
			Allow   bool     `json:"allow"`
		} `json:"scope_matrix"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	if KeyRe.String() != doc.KeyRe || KidRe.String() != doc.KidRe {
		t.Fatal("key shape laws drifted from goldens")
	}
	if Iters != doc.Iters || Domain != doc.Domain {
		t.Fatal("kdf constants drifted from goldens")
	}
	got, err := KDF(doc.KDF.Key, doc.KDF.Salt)
	if err != nil || got != doc.KDF.Hash {
		t.Fatalf("kdf drifted: %v (want %s)", err, doc.KDF.Hash)
	}
	if bad, _ := KDF("atl_cccccccccccccccccccccccccccccccc", doc.KDF.Salt); bad == doc.KDF.Hash {
		t.Fatal("wrong key must never verify")
	}
	for _, row := range doc.Scope {
		if ScopeOK(row.Tenants, row.Request) != row.Allow {
			t.Errorf("scope misscored: %+v", row)
		}
	}
}

func TestKeyLifecycle(t *testing.T) {
	home := t.TempDir()
	if !Empty(home) {
		t.Fatal("fresh home bootstraps open")
	}
	key, rec, err := Create(home, "op", []string{"*"}, "bootstrap")
	if err != nil {
		t.Fatal(err)
	}
	if !KeyRe.MatchString(key) || !KidRe.MatchString(rec.ID) {
		t.Fatalf("minted shapes break law: %q %q", key, rec.ID)
	}
	if Empty(home) {
		t.Fatal("minted store is no longer empty")
	}
	got, err := Verify(home, key)
	if err != nil || got.ID != rec.ID {
		t.Fatalf("minted key must verify: %v", err)
	}
	if _, err := Verify(home, "atl_ffffffffffffffffffffffffffffffff"); err == nil {
		t.Fatal("stranger keys refuse")
	}
	list := List(home)
	if len(list) != 1 || list[0].Hash != "" || list[0].Salt != "" {
		t.Fatalf("list must withhold secrets: %+v", list)
	}
	if err := Revoke(home, rec.ID, "op"); err != nil {
		t.Fatal(err)
	}
	if _, err := Verify(home, key); err == nil {
		t.Fatal("revoked keys refuse")
	}
	if !Empty(home) {
		t.Fatal("fully revoked store bootstraps again")
	}
	raw, err := os.ReadFile(filepath.Join(home, "state", "auth_audit.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 2 {
		t.Fatalf("audit must hold create+revoke, got %d", len(lines))
	}
	for _, l := range lines {
		var doc map[string]any
		_ = json.Unmarshal([]byte(l), &doc)
		for _, f := range []string{"ts", "action", "key_id", "tenant", "by"} {
			if _, ok := doc[f]; !ok {
				t.Fatalf("audit line misses %q: %s", f, l)
			}
		}
	}
}

// A KEY'S SCOPE MOVES WITHOUT THE SECRET (2026-09-28). The council's key was
// minted for one tenant while the record said two, and the only way to move
// it was to mint another and rotate the plaintext through .env. Scope replaces
// the list whole on a live key, audits it, and touches nothing else: the same
// plaintext verifies before and after, a revoked or unknown key is refused by
// id, and an empty list is refused because that would be a revocation by
// another name.
func TestAKeysScopeMovesWithoutTheSecret(t *testing.T) {
	home := t.TempDir()
	key, rec, err := Create(home, "council", []string{"research"}, "bootstrap")
	if err != nil {
		t.Fatal(err)
	}
	moved, err := Scope(home, rec.ID, []string{"research", "atlas"}, "op")
	if err != nil {
		t.Fatal(err)
	}
	if moved.ID != rec.ID || strings.Join(moved.Tenants, ",") != "research,atlas" {
		t.Fatalf("the moved record must carry the new list: %+v", moved)
	}
	if moved.Hash != "" || moved.Salt != "" {
		t.Fatalf("the answer must withhold the secret's derivation: %+v", moved)
	}
	got, err := Verify(home, key)
	if err != nil || strings.Join(got.Tenants, ",") != "research,atlas" {
		t.Fatalf("the same plaintext must verify and carry both: %+v %v", got, err)
	}
	if !ScopeOK(got.Tenants, "atlas") || ScopeOK(got.Tenants, "stranger") {
		t.Fatal("the widened key carries atlas and nothing it was not given")
	}
	// Narrowing is the same act.
	if _, err := Scope(home, rec.ID, []string{"atlas"}, "op"); err != nil {
		t.Fatal(err)
	}
	if got, _ := Verify(home, key); strings.Join(got.Tenants, ",") != "atlas" {
		t.Fatalf("narrowing must replace the list whole: %+v", got)
	}
	// The refusals.
	if _, err := Scope(home, rec.ID, nil, "op"); err == nil || !strings.Contains(err.Error(), "revoke") {
		t.Fatalf("an empty scope is refused and says what it would be: %v", err)
	}
	if _, err := Scope(home, "k-00000000", []string{"atlas"}, "op"); err == nil || !strings.Contains(err.Error(), `no live key "k-00000000"`) {
		t.Fatalf("an unknown id is refused by id: %v", err)
	}
	if err := Revoke(home, rec.ID, "op"); err != nil {
		t.Fatal(err)
	}
	if _, err := Scope(home, rec.ID, []string{"atlas"}, "op"); err == nil {
		t.Fatal("a revoked key's scope does not move")
	}
	// The audit holds every move, in order, with the list it set.
	raw, err := os.ReadFile(filepath.Join(home, "state", "auth_audit.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 4 {
		t.Fatalf("audit must hold create, scope, scope, revoke; got %d", len(lines))
	}
	var second map[string]any
	_ = json.Unmarshal([]byte(lines[1]), &second)
	if second["action"] != "scope" || second["tenant"] != "research,atlas" || second["key_id"] != rec.ID || second["by"] != "op" {
		t.Fatalf("the scope line must say what moved, to what, and by whom: %s", lines[1])
	}
}
