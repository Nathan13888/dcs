package scraper

import "testing"

func TestPing(t *testing.T) {
	validLinks := []string{
		"google.com",
		"duckduckgo.com",
		"amazon.ca",
		"youtube.com",
		"kernel.org",
	}
	invalidLinks := []string{
		"fakewebsite.does.not.exist.com",
		"aaaaaaaaaaaaaaaaaaaaaaaaaaa.com",
		"bbbbbbbbbbbbbbbbbbbbbbbbbbb.com",
	}
	for _, x := range validLinks {
		res := Ping(x)
		if !res {
			t.Errorf("Problem pinging expected link, %s", x)
		}
	}
	for _, y := range invalidLinks {
		res := Ping(y)
		if res {
			t.Errorf("Problem pinging unexpected link, %s", y)
		}
	}
}
