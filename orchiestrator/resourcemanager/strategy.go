package resourcemanager

const (
	LocalDockerStrategy = "docker"
	HelmChartStrategy   = "helm"
)

type ResourceManager interface {
	Install() error
	UnInstall() error
	Suspend() error
}

var (
	Strategies = []string{LocalDockerStrategy, HelmChartStrategy}

	nameToStrategy = map[string]ResourceManager{
		LocalDockerStrategy: &DockerService{},
		HelmChartStrategy:   &HelmService{},
	}

	DefaultStrategy = LocalDockerStrategy
)

// wybranie odpowiedzniej strategii, potrzebuje dodstać jakiś request params
// -- jeśli dał coś z k8s to helm
func GetResourceManager() (ResourceManager, error) {
	// switch {
	// }
	return nameToStrategy[DefaultStrategy], nil
}
