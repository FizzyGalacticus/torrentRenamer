package torrentRenamer

import (
	"strings"
	"torrentRenamer/config"
	"torrentRenamer/util"

	torrentParser "github.com/middelink/go-parse-torrent-name"
)

type Video interface {
	IsMovie() bool
	IsShow() bool
	IsValid() bool
	SetExt(string)
	GetNewName() string
	GetNewPath() string
}

type Movie struct {
	Name string `json:"name"`
	Year int    `json:"year"`
	Ext  string `json:"ext"`
}

func (m *Movie) IsMovie() bool {
	return true
}

func (m *Movie) IsShow() bool {
	return false
}

func (m *Movie) IsValid() bool {
	return m.Name != ""
}

func (m *Movie) SetExt(ext string) {
	m.Ext = ext
}

func (m *Movie) GetNewName() string {
	config := config.GetConfig()

	name, err := util.InsertTemplateData(config.RenameTemplates.Movies, m)
	if err != nil {
		return ""
	}

	return name
}

func (m *Movie) GetNewPath() string {
	config := config.GetConfig()

	path, err := util.InsertTemplateData(config.DefaultDirectories.Movies, m)
	if err != nil {
		return ""
	}

	return util.JoinPaths(path, m.GetNewName())
}

type Show struct {
	Name    string `json:"name"`
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Title   string `json:"title"`
	Ext     string `json:"ext"`
}

func (s *Show) IsMovie() bool {
	return false
}

func (s *Show) IsShow() bool {
	return true
}

func (s *Show) IsValid() bool {
	return s.Name != ""
}

func (s *Show) SetExt(ext string) {
	s.Ext = ext
}

func (s *Show) GetNewName() string {
	config := config.GetConfig()

	name, err := util.InsertTemplateData(config.RenameTemplates.Shows, s)
	if err != nil {
		return ""
	}

	return name
}

func (s *Show) GetNewPath() string {
	config := config.GetConfig()

	path, err := util.InsertTemplateData(config.DefaultDirectories.Shows, s)
	if err != nil {
		return ""
	}

	return util.JoinPaths(path, s.GetNewName())
}

func ParseTorrentName(name string) (Video, error) {
	var ret Video
	parsed, err := torrentParser.Parse(name)
	if err != nil {
		return ret, err
	}

	if parsed.Title[len(parsed.Title)-2:] == " -" {
		parsed.Title = strings.ReplaceAll(parsed.Title, " -", "")
	}

	parsed.Title = config.ApplyRenameOverrides(parsed.Title)

	parsed.Title = util.CapitalizeFirstAll(parsed.Title)

	if parsed.Season == 0 {
		ret = &Movie{
			Name: parsed.Title,
			Year: parsed.Year,
		}
	} else {
		ret = &Show{
			Name:    parsed.Title,
			Season:  parsed.Season,
			Episode: parsed.Episode,
		}
	}

	return ret, nil
}
