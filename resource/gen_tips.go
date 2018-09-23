package resource

import (
	"fmt"
	"sort"
	"strings"

	"github.com/st-l10n/etree"
)

// reference {KEY:InventorySelect}
type reference struct {
	Type string // like "Key"
	Name string // like "InventorySelect"
}

func parseReferences(raw string) []reference {
	ref := make([]reference, 0, 10)
	for {
		startIndex := strings.Index(raw, "{")
		if startIndex < 0 {
			break
		}
		raw = raw[startIndex:]
		endIndex := strings.Index(raw, "}")
		if endIndex < 0 {
			break
		}
		kv := strings.SplitN(raw[:endIndex-1], ":", 2)
		if len(kv) < 2 {
			break
		}
		ref = append(ref, reference{
			Name: kv[0],
			Type: kv[1],
		})
		raw = raw[endIndex+1:]
	}
	return ref
}

type tip struct {
	raw           string
	translation   string
	refs          []reference
	codeReference string
}

func (t tip) equal(b []reference) bool {
	if len(b) != len(t.refs) {
		return false
	}
	for i, r := range b {
		if t.refs[i] != r {
			return false
		}
	}
	return true
}

func genTips(eng, d *etree.Document) (Entries, error) {
	var tips []tip
	for _, part := range eng.SelectElement("Language").ChildElements() {
		switch part.Tag {
		case "Name", "Code", "Font":
			continue
		}
		for i, c := range part.ChildElements() {
			dPath := c.GetPath() + fmt.Sprintf("[%d]", i)
			tips = append(tips, tip{
				codeReference: dPath,
				raw:           c.Text(),
				refs:          parseReferences(c.Text()),
			})
		}
	}
	if d != nil && d.SelectElement("Language") != nil {
		for _, part := range d.SelectElement("Language").ChildElements() {
			switch part.Tag {
			case "Name", "Code", "Font":
				continue
			}
			for _, c := range part.ChildElements() {
				raw := c.Text()
				refs := parseReferences(raw)
				for i, tip := range tips {
					if tip.equal(refs) {
						tip.translation = raw
						tips[i] = tip
						break
					}
				}
			}
		}
	}
	var entries Entries
	for _, tip := range tips {
		entries = append(entries, Entry{
			File:      "Tips",
			Reference: tip.codeReference,
			ID:        tip.raw,
			Str:       tip.translation,
		})
	}
	sort.Sort(entries)
	return entries, nil
}
