package functions

import (
	"os"
	"strings"
	"testing"
)

func buildBannerFile(t *testing.T) string {
	t.Helper()

	var sb strings.Builder
	sb.WriteString("\n") 

	for i := 0; i < 95; i++ {
		for row := 0; row < 8; row++ {
			sb.WriteString(formatRow(i, row))
			sb.WriteString("\n")
		}
		if i < 94 {
			sb.WriteString("\n") 
		}
	}

	f, err := os.CreateTemp(t.TempDir(), "banner_*.txt")
	if err != nil {
		t.Fatalf("could not create temp banner file: %v", err)
	}
	if _, err := f.WriteString(sb.String()); err != nil {
		t.Fatalf("could not write banner file: %v", err)
	}
	f.Close()
	return f.Name()
}

func formatRow(i, r int) string {
	return strings.Replace(
		strings.Replace("CH__R_",
			"__", padTwo(i), 1),
		"_", padOne(r), 1)
}

func padTwo(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	tens := n / 10
	ones := n % 10
	return string(rune('0'+tens)) + string(rune('0'+ones))
}

func padOne(n int) string { return string(rune('0' + n)) }



func TestLoadBanner_ValidFile(t *testing.T) {
	path := buildBannerFile(t)
	chars, err := LoadBanner(path)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(chars) != 95 {
		t.Fatalf("expected 95 chars, got %d", len(chars))
	}
}

