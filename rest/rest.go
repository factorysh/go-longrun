package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/factorysh/go-longrun/run"
	"github.com/google/uuid"
)

type Handler struct {
	runs *run.Runs
}

func getId(root, path string) (uuid.UUID, error) {
	if !strings.HasPrefix(path, root) {
		return uuid.Nil, errors.New("root path doesn't match")
	}
	return uuid.Parse(path[len(root)+1 : len(path)])
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodHead:
		resp.WriteHeader(200)
	case http.MethodPost:
		b, err := ioutil.ReadAll(req.Body)
		raw := new(json.RawMessage)
		err = json.Unmarshal(b, raw)
		fmt.Println(err)
	case http.MethodGet:

	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}
