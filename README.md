# Long run

Golang implementation for _long run_

_Long run_ is just a _Server Sent Event_ normalization for REST API that stream response.

## Server

```golang
mux := http.NewServeMux()

runs := run.New(30 * time.Second)
h := rest.NewHandler(runs, "/user", func(r *run.Run, req *http.Request, arg map[string]interface{}) {
	r.Run(arg["name"])
	time.Sleep(100 * time.Millisecond)
	r.Success(nil)
})
mux.Handle("/user/", h)
```

You action are is a function : `func(r *run.Run, req *http.Request, arg map[string]interface{}`. Arguments came from JSON, and you can stream the response with the `run.Run` instance and JSON values.

## Client

Its a POST request, with a redirection to a GET SSE response.

```golang
req, err := http.NewRequest("POST", server.URL+"/user/", bytes.NewBuffer([]byte(`{"name": "Charly"}`)))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("accept", "text/event-stream") // Explicit event-stream
resp, err := client.Do(req) // golang http client handles the redirection
eventsReader, err := Longrun(resp)
for {
	evt, err := eventsReader.Read() // Reading an event
	if err == io.EOF { // Until the end
		break
	}
	fmt.Println("Read event: ", evt)
}
```

## Licence

3 terms BSD licence. © 2019 Mathieu Lecarme.
