package main

import (
	"fmt"
	"path"
	"torrentRenamer"
	"torrentRenamer/config"
	"torrentRenamer/services"
	"torrentRenamer/util"
)

func main() {
	files := config.GetPositionalArgs()

	videos := make(map[string]torrentRenamer.Video, len(files))

	services := []services.Service{services.OMDBService{}}

	for _, file := range files {
		base := path.Base(file)
		video, err := torrentRenamer.ParseTorrentName(base)
		if err == nil {
			ext := path.Ext(file)[1:]

			videos[path.Clean(file)] = video
			videos[path.Clean(file)].SetExt(ext)
		}
	}

	for src, video := range videos {
		config := config.GetConfig()
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
	}
}
