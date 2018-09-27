package resource

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func readAll(name string) ([]byte, error) {
	oF, oErr := os.Open(name)
	defer oF.Close()
	if oErr != nil {
		return nil, oErr
	}
	data, err := ioutil.ReadAll(oF)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Merge merges original .po file with updated template .po file.
func Merge(original, template string) error {
	var (
		err  error
		orig []byte
	)
	if orig, err = readAll(original); err != nil {
		return err
	}
	p := &parser{
		orig: orig,
	}
	p.populate()
	b := new(bytes.Buffer)

	// Running gettext merge to add new msgid's to ".po" file.
	cmd := exec.Command("msgmerge",
		"-U", "--no-wrap",
		"--backup=off",
		original, template,
	)
	cmd.Stdout = b
	cmd.Stderr = b
	if err = cmd.Run(); err != nil {
		return fmt.Errorf("failed to merge: %s (%v)", b, err)
	}

	// Merge.
	if p.merged, err = readAll(original); err != nil {
		return err
	}
	p.merge()

	f, createErr := os.Create(original)
	if createErr != nil {
		return createErr
	}
	if _, err = io.Copy(f, p.result); err != nil {
		return err
	}
	return f.Close()
}

// parser for .po files specializing on merges.
type parser struct {
	orig   []byte
	merged []byte
	result *bytes.Buffer

	origID map[string][]string // map[ref] ->
}

func (p *parser) merge() {
	r := bytes.NewReader(p.merged)
	s := bufio.NewScanner(r)
	p.result = new(bytes.Buffer)

	var (
		lines   []string
		ref     string
		isFuzzy bool
	)
	for s.Scan() {
		l := s.Text() // one line
		if strings.HasPrefix(l, "#:") {
			ref = strings.TrimSpace(l)
			ref = strings.TrimSpace(strings.TrimPrefix(ref, "#:"))
		}
		if strings.HasPrefix(l, "#, fuzzy") {
			isFuzzy = true
		}
		if len(strings.TrimSpace(l)) == 0 {
			if len(ref) > 0 && isFuzzy && len(p.origID[ref]) > 0 {
				orig := p.origID[ref]
				// Replace everything after "msgstr".
				newLines := make([]string, 0, len(lines))
				var (
					oldMsgStr []string
				)
				msgStr := false
				for _, line := range orig {
					if strings.HasPrefix(line, "msgstr") {
						msgStr = true
					}
					if msgStr {
						oldMsgStr = append(oldMsgStr, line)
					}
				}
				newLines = newLines[:0]
				for _, line := range lines {
					if strings.HasPrefix(line, "msgstr") {
						break
					}
					newLines = append(newLines, line)
				}
				lines = append(newLines, oldMsgStr...)
			}
			if p.result.Len() != 0 {
				p.result.WriteRune('\n')
			}
			for _, l := range lines {
				fmt.Fprintln(p.result, l)
			}
			ref = ""
			lines = lines[:0]
			isFuzzy = false
		} else {
			isFuzzy = false
			lines = append(lines, l)
		}
	}
	if len(lines) != 0 {
		// Write last entry.
		p.result.WriteRune('\n')
		for _, l := range lines {
			fmt.Fprintln(p.result, l)
		}
	}
}

func (p *parser) populate() {
	p.origID = make(map[string][]string)

	// parse "orig", populating all entries.
	r := bytes.NewReader(p.orig)
	s := bufio.NewScanner(r)
	var (
		lines []string
		ref   string
	)
	for s.Scan() {
		l := s.Text() // one line
		if strings.HasPrefix(l, "#:") {
			ref = strings.TrimSpace(l)
			ref = strings.TrimSpace(strings.TrimPrefix(ref, "#:"))
		}
		if len(strings.TrimSpace(l)) == 0 {
			if len(ref) > 0 {
				p.origID[ref] = append(p.origID[ref][:0], lines...)
			}
			ref = ""
			lines = lines[:0]
		} else {
			lines = append(lines, l)
		}
	}
}
