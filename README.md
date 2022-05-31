# Replacement go-ethereum/core types

This repository defines the modified `go-ethereum/core` types
to be synced accros the various go-ethereum forks of shutter-network
(cannon/minigeth, optimism/l2geth).

It is possible to use the transaction types by importing this respository,
since all necessary base-types are included here.
This is useful e.g. for RLP encoding the Shutter transaction types in
rolling shutter.

## Updating shutter-network/go-ethereum type definitions

With the help of the purpose built `shtypecopy` tool,
it is possible to copy a subset of modified files 
to the target directory, and eventually replace the 
original packes import paths.

Workflow how to update go-ethereum type definitions accross code-bases:

1) Merge the type changes in https://github.com/shutter-network/go-ethereum/ 

2) Sync code in the `cannon` repo

```
cd cannon/minigeth/core/types
./shtypecopy --in core/types --out . --rules rules.shtypecopy
git add . -u
git commit -m "Sync types with github.com/shutter-network/go-ethereum"
```
And push and merge upstream.

3) Sync code in the `txtypes` repo
```
cd txtypes/types
./shtypecopy --in core/types  --out . --rules rules.shtypecopy
git add . -u
git commit -m "Sync types with github.com/shutter-network/go-ethereum"
```

4) Push a semver tag
```
git tag v0.0.42
git push --tags
```

5) Update txypes version in `rolling-shutter` repo
You have to wait at least 1 minute after pushing the new version in order to not poison the `proxy.golang.org` cache:

From https://proxy.golang.org/
> The new version should be available within one minute. Note that if someone requested the version before the tag was pushed, it may take up to 30 minutes for the mirror's cache to expire and fresh data about the version to become available.

```
cd rolling-shutter/rolling-shutter
go get github.com/shutter-network/txtypes@v0.0.42
git add . -u
git commit -m "Update txtypes version to v0.0.42"
```
