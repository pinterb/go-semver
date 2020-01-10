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

  -h, --help                                              Help for semver
```

### Inspirational/Interesting Links
* [Git Tags and Semantic Versioning](http://www.tugberkugurlu.com/archive/versioning-software-builds-based-on-git-tags-and-semantic-versioning-semver)
* [node-semver](https://github.com/npm/node-semver)
* [masterminds semver](https://github.com/Masterminds/semver)
* [semver-howto](https://github.com/dbrock/semver-howto)
