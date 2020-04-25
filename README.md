# Mattermost vTiger plugin

This plugin enables you to query a vTiger installation from within Mattermost. It provides the `/crm` slash command.

## Installation

Download a package from the releases page or compile the project. Currently there is not Makefile. The target executable should be named vtiger.

Extract the archive in the Mattermost plugins directory or enable it via the Mattermost system console if that setting is enabled.

## Configuration

You need to configure the following settings for the plugin to work:

- vTiger User Name: The user name of a vTiger user. Usually one creates a separate user for this.
- vTiger Access Key: The access key of the vTiger user. This is NOT the password. You can find this access key at the bottom of the user's page, under User Advanced Options.
- vTiger Base URL: The URL of your vTiger installation.
