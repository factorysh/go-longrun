package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func (h *Handler) head(resp http.ResponseWriter, req *http.Request) {
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
}

func (h *Handler) post(resp http.ResponseWriter, req *http.Request) {
	if req.RequestURI != h.root+"/" {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("content-type") != "application/json" {
		resp.WriteHeader(400)
		return
	}
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(400)
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

	if parseAcceptContains(req.Header.Get("accept"), "text/event-stream") {
		resp.Header().Add("Location", h.root+"/"+run.Id().String())
		resp.WriteHeader(303)
		resp.Write(r.Bytes())
	} else {
		resp.WriteHeader(201)
		resp.Write(r.Bytes())
	}
}

func (h *Handler) get(resp http.ResponseWriter, req *http.Request) {
	id, err := getId(h.root, req.RequestURI)
	if err != nil {
		resp.WriteHeader(400)
		return
	}
	run, ok := h.runs.GetRun(id)
	if !ok {
		resp.WriteHeader(404)
		return
	}
	lei := req.Header.Get("Last-Event-ID")
	if lei == "" {
		lei = req.URL.Query().Get("last-event-id")
	}
	var since int
	if lei != "" {
		since, err = strconv.Atoi(lei)
		if err != nil {
			resp.WriteHeader(400)
			return
		}
	}
	if parseAcceptContains(req.Header.Get("accept"), "text/event-stream") {
		resp.Header().Set("content-type", "text/event-stream")
		resp.WriteHeader(200)
	} else {
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(200)
		resp.Write([]byte(`[`))
		events := run.Events(since)
		for i, event := range events {
			v, err := json.Marshal(event.Value)
			if err != nil {
				v = []byte("null")
			}
			fmt.Fprintf(resp, `{"id":%d,"state":"%s","value":%s}`,
				i+since, event.State, string(v))
			if i < len(events)-1 {
				resp.Write([]byte(","))
			}
		}
		resp.Write([]byte(`]`))
	}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodHead:
		h.head(resp, req)
	case http.MethodPost:
		h.post(resp, req)
	case http.MethodGet:
		h.get(resp, req)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}
