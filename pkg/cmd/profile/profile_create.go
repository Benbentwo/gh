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

const GLOBAL_GIT_CONFIG_FILE = "~/.gitconfig"

// options for the command
type ProfileCreateOptions struct {
	*opts.CommonOptions
	batch bool

	UseGlobal  bool
	UseLocal   bool
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

// Could do:

// Found local, create profile from that? (if doesn't exist)
// found global, (if doesn't exist already) create profile from global?
// ask questions?
// ----
func (o *ProfileCreateOptions) Run() error {

	err := CreateFromLocalGit(o)
	utils.Check(err)

	err = CreateFromGlobalGit(o)
	utils.Check(err)

	argsPass, err := CheckRequiredArgs(o.Name, o.Alias, o.Email, o.ServerName, o.ServerUrl)
	// utils.Check(err)

	if o.CommonOptions.BatchMode {
		if !argsPass {
			return errors.Wrap(err, "Missing required arguments to run in batch mode")
		}
	}

	if o.Name == "" {
		AskForString(&o.Name, "What is your Git Name", "",
			true, "Git Name", *o.CommonOptions)
	}
	if o.Email == "" {
		AskForString(&o.Email, "What is the Email Address for this git profile", "",
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

	return o.SaveConfig()
}
func (o *ProfileCreateOptions) SaveConfig() error {

	replacer := strings.NewReplacer("~", os.Getenv("HOME"))

	if o.Dir == "" {
		o.Dir = ".gh"
	}
	if o.FileName == "" {
		o.FileName = "gitAuth.yaml"
	}
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
	err = gitAuth.FlushOutGitAuth(*o.CommonOptions)
	utils.Check(err)
	utils.Info("GitAuth: %s", gitAuth)

	authConfig.Servers = append(authConfig.Servers, &gitAuth)
	err = fileAuthConfigSaver.SaveConfig(authConfig)
	utils.Check(err)

	return nil
}

/*

Thinking about this more I think we should have a function as a property of the ProfileCreateOptions,
profileCreateOptions should have batch, use local, or use global.
the function should flush out the rest of the required info - name, email, alias etc. and return a Profile struct - the stuff actually in the
yaml file. The profile create options should just be applied to the command not read or written to the file.
*/

func CreateFromLocalGit(o *ProfileCreateOptions) error {
	var name, email = "", ""
	batch := o.BatchMode
	localGitDirExists, err := util.DirExists("./.git")
	utils.Check(err)
	if (batch && localGitDirExists && o.UseLocal) || (!batch && localGitDirExists) {
		utils.Debug("Local git directory found")
		name, email = checkForNameInGitConfig("./.git/config")
	}

	if !batch && name != "" {
		utils.Debug("Git name found in local .git dir")
		utils.Info("Local Dir Name = " + name)
		useConfig := Confirm("Use above name?", false, "Use the name: "+name, o.In, o.Out, o.Err)
		if useConfig {
			o.Name = name
		}
	}
	if !batch && email != "" {
		utils.Debug("Git email found in local .git dir")
		utils.Info("Local Dir Email = " + name)
		useConfig := Confirm("Use above email?", false, "Use the email: "+email, o.In, o.Out, o.Err)
		if useConfig {
			o.Email = email
		}
	}
	if batch && o.UseLocal && name != "" && email != "" { // something found in local
		utils.Debug("Using local config with " + name + " " + email)
		o.Name = name
		o.Email = email
	}
	return nil
}

func CreateFromGlobalGit(o *ProfileCreateOptions) error {
	var name, email = "", ""
	replacer := strings.NewReplacer("~", os.Getenv("HOME"))
	gitconfigFile := replacer.Replace(GLOBAL_GIT_CONFIG_FILE)
	globalConfigExists, err := util.FileExists(gitconfigFile)
	utils.Check(err)

	if (o.BatchMode && globalConfigExists && o.UseGlobal && !o.UseLocal) || (!o.BatchMode && !o.UseLocal && globalConfigExists) {
		utils.Debug("Global Found and looking at it")
		name, email = checkForNameInGitConfig(gitconfigFile)
	}

	if !o.BatchMode && name != "" {
		utils.Debug("Git name found in global git config")
		utils.Info("Global Name = " + name)
		useConfig := Confirm("Use above name?", false, "Use the name: "+name, o.In, o.Out, o.Err)
		if useConfig {
			o.Name = name
		}
	}
	if !o.BatchMode && email != "" {
		utils.Debug("Git email found in global git config")
		utils.Info("Global Email = " + email)
		useConfig := Confirm("Use above email?", false, "Use the email: "+email, o.In, o.Out, o.Err)
		if useConfig {
			o.Email = email
		}
	}
	if o.BatchMode && o.UseGlobal && name != "" && email != "" { // something found in local
		utils.Debug("Using local config with " + name + " " + email)
		o.Name = name
		o.Email = email
	}
	return nil
}
func checkForNameInGitConfig(file string) (string, string) {

	var name, email = "", ""
	exists, lineNum, err := utils.DoesFileContainString("name =", file)
	utils.Check(err)
	utils.Debug("Line Num set to %d", lineNum)
	if exists {
		utils.Debug("Found a profile, defaulting to that.")
		bytes, err := ioutil.ReadFile(file)
		utils.Check(err)
		keys := make([]string, 0)
		keys = append(keys, "name", "email")
		kvPairs, err := GetValueFromGitConfig(bytes, keys)
		utils.Check(err)

		name = kvPairs["name"]
		email = kvPairs["email"]
	}
	return name, email
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

func (gitAuth *GitAuth) FlushOutGitAuth(o opts.CommonOptions) error {
	if gitAuth.Name == "" {
		AskForString(&gitAuth.Name, "What is your Git Name", "",
			true, "Git Name", o)
	}
	if gitAuth.Email == "" {
		AskForString(&gitAuth.Email, "What is the Email Address for this git profile", "",
			true, "what email address is tied to this account", o)
	}
	if gitAuth.Alias == "" {
		AskForString(&gitAuth.Alias, "What is the Alias for this profile", "",
			true, "Name the profile something unique", o)
	}
	if gitAuth.GitServer.ServerName == "" {
		AskForString(&gitAuth.GitServer.ServerName, "What is the Server Name for this profile", "",
			true, "Name the server something unique, like GHE_Benbentwo", o)
	}
	if gitAuth.GitServer.ServerUrl == "" {
		AskForString(&gitAuth.GitServer.ServerUrl, "What is the Server Url for this profile", "https://github.com",
			true, "Name the profile something unique", o)
	}
	if gitAuth.ApiToken == "" {
		AskForPassword(&gitAuth.ApiToken, "What is the ApiToken for this profile", false,
			"Enter your api token, it will be hidden to the console", o)
	}
	return nil
}

func GetValueFromGitConfig(bytes []byte, keys []string) (map[string]string, error) {
	kvPairs := make(map[string]string)
	content := strings.Split(string(bytes), "\n") // get an array of lines of the file
	for _, key := range keys {
		for _, line := range content {
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
