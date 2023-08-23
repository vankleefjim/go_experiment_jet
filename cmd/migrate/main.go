package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
	"jvk.com/things/internal/config"
	"jvk.com/things/internal/migrate"
)

func main() {
	cfg := config.DB{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{Use: "migrate"}
	rootCmd.AddCommand(migrate.Up(cfg))

	err = rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
