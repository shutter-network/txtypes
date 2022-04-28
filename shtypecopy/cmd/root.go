package cmd

import (
	"os"

	"github.com/shutter-network/txtypes/shtypecopy/tool"
	"github.com/spf13/cobra"
)

const (
	inDirFlagName    string = "in"
	outDirFlagName   string = "out"
	ruleFileFlagName string = "rules"
)

const help string = `       .__     __                         __                .__   
  _____|  |___/  |_ ___.__.______   _____/  |_  ____   ____ |  |  
 /  ___/  |  \   __<   |  |\____ \_/ __ \   __\/  _ \ /  _ \|  |  
 \___ \|   Y  \  |  \___  ||  |_> >  ___/|  | (  <_> |  <_> )  |__
/____  >___|  /__|  / ____||   __/ \___  >__|  \____/ \____/|____/
     \/     \/      \/     |__|        \/ 
Keeping modified go-ethereum type-redefintions in sync with multiple forks

This is a simple helper tool to distribute type-redefinitions
to various forks that can't import the code, but need to
have the source-code located in their packages.

Currently this is a simple copy tool that has a 
copies files from the source directory to the target directory 
based on rules defined in a rule file.

Additionally, the import paths in the target files can be replaced
based on a simple pattern matching algorithm.

Syntax of the rule file:
=========================


Adding files that don't exist in the output directory:
--------------------------------------------------------

'new: <filename input dir> => <filename output dir>'

or

'new: <filename input dir>'
Here the output filename will be the same as the input filename


Replacing type files that exist in the output directory:
--------------------------------------------------------

'replace: <filename input dir> => <filename output dir>'
or 
'replace: <filename input dir>'
Here the output filename will be the same as the input filename


Replacing import statements:
--------------------------------------------------------

'import-replace: github.com/foo/repo => github.com/bar/repo'
`

var (
	inDir    string
	outDir   string
	ruleFile string
)

var rootCmd = &cobra.Command{
	Use:   "shtypetool",
	Short: "A brief description of your application",
	Long:  help,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tool.Run(inDir, outDir, ruleFile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&ruleFile, ruleFileFlagName, "", "rule file path")
	rootCmd.PersistentFlags().StringVar(&inDir, inDirFlagName, "", "input dir with type definitions")
	rootCmd.PersistentFlags().StringVar(&outDir, outDirFlagName, "", "output dir with to be replaced type definitions")
	rootCmd.MarkFlagFilename(ruleFileFlagName)
	rootCmd.MarkFlagDirname(inDirFlagName)
	rootCmd.MarkFlagDirname(outDirFlagName)
	rootCmd.MarkFlagRequired(inDirFlagName)
	rootCmd.MarkFlagRequired(outDirFlagName)
	rootCmd.MarkFlagRequired(ruleFileFlagName)
}
