package rialto

type WellKnownServiceProperty struct {
    Name string `json:"name"`
    DefaultValue string `json:"defaultValue"`
}

type WellKnownService struct  {
    Type string `json:"type"`
    Properties []WellKnownServiceProperty `json:"properties"`
}

type WellKnownServiceConfig struct {
    Interfaces []WellKnownService `json:"interfaces"`
}

type ServiceInterface struct {
    Name string `json:"name"`
    WellKnownTypeName string `json:"wellKnownType"`
    Properties []ServiceConfigProperty `json:"properties"`
}

type ServiceConfigProperty struct {
    Name string `json:"name"`
    DefaultValue string `json:"defaultValue"`
    WellKnownPropertyName string `json:"wellKnownName"`
}

type ServiceConfig struct {
    Name string `json:"name"`
    Exposes []ServiceInterface `json:"exposes"`
    DependsOn []ServiceInterface `json:"dependsOn"`
}

type ServiceInstance struct {
    Name string `json:"name"`
    ServiceName string `json:"serviceName"`
    ChartSource ChartSource `json:"chartSource"`
    Exposes map[string]ServiceInstanceProperties
}

type ServiceInstanceProperties struct {
    Properties map[string]string
}

type InstallConfig struct {
    Instances []ServiceInstance `json:"instances"`
}

type EnvironmentInstance struct {
    ExternalIPs []string
}

type Deployment struct {
    Environment EnvironmentInstance
    ServiceInstance ServiceInstance
}

type ChartSource struct {
    HelmRepo string `json:"helmRepo"`
    ChartName string `json:"chartName"`
    Url string `json:"url"`
}

