package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"torrentRenamer"
	"torrentRenamer/config"
	"torrentRenamer/fetch"
	"torrentRenamer/util"
)

const (
	apiURL = "http://www.omdbapi.com/"
)

type omdbResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Season   string `json:"Season"`
	Episode  string `json:"Episode"`
	Response string `json:"Response"`
	SeriesID string `json:"seriesID"`
}

type OMDBService struct{}

func (o *OMDBService) getOMDBResponse(query *url.Values) (omdbResponse, error) {
	var ret omdbResponse
	requestURL := o.getURLWithQuery(query)

	bytes, err := fetch.Get(requestURL)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(bytes, &ret)

	return ret, err
}

func (o *OMDBService) responseToMovie(r *omdbResponse) torrentRenamer.Movie {
	year, _ := strconv.Atoi(r.Year)
	return torrentRenamer.Movie{
		Name: config.ApplyRenameOverrides(r.Title),
		Year: year,
	}
}

func (o *OMDBService) responseToShow(r *omdbResponse) torrentRenamer.Show {
	season, _ := strconv.Atoi(r.Season)
	episode, _ := strconv.Atoi((r.Episode))
	return torrentRenamer.Show{
		Title:   r.Title,
		Season:  season,
		Episode: episode,
	}
}

func (o *OMDBService) getCommonQuery() url.Values {
	config := config.GetConfig()

	query := url.Values{}
	query.Add("apikey", config.Services.Omdb.ApiKey)
	query.Add("r", "json")

	return query
}

func (o *OMDBService) getURLWithQuery(q *url.Values) string {
	return fmt.Sprintf("%s?%s", apiURL, q.Encode())
}

func (o *OMDBService) searchMovie(m *torrentRenamer.Movie) (torrentRenamer.Movie, error) {
	var ret torrentRenamer.Movie
	query := o.getCommonQuery()
	query.Add("type", "movie")
	query.Add("t", m.Name)

	if m.Year != 0 {
		query.Add("y", string(m.Year))
	}

	res, err := o.getOMDBResponse(&query)
	if err == nil {
		if res.Response != "True" {
			return ret, fmt.Errorf("Could not find in OMDB: %s", m.GetNewName())
		}

		ret = o.responseToMovie(&res)
		ret.Ext = m.Ext
	}

	return ret, err
}

func (o *OMDBService) searchShowNameFromID(showID string) string {
	var title string

	query := o.getCommonQuery()
	query.Add("i", showID)
	query.Add("type", "series")

	res, err := o.getOMDBResponse(&query)
	if err == nil {
		title = res.Title
	}

	return config.ApplyRenameOverrides(title)
}

func (o *OMDBService) searchShow(s *torrentRenamer.Show) (torrentRenamer.Show, error) {
	var ret torrentRenamer.Show
	query := o.getCommonQuery()
	query.Add("type", "episode")
	query.Add("t", s.Name)
	query.Add("Season", fmt.Sprintf("%d", s.Season))
	query.Add("Episode", fmt.Sprintf("%d", s.Episode))

	res, err := o.getOMDBResponse(&query)
	if err != nil {
		return ret, err
	}

	if res.Response != "True" {
		return ret, errors.New("Could not find in OMDB")
	}

	ret = o.responseToShow(&res)
	ret.Name = o.searchShowNameFromID(res.SeriesID)
	ret.Ext = s.Ext

	return ret, err
}

func (o OMDBService) Name() string {
	return "OMDB"
}

func (o OMDBService) Search(video torrentRenamer.Video) (torrentRenamer.Video, error) {
	if video.IsMovie() {
		movie, ok := video.(*torrentRenamer.Movie)
		if ok {
			movie, err := o.searchMovie(movie)
			return &movie, err
		}
	}

	show, ok := video.(*torrentRenamer.Show)
	if ok {
		show, err := o.searchShow(show)
		return &show, err
	}

	return &torrentRenamer.Show{}, nil
}

func (o OMDBService) IsAvailable() bool {
	config := config.GetConfig()

	return config.Services.Omdb.ApiKey != ""
}

func (o OMDBService) GetNewName(video torrentRenamer.Video) (string, error) {
	result, err := o.Search(video)
	if err != nil {
		return "", err
	}

	config := config.GetConfig()

	if !result.IsValid() {
		return "", fmt.Errorf("could not find valid video results")
	}

	if result.IsMovie() {
		return util.InsertTemplateData(config.Services.Omdb.RenameTemplates.Movies, result)
	}

	return util.InsertTemplateData(config.Services.Omdb.RenameTemplates.Shows, result)
}

func (o OMDBService) IsDefault() bool {
	config := config.GetConfig()

	return config.DefaultService == "OMDB"
}
