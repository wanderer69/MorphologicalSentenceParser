package natashaclient

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type NatashaClient struct {
	isInit  bool
	cli     *client.Client
	address string
	port    int
}

const (
	imageName     = "natasha1:latest"
	internalPort  = "8888"
	containerName = "natasha1_"
)

func NewNatashaClient() *NatashaClient {
	return &NatashaClient{}
}

func (nc *NatashaClient) Init() error {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return err
	}
	nc.cli = cli

	ctx := context.Background()
	// проверяем, что имейж есть
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return err
	}
	isHasImage := false
	for i := range images {
		fmt.Printf("%#v\r\n", images[i])
		if len(images[i].RepoTags) == 0 {
			responses, err := cli.ImageRemove(ctx, images[i].ID, image.RemoveOptions{})
			if err != nil {
				return err
			}
			for j := range responses {
				fmt.Printf("%#v\r\n", responses[j])
			}
			continue
		}
		if len(images[i].RepoTags) == 1 {
			if images[i].RepoTags[0] == imageName {
				isHasImage = true
			}
		}
	}
	if !isHasImage {
		return fmt.Errorf("image %v not found. Please execute 'make_image.sh'", imageName)
	}

	// проверяем, что контейнер запущен, если не запущен - запускаем (план Б - гасим контейней и рестартуем его заново)
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	for i := range containers {
		fmt.Printf("%#v\r\n", containers[i])
		if containers[i].Image == imageName {
			err = cli.ContainerStop(ctx, containers[i].ID, container.StopOptions{})
			if err != nil {
				return err
			}
			err = cli.ContainerRemove(ctx, containers[i].ID, container.RemoveOptions{})
			if err != nil {
				return err
			}
		}
	}
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		ExposedPorts: nat.PortSet{internalPort: struct{}{}},
	}, &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("8188"): {{HostIP: "127.0.0.1", HostPort: internalPort}}},
	}, &network.NetworkingConfig{}, nil, "test1")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	containers, err = cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	for i := range containers {
		fmt.Printf("%#v\r\n", containers[i])
		if containers[i].Image == imageName {
			info, err := cli.ContainerInspect(ctx, containers[i].ID)
			if err != nil {
				return err
			}
			fmt.Printf("info %v\r\n", info)
			nc.address = info.NetworkSettings.IPAddress
		}
	}
	nc.isInit = true
	return nil
}

func (nc *NatashaClient) Close() error {
	if !nc.isInit {
		return nil
	}
	ctx := context.Background()
	containers, err := nc.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return err
	}

	for i := range containers {
		fmt.Printf("%#v\r\n", containers[i])
		if containers[i].Image == imageName {
			err = nc.cli.ContainerStop(ctx, containers[i].ID, container.StopOptions{})
			if err != nil {
				return err
			}
			err = nc.cli.ContainerRemove(ctx, containers[i].ID, container.RemoveOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (nc *NatashaClient) ParsePhrase(ctx context.Context, phrase string) (string, error) {
	if !nc.isInit {
		return "", nil
	}
	transport := http.Transport{}
	path := fmt.Sprintf("http://%v:%v/phrase-translate", nc.address, internalPort)
	type Request struct {
		Phrase string `json:"phrase"`
	}
	request := Request{
		Phrase: phrase,
	}
	data, err := json.Marshal(&request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	client := http.Client{
		Transport: &transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("`%v`\r\n", string(b))

	type Response struct {
		Result string `json:"result"`
	}
	var response Response
	err = json.Unmarshal(b, &response)
	if err != nil {
		return "", err
	}

	dataRaw, err := base64.StdEncoding.DecodeString(response.Result[2 : len(response.Result)-1])
	if err != nil {
		return "", err
	}

	fmt.Printf("%v\r\n", string(dataRaw))

	return string(dataRaw), nil
}
