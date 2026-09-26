// Package rbac provides role-based access control for Atlas tenants.
//
// Roles are scoped to tenants: each tenant defines its own role hierarchy
// and permission set. The operator's role is absolute — approval lives in
// the hand alone, and rbac never grants it.
//
// THE ROLE MODEL (2026-09-26; the operator's ruling: "kinds from the tool's
// own declaration"). A role's permissions are keyed three ways, and Can asks
// them in this order: the TOOL BY NAME (`git_tag`), the WILDCARD (`*`), then
// the KINDS the call carries. A kind is not a word the policy invents. It is
// the one thing the door already declares about a tool beyond its name --
// whether it writes -- read as `edit` for Writes:true and `read` otherwise,
// with `tools` (calling at all) carried by every call. Until this ruling the
// shipped roles spoke only in kinds (`read`, `edit`, `bash`, `net`, `tools`)
// while Can looked up tool names and the wildcard, so assigning ANY shipped
// role to a key denied it every tool: the policy and its reader disagreed,
// and nothing had ever assigned a shipped role to find out. `bash` and `net`
// named nothing the door declares, and are gone from the shipped roles; a
// policy on disk that still carries them is read, and those keys are simply
// never asked.
package rbac

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Role is a named set of permissions within a tenant.
type Role struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Permissions map[string]string `json:"permissions"`       // tool name, "*", or a Kind → allow/deny
	Implies     []string          `json:"implies,omitempty"` // roles this role includes
}

// Policy defines the RBAC policy for a tenant.
type Policy struct {
	Roles  map[string]Role   `json:"roles"`
	Assign map[string]string `json:"assign"` // agent_id → role_name
}

// Kind is what the door declares of a tool, and what a role is asked about
// when it names neither the tool nor the wildcard. Exactly these three exist,
// because the door declares exactly one thing about a tool beyond its name.
type Kind string

const (
	// KindTools is calling at all. Every call carries it.
	KindTools Kind = "tools"
	// KindRead is a tool that declares Writes: false.
	KindRead Kind = "read"
	// KindEdit is a tool that declares Writes: true.
	KindEdit Kind = "edit"
)

// Kinds are the three, in the order a role is asked about them.
var Kinds = []Kind{KindTools, KindRead, KindEdit}

// KindsOf are the kinds one call carries, from the declaration: `tools`
// always, then `edit` or `read` by whether the tool writes.
func KindsOf(writes bool) []Kind {
	if writes {
		return []Kind{KindTools, KindEdit}
	}
	return []Kind{KindTools, KindRead}
}

// IsKind reports whether a permission key is one of the three kinds.
func IsKind(key string) bool {
	for _, k := range Kinds {
		if string(k) == key {
			return true
		}
	}
	return false
}

// DefaultPolicy returns the default Atlas RBAC policy, spoken in the kinds
// the door declares (and nothing else, which TestEveryShippedPermission-
// IsAKindTheDoorDeclares holds it to):
//   - operator: reads and writes (can_approve is still structurally false)
//   - steward:  reads and writes; implies agent
//   - agent:    reads; every tool that writes is denied by kind
//   - guest:    denied calling at all
func DefaultPolicy() Policy {
	return Policy{
		Roles: map[string]Role{
			"operator": {
				Name:        "operator",
				Description: "the hand; holds the gate",
				Permissions: map[string]string{
					"tools": "allow",
					"read":  "allow",
					"edit":  "allow",
				},
			},
			"steward": {
				Name:        "steward",
				Description: "plans, specs, keeps THE_ROAD",
				Permissions: map[string]string{
					"tools": "allow",
					"read":  "allow",
					"edit":  "allow",
				},
				Implies: []string{"agent"},
			},
			"agent": {
				Name:        "agent",
				Description: "a declared seat with limited scope",
				Permissions: map[string]string{
					"tools": "allow",
					"read":  "allow",
					"edit":  "deny",
				},
			},
			"guest": {
				Name:        "guest",
				Description: "unauthenticated; deny all",
				Permissions: map[string]string{
					"tools": "deny",
					"read":  "deny",
					"edit":  "deny",
				},
			},
		},
		Assign: map[string]string{},
	}
}

