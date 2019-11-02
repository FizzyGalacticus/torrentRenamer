package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"torrentRenamer"
	"torrentRenamer/config"
	"torrentRenamer/exec"
	"torrentRenamer/services"
	"torrentRenamer/util"
)

func main() {
	files := config.GetPositionalArgs()
	config := config.GetConfig()

	videos := make(map[string]torrentRenamer.Video, len(files))

	services := []services.Service{services.OMDBService{}}

	for _, file := range files {
		base := path.Base(file)
		video, err := torrentRenamer.ParseTorrentName(base)
		if err == nil {
			ext := path.Ext(file)[1:]

			src, err := filepath.Abs(path.Clean(file))

			if err == nil {
				videos[src] = video
				videos[src].SetExt(ext)
			}
		}
	}

	for src, video := range videos {
		serviceResults := make([]string, 0)
		var err error
		var dest string

		for _, service := range services {
			if service.IsDefault() && service.IsAvailable() {
				result, err := service.GetNewName(video)
				if err == nil {
					serviceResults = append(serviceResults, result)
				}
			} else {
				fmt.Printf("Service '%s' unavailable\n", service.Name())
			}
		}

		if len(serviceResults) > 0 {
			if _, ok := video.(*torrentRenamer.Movie); ok {
				dest = util.JoinPaths(config.DefaultDirectories.Movies, serviceResults[0])
			} else {
				dest = util.JoinPaths(config.DefaultDirectories.Shows, serviceResults[0])
			}
		} else {
			dest = video.GetNewPath()
		}

		if path.Clean(src) != path.Clean(dest) {
			err = util.MoveFile(src, dest, !config.RenameWithoutPrompt)

			if err != nil {
				fmt.Printf("Error moving file: %s", err.Error())
			}
		}

		ext := filepath.Ext(dest)
		if ext != config.Conversion.Format && config.Conversion.AutoConvert {
			if !exec.IsCommandInPath(config.Conversion.Converter) {
				fmt.Printf("Command \"%s\" is not in path, cannot convert", config.Conversion.Converter)
				continue
			}

			old := dest
			new := old[0:len(old)-len(ext)] + config.Conversion.Format

			args, err := util.InsertTemplateData(config.Conversion.ArgsTemplate, struct {
				Old string
				New string
			}{Old: "TR-OLD", New: "TR-NEW"})

			splitArgs := strings.Split(args, " ")

			for i, arg := range splitArgs {
				if arg == "TR-OLD" {
					splitArgs[i] = old
				} else if arg == "TR-NEW" {
					splitArgs[i] = new
				}
			}

			if err == nil {
				exec.ExecuteCommandWithSTDOutput(config.Conversion.Converter, splitArgs...)
			}
		}
	}
}
