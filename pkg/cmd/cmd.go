package cmd

import (
	"github.com/Benbentwo/gh/pkg/cmd/initialize"
	"github.com/Benbentwo/gh/pkg/cmd/profile"
	"github.com/Benbentwo/gh/pkg/log"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"io"
	"strconv"
	"strings"
)

func NewGHCommand(in terminal.FileReader, out terminal.FileWriter, err io.Writer, args []string) *cobra.Command {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	baseCommand := &cobra.Command{
		Use:              "gh",
		Short:            "GitHub CLI tool to manage profiles",
		PersistentPreRun: setLoggingLevel,
		Run:              runHelp,
	}
	commonOpts := &opts.CommonOptions{
		In:  in,
		Out: out,
		Err: err,
	}
	commonOpts.AddBaseFlags(baseCommand)
	baseCommand.AddCommand(profile.NewCmdProfile(commonOpts))
	baseCommand.AddCommand(initialize.NewCmdInit(commonOpts))
	return baseCommand
}

func setLoggingLevel(cmd *cobra.Command, args []string) {
	verbose, err := strconv.ParseBool(cmd.Flag(opts.OptionVerbose).Value.String())
	if err != nil {
		log.Logger().Errorf("Unable to determine log level")
	}

	if verbose {
		err := log.SetLevel("debug")
		if err != nil {
			log.Logger().Errorf("Unable to set log level to debug")
		}
	} else {
		err := log.SetLevel("info")
		if err != nil {
			log.Logger().Errorf("Unable to set log level to info")
		}
	}
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
