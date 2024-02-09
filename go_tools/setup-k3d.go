package main

import (
        "fmt"
        "github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
        "github.com/pulumi/pulumi-command/sdk/v4/go/command"
        "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
        pulumi.Run(func(ctx *pulumi.Context) error {
                // Create a Docker network for the k3d cluster
                network, err := docker.NewNetwork(ctx, "k3d-network", &docker.NetworkArgs{
                        Name: pulumi.String("k3d-net"),
                })
                if err != nil {
                        return err
                }

                // Create a k3d cluster using a local command
                clusterCreateCommand := fmt.Sprintf(
                        "k3d cluster create mycluster --servers 1 --agents 1 "+
                                "--k3s-arg '--disable=traefik@server:*' "+
                                "--k3s-arg '--disable=servicelb@server:*' "+
                                "--k3s-arg '--disable=metrics-server@server:*' "+
                                "--k3s-arg '--disable-cloud-controller@server:0' "+
                                "--api-port 6443 "+
                                "--k3d-arg '--volume=/mnt/data:/mnt/data@all' "+
                                "--network %s "+
                                "--k3s-arg '--node-ip=192.168.1.183' "+
                                "--k3s-arg '--kube-apiserver-arg=service-node-port-range=8000-32767@server:0'",
                        network.Name,
                )

                k3dCluster, err := command.NewLocalCommand(ctx, "k3d-cluster", &command.LocalCommandArgs{
                        Create:   pulumi.String(clusterCreateCommand),
                        Delete:   pulumi.String("k3d cluster delete mycluster"),
                        Triggers: pulumi.StringArray([]pulumi.StringInput{network.Name}),
                })
                if err != nil {
                        return err
                }

                // Output the network and cluster information
                ctx.Export("k3dNetworkID", network.ID())
                ctx.Export("k3dNetworkName", network.Name)
                ctx.Export("k3dClusterStdout", k3dCluster.Stdout)

                return nil
        })
}
