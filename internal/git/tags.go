package git

import (
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4"
)

// rootPath isn't really getting the root path. But it does try to make sure that the path specified is a valid directory
func rootPath(path string) (string, error) {
	if path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = dir
	} else {
		fi, err := os.Stat(path)
		if err != nil {
			return "", err
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			return path, nil

		case mode.IsRegular():
			if dir := filepath.Dir(path); dir != path {
				path = dir
			}
		}
	}
	return path, nil
}

// Tags returns a list of tag values from a git repository at a known location
func Tags(path string) ([]string, error) {
	var stags []string
	apath, err := rootPath(path)
	if err != nil {
		return stags, err
	}

	options := &git.PlainOpenOptions{DetectDotGit: true}
	r, err := git.PlainOpenWithOptions(apath, options)
	if err != nil {
		return stags, err
	}

	// all tag references, both lightweight tags and annotated tags
	tags, err := r.Tags()
	if err != nil {
		return stags, err
	}

	stags = make([]string, 0)

	err = tags.ForEach(func(t *plumbing.Reference) error {
		stags = append(stags, t.Name().Short())
		return nil
	})

	return stags, err
}
