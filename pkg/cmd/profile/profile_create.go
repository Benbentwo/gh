package profile

import (
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Benbentwo/utils"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	jxutil "github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
)

// options for the command
type ProfileCreateOptions struct {
	*opts.CommonOptions
	batch bool

	Name       string
	Alias      string
	Email      string
	ApiToken   string
	ServerName string
	ServerUrl  string

	FileName string
	Dir      string
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
	commonOpts.AddBaseFlags(cmd)
	cmd.Flags().StringVarP(&options.Name, "name", "n", "",
		"The Git Name, this is part of the commit, not a name of the profile")

	alias := ""
	if options.Name != "" {
		alias = options.Name
	}
	cmd.Flags().StringVarP(&options.Alias, "alias", "a", alias,
		"The Git profile Name... this is a name for the profile")

	cmd.Flags().StringVarP(&options.Email, "email", "e", "",
		"The Git Email, this is part of the commit")

	cmd.Flags().StringVar(&options.ApiToken, "api", "",
		"The Git API Token for the profile")

	cmd.Flags().StringVarP(&options.ServerName, "servername", "s", "",
		"The Git Server Name, this is an alias for the url")

	cmd.Flags().StringVarP(&options.ServerUrl, "serverurl", "u", "https://github.com",
		"The Git Server Url, this is the url for the server")

	// --- optional for command (keeping separate so functions can be used as lib) ---
	cmd.Flags().StringVarP(&options.Dir, "dir", "d", ".gh",
		"The file to write to inside of a directory")

	cmd.Flags().StringVarP(&options.FileName, "filename", "f", "gitAuth.yaml",
		"The file to write to inside of a directory")

	return cmd
}

func (o *ProfileCreateOptions) Run() error {

	argsPass, err := CheckRequiredArgs(o.Name, o.Alias, o.Email, o.ServerName, o.ServerUrl)

	if o.CommonOptions.BatchMode {
		if !argsPass {
			return errors.Wrap(err, "Missing required arguments to run in batch mode")
		}
	} else {
		if o.Name == "" {
			AskForString(&o.Name, "What is your Git Name", "",
				true, "Git Name", o.In, o.Out, o.Err)
		}
		if o.Alias == "" {
			AskForString(&o.Alias, "What is the Alias for this profile", "",
				true, "Name the profile something unique", o.In, o.Out, o.Err)
		}
		if o.Email == "" {
			AskForString(&o.Email, "What is the Email Address for this git profile", "",
				true, "what email address is tied to this account", o.In, o.Out, o.Err)
		}
		if o.ServerName == "" {
			AskForString(&o.ServerName, "What is the Server Name for this profile", "",
				true, "Name the server something unique, like GHE_Benbentwo", o.In, o.Out, o.Err)
		}
		if o.ServerUrl == "" {
			AskForString(&o.ServerUrl, "What is the Server Url for this profile", "",
				true, "Name the profile something unique", o.In, o.Out, o.Err)
		}
		if o.ApiToken == "" {
			AskForPassword(&o.Alias, "What is the ApiToken for this profile",
				"Enter your api token, it will be hidden to the console", o.In, o.Out, o.Err)
		}

	}
	// creates ~/.gh if it doesn't exist
	path, err := utils.ConfigDir("", o.Dir)
	utils.Check(err)
	utils.Info(path)

	totalPath := jxutil.StripTrailingSlash(o.Dir) + "/" + o.FileName
	fileAuthConfigSaver := FileAuthConfigSaver{
		FileName: totalPath,
	}
	gitAuth := o.CreateGitAuth()
	err = fileAuthConfigSaver.SaveConfig(&gitAuth)
	utils.Check(err)

	return nil
}

func AskForString(response *string, message string, defaultValue string, req bool, help string, in terminal.FileReader, out terminal.FileWriter, outErr io.Writer) {
	//noinspection GoFunctionCall
	val, err := jxutil.PickValue(message, defaultValue, req, help, in, out, outErr)
	utils.Check(err)
	*response = val
}
func AskForPassword(response *string, message string, help string, in terminal.FileReader, out terminal.FileWriter, outErr io.Writer) {
	//noinspection GoFunctionCall
	val, err := jxutil.PickPassword(message, help, in, out, outErr)
	utils.Check(err)
	*response = val
}

func (o *ProfileCreateOptions) CreateGitAuth() GitAuth {
	return GitAuth{
		o.Name,
		o.Alias,
		o.Email,
		o.ApiToken,
		GitServer{
			o.ServerName,
			o.ServerUrl,
		},
	}

}