// LoadPolicy reads a tenant's RBAC policy from <home>/rbac.json.
// If absent, returns the default policy.
func LoadPolicy(home string) Policy {
	p := DefaultPolicy()
	path := filepath.Join(home, "rbac.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return p
	}
	var got Policy
	if err := json.Unmarshal(data, &got); err != nil {
		return p
	}
	// Merge: custom roles override defaults, absent roles keep defaults
	if got.Roles != nil {
		for k, v := range got.Roles {
			p.Roles[k] = v
		}
	}
	if got.Assign != nil {
		for k, v := range got.Assign {
			p.Assign[k] = v
		}
	}
	return p
}

// RoleOf is the role a policy assigns an agent, or "" and false when it
// assigns none. A hold record reads this to SAY the role beside a parked
// call; it decides nothing.
func RoleOf(p Policy, agentID string) (string, bool) {
	r, ok := p.Assign[agentID]
	return r, ok
}

// Can checks if an agent may call a tool under a policy. `writes` is the
// tool's own declaration (Tool.Writes at the door), and every caller says it:
// a kind is never guessed from a name.
//
// The assigned role is asked first, then each role it implies, breadth-first.
// Every role is asked the same three questions in the same order, and the
// first that answers decides:
//
//  1. THE TOOL BY NAME  -- `git_tag: allow|deny`. A name beats everything
//     below it, so a role may deny a kind and still allow one tool of it.
//  2. THE WILDCARD      -- `*: allow|deny`, every tool by name.
//  3. THE KINDS         -- the ones this call carries (KindsOf). A deny on
//     any carried kind denies the call; an allow needs every carried kind
//     allowed. A role that speaks to neither has not answered, and the next
//     role is asked.
//
// Nothing answering is a refusal that names what nobody allowed. Returns
// (allowed, the role that decided, reason).
func Can(p Policy, agentID, toolName string, writes bool) (bool, string, string) {
	roleName, ok := p.Assign[agentID]
	if !ok {
		return false, "guest", fmt.Sprintf("agent %q has no assigned role", agentID)
	}
	if _, ok := p.Roles[roleName]; !ok {
		return false, "guest", fmt.Sprintf("role %q not found in policy", roleName)
	}
	kinds := KindsOf(writes)
	visited := map[string]bool{}
	queue := []string{roleName}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		r, ok := p.Roles[current]
		if !ok {
			continue
		}
		label := "role"
		if current != roleName {
			label = "implied role"
		}
		if allowed, decided, why := asked(r, label, current, toolName, kinds); decided {
			return allowed, current, why
		}
		queue = append(queue, r.Implies...)
	}
	return false, roleName, fmt.Sprintf("role %q allows neither tool %q by name nor its kind (%s: %s)",
		roleName, toolName, kinds[1], describe(toolName, kinds[1]))
}

// asked puts the three questions to one role, in order. decided is false when
// the role has nothing to say about this call. A permission that is neither
// "allow" nor "deny" is no answer.
func asked(r Role, label, roleName, toolName string, kinds []Kind) (allowed, decided bool, why string) {
	switch r.Permissions[toolName] {
	case "allow":
		return true, true, ""
	case "deny":
		return false, true, fmt.Sprintf("%s %q denies tool %q", label, roleName, toolName)
	}
	switch r.Permissions["*"] {
	case "allow":
		return true, true, ""
	case "deny":
		return false, true, fmt.Sprintf("%s %q denies all tools", label, roleName)
	}
	every := true
	for _, k := range kinds {
		switch r.Permissions[string(k)] {
		case "allow":
		case "deny":
			return false, true, fmt.Sprintf("%s %q denies %s (%s)", label, roleName, k, describe(toolName, k))
		default:
			every = false
		}
	}
	if every {
		return true, true, ""
	}
	return false, false, ""
}

// describe is a kind in a reader's words, for a refusal.
func describe(toolName string, k Kind) string {
	switch k {
	case KindEdit:
		return toolName + " writes"
	case KindRead:
		return toolName + " reads"
	default:
		return "calling any tool at all"
	}
}

// AssignAgent assigns a role to an agent in the policy.
func AssignAgent(p *Policy, agentID, roleName string) error {
	if _, ok := p.Roles[roleName]; !ok {
		return fmt.Errorf("role %q does not exist", roleName)
	}
	p.Assign[strings.ToLower(agentID)] = roleName
	return nil
}

// SavePolicy writes the RBAC policy to <home>/rbac.json.
func SavePolicy(home string, p Policy) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(home, "rbac.json"), data, 0644)
}
