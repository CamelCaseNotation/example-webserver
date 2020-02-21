package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// TODO(josh): Make nameCacheSize not out of lock-step with nameURL querystring
	nameCacheSize int           = 500
	nameURL       string        = "http://uinames.com/api/?amount=500"
	jokeTimeout   time.Duration = time.Second * 10
	nameTimeout   time.Duration = time.Second * 10
)

var randomNamesCache []*RandomNameResponse

type RandomNameResponse struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
	Region  string `json:"region"`
}

type JokeResponse struct {
	Type  string `json:"type"`
	Value JokeValue
}

type JokeValue struct {
	ID         int64    `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

func randomNameJoke(w http.ResponseWriter, r *http.Request) {
	log := Log()
	// TODO(josh): This output could be more useful.
	// Also, make every handler log incoming requests maybe?
	log.Infof("incoming request: %+v", r)
	name := getRandomName()
	joke, err := getJoke(name)
	if err != nil {
		Log().Info("request to joke either returned an error or status code != 200", err)
		w.WriteHeader(http.StatusFailedDependency)
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, joke.Value.Joke)
	}
}

func getJoke(name *RandomNameResponse) (*JokeResponse, error) {
	// TODO(josh): Allow user to modify categories for jokes via an additional query string
	jokeURL := fmt.Sprintf("http://api.icndb.com/jokes/random?firstName=%s&lastName=%s&limitTo=[nerdy]",
		url.QueryEscape(name.Name), url.QueryEscape(name.Surname))
	joke := JokeResponse{}
	statusCode, body, err := fasthttp.GetTimeout(nil, jokeURL, jokeTimeout)
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, nil
	}
	if err := json.Unmarshal(body, &joke); err != nil {
		// Issue marshalling the JSON response into our struct
		Log().WithError(err).Errorf("could not unmarshall JSON %+v into JokeResponse struct")
		return nil, err
	}
	return &joke, nil
}

// getRandomName populates a cache of random names
func getRandomName() *RandomNameResponse {
	// TODO(josh): Decide how often this cache should be re-populated
	// Currently it is filled with 500 random names and then we re-use those for all
	// requests. If that is not sufficiently random, we could "pop" names out of the cache
	// as we use them, in which case every 500 requests would take a little longer while we
	// fill up our cache.
	if len(randomNamesCache) == 0 {
		Log().Info(fmt.Sprintf(
			"randomNamesCache not yet populated, filling it with %d random names",
			nameCacheSize,
		))
		statusCode, body, err := fasthttp.GetTimeout(nil, nameURL, nameTimeout)
		if err == nil && statusCode == 200 {
			json.Unmarshal(body, &randomNamesCache)
		}
	}
	return randomNamesCache[rand.Intn(nameCacheSize)]
}
