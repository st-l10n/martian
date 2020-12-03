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
	Simplified  []string // see GenOptions.Simplified
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
		simplified := false
		for _, s := range o.Simplified {
			if s == part.Tag {
				simplified = true
			}
		}
	Loop:
		for _, e := range part.ChildElements() {
			k := e.SelectElement("Key")
			if k == nil {
				// Tips.
				translated := t.Get(e.Text())
				if translated == "" || translated == e.Text() {
					part.RemoveChild(e)
					continue
				}
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
				elemSimplified := simplified
				switch elemPart.Tag {
				case "Key":
					continue
				case "Description":
					fmt.Println("ignoring description for", engPath)
					// TODO(ernado): Handle "Description" field
					// Ref: https://github.com/st-l10n/martian/issues/3
					e.RemoveChild(elemPart)
					continue
				}
				for _, s := range o.Simplified {
					if s == part.Tag+"."+elemPart.Tag {
						elemSimplified = true
					}
				}
				p := elemPart.GetRelativePath(e)
				engPart := engElem.FindElement(p)
				if engPart == nil {
					continue
				}
				engText := engPart.Text()
				var id string
				if elemPart.Tag != "Value" {
					id += "." + elemPart.Tag
				}
				if elemSimplified {
					// Using simplified relative path as ID.
					id = elemKey
					if elemPart.Tag != "Value" {
						id += "." + elemPart.Tag
					}
				} else {
					// Using original text as ID.
					id = engText
				}
				translated := t.GetC(id, part.Tag+"."+elemKey)
				if translated == "" || translated == id {
					part.RemoveChild(e)
					continue Loop
				}
				elemPart.SetText(translated)
				if translated == Blank {
					for _, child := range elemPart.Child {
						elemPart.RemoveChild(child)
					}
				}
			}
		}
	}
	d.Indent(2)
	return d.WriteToBytes()
}
