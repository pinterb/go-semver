package git

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

	tags, err := r.TagObjects()
	if err != nil {
		return stags, err
	}

	stags = make([]string, 0)

	err = tags.ForEach(func(t *object.Tag) error {
		fmt.Println(t)
		stags = append(stags, t.Name)
		return nil
	})

	return stags, nil
}
