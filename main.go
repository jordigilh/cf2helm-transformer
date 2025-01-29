package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {

	a := Application{
		Metadata: Metadata{
			Name:        "foo",
			Labels:      map[string]string{"foo": "bar"},
			Annotations: map[string]string{"bar": "foo"},
			Space:       "default",
		},
		Env: map[string]string{"CF_ROOT": "/root", "SERVER_PORT": "8080"},
		Routes: []Route{
			{
				Route:    "foo.default.cluster.io/",
				Protocol: HTTP2RouteProtocol},
		},
		Services: []Service{
			{Name: "mysql", Parameters: map[string]interface{}{"DB_NAME": "default"}},
		},
		Processes: []Process{
			{Type: Web,
				Command: []string{"/bin/sh", "echo", "hello world >index.html", "&&", "python3", "-m", "http.server", "$SERVER_PORT"},
				Memory:  "128MB",
				ReadinessCheck: &Probe{
					Endpoint: "localhost:8080/",
					Type:     string(HTTPProbeType),
				},
				HealthCheck: &Probe{
					Endpoint: "localhost:8080/",
					Type:     string(HTTPProbeType),
				},
			},
		},
		Stack: "default",
		Docker: &Docker{
			Image: "python:latest",
		},
	}
	b, err := yaml.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

}
