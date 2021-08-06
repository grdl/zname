package main

import (
	"fmt"
	"path/filepath"
	"zname/zname"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:           "zname <WORD>",
	Short:         "Zname - search through your cloud DNS records.",
	Version:       zname.Version(),
	Args:          cobra.ExactArgs(1),
	Run:           run,
	SilenceErrors: true, // We check for errors explicitly with cobra.CheckErr()
	SilenceUsage:  true, // We don't want to show usage on legit errors
}

// vip is a local instance of Viper available only inside main package. We don't want a global variable.
var vip = viper.New()

func init() {
	cmd.PersistentFlags().BoolP("help", "h", false, "Print this help and exit")
	cmd.PersistentFlags().BoolP("version", "v", false, "Print version and exit")
	cmd.PersistentFlags().BoolP("rebuild-cache", "r", false, "Rebuild the local cache")
	cmd.PersistentFlags().StringP("cache-path", "p", fmt.Sprintf("~%c%s", filepath.Separator, ".zname.cache"), "Path to the local cache file")

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	vip.AutomaticEnv()
	vip.SetEnvPrefix("ZNAME")
	vip.AddConfigPath(".")
	vip.SetConfigName(".env")
	vip.SetConfigType("env")

	err := vip.ReadInConfig()
	// Ignore error if config file is not found, default to env vars
	if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
		cobra.CheckErr(err)
	}

	err = vip.BindPFlags(cmd.PersistentFlags())
	cobra.CheckErr(err)
}

func run(_ *cobra.Command, args []string) {
	config := &zname.Config{
		RebuildCache: vip.GetBool("rebuild-cache"),
		CachePath:    vip.GetString("cache-path"),
		Word:         args[0],
	}

	app, err := zname.New(config)
	cobra.CheckErr(err)

	cobra.CheckErr(app.Run())
}

func main() {
	cobra.CheckErr(cmd.Execute())
}
