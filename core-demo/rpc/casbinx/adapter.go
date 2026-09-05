package casbinx

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/casbinrule"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

type Adapter struct {
	client *ent.Client
}

func NewAdapter(client *ent.Client) *Adapter {
	return &Adapter{client: client}
}

var _ persist.BatchAdapter = (*Adapter)(nil)

func (a *Adapter) LoadPolicy(m model.Model) error {
	rows, err := a.client.CasbinRule.Query().All(context.Background())
	if err != nil {
		return err
	}
	for _, row := range rows {
		if err := persist.LoadPolicyArray(ruleOf(row), m); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) SavePolicy(m model.Model) error {
	ctx := context.Background()
	if _, err := a.client.CasbinRule.Delete().Exec(ctx); err != nil {
		return err
	}
	lines := policyLines(m)
	if len(lines) == 0 {
		return nil
	}
	builders := make([]*ent.CasbinRuleCreate, 0, len(lines))
	for _, line := range lines {
		builders = append(builders, a.createBuilder(line.ptype, line.rule))
	}
	return a.client.CasbinRule.CreateBulk(builders...).Exec(ctx)
}

func (a *Adapter) AddPolicy(sec, ptype string, rule []string) error {
	return a.AddPolicies(sec, ptype, [][]string{rule})
}

func (a *Adapter) AddPolicies(sec, ptype string, rules [][]string) error {
	if len(rules) == 0 {
		return nil
	}
	builders := make([]*ent.CasbinRuleCreate, 0, len(rules))
	for _, rule := range rules {
		builders = append(builders, a.createBuilder(ptype, rule))
	}
	return a.client.CasbinRule.CreateBulk(builders...).Exec(context.Background())
}

func (a *Adapter) RemovePolicies(sec, ptype string, rules [][]string) error {
	for _, rule := range rules {
		if err := a.RemovePolicy(sec, ptype, rule); err != nil {
			return err
		}
	}
	return nil
}

func (a *Adapter) RemovePolicy(sec, ptype string, rule []string) error {
	q := a.client.CasbinRule.Delete().Where(casbinrule.PtypeEQ(ptype))
	q = eqV(q, 0, at(rule, 0))
	q = eqV(q, 1, at(rule, 1))
	q = eqV(q, 2, at(rule, 2))
	q = eqV(q, 3, at(rule, 3))
	q = eqV(q, 4, at(rule, 4))
	q = eqV(q, 5, at(rule, 5))
	_, err := q.Exec(context.Background())
	return err
}

func (a *Adapter) RemoveFilteredPolicy(sec, ptype string, fieldIndex int, fieldValues ...string) error {
	q := a.client.CasbinRule.Delete().Where(casbinrule.PtypeEQ(ptype))
	for i, v := range fieldValues {
		if v == "" {
			continue
		}
		q = eqV(q, fieldIndex+i, v)
	}
	_, err := q.Exec(context.Background())
	return err
}

func (a *Adapter) createBuilder(ptype string, rule []string) *ent.CasbinRuleCreate {
	return a.client.CasbinRule.Create().
		SetPtype(ptype).
		SetV0(at(rule, 0)).
		SetV1(at(rule, 1)).
		SetV2(at(rule, 2)).
		SetV3(at(rule, 3)).
		SetV4(at(rule, 4)).
		SetV5(at(rule, 5))
}

func eqV(q *ent.CasbinRuleDelete, idx int, v string) *ent.CasbinRuleDelete {
	switch idx {
	case 0:
		return q.Where(casbinrule.V0EQ(v))
	case 1:
		return q.Where(casbinrule.V1EQ(v))
	case 2:
		return q.Where(casbinrule.V2EQ(v))
	case 3:
		return q.Where(casbinrule.V3EQ(v))
	case 4:
		return q.Where(casbinrule.V4EQ(v))
	case 5:
		return q.Where(casbinrule.V5EQ(v))
	default:
		return q
	}
}

func ruleOf(row *ent.CasbinRule) []string {
	vals := []string{row.Ptype, row.V0, row.V1, row.V2, row.V3, row.V4, row.V5}
	for i := len(vals) - 1; i >= 0; i-- {
		if strings.TrimSpace(vals[i]) != "" {
			return vals[:i+1]
		}
	}
	return vals[:1]
}

type policyLine struct {
	ptype string
	rule  []string
}

func policyLines(m model.Model) []policyLine {
	var out []policyLine
	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			out = append(out, policyLine{ptype: ptype, rule: rule})
		}
	}
	return out
}

func at(rule []string, i int) string {
	if i < 0 || i >= len(rule) {
		return ""
	}
	return rule[i]
}
