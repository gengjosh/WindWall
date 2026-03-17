package rules

import "testing"

func TestExtractHostname(t *testing.T) {
	got := extractHostname("https://cdn.mangasite.com:8080/file.js")
	want := "cdn.mangasite.com"

	if got != want {
		t.Errorf("extractHostname() = %q, want %q", got, want)
	}
}

func TestExtractDomain(t *testing.T) {
	got := extractDomain("cdn.mangasite.com")
	want := "mangasite.com"

	if got != want {
		t.Errorf("extractDomain() = %q, want %q", got, want)
	}
}

func TestDetermineThirdParty_FirstParty(t *testing.T) {
	got := DetermineThirdParty(
		"cdn.mangasite.com",
		"https://mangasite.com/chapter-12",
		"https://mangasite.com",
	)
	want := false

	if got != want {
		t.Errorf("DetermineThirdParty() = %v, want %v", got, want)
	}
}

func TestDetermineThirdParty_ThirdParty(t *testing.T) {
	got := DetermineThirdParty(
		"ads.evilnetwork.com",
		"https://mangasite.com/chapter-12",
		"https://mangasite.com",
	)
	want := true

	if got != want {
		t.Errorf("DetermineThirdParty() = %v, want %v", got, want)
	}
}
