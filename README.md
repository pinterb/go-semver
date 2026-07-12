## ⚠️ IMPORTANT LICENSE CHANGE NOTICE
As of July 2026, this software is **no longer licensed under the MIT License**. 

This software is now distributed under a **Custom Source-Available License** (see `CUSTOM_LICENSE.md`). 
* **General Users:** The software remains free to use, modify, and distribute for the general public.
* **Excluded Entities:** CDW Corporation, its subsidiaries (including CDW LLC and CDW Government LLC), affiliates, and employees are **strictly prohibited** from downloading, installing, running, or utilizing this software. 

Any use by CDW Corporation or its affiliates under this version constitutes immediate copyright infringement.

---

## Superseded by semver-trust-go

**go-semver is superseded by
[`semver-trust-go`](https://github.com/semver-trust/semver-trust-go)**, the
official [SemVer-Trust](https://github.com/semver-trust/spec) reference
implementation. This repository is no longer maintained and will be archived
shortly after this notice lands. Please migrate.

### Migration

| go-semver | semver-trust |
|---|---|
| `semver <versions>` | `semver-trust list` (lenient parity incl. coercion) |
| `semver -l` | `semver-trust latest` |
| `semver -i <level> --preid X` | `semver-trust next -i <level> --preid X` |
| `semver -r -i -d` | `semver-trust next -i -d` (tagless bootstrap `0.0.1`) |
| `semver … tag` behavior | `semver-trust tag` |

### Install semver-trust

```sh
go install github.com/semver-trust/semver-trust-go/cmd/semver-trust@latest
```

Prebuilt binaries for Linux, macOS, and Windows (amd64/arm64) — with
checksums, a keyless [cosign](https://docs.sigstore.dev/) signature, and
per-archive SBOMs — are attached to each
[GitHub Release](https://github.com/semver-trust/semver-trust-go/releases).

### License history

This repository has carried a [custom source-available
license](CUSTOM_LICENSE.md) since July 2026 (see the notice above).
`semver-trust-go` is licensed under [Apache
2.0](https://github.com/semver-trust/semver-trust-go/blob/main/LICENSE).
This notice does not modify go-semver's license files.

---

# go-semver

A tool to manage semantic versioning of software

### Synopsis

The primary goal of semver is to make working with semantic versions
easier. Currently, its two primary functions are to a) validate lists
of raw versions; and b) increment a version.

A secondary goal is work well with git repositories. So while you may
pass one or more versions (as arguments) to semver, you can just as
easily use the tags from a local git repository. So semver can validate
git repository tags and perhaps most importantly, it can help manage
git tags by providing a clean interface for incrementing a current tag
to a valid next version.

### Examples
Validate a list of versions (where some versions are malformed and others are
invalid):

```
root@laptop:~/some-dir$ semver 2.1 v1.0.1 v3 4.x 5.12
1.0.1 2.1.0 3.0.0 5.12.0
```

Increment most current, valid version to a pre-patch version (where prefix
identifier is specified): 

```
root@laptop:~/some-dir$ semver 2.1 v1.0.1 v3 4.x 5.12 -i=prepatch --preid=rc
5.12.1-rc.0
```

Increment a valid pre-release version:

```
root@laptop:~/some-dir$ semver 5.12.1-rc.0 -i=prerelease --preid=rc
5.12.1-rc.1
```

Increment version on a git repository with no tags (where default increment is
patch and default version is "0.0.0"):

```
root@laptop:~/some-repo$ semver -r -i -d
0.0.1
```

### Options

```
  -i, --increment string[="patch"]                        Increment a valid version by the specified level. Level can
                                                          be one of: major, minor, patch, premajor, preminor, prepatch,
                                                          or prerelease. If more than one version is provided, then
                                                          the most current version is incremented.

      --preid string                                      Identifier to be used to prefix premajor, preminor,
                                                          prepatch or prerelease version increments.

  -r, --repo-dir string[="/current/working/directory"]    Use tags from a local git repo as source of versions.

  -d, --default string[="0.0.0"]                          Default version to use when no valid versions are provided

  -l, --latest-only                                       Only return the latest version

  -h, --help                                              Help for semver
```

### Inspirational/Interesting Links
* [Git Tags and Semantic Versioning](http://www.tugberkugurlu.com/archive/versioning-software-builds-based-on-git-tags-and-semantic-versioning-semver)
* [node-semver](https://github.com/npm/node-semver)
* [masterminds semver](https://github.com/Masterminds/semver)
* [semver-howto](https://github.com/dbrock/semver-howto)
