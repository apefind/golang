package shellutil

import (
	"testing"
)

type testPrompt struct {
	path   string
	length int
	ratio  float64
	prompt string
}

func TestPrompt(t *testing.T) {
	t.Log("testing prompt functionality")
	R := []testPrompt{
		{"/Users/test/Movies/Millenium", 32, 0.75, "/Users/test/Movies/Millenium"},
		{"/home/www.github.com/apefind/golang/u/work/go/src/www.github.com/apefind/golang", 32, 0.75, "/home/a...u/work/go/src/www.github.com/apefind/golang"},
		{"/home/www.github.com/apefind/golang/u/work/go/src/www.github.com/apefind/golang", 20, 0.75, "/hom.../src/www.github.com/apefind/golang"},
		{"/home/www.github.com/apefind/golang/u/work/go/src/www.github.com/apefind/golang", 20, 0.25, "/home/apefin...find"},
		{"/home/xzy/walter/ist/eine/taube/nuss", 15, 0.50, "/home/...e/nuss"},
		{"/home/xyz/walter/ist/eine/taube/nuss", 10, 0.75, "/.../nuss"},
	}
	for _, r := range R {
		if r.prompt != getPromptPath(r.path, r.length, r.ratio) {
			t.Error("expected", r.prompt, ", but got", getPromptPath(r.path, r.length, r.ratio))
		}
		getPrompt(r.path, r.length, r.ratio)
		GetShellPrompt(r.length, r.ratio)
	}
}
