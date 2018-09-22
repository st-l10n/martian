package resource

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	update bool
)

func init() {
	flag.BoolVar(&update, "update", false, "update golden files")
}

func create(t *testing.T, path ...string) (*os.File, func()) {
	t.Helper()
	p := []string{"_testdata"}
	p = append(p, path...)
	name := filepath.Join(p...)
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	return f, func() {
		if closeErr := f.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}
}

func open(t *testing.T, path ...string) (*os.File, func()) {
	t.Helper()
	p := []string{"_testdata"}
	p = append(p, path...)
	name := filepath.Join(p...)
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("failed to open: %v", err)
	}
	return f, func() {
		if closeErr := f.Close(); closeErr != nil {
			t.Fatal(closeErr)
		}
	}
}

func TestBakeTips(t *testing.T) {
	engF, enfFClose := open(t, "Language", "english_tips.xml")
	defer enfFClose()
	tr, trClose := open(t, "tips.po")
	defer trClose()
	translation, err := ioutil.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	original, err := ioutil.ReadAll(engF)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Bake(Options{
		Original:    original,
		Translation: translation,
		Code:        "RU",
		Name:        "Russian",
		Font:        "font_russian",
	})
	if err != nil {
		t.Fatal(err)
	}
	if update {
		out, outClose := create(t, "tips.xml")
		out.Write(result)
		outClose()
	}
	out, outClose := open(t, "tips.xml")
	defer outClose()
	goldenOut, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.EqualFold(goldenOut, result) {
		t.Error("failed")
	}
}

func TestLoadTips(t *testing.T) {
	engF, enfFClose := open(t, "Language", "english_tips.xml")
	defer enfFClose()
	f, fClose := open(t, "Language", "russian_tips.xml")
	defer fClose()
	var (
		engRaw, rusRaw []byte
		err            error
	)
	if engRaw, err = ioutil.ReadAll(engF); err != nil {
		t.Fatal(err)
	}
	if rusRaw, err = ioutil.ReadAll(f); err != nil {
		t.Fatal(err)
	}
	result, err := Gen(GenOptions{
		Language:   "ru",
		Translated: rusRaw,
		Original:   engRaw,
	})

	if update {
		out, outClose := create(t, "tips.po")
		out.Write(result)
		outClose()
	}
}

func TestBake(t *testing.T) {
	engF, enfFClose := open(t, "Language", "english_keys.xml")
	defer enfFClose()
	tr, trClose := open(t, "out.po")
	defer trClose()
	translation, err := ioutil.ReadAll(tr)
	if err != nil {
		t.Fatal(err)
	}
	original, err := ioutil.ReadAll(engF)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Bake(Options{
		Original:    original,
		Translation: translation,
		Code:        "RU",
		Name:        "Russian",
		Font:        "font_russian",
	})
	if err != nil {
		t.Fatal(err)
	}
	if update {
		out, outClose := create(t, "out.xml")
		out.Write(result)
		outClose()
	}
	out, outClose := open(t, "out.xml")
	defer outClose()
	goldenOut, err := ioutil.ReadAll(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.EqualFold(goldenOut, result) {
		t.Error("failed")
	}
}

func TestLoad(t *testing.T) {
	engF, enfFClose := open(t, "Language", "english_keys.xml")
	defer enfFClose()
	f, fClose := open(t, "Language", "russian_keys.xml")
	defer fClose()
	var (
		engRaw, rusRaw []byte
		err            error
	)
	if engRaw, err = ioutil.ReadAll(engF); err != nil {
		t.Fatal(err)
	}
	if rusRaw, err = ioutil.ReadAll(f); err != nil {
		t.Fatal(err)
	}
	result, err := Gen(GenOptions{
		Language:   "ru",
		Translated: rusRaw,
		Original:   engRaw,
	})

	if update {
		out, outClose := create(t, "out.po")
		out.Write(result)
		outClose()
	}
}
