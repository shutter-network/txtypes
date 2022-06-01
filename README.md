# Replacement go-ethereum/core types

This repository defines the modified `go-ethereum/core` types
to be synced accros the various go-ethereum forks of shutter-network
(cannon/minigeth, optimism/l2geth).

It is possible to use the transaction types by importing this respository,
since all necessary base-types are included in the `types/` package.
This is useful e.g. for RLP encoding the Shutter transaction types in
rolling shutter.

## Updating shutter-network/go-ethereum type definitions

With the help of the purpose built `shtypecopy` tool,
it is easy to copy a subset of modified files 
to the target directory, and eventually replace the 
original packages import paths.

Workflow how to update go-ethereum type definitions accross code-bases:

#### 0 - Build the `shtypecopy` tool
```sh
git clone git@github.com:shutter-network/txtypes.git
cd txtypes/shtypecopy
go build
```
Now you can use the binary `./shtypecopy` for the instructions below.

#### 1 - Make the changes to the types in the geth fork

After cloning https://github.com/shutter-network/go-ethereum/ do:
```sh
cd go-ethereum/core/types
git checkout shutter-types
echo "make your changes, for example:"
echo "// changes" >> transaction.go
git add . -u
git commit -m "Update type definition"
```
And push and merge upstream.

#### 2 - Sync code in the `cannon` repo
After cloning https://github.com/shutter-network/cannon do:
```sh
cd cannon/minigeth/core/types
./shtypecopy --in core/types --out . --rules rules.shtypecopy
git add . -u
git commit -m "Sync types with github.com/shutter-network/go-ethereum"
```
And push and merge upstream.

#### 3 - Sync code in the `txtypes` repo
After cloning https://github.com/shutter-network/txtypes do:
```sh
cd txtypes/types
./shtypecopy --in core/types  --out . --rules rules.shtypecopy
git add . -u
git commit -m "Sync types with github.com/shutter-network/go-ethereum"
```

#### 4 - Push a new semver tag
```sh
git tag v0.0.42
git push --tags
```
And push and merge upstream.

#### 5 - Update `txtypes` version in `rolling-shutter` repo
You should wait at least 1 minute after pushing the new version in order to not poison the `proxy.golang.org` cache:

> From https://proxy.golang.org/ - 
The new version should be available within one minute. Note that if someone requested the version before the tag was pushed, it may take up to 30 minutes for the mirror's cache to expire and fresh data about the version to become available.

After cloning https://github.com/shutter-network/rolling-shutter do:
```sh
cd rolling-shutter/rolling-shutter
go get github.com/shutter-network/txtypes@v0.0.42
git add . -u
git commit -m "Update txtypes version to v0.0.42"
```