func TestLoadBanner_AllCharsMapped(t *testing.T) {
	path := buildBannerFile(t)
	chars, err := LoadBanner(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i := 0; i < 95; i++ {
		ch := rune(32 + i)
		art, ok := chars[ch]
		if !ok {
			t.Errorf("char %q (ASCII %d) missing from map", ch, 32+i)
			continue
		}
		if len(art) != 8 {
			t.Errorf("char %q: expected 8 art rows, got %d", ch, len(art))
		}
	}
}

func TestLoadBanner_CorrectArtContent(t *testing.T) {
	path := buildBannerFile(t)
	chars, _ := LoadBanner(path)

	spaceArt := chars[' ']
	for row := 0; row < 8; row++ {
		want := formatRow(0, row)
		if spaceArt[row] != want {
			t.Errorf("space row %d: want %q, got %q", row, want, spaceArt[row])
		}
	}

	
	aArt := chars['A']
	for row := 0; row < 8; row++ {
		want := formatRow(33, row)
		if aArt[row] != want {
			t.Errorf("'A' row %d: want %q, got %q", row, want, aArt[row])
		}
	}
}

func TestLoadBanner_FileNotFound(t *testing.T) {
	_, err := LoadBanner("/nonexistent/path/banner.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadBanner_EmptyFile(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "empty_*.txt")
	f.Close()

	_, err := LoadBanner(f.Name())
	if err == nil {
		t.Fatal("expected error for empty file, got nil")
	}
	if err.Error() != "empty banner file" {
		t.Errorf("expected 'empty banner file', got %q", err.Error())
	}
}

func TestLoadBanner_IncompleteTooFewLines(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "short_*.txt")
	f.WriteString("only a few lines\nline2\nline3\n")
	f.Close()

	_, err := LoadBanner(f.Name())
	if err == nil {
		t.Fatal("expected error for incomplete file, got nil")
	}
	if err.Error() != "incomplete banner file" {
		t.Errorf("expected 'incomplete banner file', got %q", err.Error())
	}
}

func TestLoadBanner_WindowsLineEndings(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("\r\n")
	for i := 0; i < 95; i++ {
		for row := 0; row < 8; row++ {
			sb.WriteString(formatRow(i, row))
			sb.WriteString("\r\n")
		}
		if i < 94 {
			sb.WriteString("\r\n")
		}
	}

	f, _ := os.CreateTemp(t.TempDir(), "crlf_*.txt")
	f.WriteString(sb.String())
	f.Close()

	chars, err := LoadBanner(f.Name())
	if err != nil {
		t.Fatalf("expected no error with CRLF file, got: %v", err)
	}
	if len(chars) != 95 {
		t.Fatalf("expected 95 chars, got %d", len(chars))
	}
}


func makeChars(runes ...rune) map[rune][]string {
	m := make(map[rune][]string)
	for _, r := range runes {
		art := make([]string, 8)
		for row := 0; row < 8; row++ {
			art[row] = string(r) + "_row" + string(rune('0'+row))
		}
		m[r] = art
	}
	return m
}

func TestRender_EmptyInput(t *testing.T) {
	chars := makeChars('A')
	got := Render("", chars)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRender_LiteralBackslashN(t *testing.T) {
	chars := makeChars('A')
	got := Render(`\n`, chars)
	if got != "" {
		t.Errorf("expected empty string for literal \\n input, got %q", got)
	}
}

func TestRender_SingleChar(t *testing.T) {
	chars := makeChars('A')
	got := Render("A", chars)

	lines := strings.Split(got, "\n")

	if len(lines) != 9 {
		t.Fatalf("expected 9 parts after split, got %d", len(lines))
	}
	if lines[8] != "" {
		t.Errorf("expected trailing empty string after last newline, got %q", lines[8])
	}
	for row := 0; row < 8; row++ {
		want := "A_row" + string(rune('0'+row))
		if lines[row] != want {
			t.Errorf("row %d: want %q, got %q", row, want, lines[row])
		}
	}
}

func TestRender_MultipleCharsOnOneLine(t *testing.T) {
	chars := makeChars('A', 'B')
	got := Render("AB", chars)

	lines := strings.Split(got, "\n")
	if len(lines) != 9 {
		t.Fatalf("expected 9 parts, got %d", len(lines))
	}
	for row := 0; row < 8; row++ {
		want := "A_row" + string(rune('0'+row)) + "B_row" + string(rune('0'+row))
		if lines[row] != want {
			t.Errorf("row %d: want %q, got %q", row, want, lines[row])
		}
	}
}

func TestRender_NewlineSeparator(t *testing.T) {
	chars := makeChars('A', 'B')
	
	got := Render(`A\nB`, chars)

	lines := strings.Split(got, "\n")
	if len(lines) != 17 {
		t.Fatalf("expected 17 parts, got %d: %v", len(lines), lines)
	}
	for row := 0; row < 8; row++ {
		wantA := "A_row" + string(rune('0'+row))
		if lines[row] != wantA {
			t.Errorf("A row %d: want %q, got %q", row, wantA, lines[row])
		}
	}
	for row := 0; row < 8; row++ {
		wantB := "B_row" + string(rune('0'+row))
		if lines[8+row] != wantB {
			t.Errorf("B row %d: want %q, got %q", row, wantB, lines[8+row])
		}
	}
}

func TestRender_ConsecutiveNewlines(t *testing.T) {
	chars := makeChars('A')
	got := Render(`A\n\nA`, chars)

	lines := strings.Split(got, "\n")
	
	if len(lines) != 18 {
		t.Fatalf("expected 18 parts, got %d: %v", len(lines), lines)
	}
	if lines[8] != "" {
		t.Errorf("expected blank line at index 8, got %q", lines[8])
	}
}

func TestRender_UnknownCharSkipped(t *testing.T) {
	chars := makeChars('A') 
	got := Render("A?", chars)

	lines := strings.Split(got, "\n")
	if len(lines) != 9 {
		t.Fatalf("expected 9 parts, got %d", len(lines))
	}
	for row := 0; row < 8; row++ {
		want := "A_row" + string(rune('0'+row))
		if lines[row] != want {
			t.Errorf("row %d: want %q, got %q", row, want, lines[row])
		}
	}
}

func TestRender_OutputEndsWithNewline(t *testing.T) {
	chars := makeChars('Z')
	got := Render("Z", chars)
	if !strings.HasSuffix(got, "\n") {
		t.Errorf("render output should end with newline, got %q", got)
	}
}

func TestRender_IntegrationWithLoadBanner(t *testing.T) {
	path := buildBannerFile(t)
	chars, err := LoadBanner(path)
	if err != nil {
		t.Fatalf("LoadBanner failed: %v", err)
	}

	got := Render(" ", chars)
	lines := strings.Split(got, "\n")
	if len(lines) < 1 {
		t.Fatal("expected at least 1 line")
	}
	want := formatRow(0, 0) 
	if lines[0] != want {
		t.Errorf("integration: space row 0: want %q, got %q", want, lines[0])
	}
}