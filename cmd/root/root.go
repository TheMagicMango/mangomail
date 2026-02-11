// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/TheMagicMango/mangomail/cmd/root/send"
	"github.com/TheMagicMango/mangomail/cmd/root/serve"
	"github.com/TheMagicMango/mangomail/configs"
	"github.com/TheMagicMango/mangomail/internal/infra/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const serviceName = "mangomail"

var Cmd = &cobra.Command{
	Use:     serviceName,
	Short:   "MangoMail - Send emails via CLI or HTTP API",
	Long:    `MangoMail is a tool for sending bulk emails using HTML templates and CSV data, or exposing an HTTP API for programmatic email sending.`,
	Version: version.BuildVersion,
}

func init() {
	configs.SetDefaults()

	// Global flags that apply to all subcommands
	Cmd.PersistentFlags().String("resend-api-key", "", "Resend API key for sending emails")
	cobra.CheckErr(viper.BindPFlag(configs.RESEND_API_KEY, Cmd.PersistentFlags().Lookup("resend-api-key")))

	Cmd.PersistentFlags().String("resend-api-key-file", "", "Path to file containing Resend API key")
	cobra.CheckErr(viper.BindPFlag(configs.RESEND_API_KEY_FILE, Cmd.PersistentFlags().Lookup("resend-api-key-file")))

	Cmd.PersistentFlags().String("log-level", "info", "Log level: debug, info, warn or error")
	cobra.CheckErr(viper.BindPFlag(configs.LOG_LEVEL, Cmd.PersistentFlags().Lookup("log-level")))

	Cmd.PersistentFlags().Uint64("rate-limit", 2, "Maximum number of email requests per second")
	cobra.CheckErr(viper.BindPFlag(configs.RESEND_RATE_LIMIT, Cmd.PersistentFlags().Lookup("rate-limit")))

	// Register subcommands
	Cmd.AddCommand(send.Cmd)
	Cmd.AddCommand(serve.Cmd)

	Cmd.DisableAutoGenTag = true
}
