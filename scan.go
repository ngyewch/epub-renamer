package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jacoblockett/gosan/v3"
	"github.com/taylorskalyo/goreader/epub"
	"github.com/urfave/cli/v3"
)

func doScan(ctx context.Context, cmd *cli.Command) error {
	inputDir := cmd.StringArg(inputDirArg.Name)
	outputFile := cmd.StringArg(outputFileArg.Name)

	config, err := getConfig(ctx, cmd)
	if err != nil {
		return err
	}

	scanner, err := NewScanner(outputFile, config)
	if err != nil {
		return err
	}
	defer func(scanner *Scanner) {
		_ = scanner.Close()
	}(scanner)

	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".epub" {
			return nil
		}

		relativePath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		err = scanner.ScanFile(path, relativePath)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type Scanner struct {
	config                *Config
	f                     *os.File
	csvWriter             *csv.Writer
	filenameTemplateParts []*template.Template
}

func NewScanner(path string, config *Config) (*Scanner, error) {
	err := os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	csvWriter := csv.NewWriter(f)
	err = csvWriter.Write([]string{
		"input",
		"title",
		"creator",
		"publisher",
		"output",
	})

	var filenameTemplateParts []*template.Template
	for _, filenameTemplatePartString := range config.FilenameTemplateParts {
		filenameTemplatePart, err := template.New("").Parse(filenameTemplatePartString)
		if err != nil {
			return nil, err
		}
		filenameTemplateParts = append(filenameTemplateParts, filenameTemplatePart)
	}

	return &Scanner{
		config:                config,
		f:                     f,
		csvWriter:             csvWriter,
		filenameTemplateParts: filenameTemplateParts,
	}, nil
}

func (scanner *Scanner) Close() error {
	scanner.csvWriter.Flush()
	_ = scanner.f.Close()
	return nil
}

type TemplateData struct {
	Metadata epub.Metadata
}

func (scanner *Scanner) ScanFile(inputFile string, relativePath string) error {
	rc, err := epub.OpenReader(inputFile)
	if err != nil {
		return err
	}
	defer rc.Close()

	book := rc.Rootfiles[0]

	templateData := TemplateData{
		Metadata: book.Metadata,
	}
	var filenameParts []string
	if scanner.config.UseOriginalDirectoryLayout {
		filenameParts = append(filenameParts, filepath.Dir(relativePath))
	}
	for _, filenameTemplatePart := range scanner.filenameTemplateParts {
		buf := bytes.NewBuffer(nil)
		err := filenameTemplatePart.Execute(buf, templateData)
		if err != nil {
			return err
		}
		filenamePart := buf.String()
		sanitizedFilenamePart, err := gosan.Filename(filenamePart, nil)
		if err != nil {
			return err
		}
		filenameParts = append(filenameParts, sanitizedFilenamePart)
	}
	outputPath := filepath.Join(filenameParts...) + ".epub"

	err = scanner.csvWriter.Write([]string{
		relativePath,
		book.Metadata.Title,
		book.Metadata.Creator,
		book.Metadata.Publisher,
		outputPath,
	})
	if err != nil {
		return err
	}

	return nil
}
