package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func doRename(ctx context.Context, cmd *cli.Command) error {
	inputDir := cmd.StringArg(inputDirArg.Name)
	inputFile := cmd.StringArg(inputFileArg.Name)
	outputDir := cmd.StringArg(outputDirArg.Name)

	if inputDir == "" {
		return fmt.Errorf("input-dir is required")
	}
	if inputFile == "" {
		return fmt.Errorf("input-file is required")
	}
	if outputDir == "" {
		return fmt.Errorf("output-dir is required")
	}

	f, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	csvReader := csv.NewReader(f)
	headers, err := csvReader.Read()
	if err != nil {
		return err
	}
	inputColumn := -1
	outputColumn := -1
	for i, header := range headers {
		switch header {
		case "input":
			if inputColumn == -1 {
				inputColumn = i
			} else {
				return fmt.Errorf("multiple 'input' columns found")
			}
		case "output":
			if outputColumn == -1 {
				outputColumn = i
			} else {
				return fmt.Errorf("multiple 'output' columns found")
			}
		}
		if inputColumn == -1 {
			return fmt.Errorf("no 'input' column found")
		}
	}
	if outputColumn == -1 {
		return fmt.Errorf("no 'output' column found")
	}

	for {
		cells, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		var input string
		var output string
		if inputColumn < len(cells) {
			input = cells[inputColumn]
		}
		if outputColumn < len(cells) {
			output = cells[outputColumn]
		}
		if (input != "") && (output != "") {
			err := copyFile(filepath.Join(inputDir, input), filepath.Join(outputDir, output))
			if err != nil {
				slog.Error("failed to rename file",
					slog.Any("err", err),
					slog.String("input", input),
					slog.String("output", output),
				)
			}
		}
	}

	return nil
}

func copyFile(source string, target string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	err = os.MkdirAll(filepath.Dir(target), 0755)
	if err != nil {
		return err
	}

	w, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func(w *os.File) {
		_ = w.Close()
	}(w)

	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}

	err = os.Chtimes(target, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return err
	}

	err = os.Chmod(target, sourceInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}
