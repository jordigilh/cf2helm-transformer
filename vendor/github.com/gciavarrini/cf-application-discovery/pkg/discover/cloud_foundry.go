package discover

import (
	"encoding/json"

	"github.com/gciavarrini/cf-application-discovery/pkg/models"
)

func Discover(cfApp AppManifest, version, space string) (models.Application, error) {
	appVersion := "1"
	if version != "" {
		appVersion = version

	}
	timeout := 60
	if cfApp.Timeout != 0 {
		timeout = int(cfApp.Timeout)
	}
	var instances int = 1
	if cfApp.Instances != nil {
		instances = int(*cfApp.Instances)
	}
	services := parseServices(cfApp.Services)
	routeSpec := parseRouteSpec(cfApp.Routes, cfApp.RandomRoute, cfApp.NoRoute)
	docker := parseDocker(cfApp.Docker)
	sidecars := parseSidecars(cfApp.Sidecars)
	processes, err := parseProcesses(cfApp)
	if err != nil {
		return models.Application{}, err
	}
	var labels, annotations map[string]*string

	if cfApp.Metadata != nil {
		labels = cfApp.Metadata.Labels
		annotations = cfApp.Metadata.Annotations
	}

	return models.Application{
		Metadata: models.Metadata{
			Version:     appVersion,
			Name:        cfApp.Name,
			Labels:      labels,
			Annotations: annotations,
			Space:       space,
		},
		Timeout:    timeout,
		Instances:  instances,
		BuildPacks: cfApp.Buildpacks,
		Env:        cfApp.Env,
		Stack:      cfApp.Stack,
		Services:   services,
		Routes:     routeSpec,
		Docker:     docker,
		Sidecars:   sidecars,
		Processes:  processes,
	}, nil
}

func parseHealthCheck(cfType AppHealthCheckType, cfEndpoint string, cfInterval, cfTimeout uint) models.ProbeSpec {
	t := models.PortProbeType
	if len(cfType) > 0 {
		t = models.ProbeType(cfType)
	}
	endpoint := "/"
	if len(cfEndpoint) > 0 {
		endpoint = cfEndpoint
	}
	timeout := 1
	if cfTimeout != 0 {
		timeout = int(cfTimeout)
	}
	interval := 30
	if cfInterval > 0 {
		interval = int(cfInterval)
	}
	return models.ProbeSpec{
		Type:     t,
		Endpoint: endpoint,
		Timeout:  timeout,
		Interval: interval,
	}
}

func parseReadinessHealthCheck(cfType AppHealthCheckType, cfEndpoint string, cfInterval, cfTimeout uint) models.ProbeSpec {
	t := models.ProcessProbeType
	if len(cfType) > 0 {
		t = models.ProbeType(cfType)
	}
	endpoint := "/"
	if len(cfEndpoint) > 0 {
		endpoint = cfEndpoint
	}
	timeout := 1
	if cfTimeout != 0 {
		timeout = int(cfTimeout)
	}
	interval := 30
	if cfInterval > 0 {
		interval = int(cfInterval)
	}
	return models.ProbeSpec{
		Type:     t,
		Endpoint: endpoint,
		Timeout:  timeout,
		Interval: interval,
	}
}

func parseProcesses(cfApp AppManifest) (models.Processes, error) {
	processes := models.Processes{}
	if cfApp.Processes == nil {
		return nil, nil
	}
	for _, cfProcess := range *cfApp.Processes {
		processes = append(processes, parseProcess(cfProcess))
	}
	if cfApp.Type != "" {
		// Type is the only mandatory field for the process.
		// https://github.com/SchemaStore/schemastore/blob/c06e2183289c50bdb0816050dfec002e5ebd8477/src/schemas/json/cloudfoundry-application-manifest.json#L280
		// If it's not defined it means there is no process spec at the application field level and we should return an empty structure
		proc, err := parseInlinedProcessSpec(cfApp)
		if err != nil {
			return nil, err
		}
		processes = append(processes, parseProcess(proc))
	}
	return processes, nil
}

