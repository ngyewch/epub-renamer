package main

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func doRename(ctx context.Context, cmd *cli.Command) error {
	inputFile := cmd.StringArg(inputFileArg.Name)
	outputDir := cmd.StringArg(outputDirArg.Name)

	fmt.Println("inputFile:", inputFile)
	fmt.Println("outputDir:", outputDir)

	return nil
}
