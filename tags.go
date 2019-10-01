package main

import (
	"github.com/integralist/go-findroot/find"

	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func tags() ([]string, error) {
	var stags []string

	root, err := find.Repo()
	if err != nil {
		return stags, err
	}

	// {Name:go-semver Path:/Users/M/Projects/golang/src/github.com/pinterb/go-semver}

	// We instanciate a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(root.Path)
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
