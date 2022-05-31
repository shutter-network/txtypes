# shtypetool
Keeping modified go-ethereum type-redefintions in sync with multiple forks

This is a simple helper tool to distribute type-redefinitions
to various forks that can't import the code, but need to
have the source-code located in their packages.

Currently this is a simple copy tool that has a 
copies files from the source directory to the target directory 
based on rules defined in a rule file.

Additionally, the import paths in the target files can be replaced
based on a simple pattern matching algorithm.

While the functionality of `shtypetool` currently is very general and naive, 
the tool itself will always stay purpose build for the specific task at hand.


Expect this tool to become less general and more narrow-purpose.


## Usage
```
Usage:
  shtypetool [flags]

Flags:
  -h, --help           help for shtypetool
      --in string      input dir with type definitions
      --out string     output dir with to be replaced type definitions
      --rules string   rule file path
```

## Rule file

The rule file specifies which files should be copied from the 
input directory.

### Syntax of the rule file:


#### Adding files that don't exist in the output directory:

`new: <filename input dir> => <filename output dir>`

or

`new: <filename input dir>`

Here the output filename will be the same as the input filename


#### Replacing type files that exist in the output directory:

`replace: <filename input dir> => <filename output dir>`

or

`replace: <filename input dir>`

Here the output filename will be the same as the input filename


#### Replacing import statements:

`import-replace: github.com/foo/repo => github.com/bar/repo`


#### Sourcing from a git repository:

`source: github.com/<user>/<repo>@<branch>`

This will set the input directory relative to the specified
repository-root and will pull the copied files direclty from there.
It will also put the hash of the last common ancester-commit of `<branch>`
and the `main` branch in a file `GOETHEREUM_BRANCH_OFF` and the 
input repository URL in a file `GOETHEREUM_SOURCE`

### Example rule file:
```
source: github.com/shutter-network/go-ethereum@shutter-types
import-replace: github.com/ethereum/go-ethereum => github.com/shutter-network/txtypes
new: transaction_extension.go
new: batch_context_tx.go
new: shutter_tx.go
replace: transaction_marshalling.go
replace: access_list_tx.go
replace: bloom9.go
replace: dynamic_fee_tx.go
replace: hashing.go
replace: legacy_tx.go
replace: log.go
replace: receipt.go
replace: transaction.go
replace: transaction_signing.go
replace: block.go
```



