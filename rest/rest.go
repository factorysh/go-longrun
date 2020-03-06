package rest

import (
	"bytes"
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
	root string
}

func getId(root, path string) (uuid.UUID, error) {
	if !strings.HasPrefix(path, root) {
		return uuid.Nil, errors.New("root path doesn't match")
	}
	return uuid.Parse(path[len(root)+1 : len(path)])
}

func parseAcceptContains(accept, contains string) bool {
	for _, slug := range strings.Split(accept, ",") {
		slug = strings.TrimSpace(slug)
		n := strings.Split(slug, ";")
		fmt.Println(n, contains)
		if n[0] == contains {
			return true
		}
	}
	return false
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodHead:
		id, err := getId(h.root, req.RequestURI)
		if err != nil {
			resp.WriteHeader(400)
			return
		}
		_, ok := h.runs.GetRun(id)
		if ok {
			resp.WriteHeader(200)
		} else {
			resp.WriteHeader(404)
		}
	case http.MethodPost:
		if req.RequestURI != h.root+"/" {
			resp.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			resp.WriteHeader(500)
		}
		req.Body.Close()
		raw := new(json.RawMessage)
		err = json.Unmarshal(b, raw)
		if err != nil {
			fmt.Println(err)
			resp.WriteHeader(500)
		}
		run := h.runs.New()
		resp.Header().Set("content-type", "application/json")
		r := bytes.NewBufferString(`{"id":"`)
		r.WriteString(run.Id().String())
		r.WriteString(`"}`)

		fmt.Println(req.Header)
		if parseAcceptContains(req.Header.Get("accept"), "text/event-stream") {
			fmt.Println("Event-stream rulez")
			resp.Header().Add("Location", h.root+"/"+run.Id().String())
			resp.WriteHeader(303)
			resp.Write(r.Bytes())
		} else {
			resp.WriteHeader(201)
			resp.Write(r.Bytes())
		}
	case http.MethodGet:
		// FIXME
		fmt.Println("Get Headers", req.Header)
		id, err := getId(h.root, req.RequestURI)
		if err != nil {
			resp.WriteHeader(400)
			return
		}
		_, ok := h.runs.GetRun(id)
		if ok {
			resp.WriteHeader(200)
		} else {
			resp.WriteHeader(404)
		}
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}
