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
