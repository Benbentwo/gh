package initialize

import (
	"github.com/Benbentwo/gh/pkg/cmd/profile"
	utils "github.com/Benbentwo/go-utils"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/spf13/cobra"
)

// options for the command
type InitOptions struct {
	*opts.CommonOptions
}

func NewCmdInit(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &InitOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use: "init",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}
	options.AddInitFlags(cmd)
	// the line below (Section to...) is for the generate-function command to add a template_command to.
	// the line above this and below can be deleted.
	// DO NOT DELETE THE FOLLOWING LINE:
	// Section to add commands to:

	return cmd
}

// Run implements this command
func (o *InitOptions) Run() error {

	createOptions := &profile.ProfileCreateOptions{
		CommonOptions: o.CommonOptions,
		UseLocal:      false,
		UseGlobal:     true,
	}
	err := profile.CreateFromLocalGit(createOptions)
	utils.Check(err)

	err = profile.CreateFromGlobalGit(createOptions)
	utils.Check(err)

	err = createOptions.SaveConfig()
	return err

}

func (o *InitOptions) AddInitFlags(cmd *cobra.Command) {
	o.Cmd = cmd
}
