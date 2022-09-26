package semver

import (
	"fmt"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func TestReleaseTypes(t *testing.T) {
	tests := []struct {
		rt       ReleaseType
		expected string
	}{
		{major, "major"},
		{minor, "minor"},
		{patch, "patch"},
		{preMajor, "premajor"},
		{preMinor, "preminor"},
		{prePatch, "prepatch"},
		{preRelease, "prerelease"},
		{pre, "pre"},
	}

	for _, tc := range tests {
		sv := fmt.Sprintf("%s", tc.rt)
		if sv != tc.expected {
			t.Fatalf("expected release type of '%s', but instead got '%s'", tc.expected, sv)
		}
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		version string
		err     bool
	}{
		{"1.2.3", false},
		{"1.2.3-alpha.01", true},
		{"1.2.3+test.01", false},
		{"1.2.3-alpha.-1", false},
		{"v1.2.3", false},
		{"1.0", false},
		{"v1.0", false},
		{"1", false},
		{"v1", false},
		{"1.2.beta", true},
		{"v1.2.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"v1.2-5", false},
		{"1.2-beta.5", false},
		{"v1.2-beta.5", false},
		{"\n1.2", true},
		{"\nv1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"v1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hypen", false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", false},
		{"1.2.3-rc1-with-hypen", false},
		{"v1.2.3-rc1-with-hypen", false},
		{"1.2.3.4", true},
		{"v1.2.3.4", true},
		{"1.2.2147483648", false},
		{"1.2147483648.3", false},
		{"2147483648.3.0", false},
	}

	for _, tc := range tests {
		_, err := Valid(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}
	}
}

func TestIncrement(t *testing.T) {

	tests := []struct {
		v1              string
		rt              ReleaseType
		ident           string
		expectedVersion string
		expectedErr     error
	}{

		{"1.2.3", major, "", "2.0.0", nil},
		{"1.2.3", minor, "", "1.3.0", nil},
		{"1.2.3", patch, "", "1.2.4", nil},
		{"1.2.3tag", major, "", "", semver.ErrInvalidSemVer},
		{"1.2.3-tag", major, "", "2.0.0", nil},
		{"1.2.0-0", patch, "", "1.2.0", nil},
		{"1.2.3-4", major, "", "2.0.0", nil},
		{"1.2.3-4", minor, "", "1.3.0", nil},
		{"1.2.3-4", patch, "", "1.2.3", nil},
		{"1.2.3-alpha.0.beta", major, "", "2.0.0", nil},
		{"1.2.3-alpha.0.beta", minor, "", "1.3.0", nil},
		{"1.2.3-alpha.0.beta", patch, "", "1.2.3", nil},
		{"1.2.4", preRelease, "", "1.2.5-0", nil},
		{"1.2.3-0", preRelease, "", "1.2.3-1", nil},
		{"1.2.3-alpha.0", preRelease, "", "1.2.3-alpha.1", nil},
		{"1.2.3-alpha.1", preRelease, "", "1.2.3-alpha.2", nil},
		{"1.2.3-alpha.2", preRelease, "", "1.2.3-alpha.3", nil},
		{"1.2.3-alpha.0.beta", preRelease, "", "1.2.3-alpha.1.beta", nil},
		{"1.2.3-alpha.1.beta", preRelease, "", "1.2.3-alpha.2.beta", nil},
		{"1.2.3-alpha.2.beta", preRelease, "", "1.2.3-alpha.3.beta", nil},
		{"1.2.3-alpha.10.0.beta", preRelease, "", "1.2.3-alpha.10.1.beta", nil},
		{"1.2.3-alpha.10.1.beta", preRelease, "", "1.2.3-alpha.10.2.beta", nil},
		{"1.2.3-alpha.10.2.beta", preRelease, "", "1.2.3-alpha.10.3.beta", nil},
		{"1.2.3-alpha.10.beta.0", preRelease, "", "1.2.3-alpha.10.beta.1", nil},
		{"1.2.3-alpha.10.beta.1", preRelease, "", "1.2.3-alpha.10.beta.2", nil},
		{"1.2.3-alpha.10.beta.2", preRelease, "", "1.2.3-alpha.10.beta.3", nil},
		{"1.2.3-alpha.9.beta", preRelease, "", "1.2.3-alpha.10.beta", nil},
		{"1.2.3-alpha.10.beta", preRelease, "", "1.2.3-alpha.11.beta", nil},
		{"1.2.3-alpha.11.beta", preRelease, "", "1.2.3-alpha.12.beta", nil},
		{"1.2.0", prePatch, "", "1.2.1-0", nil},
		{"1.2.0-1", prePatch, "", "1.2.1-0", nil},
		{"1.2.0", preMinor, "", "1.3.0-0", nil},
		{"1.2.3-1", preMinor, "", "1.3.0-0", nil},
		{"1.2.0", preMajor, "", "2.0.0-0", nil},
		{"1.2.3-1", preMajor, "", "2.0.0-0", nil},
		{"1.2.0-1", minor, "", "1.2.0", nil},
		{"1.0.0-1", major, "", "1.0.0", nil},

		{"1.2.3", major, "", "2.0.0", nil},
		{"1.2.3", minor, "", "1.3.0", nil},
		{"1.2.3", patch, "", "1.2.4", nil},
		{"1.2.3tag", major, "", "", semver.ErrInvalidSemVer},
		{"1.2.3-tag", major, "", "2.0.0", nil},
		{"1.2.0-0", patch, "", "1.2.0", nil},
		{"fake", major, "", "", semver.ErrInvalidSemVer},
		{"1.2.3-4", major, "", "2.0.0", nil},
		{"1.2.3-4", minor, "", "1.3.0", nil},
		{"1.2.3-4", patch, "", "1.2.3", nil},
		{"1.2.3-alpha.0.beta", major, "", "2.0.0", nil},
		{"1.2.3-alpha.0.beta", minor, "", "1.3.0", nil},
		{"1.2.3-alpha.0.beta", patch, "", "1.2.3", nil},
		{"1.2.4", preRelease, "dev", "1.2.5-dev.0", nil},
		{"1.2.3-0", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.0", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.0", preRelease, "", "1.2.3-alpha.1", nil},
		{"1.2.3-alpha.0.beta", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.0.beta", preRelease, "", "1.2.3-alpha.1.beta", nil},
		{"1.2.3-alpha.10.0.beta", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.10.0.beta", preRelease, "", "1.2.3-alpha.10.1.beta", nil},
		{"1.2.3-alpha.10.1.beta", preRelease, "", "1.2.3-alpha.10.2.beta", nil},
		{"1.2.3-alpha.10.2.beta", preRelease, "", "1.2.3-alpha.10.3.beta", nil},
		{"1.2.3-alpha.10.beta.0", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.10.beta.0", preRelease, "", "1.2.3-alpha.10.beta.1", nil},
		{"1.2.3-alpha.10.beta.1", preRelease, "", "1.2.3-alpha.10.beta.2", nil},
		{"1.2.3-alpha.10.beta.2", preRelease, "", "1.2.3-alpha.10.beta.3", nil},
		{"1.2.3-alpha.9.beta", preRelease, "dev", "1.2.3-dev.0", nil},
		{"1.2.3-alpha.9.beta", preRelease, "", "1.2.3-alpha.10.beta", nil},
		{"1.2.3-alpha.10.beta", preRelease, "", "1.2.3-alpha.11.beta", nil},
		{"1.2.3-alpha.11.beta", preRelease, "", "1.2.3-alpha.12.beta", nil},
		{"1.2.0", prePatch, "dev", "1.2.1-dev.0", nil},
		{"1.2.0-1", prePatch, "dev", "1.2.1-dev.0", nil},
		{"1.2.0", preMinor, "dev", "1.3.0-dev.0", nil},
		{"1.2.3-1", preMinor, "dev", "1.3.0-dev.0", nil},
		{"1.2.0", preMajor, "dev", "2.0.0-dev.0", nil},
		{"1.2.3-1", preMajor, "dev", "2.0.0-dev.0", nil},
		{"1.2.0-1", minor, "dev", "1.2.0", nil},
		{"1.0.0-1", major, "dev", "1.0.0", nil},
		{"1.2.3-dev.bar", preRelease, "dev", "1.2.3-dev.0", nil},
	}

	for _, tc := range tests {
		v1, err := Increment(tc.v1, tc.rt, tc.ident)
		if err != nil && tc.expectedErr == nil {
			t.Errorf("For %s with rt of %s, expected version string=%q, encountered unexpected error: %s", tc.v1, tc.rt, tc.expectedVersion, err.Error())
		}
		if err == nil && tc.expectedErr != nil {
			t.Errorf("For %s with rt of %s, expected version string=%q, did not encounter expected error: %s", tc.v1, tc.rt, tc.expectedVersion, tc.expectedErr.Error())
		}

		if err != nil {
			if err.Error() != tc.expectedErr.Error() {
				t.Errorf("For %s, expected to get err=%s, but got err=%s", tc.v1, tc.expectedErr, err)
			}
		}

		if v1 != tc.expectedVersion {
			t.Errorf("For %s with rt of %s, expected version string=%q, but got %q", tc.v1, tc.rt, tc.expectedVersion, v1)
		}
	}
}

func TestMajor(t *testing.T) {
	tests := []struct {
		version string
		major   uint64
		err     bool
	}{
		{"1.2.3", 1, false},
		{"1.2.3-alpha.01", 0, true},
		{"1.2.3+test.01", 1, false},
		{"1.2.3-alpha.-1", 1, false},
		{"v1.2.3", 1, false},
		{"1.0", 1, false},
		{"v1.0", 1, false},
		{"1", 1, false},
		{"v1", 1, false},
		{"1.2.beta", 0, true},
		{"v1.2.beta", 0, true},
		{"foo", 0, true},
		{"1.2-5", 1, false},
		{"v1.2-5", 1, false},
		{"1.2-beta.5", 1, false},
		{"v1.2-beta.5", 1, false},
		{"\n1.2", 0, true},
		{"\nv1.2", 0, true},
		{"1.2.0-x.Y.0+metadata", 1, false},
		{"v1.2.0-x.Y.0+metadata", 1, false},
		{"1.2.0-x.Y.0+metadata-width-hypen", 1, false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", 1, false},
		{"1.2.3-rc1-with-hypen", 1, false},
		{"v1.2.3-rc1-with-hypen", 1, false},
		{"1.2.3.4", 0, true},
		{"v1.2.3.4", 0, true},
		{"1.2.2147483648", 1, false},
		{"1.2147483648.3", 1, false},
		{"2147483648.3.0", 2147483648, false},
	}

	for _, tc := range tests {
		m, err := Major(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		if tc.major != m {
			t.Fatalf("expected major version: %d, but got %d", tc.major, m)
		}
	}
}

func TestMinor(t *testing.T) {
	tests := []struct {
		version string
		minor   uint64
		err     bool
	}{
		{"1.2.3", 2, false},
		{"1.2.3-alpha.01", 0, true},
		{"1.2.3+test.01", 2, false},
		{"1.2.3-alpha.-1", 2, false},
		{"v1.2.3", 2, false},
		{"1.0", 0, false},
		{"v1.0", 0, false},
		{"1", 0, false},
		{"v1", 0, false},
		{"1.2.beta", 0, true},
		{"v1.2.beta", 0, true},
		{"foo", 0, true},
		{"1.2-5", 2, false},
		{"v1.2-5", 2, false},
		{"1.2-beta.5", 2, false},
		{"v1.2-beta.5", 2, false},
		{"\n1.2", 0, true},
		{"\nv1.2", 0, true},
		{"1.2.0-x.Y.0+metadata", 2, false},
		{"v1.2.0-x.Y.0+metadata", 2, false},
		{"1.2.0-x.Y.0+metadata-width-hypen", 2, false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", 2, false},
		{"1.2.3-rc1-with-hypen", 2, false},
		{"v1.2.3-rc1-with-hypen", 2, false},
		{"1.2.3.4", 0, true},
		{"v1.2.3.4", 0, true},
		{"1.2.2147483648", 2, false},
		{"1.2147483648.3", 2147483648, false},
		{"2147483648.3.0", 3, false},
	}

	for _, tc := range tests {
		m, err := Minor(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		if tc.minor != m {
			t.Fatalf("expected minor version: %d, but got %d", tc.minor, m)
		}
	}
}

func TestPatch(t *testing.T) {
	tests := []struct {
		version string
		patch   uint64
		err     bool
	}{
		{"1.2.3", 3, false},
		{"1.2.3-alpha.01", 0, true},
		{"1.2.3+test.01", 3, false},
		{"1.2.3-alpha.-1", 3, false},
		{"v1.2.3", 3, false},
		{"1.0", 0, false},
		{"v1.0", 0, false},
		{"1", 0, false},
		{"v1", 0, false},
		{"1.2.beta", 0, true},
		{"v1.2.beta", 0, true},
		{"foo", 0, true},
		{"1.2-5", 0, false},
		{"v1.2-5", 0, false},
		{"1.2-beta.5", 0, false},
		{"v1.2-beta.5", 0, false},
		{"\n1.2", 0, true},
		{"\nv1.2", 0, true},
		{"1.2.0-x.Y.0+metadata", 0, false},
		{"v1.2.0-x.Y.0+metadata", 0, false},
		{"1.2.0-x.Y.0+metadata-width-hypen", 0, false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", 0, false},
		{"1.2.3-rc1-with-hypen", 3, false},
		{"v1.2.3-rc1-with-hypen", 3, false},
		{"1.2.3.4", 0, true},
		{"v1.2.3.4", 0, true},
		{"1.2.2147483648", 2147483648, false},
		{"1.2147483648.3", 3, false},
		{"2147483648.3.0", 0, false},
	}

	for _, tc := range tests {
		m, err := Patch(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		if tc.patch != m {
			t.Fatalf("expected patch version: %d, but got %d", tc.patch, m)
		}
	}
}

func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestPrerelease(t *testing.T) {
	tests := []struct {
		version    string
		prerelease []string
		err        bool
	}{
		{"1.2.3", nil, false},
		{"1.2.3-alpha.01", nil, true},
		{"1.2.3+test.01", nil, false},
		{"1.2.3-alpha.-1", []string{"alpha", "-1"}, false},
		{"v1.2.3", nil, false},
		{"1.0", nil, false},
		{"v1.0", nil, false},
		{"1", nil, false},
		{"v1", nil, false},
		{"1.2.beta", nil, true},
		{"v1.2.beta", nil, true},
		{"foo", nil, true},
		{"1.2-5", []string{"5"}, false},
		{"v1.2-5", []string{"5"}, false},
		{"1.2-beta.5", []string{"beta", "5"}, false},
		{"v1.2-beta.5", []string{"beta", "5"}, false},
		{"\n1.2", nil, true},
		{"\nv1.2", nil, true},
		{"1.2.0-x.Y.0+metadata", []string{"x", "Y", "0"}, false},
		{"v1.2.0-x.Y.0+metadata", []string{"x", "Y", "0"}, false},
		{"1.2.0-x.Y.0+metadata-width-hypen", []string{"x", "Y", "0"}, false},
		{"v1.2.0-x.Y.0+metadata-width-hypen", []string{"x", "Y", "0"}, false},
		{"1.2.3-rc1-with-hypen", []string{"rc1-with-hypen"}, false},
		{"v1.2.3-rc1-with-hypen", []string{"rc1-with-hypen"}, false},
		{"1.2.3.4", nil, true},
		{"v1.2.3.4", nil, true},
		{"1.2.2147483648", nil, false},
		{"1.2147483648.3", nil, false},
		{"2147483648.3.0", nil, false},
	}

	for _, tc := range tests {
		m, err := Prerelease(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %s", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %s: %s", tc.version, err)
		}

		if tc.prerelease != nil && m == nil {
			t.Fatalf("expected version %s to produce %d prerelease elements, but got nil", tc.version, len(tc.prerelease))
		}

		if tc.prerelease == nil && m != nil {
			t.Fatalf("expected nil prerelease, but got %d prerelease elements", len(m))
		}

		if !Equal(tc.prerelease, m) {
			t.Fatalf("expected version %s to return %v, but got %v", tc.version, tc.prerelease, m)
		}
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		raw      []string
		versions []string
		err      bool
	}{
		{[]string{"1.2.3", "1.2.3-alpha.01"}, []string{"1.2.3"}, false},
		{[]string{"1.2.3-alpha.01"}, nil, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0"}, []string{"1.2.0-5", "1.0.0", "0.2.0", "3.0.1-beta.0"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4"}, []string{"1.2.0-5", "1.0.0", "0.2.0", "3.0.1-beta.0"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4", "3.0.1"}, []string{"1.2.0-5", "1.0.0", "0.2.0", "3.0.1-beta.0", "3.0.1"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "v4", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4", "3.0.1"}, []string{"1.2.0-5", "1.0.0", "4.0.0", "0.2.0", "3.0.1-beta.0", "3.0.1"}, false},
	}

	for _, tc := range tests {
		m, err := List(tc.raw)
		if tc.err && err == nil {
			t.Fatalf("expected error for raw versions: %v", tc.raw)
		} else if !tc.err && err != nil {
			t.Fatalf("error for raw version %v: %s", tc.raw, err)
		}

		if tc.versions != nil && m == nil {
			t.Fatalf("expected list %v to produce %d versions, but got nil", tc.raw, len(tc.versions))
		}

		if tc.versions == nil && m != nil {
			t.Fatalf("expected nil versions, but got %d versions", len(m))
		}

		if !Equal(tc.versions, m) {
			t.Fatalf("expected raw version %v to return %v, but got %v", tc.raw, tc.versions, m)
		}
	}
}

func TestSortedList(t *testing.T) {
	tests := []struct {
		raw      []string
		versions []string
		err      bool
	}{
		{[]string{"1.2.3", "1.2.3-alpha.01"}, []string{"1.2.3"}, false},
		{[]string{"1.2.3-alpha.01"}, nil, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0"}, []string{"0.2.0", "1.0.0", "1.2.0-5", "3.0.1-beta.0"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4"}, []string{"0.2.0", "1.0.0", "1.2.0-5", "3.0.1-beta.0"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4", "3.0.1"}, []string{"0.2.0", "1.0.0", "1.2.0-5", "3.0.1-beta.0", "3.0.1"}, false},
		{[]string{"1.2-5", "1.2.3-alpha.01", "v1", "v4", "0.2.0", "v3.0.1-beta.0", "v1.2.3.4", "3.0.1"}, []string{"0.2.0", "1.0.0", "1.2.0-5", "3.0.1-beta.0", "3.0.1", "4.0.0"}, false},
	}

	for _, tc := range tests {
		m, err := SortedList(tc.raw)
		if tc.err && err == nil {
			t.Fatalf("expected error for raw versions: %v", tc.raw)
		} else if !tc.err && err != nil {
			t.Fatalf("error for raw version %v: %s", tc.raw, err)
		}

		if tc.versions != nil && m == nil {
			t.Fatalf("expected list %v to produce %d versions, but got nil", tc.raw, len(tc.versions))
		}

		if tc.versions == nil && m != nil {
			t.Fatalf("expected nil versions, but got %d versions", len(m))
		}

		if !Equal(tc.versions, m) {
			t.Fatalf("expected raw version %v to return %v, but got %v", tc.raw, tc.versions, m)
		}
	}
}
