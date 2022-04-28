# Replacement go-ethereum/core types

This repository defines the modified `go-ethereum/core` types
to be synced accros the various go-ethereum forks of shutter-network
(cannon/minigeth, optimism/l2geth).

It is possible to use the transaction types by importing this respository,
since all necessary base-types are included here.
This is useful e.g. for RLP encoding the Shutter transaction types in
rolling shutter.


## Updating the types

Currently updating the types itself from original go-ethereum 
is a manual process.

The upstream commit hash of go-ethereum the types are 
currently modified from should be put in the 
`txtypes/types/GOETHEREUM_BASE_COMMIT` file.


## Source code syncing

With the help of the purpose built `shtypecopy` tool,
it is possible to copy a subset of modified files 
to the target directory, and eventually replace the 
original packes import paths.

E.g. for cannon minigeth:
```
>> git clone github.com/shutter-network/cannon
>> git clone github.com/shutter-network/txtypes

>> echo "replace: transaction_extension.go
       replace: transaction_marshalling.go
       replace: access_list_tx.go
       replace: batch_context_tx.go
       replace: bloom9.go
       replace: dynamic_fee_tx.go
       replace: hashing.go
       replace: legacy_tx.go
       replace: log.go
       replace: receipt.go
       replace: shutter_tx.go
       replace: transaction.go
       replace: transaction_signing.go
       replace: block.go" >> rules.txt

>> go build txypes/shtypecopy/ -o shtypecopy
>> ./shtypecopy --in txtypes/types/ --out cannon/minigeth/core/types/ --rules rules.txt

```

The rule files could be checked out in the target repos, and additional scripts
that coordinate the `shtypecopy` tool could be constructed to further optimise the
source code syncing experience.

