package git

import (
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/src-d/go-git.v4"
)

// TestNoTags verifies a repository without tags is handled correctly
func TestNoTags(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-semver-git-no-tags")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(dir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/pinterb/semver-test-3",
		Progress: os.Stdout,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	tags, err := Tags(dir)
	if err != nil {
		t.Fatal(err.Error())
	}

	ntags := len(tags)
	if ntags > 0 {
		t.Fatalf("expected no tags, found %d tags", ntags)
	}
}

// TestTags verifies a repository with tags is handled correctly
func TestTags(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-semver-git-tags")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.RemoveAll(dir)

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/pinterb/semver-test-1",
		Progress: os.Stdout,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	tags, err := Tags(dir)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := []string{
		"0.0.2",
		"0.1.0-alpha.0.beta",
		"0.1.0-alpha.01",
		"0.1.1-beta.0",
		"v0.0.1",
		"v0.1.0",
	}

	if len(tags) != len(expected) {
		t.Fatalf("expected %d tags, found %d tags", len(expected), len(tags))
	}

	for i, a := range tags {
		if expected[i] != a {
			t.Fatalf("expected tag value '%s', found tag value '%s'", expected[i], a)
		}
	}
}
