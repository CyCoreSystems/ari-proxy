package main

var version = "master"

func main() {
	if err := RootCmd.Execute(); err != nil {
		Log.Error("server died", "error", err)
	}
}
