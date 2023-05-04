package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"

	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Config struct {
	Style  string   `json:"style"`
	Filter []string `json:"filter"`
	Dirs   []string `json:"dirs"`
}

const configExample = `{
  "style": "WebKit",
  "filter": [
    "*.c",
    "*.cc",
    "*.cpp",
    "*.h",
    "*.hh",
    "*.hpp"
  ],
  "dirs": [
    "./"
  ]
}`

func main() {
	var config Config
	var configPath string
	flag.StringVar(&configPath, "c", "", "format config path. if the param is empty, would create a example config.")
	flag.Parse()

	if configPath == "" {
		fmt.Println("param config is empty, create a example config...")
		f, err := os.Create("./format.json")
		if err != nil {
			fmt.Println("create example config fail!", err)
			return
		}
		defer f.Close()
		writer := bufio.NewWriter(f)
		defer writer.Flush()
		writer.WriteString(configExample)
		fmt.Println("example config is created!")
		return
	} else {
		f, err := os.Open(configPath)
		if err != nil {
			fmt.Println("open config file fail!", err)
			return
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		decoder := json.NewDecoder(reader)
		err = decoder.Decode(&config)
		if err != nil {
			fmt.Println("decode json fail!", err)
			return
		}
	}

	if len(config.Filter) == 0 {
		config.Filter = append(config.Filter, "*")
	} else if len(config.Filter) == 1 && config.Filter[0] == "" {
		config.Filter[0] = "*"
	}

	dirs := GetDirsList(config.Dirs)
	dirs = filterConfigDirs(dirs, config.Filter)
	files := make([]string, 0)
	for _, dir := range dirs {
		for _, filter := range config.Filter {
			path := fmt.Sprintf("%s/%s", dir, filter)
			files = append(files, path)
			fmt.Println("format file: ", path)
		}
	}

	err := ClangFormat(config.Style, files)
	if err != nil {
		fmt.Println(err)
	}
}

func filterConfigDirs(dirs []string, s []string) []string {
	res := []string{}
	for _, d := range dirs {
		var count = 0
		for _, f := range s {

			var matches, err = (fs.Glob(os.DirFS("./"+d), "./"+f))
			if err != nil {
				panic(err)

			}
			count += len(matches)
		}
		if count > 0 {
			res = append(res, d)
		}
	}

	return res
}

func GetDirList(path string) []string {
	dirs := make([]string, 0)
	rd, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("read dir error!", err)
		return dirs
	}
	for _, fi := range rd {
		if fi.IsDir() {
			subDirs := GetDirList(path + "/" + fi.Name())
			dirs = append(dirs, subDirs...)
		}
	}
	dirs = append(dirs, path)
	return dirs
}

func GetDirsList(paths []string) []string {
	dirs := make([]string, 0)
	for _, path := range paths {
		dirs = append(dirs, GetDirList(path)...)
	}
	return dirs
}

func ClangFormat(style string, paths []string) error {
	files := strings.Join(paths, " ")
	command := fmt.Sprintf("clang-format -i --style=%s %s", style, files)
	fmt.Println("command:", command)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		fmt.Println("runtime os is windows, use powershell")
		cmd = exec.Command("powershell", "-command", command)
		//cmd = exec.Command("cmd.exe", "/c", "start"+command)
	} else {
		fmt.Println("runtime os is linux, use bash")
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("create stdout pipe fail!")
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	for {
		output, _, err := bufio.NewReader(stdout).ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("read output fail!")
				return err
			} else {
				break
			}
		}
		fmt.Println(output)
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
