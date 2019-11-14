package profile

import (
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

// options for the command
type ProfileOptions struct {
	*opts.CommonOptions
}

func NewCmdProfile(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &ProfileOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use: "Profile",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	options.AddProfileFlags(cmd)
	// the line below (Section to...) is for the generate-function command to add a template_command to.
	// the line above this and below can be deleted.
	// DO NOT DELETE THE FOLLOWING LINE:
	// Section to add commands to:

	return cmd
}

// Run implements this command
func (o *ProfileOptions) Run() error {
	return o.Cmd.Help()
}

func (o *ProfileOptions) AddProfileFlags(cmd *cobra.Command) {
	o.Cmd = cmd
}
