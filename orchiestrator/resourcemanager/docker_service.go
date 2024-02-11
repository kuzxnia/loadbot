package resourcemanager

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerImage struct {
	// podstawowo obraz
	Name string
}

type DockerContainer struct {
	// podstawowo obraz
	Name string
}

type DockerService struct{}

func (d *DockerService) Install() error   { return nil }
func (d *DockerService) UnInstall() error { return nil }
func (d *DockerService) Suspend() error   { return nil }

// add optional argument with chart version
func NewDockerService() (*DockerService, error) {
	// use default or fetch from internet from tag
	// todo: later add validation for type

	// create docker client
	// ctx := context.Background()
	// cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	// or don't create at all

	return &DockerService{}, nil
}

func (c *DockerService) createContainer(image DockerImage) (err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	out, err := cli.ImagePull(ctx, image.Name, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: image.Name,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return
}

func (c *DockerService) stopContainer(cont DockerContainer) (err error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	if err := cli.ContainerStop(ctx, cont.Name, container.StopOptions{}); err != nil {
		log.Printf("Unable to stop container %s: %s", cont.Name, err)
	}

	return nil
}

func (c *DockerService) removeContainer(cont DockerContainer) (err error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := cli.ContainerRemove(ctx, cont.Name, removeOptions); err != nil {
		log.Printf("Unable to remove container: %s", err)
		return err
	}

	return
}

// func runContainer(client *client.Client, imagename string, containername string, port string, inputEnv []string) error {
// 	// Define a PORT opening
// 	newport, err := natting.NewPort("tcp", port)
// 	if err != nil {
// 		fmt.Println("Unable to create docker port")
// 		return err
// 	}

// 	// Configured hostConfig:
// 	// https://godoc.org/github.com/docker/docker/api/types/container#HostConfig
// 	hostConfig := &container.HostConfig{
// 		PortBindings: natting.PortMap{
// 			newport: []natting.PortBinding{
// 				{
// 					HostIP:   "0.0.0.0",
// 					HostPort: port,
// 				},
// 			},
// 		},
// 		RestartPolicy: container.RestartPolicy{
// 			Name: "always",
// 		},
// 		LogConfig: container.LogConfig{
// 			Type:   "json-file",
// 			Config: map[string]string{},
// 		},
// 	}

// 	// Define Network config (why isn't PORT in here...?:
// 	// https://godoc.org/github.com/docker/docker/api/types/network#NetworkingConfig
// 	networkConfig := &network.NetworkingConfig{
// 		EndpointsConfig: map[string]*network.EndpointSettings{},
// 	}
// 	gatewayConfig := &network.EndpointSettings{
// 		Gateway: "gatewayname",
// 	}
// 	networkConfig.EndpointsConfig["bridge"] = gatewayConfig

// 	// Define ports to be exposed (has to be same as hostconfig.portbindings.newport)
// 	exposedPorts := map[natting.Port]struct{}{
// 		newport: {},
// 	}

// 	// Configuration
// 	// https://godoc.org/github.com/docker/docker/api/types/container#Config
// 	config := &container.Config{
// 		Image:        imagename,
// 		Env:          inputEnv,
// 		ExposedPorts: exposedPorts,
// 		Hostname:     fmt.Sprintf("%s-hostnameexample", imagename),
// 	}

// 	// Creating the actual container. This is "nil,nil,nil" in every example.
// 	cont, err := client.ContainerCreate(
// 		context.Background(),
// 		config,
// 		hostConfig,
// 		networkConfig,
// 		// platform
// 		containername,
// 	)
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	// Run the actual container
// 	client.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
// 	log.Printf("Container %s is created", cont.ID)

// 	return nil
// }

// func main() {
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		log.Fatalf("Unable to create docker client")
// 	}

// 	imagename := "imagename"
// 	containername := "containername"
// 	portopening := "8080"
// 	inputEnv := []string{fmt.Sprintf("LISTENINGPORT=%s", portopening)}
// 	err = runContainer(cli, imagename, containername, portopening, inputEnv)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	// stop and remove
// 	client, err := client.NewEnvClient()
// 	if err != nil {
// 		fmt.Printf("Unable to create docker client: %s", err)
// 	}

// 	// Stops and removes a container
// 	stopAndRemoveContainer(client, "containername")
// }
