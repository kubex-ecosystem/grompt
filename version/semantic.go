// Package version provides functionality to manage and check the version of the Grompt CLI tool.
// It includes methods to get the current version, check for the latest version, and compare versions
package version

import (
	"os"
	"path/filepath"

	gl "github.com/rafa-mori/grompt/logger"
	l "github.com/rafa-mori/logz"

	"github.com/spf13/cobra"

	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {
	if owner := os.Getenv("GITHUB_OWNER"); owner != "" {
		projectOwner = owner
	}
	moduleName = os.Getenv("GITHUB_REPOSITORY")
	if moduleName == "" {
		moduleName, err := os.Executable()
		if err == nil {
			moduleName = filepath.Base(moduleName)
		}
	}
	if moduleAliasEnv := os.Getenv("MODULE_ALIAS"); moduleAliasEnv != "" {
		moduleAlias = moduleAliasEnv
	}
}

var moduleAlias = "Grompt" // Default module alias, can be overridden by environment variable
var moduleName = "grompt"  // Default module name, can be overridden by environment variable

const currentVersionFallback = "v1.0.0"

//go:embed CLI_VERSION
var cliVersion string
var projectOwner = "rafa-mori" // Default project owner, can be overridden by environment variable
var gitModelUrl = "https://github.com/" + projectOwner + "/" + moduleName + ".git"

type Service interface {
	GetLatestVersion() (string, error)
	GetCurrentVersion() string
	IsLatestVersion() (bool, error)
}
type ServiceImpl struct {
	gitModelUrl    string
	latestVersion  string
	currentVersion string
}
type Tag struct {
	Name string `json:"name"`
}

func init() {
	l.GetLogger(moduleAlias)
}

func getLatestTag(repoURL string) (string, error) {
	apiURL := fmt.Sprintf("%s/tags", repoURL)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch tags: %s", resp.Status)
	}

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found")
	}

	return tags[0].Name, nil
}

func (v *ServiceImpl) updateLatestVersion() error {
	repoURL := strings.TrimSuffix(v.gitModelUrl, ".git")
	tag, err := getLatestTag(repoURL)
	if err != nil {
		return err
	}
	v.latestVersion = tag
	return nil
}
func (v *ServiceImpl) vrsCompare(v1, v2 []int) (int, error) {
	if len(v1) != len(v2) {
		return 0, fmt.Errorf("version length mismatch")
	}

	for idx, v2S := range v2 {
		v1S := v1[idx]
		if v1S > v2S {
			return 1, nil
		}

		if v1S < v2S {
			return -1, nil
		}
	}
	return 0, nil
}
func (v *ServiceImpl) versionAtMost(versionAtMostArg, max []int) (bool, error) {
	if comp, err := v.vrsCompare(versionAtMostArg, max); err != nil {
		return false, err
	} else if comp == 1 {
		return false, nil
	}
	return true, nil
}
func (v *ServiceImpl) parseVersion(versionToParse string) []int {
	version := make([]int, 3)
	for idx, vStr := range strings.Split(versionToParse, ".") {
		vS, err := strconv.Atoi(vStr)
		if err != nil {
			return nil
		}
		version[idx] = vS
	}
	return version
}

func (v *ServiceImpl) IsLatestVersion() (bool, error) {
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return false, err
		}
	}

	curr := v.parseVersion(v.currentVersion)
	latest := v.parseVersion(v.latestVersion)

	if curr == nil || latest == nil {
		return false, fmt.Errorf("error parsing versions")
	}

	if isLatest, err := v.versionAtMost(curr, latest); err != nil {
		return false, err
	} else if isLatest {
		return true, nil
	}
	return false, nil
}
func (v *ServiceImpl) GetLatestVersion() (string, error) {
	if v.latestVersion == "" {
		if err := v.updateLatestVersion(); err != nil {
			return "", err
		}
	}

	return v.latestVersion, nil
}
func (v *ServiceImpl) GetCurrentVersion() string { return v.currentVersion }

func NewVersionService() Service {
	return &ServiceImpl{
		gitModelUrl:    gitModelUrl,
		currentVersion: currentVersion,
		latestVersion:  "",
	}
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of " + moduleAlias,
		Long:  "Print the version number of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetVersionInfo()
		},
	}
	subLatestCmd = &cobra.Command{
		Use:   "latest",
		Short: "Print the latest version number of " + moduleAlias,
		Long:  "Print the latest version number of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetLatestVersionInfo()
		},
	}
	subCmdCheck = &cobra.Command{
		Use:   "check",
		Short: "Check if the current version is the latest version of " + moduleAlias,
		Long:  "Check if the current version is the latest version of " + moduleAlias,
		Run: func(cmd *cobra.Command, args []string) {
			GetVersionInfoWithLatestAndCheck()
		},
	}
)

//go:embed CLI_VERSION
var currentVersion string

func GetVersion() string {
	if currentVersion == "" {
		return currentVersionFallback
	}
	return currentVersion
}

func GetGitModelUrl() string {
	return gitModelUrl
}

func GetVersionInfo() string {
	gl.Log("info", "Version: "+GetVersion())
	gl.Log("info", "Git repository: "+GetGitModelUrl())
	return fmt.Sprintf("Version: %s\nGit repository: %s", GetVersion(), GetGitModelUrl())
}

func GetLatestVersionFromGit() string {
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	gitUrlWithoutGit := strings.TrimSuffix(gitModelUrl, ".git")

	response, err := netClient.Get(gitUrlWithoutGit + "/releases/latest")
	if err != nil {
		gl.Log("error", "Error fetching latest version: "+err.Error())
		gl.Log("error", gitUrlWithoutGit+"/releases/latest")
		return err.Error()
	}

	if response.StatusCode != 200 {
		gl.Log("error", "Error fetching latest version: "+response.Status)
		gl.Log("error", "Url: "+gitUrlWithoutGit+"/releases/latest")
		body, _ := io.ReadAll(response.Body)
		return fmt.Sprintf("Error: %s\nResponse: %s", response.Status, string(body))
	}

	tag := strings.Split(response.Request.URL.Path, "/")

	return tag[len(tag)-1]
}

func GetLatestVersionInfo() string {
	gl.Log("info", "Latest version: "+GetLatestVersionFromGit())
	return "Latest version: " + GetLatestVersionFromGit()
}

func GetVersionInfoWithLatestAndCheck() string {
	if GetVersion() == GetLatestVersionFromGit() {
		gl.Log("info", "You are using the latest version.")
		return fmt.Sprintf("You are using the latest version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	} else {
		gl.Log("warn", "You are using an outdated version.")
		return fmt.Sprintf("You are using an outdated version.\n%s\n%s", GetVersionInfo(), GetLatestVersionInfo())
	}
}

func CliCommand() *cobra.Command {
	versionCmd.AddCommand(subLatestCmd)
	versionCmd.AddCommand(subCmdCheck)
	return versionCmd
}
