// The first strokes on the role model.
//
// WHY THIS FILE DID NOT EXIST UNTIL 2026-09-26. This package is the door's
// judgement of who may call what -- ranked second in tests/PROVING.md among
// the packages whose silent failure would hurt most -- and it had no prover.
// The fault that earned these was found by hand on 2026-09-25 while closing
// P0-13: the shipped roles spoke in kinds (`read`, `edit`, `bash`, `net`,
// `tools`) and Can read tool names, so every shipped role denied every tool,
// and no stroke had ever assigned one to find out.
//
// Hermetic: policies are built in the stroke; the one that touches disk uses
// t.TempDir(). Nothing reads the estate.
package rbac

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// assigned is a policy with one agent, "k", holding the named role, over the
// roles given (the shipped ones when none are).
func assigned(role string, roles ...Role) Policy {
	p := DefaultPolicy()
	for _, r := range roles {
		p.Roles[r.Name] = r
	}
	p.Assign["k"] = role
	return p
}

func role(name string, perms map[string]string, implies ...string) Role {
	return Role{Name: name, Permissions: perms, Implies: implies}
}

// mustDeny fails unless the call is refused and the reason carries every word.
func mustDeny(t *testing.T, p Policy, tool string, writes bool, words ...string) {
	t.Helper()
	ok, who, why := Can(p, "k", tool, writes)
	if ok {
		t.Fatalf("%s (writes=%v) was allowed by role %q", tool, writes, who)
	}
	for _, w := range words {
		if !strings.Contains(why, w) {
			t.Fatalf("the refusal of %s must say %q: %q", tool, w, why)
		}
	}
}

func mustAllow(t *testing.T, p Policy, tool string, writes bool) string {
	t.Helper()
	ok, who, why := Can(p, "k", tool, writes)
	if !ok {
		t.Fatalf("%s (writes=%v) was denied: %s", tool, writes, why)
	}
	return who
}

// --- the wire on the policy's side ------------------------------------------

// EVERY SHIPPED PERMISSION IS A KIND THE DOOR DECLARES. This is what was
// wrong: the roles spoke `bash` and `net`, which no tool carries, beside
// `read`, `edit` and `tools`, which nothing read. A key that names no tool
// and no kind is a permission connected to nothing, and this stroke is where
// that goes red. And every shipped role speaks to all three kinds, so a
// shipped role always ANSWERS -- it never falls through to "nothing allows".
func TestEveryShippedPermissionIsAKindTheDoorDeclares(t *testing.T) {
	p := DefaultPolicy()
	if len(p.Roles) == 0 {
		t.Fatal("no shipped roles; this stroke would be vacuous")
	}
	for name, r := range p.Roles {
		if r.Name != name {
			t.Errorf("role %q carries the name %q", name, r.Name)
		}
		for key, val := range r.Permissions {
			if key != "*" && !IsKind(key) {
				t.Errorf("role %q names %q, which is neither a tool kind nor the wildcard -- "+
					"it is read by nothing", name, key)
			}
			if val != "allow" && val != "deny" {
				t.Errorf("role %q says %q of %q; a permission is allow or deny", name, val, key)
			}
		}
		for _, k := range Kinds {
			if _, ok := r.Permissions[string(k)]; !ok {
				t.Errorf("role %q says nothing of %s and would fall through to no answer", name, k)
			}
		}
		for _, imp := range r.Implies {
			if _, ok := p.Roles[imp]; !ok {
				t.Errorf("role %q implies %q, which is not shipped", name, imp)
			}
		}
	}
	if len(p.Assign) != 0 {
		t.Fatalf("the shipped policy assigns %v; it must ship open", p.Assign)
	}
}

// KindsOf is the declaration read as kinds, and nothing more.
func TestTheKindsACallCarriesComeFromTheDeclaration(t *testing.T) {
	if got := KindsOf(true); len(got) != 2 || got[0] != KindTools || got[1] != KindEdit {
		t.Fatalf("a writer carries tools+edit, got %v", got)
	}
	if got := KindsOf(false); len(got) != 2 || got[0] != KindTools || got[1] != KindRead {
		t.Fatalf("a reader carries tools+read, got %v", got)
	}
	for _, k := range []string{"tools", "read", "edit"} {
		if !IsKind(k) {
			t.Fatalf("%q is a kind", k)
		}
	}
	for _, k := range []string{"bash", "net", "*", "git_tag", ""} {
		if IsKind(k) {
			t.Fatalf("%q is not a kind", k)
		}
	}
}

