package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/slapec93/bitrise-plugins-io/configs"
	"github.com/slapec93/bitrise-plugins-io/services"
	"github.com/slapec93/bitrise-plugins-io/version"

	bitriseConfigs "github.com/bitrise-io/bitrise/configs"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/bitrise-io/go-utils/log"
	"github.com/codegangsta/cli"
	"github.com/pkg/errors"

	ver "github.com/hashicorp/go-version"
)

var commands = []cli.Command{
	cli.Command{
		Name:   "add-auth-token",
		Usage:  "Add API authentication token",
		Action: addAuthToken,
	},
	cli.Command{
		Name:   "apps",
		Usage:  "Get apps for user",
		Action: apps,
	},
}

//=======================================
// Functions
//=======================================

func printVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v\n", c.App.Version)
}

func before(c *cli.Context) error {
	configs.DataDir = os.Getenv(plugins.PluginInputDataDirKey)
	configs.IsCIMode = (os.Getenv(bitriseConfigs.CIModeEnvKey) == "true")

	return nil
}

func ensureFormatVersion(pluginFormatVersionStr, hostBitriseFormatVersionStr string) (string, error) {
	if hostBitriseFormatVersionStr == "" {
		return "This io plugin version would need bitrise-cli version >= 1.6.0 to access Bitrise IO", nil
	}

	hostBitriseFormatVersion, err := ver.NewVersion(hostBitriseFormatVersionStr)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to parse bitrise format version (%s)", hostBitriseFormatVersionStr)
	}

	pluginFormatVersion, err := ver.NewVersion(pluginFormatVersionStr)
	if err != nil {
		return "", errors.Errorf("Failed to parse io plugin format version (%s), error: %s", pluginFormatVersionStr, err)
	}

	if pluginFormatVersion.LessThan(hostBitriseFormatVersion) {
		return "Outdated io plugin, used format version is lower then host bitrise-cli's format version, please update the plugin", nil
	} else if pluginFormatVersion.GreaterThan(hostBitriseFormatVersion) {
		return "Outdated bitrise-cli, used format version is lower then the io plugin's format version, please update the bitrise-cli", nil
	}

	return "", nil
}

func addAuthToken(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mSet authentication token...\x1b[0m")

	args := c.Args()
	if len(args) != 1 {
		log.Errorf("Failed to set authentication token, error: %s", errors.New("invalid number of arguments"))
		os.Exit(1)
	}

	if err := configs.SetAPIToken(args[0]); err != nil {
		log.Errorf("Failed to set authentication token, error: %s", err)
		os.Exit(1)
	}

	log.Infof("\x1b[32;1mAuthentication token set successfully...\x1b[0m")
}

func apps(c *cli.Context) {
	log.Infof("")
	log.Infof("\x1b[34;1mGet user apps...\x1b[0m")

	services.GetBitriseAppsForUser()
}

//=======================================
// Main
//=======================================

// Run ...
func Run() {
	// Parse cl
	cli.VersionPrinter = printVersion

	app := cli.NewApp()

	app.Name = path.Base(os.Args[0])
	app.Usage = "Bitrise IO plugin"
	app.Version = version.VERSION

	app.Author = ""
	app.Email = ""

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "loglevel, l",
			Usage:  "Log level (options: debug, info, warn, error, fatal, panic).",
			EnvVar: "LOGLEVEL",
		},
	}
	app.Before = before
	app.Commands = commands
	// app.Action = action

	if err := app.Run(os.Args); err != nil {
		log.Errorf("Finished with Error: %s", err)
		os.Exit(1)
	}
}
