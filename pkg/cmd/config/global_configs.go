package config

import (
	"fmt"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	gitcfg "gopkg.in/src-d/go-git.v4/config"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
)

const (
	DefaultWritePermissions = 0760
	ConfigDir               = ".gh"
	GitAuthConfigFile       = "gitAuth.yaml"
)

// -------------- Section for utils

// FileConfigSaver is a ConfigSaver that saves its config to the local filesystem
type FileConfigSaver struct {
	FileName string
}

type AuthConfig struct {
	Servers []*AuthServer `json:"servers"`

	DefaultUsername string `json:"defaultusername"`
	CurrentServer   string `json:"currentserver"`
}

type AuthServer struct {
	URL   string      `json:"url"`
	Users []*UserAuth `json:"users"`
	Name  string      `json:"name"`
	Kind  string      `json:"kind"`

	CurrentUser string `json:"currentuser"`
}

type UserAuth struct {
	Username    string `json:"username"`
	ApiToken    string `json:"apitoken"`
	BearerToken string `json:"bearertoken"`
	Password    string `json:"password,omitempty"`
}
type IOFileHandles struct {
	Err io.Writer
	In  terminal.FileReader
	Out terminal.FileWriter
}

type GitServer struct {
	ServerName string `json:"servername"`
	ServerUrl  string `json:"serverurl"`
	// Users		*[]gitAuth	`json:"users"`
}
type GitAuth struct {
	Name      string    `json:"name"`
	Alias     string    `json:"alias,omitempty"`
	Email     string    `json:"email"`
	ApiToken  string    `json:"apitoken"`
	GitServer GitServer `json:"gitserver"`
}

// type Config interface {
// 	// GetStruct() *struct
// 	type Struct *struct
// }

// TODO figure out how to abstract SaveConfig(some struct)
// func (s *FileConfigSaver) SaveConfig(config *struct{}) error {
func SaveGitAuthConfig(fileConfigSaver FileConfigSaver, config *GitAuth) error {
	fileName := fileConfigSaver.FileName
	if fileName == "" {
		return fmt.Errorf("no filename defined")
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, util.DefaultWritePermissions)
}

// This is used
func CreateGitAuthConfigSaver() FileConfigSaver {
	totalPath, err := GetGitAuthConfigFile()
	if err != nil {
		panic(err)
	}

	fileAuthConfigSaver := FileConfigSaver{
		FileName: totalPath,
	}
	return fileAuthConfigSaver
}

func CheckRequiredArgs(args ...interface{}) (bool, error) {
	for key, value := range args {
		if value == nil || value == "" {
			return false, errors.Wrapf(nil, "Error: missing arg %s", key)
		}
	}

	return true, nil

}

func GetConfigDir() (string, error) {
	h := util.HomeDir()
	path := filepath.Join(h, ConfigDir)
	err := os.MkdirAll(path, DefaultWritePermissions)
	if err != nil {
		return "", err
	}
	return path, nil
}

func GetGitAuthConfigFile() (string, error) {
	config, err := GetConfigDir()
	if err != nil {
		return "", errors.Errorf("couldn't get the ConfigDir: %s", err)
	}
	return filepath.Join(config, GitAuthConfigFile), nil
}
func ReadGitProfile() {
	gitcfg.NewConfig()
}

// // PickValue gets an answer to a prompt from a user's free-form input
// func PickValue(message string, defaultValue string, required bool, help string, handles IOFileHandles) (string, error) {
// 	answer := ""
// 	prompt := &survey.Input{
// 		Message: message,
// 		Default: defaultValue,
// 		Help:    help,
// 	}
// 	validator := survey.Required
// 	if !required {
// 		validator = nil
// 	}
// 	surveyOpts := survey.WithStdio(handles.In, handles.Out, handles.Err)
// 	err := survey.AskOne(prompt, &answer, validator, surveyOpts)
// 	if err != nil {
// 		return "", err
// 	}
// 	return answer, nil
// }
//
// // PickPassword gets a password (via hidden input) from a user's free-form input
// func PickPassword(message string, help string, handles IOFileHandles) (string, error) {
// 	answer := ""
// 	prompt := &survey.Password{
// 		Message: message,
// 		Help:    help,
// 	}
// 	validator := survey.Required
// 	surveyOpts := survey.WithStdio(handles.In, handles.Out, handles.Err)
// 	err := survey.AskOne(prompt, &answer, validator, surveyOpts)
// 	if err != nil {
// 		return "", err
// 	}
// 	return strings.TrimSpace(answer), nil
// }