// --- the shipped roles, through the model -----------------------------------

// A SHIPPED ROLE DECIDES BY WHAT THE TOOL DECLARES. Before 2026-09-26 every
// row of this table was DENIED.
func TestTheShippedRolesDecideByTheDeclaration(t *testing.T) {
	for _, c := range []struct {
		role       string
		read, edit bool
		whyRead    string
		whyEdit    string
	}{
		{"operator", true, true, "", ""},
		{"steward", true, true, "", ""},
		{"agent", true, false, "", `role "agent" denies edit (git_tag writes)`},
		{"guest", false, false, `role "guest" denies tools (calling any tool at all)`,
			`role "guest" denies tools (calling any tool at all)`},
	} {
		p := assigned(c.role)
		for _, call := range []struct {
			tool   string
			writes bool
			want   bool
			why    string
		}{{"muster", false, c.read, c.whyRead}, {"git_tag", true, c.edit, c.whyEdit}} {
			ok, who, why := Can(p, "k", call.tool, call.writes)
			if ok != call.want {
				t.Fatalf("%s calling %s (writes=%v): allowed=%v, wanted %v (%s)",
					c.role, call.tool, call.writes, ok, call.want, why)
			}
			if who != c.role {
				t.Fatalf("%s calling %s was decided by %q", c.role, call.tool, who)
			}
			if !ok && why != call.why {
				t.Fatalf("%s calling %s: reason %q, wanted %q", c.role, call.tool, why, call.why)
			}
			if ok && why != "" {
				t.Fatalf("an allow carries no reason: %q", why)
			}
		}
	}
}

// --- the order of the three questions ---------------------------------------

// A NAME BEATS THE WILDCARD, AND THE WILDCARD BEATS A KIND. Each pair below
// sets two questions against each other and the earlier one decides.
func TestANameBeatsTheWildcardBeatsAKind(t *testing.T) {
	// name over kind: a role that denies writing may still allow one writer.
	p := assigned("x", role("x", map[string]string{
		"tools": "allow", "read": "allow", "edit": "deny", "git_tag": "allow"}))
	mustAllow(t, p, "git_tag", true)
	mustDeny(t, p, "git_commit", true, `role "x" denies edit (git_commit writes)`)

	// name over kind, the other way: a reader denied by name under an open kind.
	p = assigned("x", role("x", map[string]string{
		"tools": "allow", "read": "allow", "edit": "allow", "muster": "deny"}))
	mustDeny(t, p, "muster", false, `role "x" denies tool "muster"`)
	mustAllow(t, p, "records", false)

	// wildcard over kind, both ways.
	p = assigned("x", role("x", map[string]string{"*": "deny", "tools": "allow", "read": "allow"}))
	mustDeny(t, p, "muster", false, `role "x" denies all tools`)
	p = assigned("x", role("x", map[string]string{"*": "allow", "edit": "deny"}))
	mustAllow(t, p, "git_tag", true)

	// name over wildcard.
	p = assigned("x", role("x", map[string]string{"*": "deny", "muster": "allow"}))
	mustAllow(t, p, "muster", false)
	mustDeny(t, p, "records", false, `denies all tools`)
}

// AN ALLOW NEEDS EVERY KIND THE CALL CARRIES; A DENY NEEDS ONE. A role that
// speaks to neither has not answered, and the refusal says what it named.
func TestAnAllowNeedsEveryKindTheCallCarries(t *testing.T) {
	// `tools` alone allows nothing: the call also carries read or edit.
	p := assigned("x", role("x", map[string]string{"tools": "allow"}))
	mustDeny(t, p, "muster", false, `role "x" allows neither tool "muster" by name nor its kind (read: muster reads)`)
	mustDeny(t, p, "git_tag", true, `nor its kind (edit: git_tag writes)`)

	// `read` alone allows nothing either: every call carries `tools`.
	p = assigned("x", role("x", map[string]string{"read": "allow"}))
	mustDeny(t, p, "muster", false, `allows neither`)

	// both carried kinds allowed: the call goes; a kind not spoken to does not.
	p = assigned("x", role("x", map[string]string{"tools": "allow", "read": "allow"}))
	mustAllow(t, p, "muster", false)
	mustDeny(t, p, "git_tag", true, `allows neither tool "git_tag"`)

	// one deny is enough, whatever else is allowed.
	p = assigned("x", role("x", map[string]string{"tools": "deny", "read": "allow", "edit": "allow"}))
	mustDeny(t, p, "muster", false, `role "x" denies tools (calling any tool at all)`)
	p = assigned("x", role("x", map[string]string{"tools": "allow", "read": "deny", "edit": "allow"}))
	mustDeny(t, p, "muster", false, `role "x" denies read (muster reads)`)
	mustAllow(t, p, "git_tag", true)

	// a word that is neither allow nor deny is no answer.
	p = assigned("x", role("x", map[string]string{"tools": "allow", "read": "ask"}))
	mustDeny(t, p, "muster", false, `allows neither`)
}

