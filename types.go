package main

type Application struct {
	// Metadata captures the name, labels and annotations in the application.
	Metadata Metadata `yaml:",inline" validate:"required"`
	// Env captures the `env` field values in the CF application manifest.
	Env map[string]string `yaml:"env,omitempty"`
	// Routes represent the routes that are made available by the application.
	Routes []Route `yaml:"route,omitempty"`
	// Services captures the `services` field values in the CF application manifest.
	Services []Service `yaml:"service,omitempty"`
	// Processes captures the `processes` field values in the CF application manifest.
	Processes []Process `yaml:"process,omitempty"`
	// Sidecars captures the `sidecars` field values in the CF application manifest.
	Sidecars []Sidecar `yaml:"sidecar,omitempty"`
	// Stack represents the `stack` field in the application manifest.
	// The value is captured for information purposes because it has no relevance
	// in Kubernetes.
	Stack string `yaml:"stack,omitempty"`
	// Timeout specifies the maximum time allowed for an application to
	// respond to readiness or health checks during startup.
	// If the application does not respond within this time, the platform will mark
	// the deployment as failed. The default value is 60 seconds and maximum to 180 seconds, but both values can be changed in the Cloud Foundry Controller.
	// https://github.com/cloudfoundry/docs-dev-guide/blob/96f19d9d67f52ac7418c147d5ddaa79c957eec34/deploy-apps/large-app-deploy.html.md.erb#L35
	// Default is 60 (seconds).
	Timeout int `yaml:"timeout" validate:"min=0,max=180"`
	// BuildPacks capture the buildpacks defined in the CF application manifest.
	BuildPacks []string `yaml:"buildPacks,omitempty"`
	// Docker captures the Docker specification in the CF application manifest.
	Docker Docker `yaml:"docker,omitempty"`
	// Instances captures the number of instances to run concurrently for this application. Default is 1.
	Instances int `yaml:"instances" validate:"required,min=1"`
}

type Docker struct {
	// Image represents the pullspect where the container image is located.
	Image string `yaml:"image" validate:"required"`
	// Username captures the username to authenticate against the container registry.
	Username string `yaml:"username,omitempty"`
}

type Sidecar struct {
	// Name represents the name of the Sidecar
	Name string `yaml:"name" validate:"required"`
	// ProcessTypes captures the different process types defined for the sidecar.
	// Compared to a Process, which has only one type, sidecar processes can
	// accumulate more than one type.
	ProcessTypes []ProcessType `yaml:"processType" validate:"required,oneof=worker web"`
	// Command captures the command to run the sidecar
	Command string `yaml:"command" validate:"required"`
	// Memory represents the amount of memory to allocate to the sidecar.
	// It's an optional field.
	Memory string `yaml:"memory,omitempty"`
}

type Service struct {
	// Name represents the name of the Cloud Foundry service required by the
	// application. This field represents the runtime name of the service, captured
	// from the 3 different cases where the service name can be listed.
	// For more information check https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#services-block
	Name string `yaml:"name" validate:"required"`
	// Parameters contain the k/v relationship for the aplication to bind to the service
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
	// BindingName captures the name of the service to bind to.
	BindingName string `yaml:"bindingName,omitempty"`
}

type Metadata struct {
	// Name capture the `name` field int CF application manifest
	Name string `yaml:"name" validate:"required"`
	// Space captures the `space` where the CF application is deployed at runtime. The field is empty if the
	// application is discovered directly from the CF manifest. It is equivalent to a Namespace in Kubernetes.
	Space string `yaml:"space,omitempty"`
	// Labels capture the labels as defined in the `annotations` field in the CF application manifest
	Labels map[string]string `yaml:"labels,omitempty"`
	// Annotations capture the annotations as defined in the `labels` field in the CF application manifest
	Annotations map[string]string `yaml:"annotations,omitempty"`
	// Version captures the version of the manifest containing the resulting CF application manifests list retrieved via REST API.
	// Only version 1 is supported at this moment See https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#manifest-schema-version
	// Defaults to 1
	Version string `yaml:"version"`
}

type Process struct {
	// Type captures the `type` field in the Process specification.
	// Accepted values are `web` or `worker`
	Type ProcessType `yaml:"type" validate:"required,oneof=web worker"`
	// Command represents the command used to run the process.
	Command string `yaml:"command,omitempty"`
	// DiskQuota represents the amount of persistent disk requested by the process.
	DiskQuota string `yaml:"disk,omitempty"`
	// Memory represents the amount of memory requested by the process.
	Memory string `yaml:"memory" validate:"required"`
	// HealthCheck captures the health check information
	HealthCheck Probe `yaml:"healthCheck"`
	// ReadinessCheck captures the readiness check information.
	ReadinessCheck Probe `yaml:"readinessCheck"`
	// Instances represents the number of instances for this process to run.
	Instances int `yaml:"instances" validate:"required,min=1"`
	// LogRateLimit represents the maximum amount of logs to be captured per second. Defaults to `16K`
	LogRateLimit string `yaml:"logRateLimit" validate:"required"`
	// Lifecycle captures the value fo the lifecycle field in the CF application manifest.
	// Valid values are `buildpack`, `cnb`, and `docker`. Defaults to `buildpack`
	Lifecycle LifecycleType `yaml:"lifecycle" validate:"required,oneof=buildpack cnb docker"`
}

type LifecycleType string

const (
	BuildPackLifecycleType LifecycleType = "buildpack"
	CNBLifecycleType       LifecycleType = "cnb"
	DockerLifecycleType    LifecycleType = "docker"
)

type ProcessType string

const (
	// Web represents a `web` application type
	Web ProcessType = "web"
	// Worker represents a `worker` application type
	Worker ProcessType = "worker"
)

type Probe struct {
	// Endpoint represents the URL location where to perform the probe check.
	Endpoint string `yaml:"endpoint" validate:"required"`
	// Timeout represents the number of seconds in which the probe check can be considered as timedout.
	// https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#timeout
	Timeout int `yaml:"timeout" validate:"required,min=0"`
	// Interval represents the number of seconds between probe checks.
	Interval int `yaml:"interval" validate:"required,min=0"`
	// Type specifies the type of health check to perform.
	Type ProbeType `yaml:"type" validate:"required,oneof=http process port"`
}

type ProbeType string

const (
	HTTPProbeType    ProbeType = "http"
	TCPProbeType     ProbeType = "tcp"
	ProcessProbeType ProbeType = "process"
)

type Route struct {
	// Route captures the domain name, port and path of the route.
	Route string `yaml:"route" validate:"required"`
	// Protocol captures the protocol type: http, http2 or tcp. Note that the CF `protocol` field is only available
	// for CF deployments that use HTTP/2 routing.
	Protocol RouteProtocol `yaml:"protocol" validate:"required,oneof=http http2 tcp"`
}

type RouteProtocol string

const (
	HTTPRouteProtocol  RouteProtocol = "http"
	HTTP2RouteProtocol RouteProtocol = "http2"
	TCPRouteProtocol   RouteProtocol = "tcp"
)
