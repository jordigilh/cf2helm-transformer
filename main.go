package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

func main() {
	validate := validator.New(validator.WithRequiredStructEnabled())
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
				Command: []string{"/usr/bin/echo hello world>index.html; /usr/local/bin/python3 -m http.server $SERVER_PORT"},
				Memory:  "128Mi",
				ReadinessCheck: Probe{
					Endpoint: "localhost:8080/",
					Type:     HTTPProbeType,
					Timeout:  30,
					Interval: 30,
				},
				HealthCheck: Probe{
					Endpoint: "localhost:8080/",
					Type:     HTTPProbeType,
					Timeout:  30,
					Interval: 30,
				},
				Instances:    1,
				LogRateLimit: "16K",
				Lifecycle:    "docker",
			},
		},
		Stack: "default",
		Docker: Docker{
			Image: "python:3",
		},
		Instances: 1,
	}
	err := validate.Struct(a)
	if err != nil {
		log.Fatal(err)
	}
	b, err := yaml.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("sample/values.yaml", b, os.ModeAppend); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)

}
