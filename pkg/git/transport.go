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

package git

import (
	"strings"

	"github.com/fluxcd/pkg/ssh/knownhosts"
	sourcev1 "github.com/fluxcd/source-controller/api/v1alpha1"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// Repository defines the desired state of a Git repository.
type Repository struct {
	Name       string `json:"name"`
	URL        string `json:"url"`
	Password   string `json:"password"`
	Username   string `json:"username"`
	Identity   []byte `json:"identity"`
	KnownHosts []byte `json:"known_hosts"`

	// The git reference to checkout and monitor for changes, defaults to
	// master branch.
	// +optional
	Reference *sourcev1.GitRepositoryRef `json:"ref,omitempty"`
}

// AuthSecretStrategyForURL returns the type of authorization based on url
func AuthSecretStrategyForURL(url string) AuthSecretStrategy {
	switch {
	case strings.HasPrefix(url, "http"):
		basicAuth := &BasicAuth{}
		return basicAuth
	case strings.HasPrefix(url, "ssh"):
		return &PublicKeyAuth{}
	}
	return nil
}

// AuthSecretStrategy defines a method to get authentication parameters
type AuthSecretStrategy interface {
	Method(repository Repository) (transport.AuthMethod, error)
}

// BasicAuth returns AuthSecretStrategy
type BasicAuth struct {
	AuthSecretStrategy
}

// Method returns transport.AuthMethod
func (s *BasicAuth) Method(repository Repository) (transport.AuthMethod, error) {
	auth := &http.BasicAuth{}
	auth.Username = repository.Username
	auth.Password = repository.Password
	return auth, nil
}

// PublicKeyAuth returns AuthSecretStrategy
type PublicKeyAuth struct {
	AuthSecretStrategy
}

// Method returns transport.AuthMethod
func (s *PublicKeyAuth) Method(repository Repository) (transport.AuthMethod, error) {
	identity := repository.Identity
	knownHosts := repository.KnownHosts

	pk, err := ssh.NewPublicKeys("git", identity, "")
	if err != nil {
		return nil, err
	}

	callback, err := knownhosts.New(knownHosts)
	if err != nil {
		return nil, err
	}
	pk.HostKeyCallback = callback
	return pk, nil
}
