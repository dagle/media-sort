package mediasearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type imdbID string

//imdbGet uses movie DB because it accepts IMDB IDs
func imdbGet(id imdbID, mediatype MediaType) (Result, error) {
	v := url.Values{}
	v.Set("external_source", "imdb_id")
	resp, err := movieDBRequest("/find/"+string(id), v)
	if err != nil {
		return Result{}, err
	}
	defer resp.Body.Close()

	data := movieDBData{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return Result{}, fmt.Errorf("omdb Get: Failed to decode: %s", err)
	}
	if debugMode {
		log.Printf("Fetch IMDB entry %s -> %+v", id, data)
	}
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("movieDB error: %s: %s", id, data.StatusMessage)
	}
	//pick first result
	if mediatype == Series {
		for _, series := range data.TVResults {
			return series.toResult()
		}
	} else if mediatype == Movie {
		for _, movie := range data.MovieResults {
			return movie.toResult()
		}
	}
	return Result{}, fmt.Errorf("movieDB error: no match for %s", id)
}
