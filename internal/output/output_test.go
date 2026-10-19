package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestTableAlignment(t *testing.T) {
	var out bytes.Buffer
	p := &Printer{Format: Table, Out: &out, Err: &out}
	items := []map[string]any{
		{"id": "1", "name": "élégant", "price": 4.5, "hidden": "x"},
		{"id": "22", "name": "b", "price": nil},
	}
	err := p.List(items, items, []Column{
		{Header: "ID", Path: "id"},
		{Header: "Nom", Path: "name"},
		{Header: "Prix", Path: "price", Format: Money("EUR")},
		{Header: "Caché", Path: "hidden", Wide: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := "ID   NOM       PRIX\n" +
		"1    élégant   4,50 €\n" +
		"22   b         -\n"
	if out.String() != want {
		t.Fatalf("table :\n%q\nattendu :\n%q", out.String(), want)
	}
	out.Reset()
	p.Format = Wide
	_ = p.List(items, items, []Column{{Header: "ID", Path: "id"}, {Header: "Caché", Path: "hidden", Wide: true}})
	if !strings.Contains(out.String(), "CACHÉ") {
		t.Fatalf("colonne wide absente : %q", out.String())
	}
}

func TestJSONIsRaw(t *testing.T) {
	var out bytes.Buffer
	p := &Printer{Format: JSON, Out: &out}
	raw := map[string]any{"items": []any{1}, "hasMore": false}
	if err := p.List(raw, nil, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"hasMore": false`) {
		t.Fatalf("json = %s", out.String())
	}
	out.Reset()
	p.Format = YAML
	_ = p.Raw(raw)
	if !strings.Contains(out.String(), "hasMore: false") {
		t.Fatalf("yaml = %s", out.String())
	}
}

func TestObjectSkipsEmpty(t *testing.T) {
	var out bytes.Buffer
	p := &Printer{Format: Table, Out: &out}
	obj := map[string]any{"a": "x", "b": nil, "address": map[string]any{"ip": "1.2.3.4"}}
	_ = p.Object(obj, obj, []Field{{Header: "A", Path: "a"}, {Header: "B", Path: "b"}, {Header: "IP", Path: "address.ip"}})
	if out.String() != "A    x\nIP   1.2.3.4\n" {
		t.Fatalf("objet : %q", out.String())
	}
}

func TestFormats(t *testing.T) {
	cases := [][2]string{
		{Text(true), "oui"},
		{Text(3.0), "3"},
		{Text(2.5), "2.50"},
		{Text([]any{"a", "b"}), "a, b"},
		{MB(2048.0), "2.0 Go"},
		{MB(512.0), "512 Mo"},
		{Bytes(1536.0), "1.5 Kio"},
		{Percent(12.34), "12.3 %"},
		{Seconds(90061.0), "1j 1h"},
		{Money("EUR")(1234.5), "1234,50 €"},
		{Money("USD")(1.0), "1,00 USD"},
		{Date("not a date"), "not a date"},
		{truncate("abcdefghij", 5), "abcd…"},
	}
	for _, c := range cases {
		if got, want := c[0], c[1]; got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}
