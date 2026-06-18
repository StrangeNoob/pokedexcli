package tui

import "testing"

func TestArtStoreRequestDedup(t *testing.T) {
	d := testDeps()
	s := NewArtStore()

	if s.request(d, "pikachu") == nil {
		t.Fatal("first request should return a command")
	}
	if s.request(d, "pikachu") != nil {
		t.Fatal("second request while pending should return nil")
	}

	s.handle(artLoadedMsg{name: "pikachu", art: "ART"})
	if s.get("pikachu") != "ART" {
		t.Fatalf("get = %q, want ART", s.get("pikachu"))
	}
	if s.request(d, "pikachu") != nil {
		t.Fatal("request after cached should return nil")
	}
}

func TestArtStoreEmptyName(t *testing.T) {
	if NewArtStore().request(testDeps(), "") != nil {
		t.Fatal("empty name should return nil")
	}
}
