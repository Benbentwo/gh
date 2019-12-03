package profile

import (
	utils "github.com/Benbentwo/go-utils"
	"github.com/jenkins-x/jx/pkg/cmd/helper"
	"github.com/jenkins-x/jx/pkg/cmd/opts"
	"github.com/jenkins-x/jx/pkg/cmd/templates"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

const GIT_CONFIG_FILE = "~/.gitconfig"
// options for the command
type ProfileCreateOptions struct {
	*opts.CommonOptions
	batch bool

	useGlobal	bool
	UseLocal	bool
	Name       	string
	Alias      	string
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

	cmd.Flags().StringVarP(&options.ServerUrl, "serverurl", "u", "",
		"The Git Server Url, this is the url for the server")

	cmd.Flags().BoolVarP(&options.UseLocal, "local", "l", false, "Use the local git directory to set values")
	cmd.Flags().BoolVarP(&options.UseLocal, "global", "g", true, "Use the global git config to set values")
	// --- optional for command (keeping separate so functions can be used as lib) ---
	// cmd.Flags().StringVarP(&options.Dir, "dir", "d", "",
	// 	"The file to write to inside of a directory")
	//
	// cmd.Flags().StringVarP(&options.FileName, "filename", "f", "gitAuth.yaml",
	// 	"The file to write to inside of a directory")
	options.Dir = ".gh"
	options.FileName = "gitAuth.yaml"

	return cmd
}

func (o *ProfileCreateOptions) Run() error {
	var name, email = "", ""
	batch := o.batch
	argsPass, err := CheckRequiredArgs(o.Name, o.Alias, o.Email, o.ServerName, o.ServerUrl)
	// utils.Check(err)
	replacer := strings.NewReplacer("~", os.Getenv("HOME"))

	if o.CommonOptions.BatchMode {
		if !argsPass {
			return errors.Wrap(err, "Missing required arguments to run in batch mode")
		}
	}
	gitconfigFile := replacer.Replace(GIT_CONFIG_FILE)
	globalConfigExists, err := util.FileExists(gitconfigFile)
	utils.Check(err)

	localGitDirExists, err := util.DirExists("./.git")
	utils.Check(err)
	if (batch && localGitDirExists && o.UseLocal) || (!batch && localGitDirExists) {
		utils.Debug("Local git directory found, looking at config for name or email")
		exists, lineNum, err := utils.DoesFileContainString("name =", "./.git/config")
		utils.Check(err)
		utils.Debug("Line Num set to %d", lineNum)
		if exists {
			utils.Info("Found a local profile, defaulting to that.")
			bytes, err := ioutil.ReadFile("./.git/config")
			utils.Check(err)
			keys := make([]string, 0)
			keys = append(keys, "name", "email")
			kvPairs, err := GetValueFromGitConfig(bytes, keys)
			utils.Check(err)

			name = kvPairs["name"]
			email = kvPairs["email"]
		}
	}
	if name != "" || email != "" { // something found in local

	}
	utils.Info("Git Config Found in local Dire")
	util.Confirm("Git Config found in Current directory")

	if globalConfigExists {
		utils.Debug("Looking at ~/.gitconfig for Defaults")
		bytes, err := ioutil.ReadFile(gitconfigFile)
		utils.Check(err)
		utils.Debug("Global Config file bytes: %s", string(bytes))

	}
	if o.Name == "" {
		AskForString(&o.Name, "What is your Git Name", name,
			true, "Git Name", *o.CommonOptions)
	}
	if o.Email == "" {
		AskForString(&o.Email, "What is the Email Address for this git profile", email,
			true, "what email address is tied to this account", *o.CommonOptions)
	}
	if o.Alias == "" {
		AskForString(&o.Alias, "What is the Alias for this profile", "",
			true, "Name the profile something unique", *o.CommonOptions)
	}
	if o.ServerName == "" {
		AskForString(&o.ServerName, "What is the Server Name for this profile", "",
			true, "Name the server something unique, like GHE_Benbentwo", *o.CommonOptions)
	}
	if o.ServerUrl == "" {
		AskForString(&o.ServerUrl, "What is the Server Url for this profile", "https://github.com",
			true, "Name the profile something unique", *o.CommonOptions)
	}
	if o.ApiToken == "" {
		AskForPassword(&o.ApiToken, "What is the ApiToken for this profile", false,
			"Enter your api token, it will be hidden to the console", *o.CommonOptions)
	}
	// if o.Dir == "" {
	// 	AskForString(&o.Dir, "What is Dir would you like to place this in","~/.gh", false,
	// 		"Enter a path, or hit enter for default", *o.CommonOptions)
	// }

	// creates ~/.gh if it doesn't exist
	// exists, err := util.DirExists(o.Dir)
	// utils.Check(err)
	path, err := utils.ConfigDir("", o.Dir)
	utils.Check(err)
	utils.Info(path)

	totalPath := "~/" + util.StripTrailingSlash(o.Dir) + "/" + o.FileName
	totalPath = replacer.Replace(totalPath)
	utils.Info("Total Path: %s", totalPath)

	fileAuthConfigSaver := FileAuthConfigSaver{
		FileName: totalPath,
	}

	authConfig, err := fileAuthConfigSaver.LoadConfig()
	if err != nil {
		return err
	}
	utils.Info("Total Path save: %s", fileAuthConfigSaver.FileName)
	gitAuth := o.CreateGitAuth()
	utils.Info("GitAuth: %s", gitAuth)

	authConfig.Servers = append(authConfig.Servers, &gitAuth)
	err = fileAuthConfigSaver.SaveConfig(authConfig)
	utils.Check(err)

	return nil
}

func AskForString(response *string, message string, defaultValue string, req bool, help string, o opts.CommonOptions) {
	val, err := PickValue(message, defaultValue, req, help, o.In, o.Out, o.Err)
	utils.Check(err)
	*response = val
}
func AskForPassword(response *string, message string, required bool, help string, o opts.CommonOptions) {
	val, err := PickPasswordNotReq(message, required, help, o.In, o.Out, o.Err)
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

func GetValueFromGitConfig(bytes []byte, keys []string) (map[string]string, error) {
	kvPairs := make(map[string]string)
	content := strings.Split(string(bytes), "\n") // get an array of lines of the file
	for _, key := range keys {
		for _,line := range content {
			if strings.Contains(line, key) {
				value := strings.SplitAfter(line, " = ")[1]
				utils.Debug("KV-Pair: %s = %s", key, value)
				kvPairs[key] = value
			}
		}
	}
	if len(kvPairs) == 0 {
		return kvPairs, errors.New("No Keys Found")
	}
	return kvPairs, nil
}