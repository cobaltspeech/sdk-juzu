# SDK for Juzu (Cobalt's Speech Diarization Engine)

This repository contains the SDK for Cobalt's Juzu Speech Diarization Engine.

## Network API (using GRPC)

The `grpc` folder at the top level of this repository contains code for Juzu's
GRPC API.  The `grpc/juzu.proto` file is the authoritative service definition of
the API and is used for auto generating SDK code in multiple languages.

### Auto-generated code and documentation

The `grpc` folder contains auto-generated code in several languages.  In order
to generate the code again, you should run `make`.  Generated code is checked
in, and you must make sure it is up to date when you push commits to this
repository.

Code generation has the following dependencies:
  - The protobuf compiler itself (protoc)
  - The protobuf documentation generation plugin (protoc-gen-doc)
  - The golang plugins (protoc-gen-go and protoc-gen-grpc-gateway)

A few system dependencies are required:
  - Go 1.12
  - git
  - wget

The top level Makefile can set up all other dependencies.

To generate the code and documentation, run `make`.  

### Tagging New Versions

This repository has several components, and they need more than just a "vX.Y.Z"
tag on the git repo.  In particular, this repository has two go modules, one of
which depends on the other, and in order to make sure correct versions are used,
we need to follow a few careful steps to release new versions on this
repository.

Step 1: Make sure all generated code and documentation is up to date.

```
pushd grpc && make && popd
```

Step 2: Update the version number.

In addition to the git tags, we also save the version string in a few places in
our sources.  These strings should all be updated and a new commit created.  The
git tags should then be placed on that commit once merged to master.

Decide which version you'd like to tag. For this README, let's say the next
version to tag is `1.0.1`.

Step 3: Add version tags to the sources.

```
NEW_VERSION="1.0.1"

git checkout master
git checkout -b version-update-v$NEW_VERSION

sed -i 's|grpc/go-juzu v[0-9.]*|grpc/go-juzu v'$NEW_VERSION'|g' grpc/go-juzu/juzupb/gw/go.mod
sed -i 's|<Version>[0-9.]*</Version>|<Version>'$NEW_VERSION'</Version>|g' grpc/csharp-juzu/juzu.csproj
sed -i 's|^VERSION="[0-9.]*"|VERSION="'$NEW_VERSION'"|g' grpc/Makefile

git commit -m "Update version to v$NEW_VERSION"
git push origin version-update-v$NEW_VERSION
```

Step 4: Create a pull request and get changes merged to master.

Step 5: Create version tags on the latest master branch:

```
git checkout master
git pull origin master
git tag -a v$NEW_VERSION -m ''
git tag -a grpc/go-juzu/v$NEW_VERSION -m ''
git tag -a grpc/go-juzu/juzupb/gw/v$NEW_VERSION -m ''
git push origin --tags
```
