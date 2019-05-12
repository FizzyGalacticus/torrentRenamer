package fetch

import (
	"io/ioutil"
	"net/http"
)

func Get(url string) ([]byte, error) {
	var response []byte

	req, err := http.Get(url)
	if err != nil {
		return response, err
	}

	defer req.Body.Close()

	return ioutil.ReadAll(req.Body)
}
