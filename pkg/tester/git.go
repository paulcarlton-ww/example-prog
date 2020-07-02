/*
Copyright 2020 The Flux CD contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tester

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	sourcev1 "github.com/fluxcd/source-controller/api/v1alpha1"
	"github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/paulcarlton-ww/example-prog/pkg/git"
)

// GitRepository defines the desired state of a Git repository.
type GitRepository struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Password string `json:"password"`
	Username string `json:"username"`

	// The git reference to checkout and monitor for changes, defaults to
	// master branch.
	// +optional
	Reference *sourcev1.GitRepositoryRef `json:"ref,omitempty"`
}

func Do(repository GitRepository) error {
	ctx := context.Background()
	// create tmp dir for the Git clone
	tmpGit, err := ioutil.TempDir("", repository.Name)
	if err != nil {
		return fmt.Errorf("tmp dir error: %w", err)
	}
	defer os.RemoveAll(tmpGit)

	// determine auth method
	var auth transport.AuthMethod
	checkoutStrategy := git.CheckoutStrategyForRef(repository.Reference)
	commit, revision, err := checkoutStrategy.Checkout(ctx, tmpGit, repository.URL, auth)
	if err != nil {
		return err
	}
	fmt.Printf("commit: %s, revision: %s", commit, revision)
	return nil
}
