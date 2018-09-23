package resource

import (
	"errors"
	"fmt"

	"github.com/leonelquinteros/gotext"
	"github.com/st-l10n/etree"
)

type Options struct {
	Original    []byte   // "xml"
	Translation [][]byte // ".po" files
	Code        string
	Name        string
	Font        string
}

const Blank = "{BLANK}"

// Bake generates new translation file.
// Original is original english xml file, translation is po-formatted file.
// Returns new xml.
func Bake(o Options) ([]byte, error) {
	if o.Code == "" {
		return nil, errors.New("no code provided")
	}
	translations, original := o.Translation, o.Original
	t := &gotext.Po{}
	for _, translation := range translations {
		t.Parse(translation)
	}
	eng := etree.NewDocument()
	if err := eng.ReadFromBytes(original); err != nil {
		return original, fmt.Errorf("failed to parse original: %v", err)
	}
	d := eng.Copy()
	e := d.SelectElement("Language")
	name := e.FindElement("Name")
	if name != nil {
		name.SetText(o.Name)
	}
	if o.Name == "" && name != nil {
		e.RemoveChild(name)
	}
	e.FindElement("Code").SetText(o.Code)
	f := e.FindElement("Font")
	if len(o.Font) > 0 {
		if f == nil {
			f = e.CreateElement("Font")
		}
		f.SetText(o.Font)
	} else if f != nil {
		e.RemoveChild(f)
	}
	for _, part := range e.ChildElements() {
		switch part.Tag {
		case "Name", "Code", "Font":
			continue
		}
		for _, e := range part.ChildElements() {
			k := e.SelectElement("Key")
			if k == nil {
				// Tips.
				translated := t.Get(e.Text())
				e.SetText(translated)
				continue
			}
			elemKey := k.Text()
			engPath := e.GetPath() + "[Key='" + elemKey + "']"
			engElem := eng.FindElement(engPath)
			if engElem == nil {
				continue
			}
			for _, elemPart := range e.ChildElements() {
				switch elemPart.Tag {
				case "Key":
					continue
				}
				p := elemPart.GetRelativePath(e)
				engPart := engElem.FindElement(p)
				engText := engPart.Text()
				id := elemKey
				if elemPart.Tag != "Value" {
					id += "." + elemPart.Tag
				}
				translated := t.GetC(id, part.Tag)
				elemPart.SetText(translated)
				if translated == "" {
					translated = engText
				}
				if translated == Blank {
					for _, child := range elemPart.Child {
						elemPart.RemoveChild(child)
					}
				}
			}
		}
	}
	d.WriteSettings.WhitespaceEndTags = true
	return d.WriteToBytes()
}
