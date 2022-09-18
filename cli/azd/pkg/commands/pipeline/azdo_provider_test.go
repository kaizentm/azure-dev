// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package pipeline

import (
	"context"
	"testing"

	"github.com/azure/azure-dev/cli/azd/pkg/azdo"
	"github.com/azure/azure-dev/cli/azd/pkg/environment"
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
		provider := &AzdoHubScmProvider{}
		ctx := context.Background()

		//act
		details, e := provider.gitRepoDetails(ctx, "https://github.com/Azure/azure-dev.git")

		//asserts
		require.Error(t, e, ErrRemoteHostIsNotAzDo)
		require.EqualValues(t, (*gitRepositoryDetails)(nil), details)
	})

	t.Run("non azure devops git remote", func(t *testing.T) {
		//arrange
		provider := &AzdoHubScmProvider{}
		ctx := context.Background()

		//act
		details, e := provider.gitRepoDetails(ctx, "git@github.com:Azure/azure-dev.git")

		//asserts
		require.Error(t, e, ErrRemoteHostIsNotAzDo)
		require.EqualValues(t, (*gitRepositoryDetails)(nil), details)
	})
}

func getAzdoScmProviderTestHarness() *AzdoHubScmProvider {
	return &AzdoHubScmProvider{
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