func parseInlinedProcessSpec(cfApp AppManifest) (AppManifestProcess, error) {
	cfProc := AppManifestProcess{}
	b, err := json.Marshal(cfApp)
	if err != nil {
		return cfProc, err
	}
	err = json.Unmarshal(b, &cfProc)
	return cfProc, err
}

func parseProcess(cfProcess AppManifestProcess) models.ProcessSpec {
	memory := "1G"
	if len(cfProcess.Memory) != 0 {
		memory = cfProcess.Memory
	}
	instances := 1
	if cfProcess.Instances != nil {
		instances = int(*cfProcess.Instances)
	}
	logRateLimit := "16K"
	if len(cfProcess.LogRateLimitPerSecond) > 0 {
		logRateLimit = cfProcess.LogRateLimitPerSecond
	}
	p := models.ProcessSpec{
		Type:           models.ProcessType(cfProcess.Type),
		Command:        cfProcess.Command,
		DiskQuota:      cfProcess.DiskQuota,
		Memory:         memory,
		HealthCheck:    parseHealthCheck(cfProcess.HealthCheckType, cfProcess.HealthCheckHTTPEndpoint, cfProcess.HealthCheckInterval, cfProcess.HealthCheckInvocationTimeout),
		ReadinessCheck: parseReadinessHealthCheck(cfProcess.ReadinessHealthCheckType, cfProcess.ReadinessHealthCheckHttpEndpoint, cfProcess.ReadinessHealthCheckInterval, cfProcess.ReadinessHealthInvocationTimeout),
		Instances:      instances,
		LogRateLimit:   logRateLimit,
		Lifecycle:      models.LifecycleType(cfProcess.Lifecycle),
	}
	return p
}

func parseProcessTypes(cfProcessTypes []AppProcessType) []models.ProcessType {
	types := []models.ProcessType{}
	for _, cfType := range cfProcessTypes {
		types = append(types, models.ProcessType(cfType))
	}
	return types

}
func parseSidecars(cfSidecars *AppManifestSideCars) models.Sidecars {
	sidecars := models.Sidecars{}
	if cfSidecars == nil {
		return nil
	}
	for _, cfSidecar := range *cfSidecars {
		pt := parseProcessTypes(cfSidecar.ProcessTypes)
		s := models.SidecarSpec{
			Name:         cfSidecar.Name,
			Command:      cfSidecar.Command,
			ProcessTypes: pt,
			Memory:       cfSidecar.Memory,
		}
		sidecars = append(sidecars, s)
	}
	return sidecars
}

func parseDocker(cfDocker *AppManifestDocker) models.Docker {
	if cfDocker == nil {
		return models.Docker{}
	}
	return models.Docker{
		Image:    cfDocker.Image,
		Username: cfDocker.Username,
	}
}
func parseServices(cfServices *AppManifestServices) models.Services {
	services := models.Services{}
	if cfServices == nil {
		return nil
	}
	for _, svc := range *cfServices {
		s := models.ServiceSpec{
			Name:        svc.Name,
			Parameters:  svc.Parameters,
			BindingName: svc.BindingName,
		}
		services = append(services, s)
	}
	return services
}

func parseRouteSpec(cfRoutes *AppManifestRoutes, randomRoute, noRoute bool) models.RouteSpec {
	if noRoute {
		return models.RouteSpec{
			NoRoute: noRoute,
		}
	}
	routeSpec := models.RouteSpec{
		RandomRoute: randomRoute,
	}

	if cfRoutes == nil {
		return routeSpec
	}

	routeSpec.Routes = parseRoutes(*cfRoutes)
	return routeSpec
}

func parseRoutes(cfRoutes AppManifestRoutes) models.Routes {
	if cfRoutes == nil {
		return nil
	}
	routes := models.Routes{}
	for _, cfRoute := range cfRoutes {
		options := models.RouteOptions{}
		if cfRoute.Options != nil {
			options.LoadBalancing = models.LoadBalancingType(cfRoute.Options.LoadBalancing)
		}
		r := models.Route{
			Route:    cfRoute.Route,
			Protocol: models.RouteProtocol(cfRoute.Protocol),
			Options:  options,
		}
		routes = append(routes, r)
	}
	return routes
}
