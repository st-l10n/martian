// Package resource implements utilities to generate po-files from
// language files and back.
package resource

import (
	"bytes"
	"errors"
	"fmt"
	"path"

	"github.com/leonelquinteros/gotext"
	"github.com/st-l10n/etree"
)

type poEntry struct {
	TranslatorComment string
	Reference         string
	ID                string
	Str               string
	Context           string
}

type GenOptions struct {
	Original   []byte
	Translated []byte
	Language   string
}

func writePO(po poEntry, b *bytes.Buffer) {
	b.WriteRune('\n')
	if len(po.TranslatorComment) > 0 {
		fmt.Fprintf(b, "# %s\n", po.TranslatorComment)
	}
	if len(po.Reference) > 0 {
		fmt.Fprintf(b, "#: %s\n", po.Reference)
	}
	if len(po.Context) > 0 {
		fmt.Fprintf(b, `msgctxt %q`, po.Context)
		b.WriteRune('\n')
	}
	fmt.Fprintf(b, `msgid %q`, po.ID)
	b.WriteRune('\n')
	fmt.Fprintf(b, `msgstr %q`, po.Str)
	b.WriteRune('\n')
}

func Gen(o GenOptions) ([]byte, error) {
	eng := etree.NewDocument()
	if err := eng.ReadFromBytes(o.Original); err != nil {
		return nil, err
	}
	d := etree.NewDocument()
	if err := d.ReadFromBytes(o.Translated); err != nil {
		return nil, err
	}
	b := new(bytes.Buffer)
	fmt.Fprintf(b, `# Stationeers translation
#
#, fuzzy
msgid ""
msgstr ""
"Project-Id-Version: v0.1n"
"Report-Msgid-Bugs-To: \n"
"Language: %s\n"
"MIME-Version: 1.0\n"
"Content-Type: text/plain; charset=UTF-8\n"
`, o.Language)
	for _, part := range eng.SelectElement("Language").ChildElements() {
		switch part.Tag {
		case "Name", "Code", "Font":
			continue
		}
		for i, c := range part.ChildElements() {
			k := c.SelectElement("Key")
			if k == nil {
				dPath := c.GetPath() + fmt.Sprintf("[%d]", i)
				dElem := d.FindElement(dPath)
				po := poEntry{
					Context:   part.Tag,
					ID:        c.Text(),
					Str:       c.Text(),
					Reference: fmt.Sprintf("%s/%s:%d", part.Tag, c.Tag, i),
				}
				if dElem != nil {
					po.Str = dElem.Text()
				}
				writePO(po, b)
				continue
			}
			elemKey := c.SelectElement("Key").Text()
			dPath := c.GetPath() + "[Key='" + elemKey + "']"
			dElem := d.FindElement(dPath)
			if dElem == nil {
				return nil, fmt.Errorf("path: %q failed", dPath)
			}
			for _, elemPart := range c.ChildElements() {
				switch elemPart.Tag {
				case "Key":
					continue
				}
				p := elemPart.GetRelativePath(c)
				dPart := dElem.FindElement(p)
				simplePath := path.Join(part.Tag, elemKey)
				if elemPart.Tag != "Value" {
					simplePath = path.Join(simplePath, elemPart.Tag)
				}
				po := poEntry{
					Context:           simplePath,
					Str:               dPart.Text(),
					ID:                elemPart.Text(),
					TranslatorComment: simplePath,
				}
				if po.ID == "" {
					po.TranslatorComment += " (Blank)"
				}
				writePO(po, b)
			}
		}
	}
	return b.Bytes(), nil
}

type Options struct {
	Original    []byte
	Translation []byte
	Code        string
	Name        string
	Font        string
}

// Bake generates new translation file.
// Original is original english xml file, translation is po-formatted file.
// Returns new xml.
func Bake(o Options) ([]byte, error) {
	if o.Name == "" || o.Code == "" {
		return nil, errors.New("no code or name provided")
	}
	translation, original := o.Translation, o.Original
	t := &gotext.Po{}
	t.Parse(translation)
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
				translated := t.GetC(e.Text(), part.Tag)
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
				simplePath := path.Join(part.Tag, elemKey)
				if elemPart.Tag != "Value" {
					simplePath = path.Join(simplePath, elemPart.Tag)
				}
				translated := t.GetC(engText, simplePath)
				elemPart.SetText(translated)
				if translated == "" {
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