// --- implied roles ------------------------------------------------------------

// THE ASSIGNED ROLE IS ASKED FIRST, THEN WHAT IT IMPLIES, and the role that
// answered is the one named back. A cycle terminates; a missing role is
// skipped.
func TestAnImpliedRoleIsAskedAfterTheAssignedOne(t *testing.T) {
	// a role that says nothing hands the question to what it implies.
	p := assigned("lead", role("lead", map[string]string{}, "agent"))
	if who := mustAllow(t, p, "muster", false); who != "agent" {
		t.Fatalf("the implied role answered and %q was named", who)
	}
	mustDeny(t, p, "git_tag", true, `implied role "agent" denies edit (git_tag writes)`)

	// the assigned role's own answer stands over what it implies.
	p = assigned("lead", role("lead", map[string]string{
		"tools": "allow", "read": "allow", "edit": "allow"}, "agent"))
	if who := mustAllow(t, p, "git_tag", true); who != "lead" {
		t.Fatalf("the assigned role answered and %q was named", who)
	}

	// a cycle ends; a role nobody defined is skipped; nothing answers.
	p = assigned("a", role("a", map[string]string{}, "b", "ghost"), role("b", map[string]string{}, "a"))
	mustDeny(t, p, "muster", false, `role "a" allows neither`)
	ok, who, _ := Can(p, "k", "muster", false)
	if ok || who != "a" {
		t.Fatalf("an unanswered call names the assigned role back: %v %q", ok, who)
	}
}

// --- no role, no entry ----------------------------------------------------------

func TestNoRoleNoEntry(t *testing.T) {
	p := DefaultPolicy()
	ok, who, why := Can(p, "stranger", "muster", false)
	if ok || who != "guest" || !strings.Contains(why, `agent "stranger" has no assigned role`) {
		t.Fatalf("an unassigned agent: %v %q %q", ok, who, why)
	}
	p.Assign["k"] = "nobody-defined-this"
	ok, who, why = Can(p, "k", "muster", false)
	if ok || who != "guest" || !strings.Contains(why, `role "nobody-defined-this" not found`) {
		t.Fatalf("a role the policy lacks: %v %q %q", ok, who, why)
	}
	if r, has := RoleOf(p, "k"); !has || r != "nobody-defined-this" {
		t.Fatalf("RoleOf says the assignment as written, decides nothing: %q %v", r, has)
	}
	if r, has := RoleOf(p, "stranger"); has || r != "" {
		t.Fatalf("RoleOf on the unassigned: %q %v", r, has)
	}
}

// --- a policy on disk that still speaks the old words -----------------------------

// THE ESTATE'S OWN atlas/rbac.json CARRIES `bash` AND `net` (written by the
// old defaults). LoadPolicy merges its roles over the shipped ones, and the
// dead keys are never asked -- so the file keeps deciding exactly as the
// shipped roles do, and nothing has to be rewritten for the door to be right.
func TestAPolicyOnDiskWithDeadKindsStillDecidesByTheLiveOnes(t *testing.T) {
	home := t.TempDir()
	old := `{"roles": {"agent": {"name": "agent", "permissions":
	          {"bash": "deny", "edit": "deny", "net": "deny", "read": "allow", "tools": "allow"}}},
	         "assign": {"k": "agent"}}`
	if err := os.WriteFile(filepath.Join(home, "rbac.json"), []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	p := LoadPolicy(home)
	if _, has := p.Roles["agent"].Permissions["bash"]; !has {
		t.Fatal("the file's own role must be the one loaded")
	}
	mustAllow(t, p, "muster", false)
	mustDeny(t, p, "git_tag", true, `role "agent" denies edit (git_tag writes)`)
	// and the other shipped roles are still there beneath it.
	for _, r := range []string{"operator", "steward", "guest"} {
		if _, ok := p.Roles[r]; !ok {
			t.Fatalf("shipped role %q was lost in the merge", r)
		}
	}
}
