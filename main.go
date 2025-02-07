package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/gciavarrini/cf-application-discovery/pkg/discover"
	"gopkg.in/yaml.v2"
)

func main() {

	if err := generateCFManifest(); err != nil {
		log.Fatal(err)
	}

	b, err := os.ReadFile("manifest.yaml")
	if err != nil {
		log.Fatal(err)
	}
	ma := discover.Manifest{Applications: []*discover.AppManifest{}}
	err = yaml.Unmarshal(b, &ma)
	if err != nil {
		log.Fatal(err)
	}
	a, err := discover.Discover(*ma.Applications[0], "1", "default")
	if err != nil {
		log.Fatal(err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(a)
	if err != nil {
		log.Fatal(err)
	}
	b, err = yaml.Marshal(a)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("sample/values.yaml", b, os.ModeAppend); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", b)
}

func ptrTo[T comparable](s T) *T {
	return &s
}

func generateCFManifest() error {

	app := discover.AppManifest{
		Name: "foo",
		Metadata: &discover.Metadata{
			Labels:      map[string]*string{"foo": ptrTo("bar")},
			Annotations: map[string]*string{"bar": ptrTo("foo")},
		},
		Env: map[string]string{"CF_ROOT": "/root", "SERVER_PORT": "8080"},
		Routes: &discover.AppManifestRoutes{
			{
				Route:    "foo.default.cluster.io/",
				Protocol: discover.HTTP2,
			},
		},
		Services: &discover.AppManifestServices{
			{
				Name:       "mysql",
				Parameters: map[string]interface{}{"DB_NAME": "default"},
			},
		},
		Processes: &discover.AppManifestProcesses{
			{
				Type:                             discover.Web,
				Command:                          "/usr/bin/echo hello world>index.html; /usr/local/bin/python3 -m http.server $SERVER_PORT",
				Memory:                           "128Mi",
				HealthCheckType:                  discover.Http,
				HealthCheckHTTPEndpoint:          "localhost:8080/",
				HealthCheckInvocationTimeout:     30,
				HealthCheckInterval:              30,
				LogRateLimitPerSecond:            "16K",
				Lifecycle:                        "docker",
				ReadinessHealthCheckType:         discover.Http,
				ReadinessHealthCheckHttpEndpoint: "localhost:8080/",
				ReadinessHealthInvocationTimeout: 30,
				ReadinessHealthCheckInterval:     30,
			},
		},
		Docker: &discover.AppManifestDocker{
			Image: "python:3",
		},
		Stack: "default",
	}
	m := discover.NewManifest("default", &app)
	mb, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile("manifest.yaml", mb, os.ModeAppend)
}
