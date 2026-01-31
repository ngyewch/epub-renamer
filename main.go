package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string

	configFileFlag = &cli.StringFlag{
		Name:    "config-file",
		Usage:   "config file",
		Value:   "config.yml",
		Sources: cli.EnvVars("CONFIG_FILE"),
	}

	inputFileArg = &cli.StringArg{
		Name:      "input-file",
		UsageText: "(input-file)",
	}
	inputDirArg = &cli.StringArg{
		Name:      "input-dir",
		UsageText: "(input-dir)",
	}
	outputFileArg = &cli.StringArg{
		Name:      "output-file",
		UsageText: "(output-file)",
	}
	outputDirArg = &cli.StringArg{
		Name:      "output-dir",
		UsageText: "(output-dir)",
	}

	app = &cli.Command{
		Name:    "epub-renamer",
		Usage:   "epub renamer",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "scan",
				Usage:  "scan",
				Action: doScan,
				Flags: []cli.Flag{
					configFileFlag,
				},
				Arguments: []cli.Argument{
					inputDirArg,
					outputFileArg,
				},
			},
			{
				Name:   "rename",
				Usage:  "rename",
				Action: doRename,
				Arguments: []cli.Argument{
					inputDirArg,
					inputFileArg,
					outputDirArg,
				},
			},
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		slog.Error("error",
			slog.Any("err", err),
		)
	}
}
