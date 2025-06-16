package cmd

import (
	"builder/internal/build"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build paths_to_buildjson",
	Short: "Build your project.",
	Long:  "Build your project according to `build.json` file.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			config, err := build.LoadConfig("./build.json", false)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = build.Build(config)
			if err != nil {
				fmt.Println(err)
				return
			}

			return
		}

		for _, configPath := range args {
			config, err := build.LoadConfig(configPath, false)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = os.Chdir(filepath.Dir(configPath))
			if err != nil {
				fmt.Println(err)
				return
			}

			err = build.Build(config)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
