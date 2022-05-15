package status

import (
	"context"
	"io"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type MemoryStat struct {
	Id       string
	Name     string
	MemUsage int64
	MemLimit int64
	Usage    float64
}

func GetMemUsage() ([]MemoryStat, error) {
	var ms []MemoryStat
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// for _, container := range containers {
	for i := 0; i < len(containers); i++ {
		var m MemoryStat
		m.Id = containers[i].ID
		m.Name = strings.Split(containers[i].Names[0], "/")[1]

		s, err := cli.ContainerStats(ctx, containers[i].ID, false)
		defer s.Body.Close()

		if err != nil {
			return nil, err
		}

		c, err := io.ReadAll(s.Body)
		if err != nil {
			return nil, err
		}
		js, err := simplejson.NewJson(c)
		if err != nil {
			return nil, err
		}

		memUsage := js.Get("memory_stats").Get("usage").MustInt64()
		m.MemUsage = memUsage

		memLimit := js.Get("memory_stats").Get("limit").MustInt64()
		m.MemLimit = memLimit

		m.Usage = float64(memUsage) / float64(memLimit) * 100

		ms = append(ms, m)
	}

	return ms, nil
}
