package docker_compose_client

import (
	"context"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	commands "github.com/docker/compose/v2/cmd/compose"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"sync"
)

var (
	once           sync.Once
	composeCommand *cobra.Command
)

func newDockerCli() (*command.DockerCli, error) {
	initOpt := command.WithInitializeClient(func(dockerCli *command.DockerCli) (client.APIClient, error) {
		return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	})
	cli, err := command.NewDockerCli()
	if err != nil {
		return nil, err
	}
	err = cli.Initialize(flags.NewClientOptions(), initOpt)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func initComposeCommand() error {
	var err error
	once.Do(func() {
		dockerCli, err := newDockerCli()
		if err != nil {
			return
		}
		lazyInit := api.NewServiceProxy()
		lazyInit.WithService(compose.NewComposeService(dockerCli))
		composeCommand = commands.RootCommand(dockerCli, lazyInit)
	})
	return err
}

func Compose(ctx context.Context, args []string) error {
	err := initComposeCommand()
	if err != nil {
		return err
	}
	composeCommand.SetContext(ctx)
	composeCommand.SetArgs(args)
	err = composeCommand.Execute()
	if err != nil {
		return err
	}
	return nil
}
