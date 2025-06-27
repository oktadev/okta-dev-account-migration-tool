# envsync 

> ⚠️ Envsync is currently in BETA. Let us know about any bugs by filing a GitHub issue! 

This tool helps you back up and restore configuration from Okta developer accounts. 
It extends the [Okta CLI Client](https://github.com/okta/okta-cli-client) to handle batch operations on all supported API endpoints which may have data to back up. 

## Alternatives

* use [Terraformer](https://github.com/GoogleCloudPlatform/terraformer/blob/master/docs/okta.md) to capture configuration that you plan to manage in Terraform moving forward
* Consider using [Okta's PowerShell Module](https://github.com/okta/okta-powershell-cli) or the [Okta CLI Client](https://github.com/okta/okta-cli-client)

## Setup

You will need: 
* Existing free dev org (for example, https://dev-123456.okta.com)
* New Integrator free plan org (for example, https://trial-123456.okta.com)
* Go language installed
* Okta CLI Client (installation instructions below)
* envsync (installation instructions below) 

### golang setup

Follow [these instructions](https://go.dev/doc/install) to set up golang on your system.

To test that you have golang set up, run `go version`. The output should show that go is successfully installed, for example: 

```
$ go version
go version go1.24.2 linux/amd64
```

### okta-cli-client setup

> ⚠️ Setup will change to using the official okta-cli-client when [this PR](https://github.com/okta/okta-cli-client/pull/18) lands

To set up the supported `okta-cli-client`:

```
git clone git@github.com:edunham/okta-cli-client.git
cd okta-cli-client
git checkout file-handling
make install
```
See the `okta-cli-client` README for setup help.

Follow the `okta-cli-client` [configuration instructions](https://github.com/okta/okta-cli-client?tab=readme-ov-file#configuration) to set up access to your developer account. 

Test the okta-cli-client installation by running `okta-cli-client user lists -h`. If the install was successful, you will see this output: 

```
$ okta-cli-client user lists -h
List all Users

Usage:
  okta-cli-client user lists [flags]

Flags:
      --all                 Fetch all items by following pagination automatically
  -d, --backup-dir string   Directory to save backups
  -b, --batch-backup        Backup multiple Users to a directory
  -h, --help                help for lists
  -l, --limit int32         Maximum number of items to return per page
  -p, --page string         Page to fetch (if supported)
  -q, --quiet               Suppress normal output

Global Flags:
      --config string   config file (default is $HOME/.okta/okta.yaml)
```

If the "Global Flags" section is missing, retry the install process. You may need to adjust your local settings to ensure that the binary which Golang created at the `make install` step was correctly added to your system path. 

### envsync setup
Clone the envsync repo (okta-dev-account-migration-tool) and open the okta-dev-account-migration-tool directory:

```
git clone https://github.com/oktadev/okta-dev-account-migration-tool.git
cd okta-dev-account-migration-tool
```

In the directory where you cloned the repository, run the install command: 

```
make install
```

Test the installation by running `envsync -h`. The output should look like: 

```
$ envsync -h
envsync is a tool for backing up and restoring Okta developer environments.
It is only designed and tested for use with Okta developer accounts.

Usage:
  envsync [command]

Available Commands:
  backup      Backup an Okta developer environment
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  restore     Restore an Okta developer environment

Flags:
  -h, --help   help for envsync

Use "envsync [command] --help" for more information about a command.
```

## Usage

### Backup

Envsync's backup behavior will attempt to copy everything in the specified dev org into files on your system. By default, these files will go into a subdirectory of `~/.okta/`, named after the org that was backed up. 

Ensure your `~/.okta.yaml` file is referencing your free dev org.

From the okta-dev-account-migration-tool directory, run a backup of your free dev org with the following command:

```
envsync backup
```

	Note: You can also use the `-o` flag to specify an output directory of your backup files.

The backup tool runs and prints to the console the content that is backed up from the free dev org you specified in the `okta.yaml` file. If you didn't specify a backup directory, the backup files are stored in the same folder as your `okta.yaml` file under the domain name of your free dev org, for example `dev-1234567`.

### Restore

Note that this tool currently supports partial restore. Not all objects can be migrated, due to API and tooling limitations. 

Migrate the backup data to your Integrator Free Plan org

Modify your `~/.okta.yaml` file to reference your Integrator Free Plan org domain and API key or OAuth 2.0 parameters.

Run the `envsync restore` command with the following flags:

```
-i (path to the directory containing your backup files)
-c (path to your okta.yaml file configuration data)
```

For example, this command might look like:

```
envsync restore -i \Users\admin\.okta\dev-1234567 -c \Users\admin\.okta\okta.yaml
```

The restore tool runs and prints to the console the content that is migrated to your new Integrator Free Plan org you specified in the `okta.yaml` file. Sign in to your Integrator Free Plan Org to review the migrated users, groups, and so on from your free dev org.

## Feedback

Please create an issue on this repo if you have feedback or feature requests!
