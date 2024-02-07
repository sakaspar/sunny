package main

import (
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	networkingv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/networking/v1"
	storagev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := corev1.NewPersistentVolume(ctx, "pv", &corev1.PersistentVolumeArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.StringPtr("pv"),
			},
			Spec: &corev1.PersistentVolumeSpecArgs{
				Capacity: pulumi.StringMap{
					"storage": pulumi.String("10Gi"),
				},
				AccessModes: pulumi.StringArray{
					pulumi.String("ReadWriteOnce"),
				},
				PersistentVolumeSource: corev1.PersistentVolumeSourceArgs{
					Nfs: &corev1.NFSVolumeSourceArgs{
						Server: pulumi.String("192.168.1.183"),
						Path:   pulumi.String("/mnt/data"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = corev1.NewPersistentVolumeClaim(ctx, "pvc", &corev1.PersistentVolumeClaimArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.StringPtr("pvc"),
			},
			Spec: &corev1.PersistentVolumeClaimSpecArgs{
				AccessModes: pulumi.StringArray{
					pulumi.String("ReadWriteOnce"),
				},
				Resources: corev1.ResourceRequirementsArgs{
					Requests: pulumi.StringMap{
						"storage": pulumi.String("10Gi"),
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = networkingv1.NewNetworkPolicy(ctx, "networkPolicy", &networkingv1.NetworkPolicyArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.StringPtr("networkPolicy"),
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

		_, err = storagev1.NewStorageClass(ctx, "storageClass", &storagev1.StorageClassArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name: pulumi.StringPtr("storageClass"),
			},
			Provisioner: pulumi.String("kubernetes.io/no-provisioner"),
			VolumeBindingMode: pulumi.StringPtr("WaitForFirstConsumer"),
		})
		if err != nil {
			return err
		}
		
		return nil
	})
}

