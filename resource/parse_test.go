package resource

import (
	"bytes"
	"encoding/json"
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

func TestBake(t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
		engF, enfFClose := open(t, "Language", "english_keys.xml")
		defer enfFClose()
		tr, trClose := open(t, "Keys-RU.po")
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
			Translation: [][]byte{translation},
			Code:        "RU",
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
	})
	t.Run("Tips", func(t *testing.T) {
		engF, enfFClose := open(t, "Language", "english_tips.xml")
		defer enfFClose()
		tr, trClose := open(t, "Tips-RU.po")
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
			Translation: [][]byte{translation},
			Code:        "RU",
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
	})
	t.Run("Default", func(t *testing.T) {
		engF, enfFClose := open(t, "Language", "english.xml")
		defer enfFClose()
		var translations [][]byte
		for _, name := range []string{
			"Actions",
			"Colors",
			"Gases",
			"Interactables",
			"Interface",
			"Mineables",
			"Reagents",
			"Slots",
			"Things",
		} {
			tr, trClose := open(t, name+"-RU.po")
			translation, err := ioutil.ReadAll(tr)
			if err != nil {
				t.Fatal(err)
			}
			trClose()
			translations = append(translations, translation)
		}
		original, err := ioutil.ReadAll(engF)
		if err != nil {
			t.Fatal(err)
		}
		result, err := Bake(Options{
			Original:    original,
			Translation: translations,
			Code:        "RU",
			Name:        "Russian",
			Font:        "font_russian",
		})
		if err != nil {
			t.Fatal(err)
		}
		if update {
			out, outClose := create(t, "default.xml")
			out.Write(result)
			outClose()
		}
		out, outClose := open(t, "default.xml")
		defer outClose()
		goldenOut, err := ioutil.ReadAll(out)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.EqualFold(goldenOut, result) {
			t.Error("failed")
		}
	})
}

func TestLoad(t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
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
		result, err := Gen(engRaw, rusRaw)
		if err != nil {
			t.Fatal(err)
		}
		if update {
			out, outClose := create(t, "keys.json")
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			if err = enc.Encode(result); err != nil {
				t.Fatal(err)
			}
			outClose()
			for _, f := range result.Files() {
				out, outClose = create(t, f+"-RU.po")
				if err = result.WriteFile(f, out); err != nil {
					t.Fatal(err)
				}
				outClose()
			}
		}
		if len(result) == 0 {
			t.Error("unexpected blank result")
		}
	})
	t.Run("Tips", func(t *testing.T) {
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
		result, err := Gen(engRaw, rusRaw)
		if err != nil {
			t.Fatal(err)
		}
		if update {
			out, outClose := create(t, "tips.json")
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			if err = enc.Encode(result); err != nil {
				t.Fatal(err)
			}
			outClose()
			for _, f := range result.Files() {
				out, outClose = create(t, f+"-RU.po")
				if err = result.WriteFile(f, out); err != nil {
					t.Fatal(err)
				}
				outClose()
			}
		}
		if len(result) == 0 {
			t.Error("unexpected blank result")
		}
		out, outClose := open(t, "tips.json")
		defer outClose()
		dec := json.NewDecoder(out)
		var expectedResult []Entry
		if err = dec.Decode(&expectedResult); err != nil {
			t.Fatal(err)
		}
		if len(expectedResult) != len(result) {
			t.Fatal("unexpected result length")
		}
		for i, r := range expectedResult {
			got := result[i]
			if got != r {
				t.Errorf("%+v (got) != %+v (expected)", got, r)
			}
		}
	})
	t.Run("Default", func(t *testing.T) {
		engF, enfFClose := open(t, "Language", "english.xml")
		defer enfFClose()
		f, fClose := open(t, "Language", "russian.xml")
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
		result, err := Gen(engRaw, rusRaw)
		if err != nil {
			t.Fatal(err)
		}
		if update {
			out, outClose := create(t, "default.json")
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			if err = enc.Encode(result); err != nil {
				t.Fatal(err)
			}
			outClose()
			for _, f := range result.Files() {
				out, outClose = create(t, f+"-RU.po")
				if err = result.WriteFile(f, out); err != nil {
					t.Fatal(err)
				}
				outClose()
			}
		}
		if len(result) == 0 {
			t.Error("unexpected blank result")
		}
		out, outClose := open(t, "default.json")
		defer outClose()
		dec := json.NewDecoder(out)
		var expectedResult []Entry
		if err = dec.Decode(&expectedResult); err != nil {
			t.Fatal(err)
		}
		if len(expectedResult) != len(result) {
			t.Fatal("unexpected result length")
		}
		for i, r := range expectedResult {
			got := result[i]
			if got != r {
				t.Errorf("%+v (got) != %+v (expected)", got, r)
			}
		}

	})
	t.Run("DefaultEnglish", func(t *testing.T) {
		engF, enfFClose := open(t, "Language", "english.xml")
		defer enfFClose()
		f, fClose := open(t, "Language", "english.xml")
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
		result, err := Gen(engRaw, rusRaw)
		if err != nil {
			t.Fatal(err)
		}
		if update {
			out, outClose := create(t, "default-english.json")
			enc := json.NewEncoder(out)
			enc.SetIndent("", "  ")
			if err = enc.Encode(result); err != nil {
				t.Fatal(err)
			}
			outClose()
		}
		if len(result) == 0 {
			t.Error("unexpected blank result")
		}
		out, outClose := open(t, "default-english.json")
		defer outClose()
		dec := json.NewDecoder(out)
		var expectedResult []Entry
		if err = dec.Decode(&expectedResult); err != nil {
			t.Fatal(err)
		}
		if len(expectedResult) != len(result) {
			t.Fatal("unexpected result length")
		}
		for i, r := range expectedResult {
			got := result[i]
			if got != r {
				t.Errorf("%+v (got) != %+v (expected)", got, r)
			}
		}
	})
}
