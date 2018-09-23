package resource

import (
	"errors"
	"fmt"
	"sort"

	"github.com/st-l10n/etree"
)

type GenOptions struct {
	Original   []byte
	Translated []byte

	// If simplified, the Entry.ID value of translation is set to simplified
	// relative path of translated element, not the original text.
	//
	// Like for Reagents/RecordReagent with Key=Flour, the id for Unit will be "Flour.Unit"
	// instead of "g".
	//
	// The "Tips" part is always assumed as non-simplified.
	SimplifiedParts []string
}

// Gen generates .po entry list from original xml, trying to apply translations
// from translated xml.
func Gen(o GenOptions) (Entries, error) {
	var entries Entries
	eng := etree.NewDocument()
	if err := eng.ReadFromBytes(o.Original); err != nil {
		return nil, err
	}
	d := etree.NewDocument()
	if len(o.Translated) > 0 {
		if err := d.ReadFromBytes(o.Translated); err != nil {
			return nil, err
		}
	}
	l := eng.SelectElement("Language")
	if l == nil {
		return nil, errors.New("no language elem")
	}
	g := l.SelectElement("GameTip")
	if g != nil && len(g.Child) != 0 {
		return genTips(eng, d)
	}
	for _, part := range l.ChildElements() {
		switch part.Tag {
		case "Name", "Code", "Font":
			continue
		}
		simplified := false
		for _, s := range o.SimplifiedParts {
			if s == part.Tag {
				simplified = true
			}
		}
		for _, c := range part.ChildElements() {
			elemKey := c.SelectElement("Key").Text()
			dPath := c.GetPath() + "[Key='" + elemKey + "']"
			dElem := d.FindElement(dPath)
			for _, elemPart := range c.ChildElements() {
				switch elemPart.Tag {
				case "Key":
					continue
				}
				p := elemPart.GetRelativePath(c)
				var dPart *etree.Element
				entry := Entry{
					Context:   part.Tag,
					File:      part.Tag,
					Reference: dPath,
					Original:  elemPart.Text(),
				}
				if entry.Original == "" {
					entry.Original = Blank
				}
				if dElem != nil {
					dPart = dElem.FindElement(p)
				}
				if dPart != nil {
					entry.Str = dPart.Text()
					if entry.Str == "" {
						entry.Str = Blank
					}
				}
				if simplified {
					// Using simplified relative path as ID.
					entry.ID = elemKey
					if elemPart.Tag != "Value" {
						entry.ID += "." + elemPart.Tag
					}
				} else {
					// Using original text as ID.
					entry.ID = entry.Original
				}
				if entry.ID != entry.Original {
					entry.TranslatorComment = fmt.Sprintf("Original: %q", entry.Original)
				}
				entries = append(entries, entry)
			}
		}
	}
	sort.Sort(entries)
	return entries, nil
}
