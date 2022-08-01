package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/CyCoreSystems/ari-proxy/v5/server"
	"github.com/CyCoreSystems/ari/v5/client/native"

	"github.com/inconshreveable/log15"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Log is the package logger
var Log log15.Logger

// RootCmd is the Cobra root command descriptor
var RootCmd = &cobra.Command{
	Use:   "ari-proxy",
	Short: "Proxy for the Asterisk REST interface.",
	Long: `ari-proxy is a proxy for working the Asterisk daemon over NATS/RabbitMQ.
	ARI commands are exposed over NATS/RabbitMQ for operation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if ok, _ := cmd.PersistentFlags().GetBool("version"); ok { // nolint: gas
			fmt.Println(version)
			os.Exit(0)
		}

		handler := log15.StdoutHandler
		if viper.GetBool("verbose") {
			Log.Info("verbose logging enabled")
			handler = log15.LvlFilterHandler(log15.LvlDebug, handler)
		} else {
			handler = log15.LvlFilterHandler(log15.LvlInfo, handler)
		}
		Log.SetHandler(handler)

		native.Logger.SetHandler(handler)

		return runServer(ctx, Log)
	},
}

var cfgFile string

func init() {
	Log = log15.New()

	cobra.OnInitialize(readConfig)

	p := RootCmd.PersistentFlags()

	p.BoolP("version", "V", false, "Print version information and exit")

	p.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ari-proxy.yaml)")
	p.BoolP("verbose", "v", false, "Enable verbose logging")

	p.String("nats.url", nats.DefaultURL, "URL for connecting to the NATS cluster") //backward compatibility
	p.String("messagebus.url", nats.DefaultURL, "URL for connecting to the Message Bus cluster")
	p.String("ari.application", "", "ARI Stasis Application")
	p.String("ari.username", "", "Username for connecting to ARI")
	p.String("ari.password", "", "Password for connecting to ARI")
	p.String("ari.http_url", "http://localhost:8088/ari", "HTTP Base URL for connecting to ARI")
	p.String("ari.websocket_url", "ws://localhost:8088/ari/events", "Websocket URL for connecting to ARI")

	for _, n := range []string{"verbose", "nats.url", "messagebus.url", "ari.application", "ari.username", "ari.password", "ari.http_url", "ari.websocket_url"} {
		err := viper.BindPFlag(n, p.Lookup(n))
		if err != nil {
			panic("failed to bind flag " + n)
		}
	}
}

// readConfig reads in config file and ENV variables if set.
func readConfig() {
	viper.SetConfigName(".ari-proxy") // name of config file (without extension)
	viper.AddConfigPath("$HOME")      // adding home directory as first search path

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// Load from the environment, as well
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		Log.Debug("read configuration from file")
	}
}

func runServer(ctx context.Context, log log15.Logger) error {
	messagebusURL := viper.GetString("messagebus.url")
	if messagebusURL == "" {
		messagebusURL = viper.GetString("nats.url") //backward compatibility
	}
	if os.Getenv("NATS_SERVICE_HOST") != "" {
		messagebusURL = "nats://" + os.Getenv("NATS_SERVICE_HOST") + ":" + os.Getenv("NATS_SERVICE_PORT_CLIENT")
	}

	srv := server.New()
	srv.Log = log

	log.Info("starting ari-proxy server", "version", version)
	return srv.Listen(ctx, &native.Options{
		Application:  viper.GetString("ari.application"),
		Username:     viper.GetString("ari.username"),
		Password:     viper.GetString("ari.password"),
		URL:          viper.GetString("ari.http_url"),
		WebsocketURL: viper.GetString("ari.websocket_url"),
	}, messagebusURL)
}
