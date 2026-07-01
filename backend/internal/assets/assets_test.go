package assets

import (
	"io"
	"testing"
)

func TestDist_ReturnsNonNilFS(t *testing.T) {
	dist := Dist()
	if dist == nil {
		t.Fatal("Dist() returned nil")
	}
}

func TestDist_ContainsIndexHTML(t *testing.T) {
	dist := Dist()
	f, err := dist.Open("index.html")
	if err != nil {
		t.Fatalf("expected index.html to be present: %v", err)
	}
	defer func() { _ = f.Close() }()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to read index.html: %v", err)
	}

	if len(content) == 0 {
		t.Fatal("index.html is empty")
	}
}
