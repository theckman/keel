package registry

import (
	"errors"

	"github.com/rusenask/docker-registry-client/registry"

	log "github.com/Sirupsen/logrus"
)

// errors
var (
	ErrTagNotSupplied = errors.New("tag not supplied")
)

// Repository - holds repository related info
type Repository struct {
	Name string
	Tags []string // available tags
}

type Client interface {
	Get(opts Opts) (*Repository, error)
	Digest(opts Opts) (digest string, err error)
}

func New() *DefaultClient {
	return &DefaultClient{}
}

type DefaultClient struct {
}

type Opts struct {
	Registry, Name, Tag string
	Username, Password  string // if "" - anonymous
}

// LogFormatter - formatter callback passed into registry client
func LogFormatter(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Get - get repository
func (c *DefaultClient) Get(opts Opts) (*Repository, error) {

	repo := &Repository{}
	hub, err := registry.New(opts.Registry, opts.Username, opts.Password)
	if err != nil {
		return nil, err
	}
	hub.Logf = LogFormatter

	tags, err := hub.Tags(opts.Name)
	if err != nil {
		return nil, err
	}
	repo.Tags = tags

	return repo, nil
}

// Digest - get digest for repo
func (c *DefaultClient) Digest(opts Opts) (digest string, err error) {
	if opts.Tag == "" {
		return "", ErrTagNotSupplied
	}

	log.WithFields(log.Fields{
		"registry":   opts.Registry,
		"repository": opts.Name,
		"tag":        opts.Tag,
	}).Info("registry client: getting digest")

	hub, err := registry.New(opts.Registry, opts.Username, opts.Password)
	if err != nil {
		return
	}
	hub.Logf = LogFormatter

	manifestDigest, err := hub.ManifestDigest(opts.Name, opts.Tag)
	if err != nil {
		return
	}

	return manifestDigest.String(), nil
}
