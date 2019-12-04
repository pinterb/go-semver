package semver

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type ReleaseType int

const (
	major ReleaseType = iota
	minor
	patch
	preMajor
	preMinor
	prePatch
	preRelease
	pre // should not be referenced externally
)

var (
	// ErrInternalOnlyReleaseType is returned when a requested release type is for internal use only
	ErrInternalOnlyReleaseType = errors.New("release type is for internal use only")
	// ErrUnknownReleaseType is returned when a requested release type is unknown
	ErrUnknownReleaseType = errors.New("unknown release type")
)

// String is the string representation of a ReleaseType
func (t ReleaseType) String() string {
	return [...]string{"major", "minor", "patch", "premajor", "preminor", "prepatch", "prerelease", "pre"}[t]
}

// ToReleaseType is a convenience function for getting a valid ReleaseType
func ToReleaseType(rt string) (ReleaseType, error) {
	var rtn ReleaseType
	var err error

	trm := strings.TrimSpace(rt)
	switch strings.ToLower(trm) {
	case "major":
		rtn = major
	case "minor":
		rtn = minor
	case "patch":
		rtn = patch
	case "premajor":
		rtn = preMajor
	case "preminor":
		rtn = preMinor
	case "prepatch":
		rtn = prePatch
	case "prerelease":
		rtn = preRelease
	case "pre":
		err = ErrInternalOnlyReleaseType
	default:
		err = ErrUnknownReleaseType
	}

	return rtn, err
}

// Valid returns a parsed version or an error if it's not valid
func Valid(in string) (string, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return "", err
	}

	return v.String(), nil
}

// Increment returns the version incremented by the release type. This function
// largely mimics the increment logic found in https://github.com/npm/node-semver
func Increment(in string, rt ReleaseType, ident string) (string, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return "", err
	}

	prerelease, _ := Prerelease(in)
	var rtn string

	switch rt {
	case preMajor:
		v2 := v.IncMajor()
		rtn, _ = Increment(v2.String(), pre, ident)

	case preMinor:
		v2 := v.IncMinor()
		rtn, _ = Increment(v2.String(), pre, ident)

	case prePatch:
		v2, err := v.SetPrerelease("")
		if err != nil {
			return "", nil
		}
		v3 := v2.IncPatch()
		rtn, _ = Increment(v3.String(), pre, ident)

	case preRelease:
		// if the input is a non-prerelease version, this acts the same as prepatch
		if len(prerelease) == 0 {
			in, _ = Increment(in, patch, ident)
		}
		rtn, _ = Increment(in, pre, ident)

	case major:
		// if this is a pre-major version, bump up to the same major version.
		// Otherwise, increment major
		// 1.0.0-5 bumps to 1.0.0
		// 1.1.0 bumps to 2.0.0
		if v.Minor() != 0 || v.Patch() != 0 || len(prerelease) == 0 {
			rtn = v.IncMajor().String()
		} else {
			v2, err := semver.NewVersion(fmt.Sprintf("%s.0.0", strconv.Itoa(int(v.Major()))))
			if err != nil {
				return "", err
			}
			rtn = v2.String()
		}

	case minor:
		// If this is a pre-minor version, bump up to the same minor version.
		// Otherwise increment minor.
		// 1.2.0-5 bumps to 1.2.0
		// 1.2.1 bumps to 1.3.0
		if v.Patch() != 0 || len(prerelease) == 0 {
			rtn = v.IncMinor().String()
		} else {
			v2, err := semver.NewVersion(fmt.Sprintf("%s.%s.0", strconv.Itoa(int(v.Major())), strconv.Itoa(int(v.Minor()))))
			if err != nil {
				return "", err
			}
			rtn = v2.String()
		}

	case patch:
		rtn = v.IncPatch().String()

	case pre:
		if len(prerelease) == 0 {
			prerelease = []string{"0"}
		} else {
			var i int
			for i = len(prerelease) - 1; i >= 0; i-- {
				if nv, err := strconv.Atoi(prerelease[i]); err == nil {
					prerelease[i] = strconv.Itoa(nv + 1)
					i = -2
				}
			}

			// didn't increment anything
			if i == -1 {
				prerelease = append(prerelease, "0")
			}
		} // len of prerelease

		if ident != "" {
			// 1.2.0-beta.1 bumps to 1.2.0-beta.2,
			// 1.2.0-beta.fooblz or 1.2.0-beta bumps to 1.2.0-beta.0
			if len(prerelease) > 0 {
				if ident == prerelease[0] {
					if _, err := strconv.Atoi(prerelease[1]); err != nil {
						prerelease = []string{ident, "0"}
					}
				} else {
					prerelease = []string{ident, "0"}
				}
			} else {
				prerelease = []string{ident, "0"}
			}
		} // ident was set

		if len(prerelease) > 0 {
			nv, _ := v.SetPrerelease(strings.Join(prerelease, "."))
			rtn = nv.String()
		}
	}

	return rtn, nil
}

// Major returns the major version number
func Major(in string) (uint64, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return 0, err
	}

	return v.Major(), nil
}

// Minor returns the minor version number
func Minor(in string) (uint64, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return 0, err
	}

	return v.Minor(), nil
}

// Patch returns the patch version number
func Patch(in string) (uint64, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return 0, err
	}

	return v.Patch(), nil
}

// Prerelease returns an array of prerelease components or nil if none exist
func Prerelease(in string) ([]string, error) {
	v, err := semver.NewVersion(in)
	if err != nil {
		return nil, err
	}

	var eparts []string
	pre := v.Prerelease()
	if len(pre) > 0 {
		eparts = strings.Split(pre, ".")
	}
	return eparts, nil
}

// List takes a collection of raw version values and returns a list of valid versions
func list(in []string) []*semver.Version {
	vs := make([]*semver.Version, 0)
	for _, r := range in {
		v, err := semver.NewVersion(r)
		if err != nil {
			continue
		}
		vs = append(vs, v)
	}
	return vs
}

// List takes a collection of raw version values and returns a list of valid versions
func List(in []string) ([]string, error) {
	vs := list(in)
	if len(vs) == 0 {
		return nil, nil
	}

	r := make([]string, len(vs))
	for i, v := range vs {
		r[i] = v.String()
	}
	return r, nil
}

// SortedList takes a collection of raw version values and returns a sorted list of valid versions
func SortedList(in []string) ([]string, error) {
	vs := list(in)
	if len(vs) == 0 {
		return nil, nil
	}

	sort.Sort(semver.Collection(vs))
	r := make([]string, len(vs))
	for i, v := range vs {
		r[i] = v.String()
	}
	return r, nil
}
