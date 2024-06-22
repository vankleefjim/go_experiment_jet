package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/spf13/cobra"
	"github.com/vankleefjim/go_experiment_jet/internal/migrate"
	"github.com/vankleefjim/go_experiment_jet/pkg/dbconn"
)

func main() {
	cfg := dbconn.Config{}
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
