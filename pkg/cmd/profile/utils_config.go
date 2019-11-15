package profile

import (
	"fmt"
	"github.com/AlecAivazis/survey/terminal"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"io"
	"io/ioutil"
	"sigs.k8s.io/yaml"
	"strings"
)

// -------------- Section for utils

// FileAuthConfigSaver is a ConfigSaver that saves its config to the local filesystem
type FileAuthConfigSaver struct {
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

func (s *FileAuthConfigSaver) SaveConfig(config *GitAuth) error {
	fileName := s.FileName
	if fileName == "" {
		return fmt.Errorf("no filename defined")
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, data, util.DefaultWritePermissions)
}

func CheckRequiredArgs(args ...interface{}) (bool, error) {
	for key, value := range args {
		if value == nil || value == "" {
			return false, errors.Wrapf(nil, "Error: missing arg %s", key)
		}
	}

	return true, nil

}

// PickValue gets an answer to a prompt from a user's free-form input
func PickValue(message string, defaultValue string, required bool, help string, handles IOFileHandles) (string, error) {
	answer := ""
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
		Help:    help,
	}
	validator := survey.Required
	if !required {
		validator = nil
	}
	surveyOpts := survey.WithStdio(handles.In, handles.Out, handles.Err)
	err := survey.AskOne(prompt, &answer, validator, surveyOpts)
	if err != nil {
		return "", err
	}
	return answer, nil
}

// PickPassword gets a password (via hidden input) from a user's free-form input
func PickPassword(message string, help string, handles IOFileHandles) (string, error) {
	answer := ""
	prompt := &survey.Password{
		Message: message,
		Help:    help,
	}
	validator := survey.Required
	surveyOpts := survey.WithStdio(handles.In, handles.Out, handles.Err)
	err := survey.AskOne(prompt, &answer, validator, surveyOpts)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}
