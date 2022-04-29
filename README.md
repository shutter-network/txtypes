# Replacement go-ethereum/core types

This repository defines the modified `go-ethereum/core` types
to be synced accros the various go-ethereum forks of shutter-network
(cannon/minigeth, optimism/l2geth).

It is possible to use the transaction types by importing this respository,
since all necessary base-types are included here.
This is useful e.g. for RLP encoding the Shutter transaction types in
rolling shutter.


## Source code syncing

With the help of the purpose built `shtypecopy` tool,
it is possible to copy a subset of modified files 
to the target directory, and eventually replace the 
original packes import paths.


```
>> git clone github.com/shutter-network/txtypes

>> echo "source-url: github.com/shutter-network/go-ethereum@shutter-types
       replace: transaction_extension.go
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
       replace: block.go" >> rules.shtype

>> go build txypes/shtypecopy/ -o shtypecopy
>> ./shtypecopy --in core/types/ --out txtypes/types/ --rules rules.shtype

```

The rule files could be checked out in the target repos, and additional scripts
that coordinate the `shtypecopy` tool could be constructed to further optimise the
source code syncing experience.

