// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pipeline

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"testing"

	"github.com/azure/azure-dev/cli/azd/pkg/azdo"
	"github.com/azure/azure-dev/cli/azd/pkg/environment"
	"github.com/azure/azure-dev/cli/azd/pkg/input"
	"github.com/stretchr/testify/require"
)

func Test_azdo_provider_getRepoDetails(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		// arrange
		provider := getAzdoScmProviderTestHarness()
		testOrgName := provider.Env.Values[azdo.AzDoEnvironmentOrgName]
		testRepoName := provider.Env.Values[azdo.AzDoEnvironmentRepoName]
		ctx := context.Background()

		// act
		details, e := provider.gitRepoDetails(ctx, "https://fake_org@dev.azure.com/fake_org/repo1/_git/repo1")

		// assert
		require.NoError(t, e)
		require.EqualValues(t, testOrgName, details.owner)
		require.EqualValues(t, testRepoName, details.repoName)
		require.EqualValues(t, false, details.pushStatus)
	})

	t.Run("ssh", func(t *testing.T) {
		// arrange
		provider := getAzdoScmProviderTestHarness()
		testOrgName := provider.Env.Values[azdo.AzDoEnvironmentOrgName]
		testRepoName := provider.Env.Values[azdo.AzDoEnvironmentRepoName]
		ctx := context.Background()

		// act
		details, e := provider.gitRepoDetails(ctx, "git@ssh.dev.azure.com:v3/fake_org/repo1/repo1")

		// assert
		require.NoError(t, e)
		require.EqualValues(t, testOrgName, details.owner)
		require.EqualValues(t, testRepoName, details.repoName)
		require.EqualValues(t, false, details.pushStatus)
	})

	t.Run("non azure devops https remote", func(t *testing.T) {
		//arrange
		provider := &AzdoScmProvider{}
		ctx := context.Background()

		//act
		details, e := provider.gitRepoDetails(ctx, "https://github.com/Azure/azure-dev.git")

		//asserts
		require.Error(t, e, ErrRemoteHostIsNotAzDo)
		require.EqualValues(t, (*gitRepositoryDetails)(nil), details)
	})

	t.Run("non azure devops git remote", func(t *testing.T) {
		//arrange
		provider := &AzdoScmProvider{}
		ctx := context.Background()

		//act
		details, e := provider.gitRepoDetails(ctx, "git@github.com:Azure/azure-dev.git")

		//asserts
		require.Error(t, e, ErrRemoteHostIsNotAzDo)
		require.EqualValues(t, (*gitRepositoryDetails)(nil), details)
	})
}

func Test_azdo_provider_preConfigureCheck(t *testing.T) {
	t.Run("accepts a PAT via system environment variables", func(t *testing.T) {
		// arrange
		testPat := "12345"
		provider := getEmptyAzdoScmProviderTestHarness()
		os.Setenv(azdo.AzDoEnvironmentOrgName, "testOrg")
		os.Setenv(azdo.AzDoPatName, testPat)
		testConsole := &circularConsole{}
		ctx := context.Background()

		// act
		e := provider.preConfigureCheck(ctx, testConsole)

		// assert
		require.NoError(t, e)

		//cleanup
		os.Unsetenv(azdo.AzDoPatName)
	})

	t.Run("returns an error if no pat is provided", func(t *testing.T) {
		// arrange
		os.Unsetenv(azdo.AzDoPatName)
		os.Setenv(azdo.AzDoEnvironmentOrgName, "testOrg")
		provider := getEmptyAzdoScmProviderTestHarness()
		testConsole := &configurablePromptConsole{}
		testPat := "testPAT12345"
		testConsole.promptResponse = testPat
		ctx := context.Background()

		// act
		e := provider.preConfigureCheck(ctx, testConsole)

		// assert
		require.Nil(t, e)
		require.EqualValues(t, testPat, provider.Env.Values[azdo.AzDoPatName])
	})

}

func Test_saveEnvironmentConfig(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("saves to environment file", func(t *testing.T) {
		// arrange
		key := "test"
		value := "12345"
		provider := getEmptyAzdoScmProviderTestHarness()
		envPath := path.Join(tempDir, ".test.env")
		provider.Env = environment.EmptyWithFile(envPath)
		// act
		e := provider.saveEnvironmentConfig(key, value)
		// assert
		writtenEnv, err := environment.FromFile(envPath)
		require.NoError(t, err)

		require.EqualValues(t, writtenEnv.Values[key], value)
		require.NoError(t, e)
	})

}
func getEmptyAzdoScmProviderTestHarness() *AzdoScmProvider {
	return &AzdoScmProvider{
		Env: &environment.Environment{
			Values: map[string]string{},
		},
	}
}

func getAzdoScmProviderTestHarness() *AzdoScmProvider {
	return &AzdoScmProvider{
		Env: &environment.Environment{
			Values: map[string]string{
				azdo.AzDoEnvironmentOrgName:       "fake_org",
				azdo.AzDoEnvironmentProjectName:   "project1",
				azdo.AzDoEnvironmentProjectIdName: "12345",
				azdo.AzDoEnvironmentRepoName:      "repo1",
				azdo.AzDoEnvironmentRepoIdName:    "9876",
				azdo.AzDoEnvironmentRepoWebUrl:    "https://repo",
			},
		},
	}
}

// For tests where the console won't matter at all
type configurablePromptConsole struct {
	promptResponse string
}

func (console *configurablePromptConsole) Message(ctx context.Context, message string) {}
func (console *configurablePromptConsole) Prompt(ctx context.Context, options input.ConsoleOptions) (string, error) {
	return console.promptResponse, nil
}
func (console *configurablePromptConsole) Select(ctx context.Context, options input.ConsoleOptions) (int, error) {
	return 0, nil
}
func (console *configurablePromptConsole) Confirm(ctx context.Context, options input.ConsoleOptions) (bool, error) {
	return true, nil
}
func (console *configurablePromptConsole) SetWriter(writer io.Writer) {}

// For tests where console.prompt returns values provided in its internal []string
type circularConsole struct {
	selectReturnValues []int
	index              int
}

func (console *circularConsole) Message(ctx context.Context, message string) {}
func (console *circularConsole) Prompt(ctx context.Context, options input.ConsoleOptions) (string, error) {
	return "", nil
}

func (console *circularConsole) Select(ctx context.Context, options input.ConsoleOptions) (int, error) {
	// If no values where provided, return error
	arraySize := len(console.selectReturnValues)
	if arraySize == 0 {
		return 0, errors.New("no values to return")
	}

	// Reset index when it reaches size (back to first value)
	if console.index == arraySize {
		console.index = 0
	}

	returnValue := console.selectReturnValues[console.index]
	console.index += 1
	return returnValue, nil
}
func (console *circularConsole) Confirm(ctx context.Context, options input.ConsoleOptions) (bool, error) {
	return false, nil
}
func (console *circularConsole) SetWriter(writer io.Writer) {}
