package main

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type configuration struct {
	UserName string
	AccessKey string
	BaseUrl string
}

func (c *configuration) Clone() *configuration {
	var clone = *c
	return &clone
}

func (p *Plugin) getConfiguration() *configuration {
	p.configurationLock.RLock()
	defer p.configurationLock.RUnlock()

	if p.configuration == nil {
		return &configuration{}
	}

	return p.configuration
}

func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()
	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}
		panic("setConfiguration called with the existing configuration")
	}
	if configuration != nil {
		configuration.BaseUrl = strings.Trim(configuration.BaseUrl, "/")
	}
	p.configuration = configuration
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)
	// Load the public configuration fields from the Mattermost server configuration.
	if err := p.API.LoadPluginConfiguration(configuration); err != nil {
		return errors.Wrap(err, "failed to load plugin configuration")
	}
	p.setConfiguration(configuration)
	return nil
}

func (c *configuration) IsValid() error {
	if c.UserName == "" {
		return fmt.Errorf("You must provide a user name.")
	}
	if c.AccessKey == "" {
		return fmt.Errorf("You must provide an access key.")
	}
	if c.BaseUrl == "" {
		return fmt.Errorf("You must provide a base url.")
	}
	return nil
}
