package profile

import (
	"fmt"
	"github.com/Benbentwo/gh/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
	"io"
	"io/ioutil"
	"sigs.k8s.io/yaml"
	"strings"
)

// -------------- Section for utils
// FileWriter provides a minimal interface for Stdin.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader provides a minimal interface for Stdout.
type FileReader interface {
	io.Reader
	Fd() uintptr
}

// FileAuthConfigSaver is a ConfigSaver that saves its config to the local filesystem
type FileAuthConfigSaver struct {
	FileName string
}

type AuthConfig struct {
	Servers []*GitAuth `json:"servers"`

	DefaultUsername string `json:"defaultusername"`
	CurrentServer   string `json:"currentserver"`
}

// type UserAuth struct {
// 	Username    string `json:"username"`
// 	ApiToken    string `json:"apitoken"`
// 	BearerToken string `json:"bearertoken"`
// 	Password    string `json:"password,omitempty"`
// }
type IOFileHandles struct {
	Err io.Writer
	In  FileReader
	Out FileWriter
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

func (s *FileAuthConfigSaver) SaveConfig(config *AuthConfig) error {
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
func PickValue(message string, defaultValue string, required bool, help string, in terminal.FileReader, out terminal.FileWriter, outErr io.Writer) (string, error) {
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
	surveyOpts := survey.WithStdio(in, out, outErr)
	var err error
	if required {
		err = survey.AskOne(prompt, &answer, validator, surveyOpts)
	} else {
		err = survey.AskOne(prompt, &answer, nil, surveyOpts)
	}
	if err != nil {
		return "", err
	}
	return answer, nil
}

// PickPassword gets a password (via hidden input) from a user's free-form input
func PickPasswordNotReq(message string, required bool, help string, in terminal.FileReader, out terminal.FileWriter, outErr io.Writer) (string, error) {
	answer := ""
	prompt := &survey.Password{
		Message: message,
		Help:    help,
	}
	surveyOpts := survey.WithStdio(in, out, outErr)
	err := survey.AskOne(prompt, &answer, nil, surveyOpts)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

// Confirm prompts the user to confirm something
func Confirm(message string, defaultValue bool, help string, in terminal.FileReader, out terminal.FileWriter, outErr io.Writer) bool {
	answer := defaultValue
	prompt := &survey.Confirm{
		Message: message,
		Default: defaultValue,
		Help:    help,
	}
	surveyOpts := survey.WithStdio(in, out, outErr)
	err := survey.AskOne(prompt, &answer, nil, surveyOpts)
	if err != nil {
		panic(err)
	}
	log.Blank()
	return answer
}

// LoadConfig loads the configuration from the users JX config directory
func (s *FileAuthConfigSaver) LoadConfig() (*AuthConfig, error) {
	config := &AuthConfig{}
	fileName := s.FileName
	if fileName != "" {
		exists, err := util.FileExists(fileName)
		if err != nil {
			return config, fmt.Errorf("Could not check if file exists %s due to %s", fileName, err)
		}
		if exists {
			data, err := ioutil.ReadFile(fileName)
			if err != nil {
				return config, fmt.Errorf("Failed to load file %s due to %s", fileName, err)
			}
			err = yaml.Unmarshal(data, config)
			if err != nil {
				return config, fmt.Errorf("Failed to unmarshal YAML file %s due to %s", fileName, err)
			}
		}
	}
	return config, nil
}
