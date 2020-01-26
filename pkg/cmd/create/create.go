package create

import (
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

// options for the command
type CreateOptions struct {
	*opts.CommonOptions
}

func NewCmdCreate(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &CreateOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	options.AddCreateFlags(cmd)
	// the line below (Section to...) is for the generate-function command to add a template_command to.
	// the line above this and below can be deleted.
	// DO NOT DELETE THE FOLLOWING LINE:
	// Section to add commands to:
	cmd.AddCommand(NewCmdCreateProfile(commonOpts))

	return cmd
}

// Run implements this command
func (o *CreateOptions) Run() error {
	return o.Cmd.Help()
}

func (o *CreateOptions) AddCreateFlags(cmd *cobra.Command) {
	o.Cmd = cmd
}
