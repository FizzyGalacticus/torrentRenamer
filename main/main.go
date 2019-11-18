package main

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"torrentRenamer"
	"torrentRenamer/config"
	"torrentRenamer/exec"
	"torrentRenamer/services"
	"torrentRenamer/util"
)

func getParsedVideosBySource(files []string) map[string]torrentRenamer.Video {
	videos := make(map[string]torrentRenamer.Video, len(files))

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

	return videos
}

func getVideoDestination(v *torrentRenamer.Video) string {
	video := *v

	var serviceResult string
	var err error

	config := config.GetConfig()

	if serviceResult, err = services.GetDefaultServiceResults(v); err == nil {
		if _, ok := video.(*torrentRenamer.Movie); ok {
			return util.JoinPaths(config.DefaultDirectories.Movies, serviceResult)
		}

		return util.JoinPaths(config.DefaultDirectories.Shows, serviceResult)
	}

	return video.GetNewPath()
}

func processConversions(possibleConversions []string) error {
	config := config.GetConfig()
	var err error

	for _, dest := range possibleConversions {
		ext := filepath.Ext(dest)
		if ext != config.Conversion.Format {
			if !config.Conversion.AutoConvert {
				if !util.GetYesOrNo(fmt.Sprintf("Do you want to convert %s to a(n) %s?", dest, config.Conversion.Format)) {
					break
				}
			}

			if !exec.IsCommandInPath(config.Conversion.Converter) {
				fmt.Printf("Command \"%s\" is not in path, cannot convert", config.Conversion.Converter)
				continue
			}

			var args string

			old := dest
			new := old[0:len(old)-len(ext)] + config.Conversion.Format

			args, err = util.InsertTemplateData(config.Conversion.ArgsTemplate, struct {
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
				err = exec.ExecuteCommandWithSTDOutput(config.Conversion.Converter, splitArgs...)
			}
		}
	}

	return err
}

func processVideoRenaming(videos *map[string]torrentRenamer.Video) ([]string, []string) {
	config := config.GetConfig()
	var wg sync.WaitGroup
	var lock sync.RWMutex

	movedVideos := make([]string, 0)
	notMovedVideos := make([]string, 0)

	for src, video := range *videos {
		wg.Add(1)
		go func(wg *sync.WaitGroup, src string, video torrentRenamer.Video) {
			dest := getVideoDestination(&video)

			if path.Clean(src) != path.Clean(dest) {
				lock.Lock()
				if err := util.MoveFile(src, dest, !config.RenameWithoutPrompt); err != nil {
					fmt.Printf("Error moving file: %s", err.Error())
					movedVideos = append(movedVideos, dest)
				} else {
					notMovedVideos = append(notMovedVideos, dest)
				}
				lock.Unlock()
			} else {
				notMovedVideos = append(notMovedVideos, dest)
			}

			wg.Done()
		}(&wg, src, video)
	}

	wg.Wait()

	return movedVideos, notMovedVideos
}

func main() {
	files := config.GetPositionalArgs()

	videos := getParsedVideosBySource(files)
	movedVideos, notMovedVideos := processVideoRenaming(&videos)

	possibleConversions := util.CombineStringArrays(movedVideos, notMovedVideos)

	if err := processConversions(possibleConversions); err != nil {
		fmt.Printf("Error converting video(s): %s", err.Error())
	}
}
