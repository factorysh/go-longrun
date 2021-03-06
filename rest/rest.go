package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/factorysh/go-longrun/run"
	"github.com/factorysh/go-longrun/sse"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	runs   *run.Runs
	root   string
	action func(r *run.Run, req *http.Request, arg map[string]interface{})
}

func NewHandler(runs *run.Runs, root string, action func(r *run.Run, req *http.Request, arg map[string]interface{})) *Handler {
	return &Handler{
		runs:   runs,
		root:   root,
		action: action,
	}
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
	defer req.Body.Close()
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(400)
		return
	}
	var arg map[string]interface{}
	err = json.Unmarshal(b, &arg)
	if err != nil {
		fmt.Println(err)
		resp.WriteHeader(500)
		return
	}
	ctx := context.TODO()
	run := h.runs.NewRun(ctx)
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

	go h.action(run, req, arg)
}

func (h *Handler) get(l *log.Entry, resp http.ResponseWriter, req *http.Request) {
	id, err := getId(h.root, req.RequestURI)
	if err != nil {
		l.WithError(err).Error()
		resp.WriteHeader(400)
		return
	}
	l = l.WithField("id", id.String())
	run, ok := h.runs.GetRun(id)
	if !ok {
		l.WithError(err).Error()
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
			l.WithError(err).Error()
			resp.WriteHeader(400)
			return
		}
	}
	l = l.WithField("since", since)
	if parseAcceptContains(req.Header.Get("accept"), "text/event-stream") {
		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()
		sse.HandleSSE(ctx, run.Events, resp, l, since)
	} else {
		events, err := run.Since(since)
		if err != nil {
			l.WithError(err).Error()
			resp.WriteHeader(500)
			return
		}
		resp.Header().Set("content-type", "application/json")
		resp.WriteHeader(200)
		resp.Write([]byte(`[`))
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
	l.Info()
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	l := log.WithField("url", req.URL)
	switch req.Method {
	case http.MethodHead:
		h.head(resp, req)
	case http.MethodPost:
		h.post(resp, req)
	case http.MethodGet:
		h.get(l, resp, req)
	default:
		resp.WriteHeader(http.StatusMethodNotAllowed)
	}
}
