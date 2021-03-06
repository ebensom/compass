package gqlschema

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/kyma-incubator/compass/components/provisioner/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRuntimeConfig_UnmarshalJSON(t *testing.T) {

	gcpClusterConfig := GCPConfig{
		Name:              util.StringPtr("name"),
		ProjectName:       util.StringPtr("project"),
		KubernetesVersion: util.StringPtr("1.17"),
		NumberOfNodes:     util.IntPtr(3),
		BootDiskSizeGb:    util.IntPtr(50),
		MachineType:       util.StringPtr("machine"),
		Region:            util.StringPtr("region"),
		Zone:              util.StringPtr("zone"),
	}

	azureProviderCfgNoZones := &AzureProviderConfig{VnetCidr: util.StringPtr("10.10.11.11/25")}
	azureProviderCfg := &AzureProviderConfig{VnetCidr: util.StringPtr("10.10.11.11/25"), Zones: []string{"az-zone-1", "az-zone-2"}}
	gcpProviderCfg := &GCPProviderConfig{Zones: []string{"gcp-zone-1", "gcp-zone-2"}}
	awsProviderCfg := &AWSProviderConfig{
		Zone:         util.StringPtr("aws zone"),
		VpcCidr:      util.StringPtr("10.10.10.11/25"),
		PublicCidr:   util.StringPtr("10.10.10.12/25"),
		InternalCidr: util.StringPtr("10.10.10.13/25"),
	}

	for _, testCase := range []struct {
		description string
		runtimeCfg  RuntimeConfig
	}{
		{
			description: "gardener cluster with Azure with no zones passed",
			runtimeCfg: RuntimeConfig{
				ClusterConfig: newGardenerClusterCfg(fixGardenerConfig("azure"), azureProviderCfgNoZones),
				KymaConfig:    &KymaConfig{Version: util.StringPtr("my precious")},
				Kubeconfig:    util.StringPtr("kubeconfig"),
			},
		},
		{
			description: "gardener cluster with Azure",
			runtimeCfg: RuntimeConfig{
				ClusterConfig: newGardenerClusterCfg(fixGardenerConfig("azure"), azureProviderCfg),
				KymaConfig:    &KymaConfig{Version: util.StringPtr("my favourite")},
				Kubeconfig:    util.StringPtr("kubeconfig"),
			},
		},
		{
			description: "gardener cluster with GCP",
			runtimeCfg: RuntimeConfig{
				ClusterConfig: newGardenerClusterCfg(fixGardenerConfig("gcp"), gcpProviderCfg),
				KymaConfig:    &KymaConfig{Version: util.StringPtr("my favourite")},
				Kubeconfig:    util.StringPtr("kubeconfig"),
			},
		},
		{
			description: "gardener cluster with AWS",
			runtimeCfg: RuntimeConfig{
				ClusterConfig: newGardenerClusterCfg(fixGardenerConfig("aws"), awsProviderCfg),
				KymaConfig:    &KymaConfig{Version: util.StringPtr("my favourite")},
				Kubeconfig:    util.StringPtr("kubeconfig"),
			},
		},
		{
			description: "GCP cluster",
			runtimeCfg: RuntimeConfig{
				ClusterConfig: &gcpClusterConfig,
				KymaConfig:    &KymaConfig{Version: util.StringPtr("my favourite")},
				Kubeconfig:    util.StringPtr("kubeconfig"),
			},
		},
	} {
		t.Run(testCase.description, func(t *testing.T) {
			// given
			marshalled, err := json.Marshal(testCase.runtimeCfg)
			require.NoError(t, err)

			var unmarshalledConfig RuntimeConfig

			// when
			err = json.NewDecoder(bytes.NewBuffer(marshalled)).Decode(&unmarshalledConfig)
			require.NoError(t, err)

			// then
			assert.Equal(t, testCase.runtimeCfg, unmarshalledConfig)
		})
	}

}

func newGardenerClusterCfg(gardenerCfg GardenerConfig, providerCfg ProviderSpecificConfig) ClusterConfig {
	gardenerCfg.ProviderSpecificConfig = providerCfg

	return &gardenerCfg
}

func fixGardenerConfig(providerName string) GardenerConfig {
	return GardenerConfig{
		Name:              util.StringPtr("name"),
		KubernetesVersion: util.StringPtr("1.16"),
		VolumeSizeGb:      util.IntPtr(50),
		MachineType:       util.StringPtr("machine"),
		Region:            util.StringPtr("region"),
		Provider:          util.StringPtr(providerName),
		Seed:              util.StringPtr("seed"),
		TargetSecret:      util.StringPtr("secret"),
		DiskType:          util.StringPtr("disk"),
		WorkerCidr:        util.StringPtr("10.10.10.10/25"),
		AutoScalerMin:     util.IntPtr(1),
		AutoScalerMax:     util.IntPtr(4),
		MaxSurge:          util.IntPtr(25),
		MaxUnavailable:    util.IntPtr(2),
	}
}
