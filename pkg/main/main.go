package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/paulcarlton-ww/example-prog/pkg/git"
	"github.com/paulcarlton-ww/example-prog/pkg/tester"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage... <items>\n"+
		"where <items> is the not currently used\n\nOptions...\n")

	flag.PrintDefaults()
}

// Config holds log level etc
type Config struct {
	FlagSet    *flag.FlagSet
	LogLevel   *string
	Verbose    *bool
	username   *string
	password   *string
	repository *string
	url        *string
	sshKey     *string
	knownHost  *string
}

func setup() *Config {
	config := &Config{}
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	config.url = flag.String("url", "", "repository url.")
	config.repository = flag.String("repository", "", "repository name.")
	config.username = flag.String("username", "", "git username.")
	config.password = flag.String("password", "", "git user password.")
	config.LogLevel = flag.String("log-level", "info", "logging level.")
	config.sshKey = flag.String("ssh-key", "", "file name containing ssh private key, defaults to $HOME/.ssh/id_rsa")
	config.knownHost = flag.String("known-host", "", "file name containing known hosts data for server.")
	config.Verbose = flag.Bool("verbose", false, "set verbose output mode, defaults to off.")

	config.FlagSet = flag.CommandLine
	flag.Usage = usage
	flag.Parse()

	level, err := log.ParseLevel(*config.LogLevel)
	if err != nil {
		log.Warnf("failed to set log level %s", err)
		log.Warnf("invalid log level: %s, using default level", *config.LogLevel)
	} else {
		log.SetLevel(level)
	}
	log.Infof("logging at level: %s", log.GetLevel().String())

	home, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("failed to get home directory: %s", err)
		os.Exit(1)
	}
	if len(*config.sshKey) == 0 {
		sshKeyFile := fmt.Sprintf("%s/.ssh/id_rsa", home)
		config.sshKey = &sshKeyFile
	}
	return config
}

func main() {
	config := setup()
	if *config.Verbose {
		fmt.Printf("Arguments: %q\n", config.FlagSet.Args())
	}

	sshKey, err := ioutil.ReadFile(*config.sshKey) // just pass the file name
	if err != nil {
		log.Errorf("failed to read identity file: %s", err)
		os.Exit(1)
	}

	repo := git.Repository{
		Name:     *config.repository,
		URL:      *config.url,
		Username: *config.username,
		Password: *config.password,
		Identity: sshKey,
	}

	if len(*config.knownHost) != 0 {
		repo.KnownHosts, err = ioutil.ReadFile(*config.knownHost) // just pass the file name
		if err != nil {
			log.Errorf("failed to read known-host file: %s", err)

			os.Exit(1)
		}
	}

	if err := tester.Do(repo); err != nil {
		log.Errorf("failed: %s", err)
	}
}
