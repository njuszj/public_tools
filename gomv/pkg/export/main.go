package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/njuszj/gomv/pkg/common"
	"gopkg.in/yaml.v2"
)

const defaultOutputFilename = "results"
const JSON = "json"
const YAML = "yaml"

type Options struct {
	Dir          string `short:"d" long:"dir" description:"The target directory to export file infomations." default:"."`
	OutputPath   string `short:"p" long:"output-path" description:"The output filepath." default:"./"`
	OutputFormat string `short:"o" long:"output-format" description:"The output format, json and yaml are supported." default:"json"`
	Verbose      bool   `short:"v" long:"verbose" description:"Output all file infomation"`
}

type File struct {
	Name       string `json:"name" yaml:"name"`
	Size       string `json:"size" yaml:"size"`
	ModifyTime string `json:"modify_time" yaml:"modify_time"`
	CreateTime string `json:"create_time" yaml:"create_time"`
}

type Directory struct {
	Path           string      `json:"path" yaml:"path"`
	SubDirectories []Directory `json:"directories,omitempty" yaml:"directories,omitempty"`
	Files          []File      `json:"files,omitempty" yaml:"files"`
}

func pathJoin(parent, dirname string) string {
	seq := "/"
	if strings.HasSuffix(parent, seq) {
		return parent + dirname
	}
	return parent + seq + dirname
}

func walkDir(dirPath string) (Directory, error) {
	var dir Directory = Directory{
		Path:           dirPath,
		SubDirectories: make([]Directory, 0),
		Files:          make([]File, 0),
	}
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return dir, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			subDir, err := walkDir(pathJoin(dirPath, entry.Name()))
			if err != nil {
				return dir, err
			}
			dir.SubDirectories = append(dir.SubDirectories, subDir)
		} else {
			fileinfo, err := entry.Info()
			if err != nil {
				fmt.Println(err)
			}
			dir.Files = append(dir.Files, File{
				Name:       fileinfo.Name(),
				Size:       common.GetSize(fileinfo.Size()),
				ModifyTime: fileinfo.ModTime().Format("2006-01-02 15:04:05"),
			})
		}
	}

	return dir, nil
}

func Run() {
	var opts Options
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrRequired {
			fmt.Println(flagsErr.Message)
			parser.WriteHelp(os.Stdout)
			return
		} else {
			fmt.Println(err)
			return
		}
	}
	dir, err := walkDir(opts.Dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	outputFilename := defaultOutputFilename + "." + opts.OutputFormat

	var res []byte

	switch opts.OutputFormat {
	case JSON:
		res, err = json.MarshalIndent(dir, "", "	")
		if err != nil {
			fmt.Println(err)
			return
		}
	case YAML:
		res, err = yaml.Marshal(dir)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		fmt.Println("Unsupported format: " + opts.OutputFormat)
	}

	file, err := os.Create(pathJoin(opts.OutputPath, outputFilename))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.Write(res)

}

func main() {
	Run()
}
