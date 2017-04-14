package main

import (
	"context"
	"os"
	"strings"

	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/CyCoreSystems/ari-proxy/server"
	"github.com/CyCoreSystems/ari/client/native"
	"github.com/nats-io/nats"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log log15.Logger

var rootCmd = &cobra.Command{
	Use:   "ari-proxy",
	Short: "Proxy for the Asterisk REST interface.",
	Long: `ari-proxy is a proxy for working the Asterisk daemon over NATS.
	ARI commands are exposed over NATS for operation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var handler = log15.StdoutHandler
		if viper.GetBool("verbose") {
			log.Info("Verbose logging enabled")
			handler = log15.LvlFilterHandler(log15.LvlDebug, handler)
		} else {
			handler = log15.LvlFilterHandler(log15.LvlInfo, handler)
		}
		log.SetHandler(handler)

		return runServer(ctx, log)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("server died", "error", err)
	}
}

var cfgFile string

func init() {
	log = log15.New()

	cobra.OnInitialize(readConfig)

	p := rootCmd.PersistentFlags()

	p.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ari-proxy.yaml)")
	p.BoolP("verbose", "v", false, "Enable verbose logging")

	p.String("nats.url", nats.DefaultURL, "URL for connecting to the NATS cluster")
	p.String("ari.application", "", "ARI Stasis Application")
	p.String("ari.username", "", "Username for connecting to ARI")
	p.String("ari.password", "", "Password for connecting to ARI")
	p.String("ari.http_url", "http://localhost:8088/ari", "HTTP Base URL for connecting to ARI")
	p.String("ari.websocket_url", "ws://localhost:8088/ari/events", "Websocket URL for connecting to ARI")

	for _, n := range []string{"verbose", "nats.url", "ari.application", "ari.username", "ari.password", "ari.http_url", "ari.websocket_url"} {
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
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}
}

func runServer(ctx context.Context, log log15.Logger) error {

	natsURL := viper.GetString("nats.url")
	if os.Getenv("NATS_SERVICE_HOST") != "" {
		natsURL = "nats://" + os.Getenv("NATS_SERVICE_HOST") + ":" + os.Getenv("NATS_SERVICE_PORT_CLIENT")
	}

	srv := server.New()
	srv.Log = log

	log.Info("Starting ari-proxy server")
	err := srv.Listen(ctx, native.Options{
		Application:  viper.GetString("ari.application"),
		Username:     viper.GetString("ari.username"),
		Password:     viper.GetString("ari.password"),
		URL:          viper.GetString("ari.http_url"),
		WebsocketURL: viper.GetString("ari.websocket_url"),
	}, natsURL)
	if err != nil {
		return err
	}

	return nil
}
