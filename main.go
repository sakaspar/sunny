package main

import (
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
    corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
    metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
    storagev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/storage/v1"
    //"os/exec"
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

		// Run k3d cluster create command using Docker
		clusterCreateCommand := pulumi.Sprintf(
			`docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
				rancher/k3s:v1.22.2-k3s1 \
				k3d cluster create mycluster --servers 1 --agents 1 \
				--k3s-arg '--disable=traefik@server:*' \
				--k3s-arg '--disable=servicelb@server:*' \
				--k3s-arg '--disable=metrics-server@server:*' \
				--k3s-arg '--disable-cloud-controller@server:0' \
				--api-port 6443 \
				--k3d-arg '--volume=/mnt/data:/mnt/data@all' \
				--network %s \
				--k3s-arg '--node-ip=192.168.1.183' \
				--k3s-arg '--kube-apiserver-arg=service-node-port-range=8000-32767@server:0'`,
			network.Name,
		)

		_, err = docker.NewContainer(ctx, "k3d-cluster", &docker.ContainerArgs{
			Image: pulumi.String("docker"),
			Command: pulumi.StringArray{
				pulumi.String("sh"),
				pulumi.String("-c"),
				clusterCreateCommand,
			},
		})
		if err != nil {
			return err
		}

		// Output the network information
		ctx.Export("k3dNetworkID", network.ID())
		ctx.Export("k3dNetworkName", network.Name)

		
		/*
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
					ctx.Export("k3dClusterStdout", k3dCluster.Stdout)*/

        // Define and provision the Kubernetes resources
        _, err = corev1.NewPersistentVolume(ctx, "pv", &corev1.PersistentVolumeArgs{
            Metadata: &metav1.ObjectMetaArgs{
                Name: pulumi.String("pv"),
            },
            Spec: &corev1.PersistentVolumeSpecArgs{
                Capacity: pulumi.StringMap{
                    "storage": pulumi.String("5Gi"),
                },
                AccessModes: pulumi.StringArray{
                    pulumi.String("ReadWriteOnce"),
                },
                Nfs: &corev1.NFSVolumeSourceArgs{
                    Server:   pulumi.String("192.168.1.183"),
                    Path:     pulumi.String("/mnt/data"),
                    ReadOnly: pulumi.Bool(false),
                },
            },
        })
        if err != nil {
            return err
        }

        _, err = corev1.NewPersistentVolumeClaim(ctx, "pvc", &corev1.PersistentVolumeClaimArgs{
            Metadata: &metav1.ObjectMetaArgs{
                Name: pulumi.String("pvc"),
            },
            Spec: &corev1.PersistentVolumeClaimSpecArgs{
                AccessModes: pulumi.StringArray{
                    pulumi.String("ReadWriteOnce"),
                },
                Resources: &corev1.ResourceRequirementsArgs{
                    Requests: pulumi.StringMap{
                        "storage": pulumi.String("10Gi"),
                    },
                },
            },
        })
        if err != nil {
            return err
        }

        // Uncomment to add additional Kubernetes resources
        /*
        _, err = networkingv1.NewNetworkPolicy(ctx, "networkPolicy", &networkingv1.NetworkPolicyArgs{
            Metadata: &metav1.ObjectMetaArgs{
                Name: pulumi.String("networkPolicy"),
            },
            Spec: &networkingv1.NetworkPolicySpecArgs{
                PodSelector: metav1.LabelSelectorArgs{
                    MatchLabels: pulumi.StringMap{"app": pulumi.String("myapp")},
                },
                Ingress: networkingv1.NetworkPolicyIngressRuleArray{
                    &networkingv1.NetworkPolicyIngressRuleArgs{
                        From: networkingv1.NetworkPolicyPeerArray{
                            &networkingv1.NetworkPolicyPeerArgs{
                                PodSelector: &metav1.LabelSelectorArgs{},
                            },
                        },
                    },
                },
            },
        })
        if err != nil {
            return err
        }
        */

        _, err = storagev1.NewStorageClass(ctx, "storageClass", &storagev1.StorageClassArgs{
            Metadata: &metav1.ObjectMetaArgs{
                Name: pulumi.String("storageClass"),
            },
            Provisioner:        pulumi.String("kubernetes.io/no-provisioner"),
            VolumeBindingMode: pulumi.StringPtr("WaitForFirstConsumer"),
        })
        if err != nil {
            return err
        }

        return nil
    })
}
