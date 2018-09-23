package resource

import (
	"errors"
	"sort"

	"github.com/st-l10n/etree"
)

// Gen generates .po entry list from original xml, trying to apply translations
// from translated xml.
func Gen(original, translated []byte) (Entries, error) {
	var entries Entries
	eng := etree.NewDocument()
	if err := eng.ReadFromBytes(original); err != nil {
		return nil, err
	}
	d := etree.NewDocument()
	if len(translated) > 0 {
		if err := d.ReadFromBytes(translated); err != nil {
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
					ID:        elemKey,
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
				if elemPart.Tag != "Value" {
					entry.ID += "." + elemPart.Tag
				}
				entries = append(entries, entry)
			}
		}
	}
	sort.Sort(entries)
	return entries, nil
}
