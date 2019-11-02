package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"torrentRenamer/util"

	flag "github.com/spf13/pflag"
)

const (
	configLocationTemplate = "{{home}}/.torrentRenamerrc"
)

type renameTemplates struct {
	Movies string `json:"movies"`
	Shows  string `json:"shows"`
}

type videoDirectories struct {
	Movies string `json:"movies"`
	Shows  string `json:"shows"`
}

type service struct {
	ApiKey          string          `json:"apiKey"`
	RenameTemplates renameTemplates `json:"renameTemplates"`
}

type services struct {
	Omdb service `json:"omdb"`
}

type conversion struct {
	AutoConvert  bool   `json:"autoConvert"`
	Format       string `json:"format"`
	Converter    string `json:"converter"`
	ArgsTemplate string `json:"commandTemplate"`
}

type Config struct {
	DefaultDirectories  videoDirectories  `json:"defaultDirectories"`
	Services            services          `json:"services"`
	DefaultService      string            `json:"defaultService"`
	RenameTemplates     renameTemplates   `json:"renameTemplates"`
	Conversion          conversion        `json:"conversion"`
	RenameOverrides     map[string]string `json:"renameOverrides"`
	RenameWithoutPrompt bool
}

var config Config

func getConfigLocation() (string, error) {
	homeDir, err := util.GetUserHomeDirectory()
	if err != nil {
		return "", err
	}

	return util.InsertTemplateData(configLocationTemplate, struct{ Home string }{homeDir})
}

func configFileExists() bool {
	configLocation, err := getConfigLocation()
	if err != nil {
		return false
	}

	_, err = os.Stat(configLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}

		panic(fmt.Errorf("Could not check config file existence: %e", err))
	}

	return true
}

func loadConfigFile() (Config, error) {
	conf := Config{}

	configLocation, err := getConfigLocation()
	if err != nil {
		return conf, err
	}

	bytes, err := ioutil.ReadFile(configLocation)
	if err == nil {
		err = json.Unmarshal(bytes, &conf)
	}

	return conf, err
}

func GetConfig() *Config {
	return &config
}

func saveConfig() {
	var configLocation string
	conf := GetConfig()

	bytes, err := json.MarshalIndent(conf, "", "\t")
	if err == nil {
		configLocation, err = getConfigLocation()
		if err == nil {
			ioutil.WriteFile(configLocation, bytes, 0644)
		}
	}
}

func addRenameOverride(override []string) bool {
	if len(override) != 2 {
		return false
	}

	conf := GetConfig()
	conf.RenameOverrides[override[0]] = override[1]

	return true
}

func removeRenameOverride(override string) bool {
	if override == "" {
		return false
	}

	conf := GetConfig()

	if _, ok := conf.RenameOverrides[override]; ok {
		delete(conf.RenameOverrides, override)
	}

	return true
}

