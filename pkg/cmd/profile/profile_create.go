package profile

import (
	"github.com/Benbentwo/utils"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	"github.com/spf13/cobra"
)

// options for the command
type ProfileCreateOptions struct {
	*opts.CommonOptions
	batch bool
}

var (
	profile_create_long = templates.LongDesc(`
initialize the config dir, and create a github profile if it doesn't exist
`)

	profile_create_example = templates.Examples(`
 
`)
)

func NewCmdProfileCreate(commonOpts *opts.CommonOptions) *cobra.Command {
	options := &ProfileCreateOptions{
		CommonOptions: commonOpts,
	}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "creates a github profile",
		Long:    profile_create_long,
		Example: profile_create_example,
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			helper.CheckErr(err)
		},
	}

	return cmd
}

// Run implements this command
func (o *ProfileCreateOptions) Run() error {
	path, err := utils.ConfigDir("", ".gh")
	utils.Check(err)
	utils.Info(path)
	return nil
}
