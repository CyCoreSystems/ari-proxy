package main

import "github.com/CyCoreSystems/ari-proxy/cmd/ari-proxy/cmd"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		cmd.Log.Error("server died", "error", err)
	}
}