func init() {
	var defaultConfig Config

	if configFileExists() {
		userConfig, err := loadConfigFile()
		if err != nil {
			panic(fmt.Errorf("Could not get config from file: %e", err))
		}

		defaultConfig = userConfig
	} else {
		userHomeDir, err := util.GetUserHomeDirectory()
		if err != nil {
			panic("Could not get current user information")
		}

		defaultConfig = Config{
			DefaultDirectories: videoDirectories{
				Movies: util.JoinPaths(userHomeDir, "Videos", "Movies"),
				Shows:  util.JoinPaths(userHomeDir, "Videos", "TV Shows"),
			},
			Services: services{
				Omdb: service{
					ApiKey: "",
					RenameTemplates: renameTemplates{
						Movies: "{{ .Name }} ({{ .Year }}).{{ .Ext }}",
						Shows:  "{{ .Name }}{{sep}}{{ .Name }} - Season {{padDigit .Season 2}}{{sep}}{{ .Name }} - S{{padDigit .Season 2 }}E{{padDigit .Episode 2 }} - {{ .Title }}.{{ .Ext }}",
					},
				},
			},
			DefaultService: "OMDB",
			RenameTemplates: renameTemplates{
				Movies: "{{ .Name }} ({{ .Year }}).{{ .Ext }}",
				Shows:  "{{ .Name }}{{sep}}{{ .Name }} - Season {{padDigit .Season 2}}{{sep}}{{ .Name }} - S{{padDigit .Season 2 }}E{{padDigit .Episode 2 }}.{{ .Ext }}",
			},
			Conversion: conversion{
				AutoConvert:  false,
				Format:       "mkv",
				Converter:    "ffmpeg",
				ArgsTemplate: "-i \"{{escapeSpaces .Old }}\" \"{{escapeSpaces .New }}\"",
			},
			RenameOverrides: make(map[string]string),
		}
	}

	// Default Directories
	moviesDir := flag.StringP("movies", "m", defaultConfig.DefaultDirectories.Movies, "The directory where movies are stored")
	showsDir := flag.StringP("shows", "s", defaultConfig.DefaultDirectories.Shows, "The directory where shows are stored")

	// Services
	omdbApiKey := flag.String("omdb-key", defaultConfig.Services.Omdb.ApiKey, "Your OMDB API key")
	omdbMovieTemplate := flag.String("omdb-movie-template", defaultConfig.Services.Omdb.RenameTemplates.Movies, "How you would like to rename movies with data from OMDB")
	omdbShowTempalte := flag.String("omdb-show-template", defaultConfig.Services.Omdb.RenameTemplates.Shows, "How you would like to rename shows with data from OMDB")

	defaultService := flag.String("service", defaultConfig.DefaultService, "The default service to use for video lookup")

	// Default rename templates
	movieTemplate := flag.String("movie-template", defaultConfig.RenameTemplates.Movies, "How you would like to rename movies")
	showTempalte := flag.String("show-template", defaultConfig.RenameTemplates.Shows, "How you would like to rename shows")

	// Conversion
	autoConvert := flag.BoolP("auto-convert", "a", defaultConfig.Conversion.AutoConvert, "Whether or not to attempt to auto-convert video file")
	convertFormat := flag.StringP("convert-format", "f", defaultConfig.Conversion.Format, "The format to which you'd like to auto-convert the video file")
	convertConverter := flag.StringP("converter", "c", defaultConfig.Conversion.Converter, "The program (command) used to run the video conversion")
	convertArgsTemplate := flag.String("convert-args", defaultConfig.Conversion.ArgsTemplate, "The Golang template for args passed to the converter")

	// Rename override options
	addOverride := flag.StringSlice("add-override", []string{}, "Add an override to parsed names")
	removeOverride := flag.String("rm-override", "", "Remove an override from parsed names")

	// Save config
	save := flag.Bool("save-config", false, "Saves your specified config as default")

	// Rename without prompt
	rename := flag.BoolP("yes", "y", false, "Will not prompt before renaming files")

	flag.Parse()

	config = Config{
		DefaultDirectories: videoDirectories{
			Movies: *moviesDir,
			Shows:  *showsDir,
		},
		Services: services{
			Omdb: service{
				ApiKey: *omdbApiKey,
				RenameTemplates: renameTemplates{
					Movies: *omdbMovieTemplate,
					Shows:  *omdbShowTempalte,
				},
			},
		},
		DefaultService: *defaultService,
		RenameTemplates: renameTemplates{
			Movies: *movieTemplate,
			Shows:  *showTempalte,
		},
		Conversion: conversion{
			AutoConvert:  *autoConvert,
			Format:       *convertFormat,
			Converter:    *convertConverter,
			ArgsTemplate: *convertArgsTemplate,
		},
		RenameOverrides:     defaultConfig.RenameOverrides,
		RenameWithoutPrompt: *rename,
	}

	exit := false

	if len(*addOverride) == 2 {
		change := addRenameOverride(*addOverride)

		if change {
			*save = true
		}
	}

	if *removeOverride != "" {
		change := removeRenameOverride(*removeOverride)

		if change {
			*save = true
		}
	}

	if *save {
		saveConfig()
	}

	if flag.NArg() < 1 {
		if !(*save) {
			flag.Usage()
		}

		exit = true
	}

	if exit {
		os.Exit(0)
	}
}

func GetPositionalArgs() []string {
	return flag.Args()
}

func ApplyRenameOverrides(str string) string {
	config := GetConfig()
	lowStr := strings.ToLower(str)

	for key, value := range config.RenameOverrides {
		lowKey := strings.ToLower(key)

		if strings.Contains(lowStr, lowKey) {
			str = util.CapitalizeFirstAll(value)
		}
	}

	return str
}
