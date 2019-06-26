// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster_test

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2018-03-31/containerservice"
	pipCluster "github.com/banzaicloud/pipeline/cluster"
	"github.com/banzaicloud/pipeline/model"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
	"github.com/banzaicloud/pipeline/pkg/cluster/aks"
	pkgErrors "github.com/banzaicloud/pipeline/pkg/errors"
	"github.com/pkg/errors"
)

func TestCreateAKSClusterFromRequest(t *testing.T) {

	testCases := []struct {
		name          string
		createRequest *pkgCluster.CreateClusterRequest
		expectedModel *model.ClusterModel
		expectedError error
	}{
		{name: "aks create with azure cni network", createRequest: aksAzureCniCreateFull, expectedModel: aksModelFullAzureCni, expectedError: nil},
		{name: "aks create with kubenet network", createRequest: aksKubenetCreateFull, expectedModel: aksModelFullKubenet, expectedError: nil},

		{name: "aks empty location", createRequest: aksEmptyLocationCreate, expectedModel: nil, expectedError: pkgErrors.ErrorLocationEmpty},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			commonCluster, err := pipCluster.CreateCommonClusterFromRequest(tc.createRequest, organizationId, userId)

			if tc.expectedError != nil {

				if err != nil {
					if !reflect.DeepEqual(tc.expectedError, err) {
						t.Errorf("Expected model: %v, got: %v", tc.expectedError, err)
					}
				} else {
					t.Errorf("Expected error: %s, but not got error!", tc.expectedError.Error())
					t.FailNow()
				}

			} else {
				if err != nil {
					t.Errorf("Error during CreateCommonClusterFromRequest: %s", err.Error())
					t.FailNow()
				}

				modelAccessor, ok := commonCluster.(interface{ GetModel() *model.ClusterModel })
				if !ok {
					t.Fatal("model cannot be accessed")
				}

				if !reflect.DeepEqual(modelAccessor.GetModel(), tc.expectedModel) {
					t.Errorf("Expected model: %v, got: %v", tc.expectedModel, modelAccessor.GetModel())
				}
			}

		})

	}
}
func TestCreateAKSClusterValidationFromRequest(t *testing.T) {

	testCases := []struct {
		name          string
		createRequest *pkgCluster.CreateClusterRequest
		expectedModel *model.ClusterModel
		expectedError error
	}{
		{name: "aks Dns Service Ip should not end with 1", createRequest: aksAzureCniCreateFullDnsServiceIpShouldNotEndWithOne, expectedModel: nil, expectedError: pkgErrors.ErrorAksNetworkDnsServiceIpInvalidValue},
		{name: "aks Dns Servive Ip not in serviceCIDR", createRequest: aksAzureCniCreateFullDnsServiceIpNotInsideServiceCidr, expectedModel: nil, expectedError: pkgErrors.ErrorAksNetworkDnsServiceIpNotInsideServiceCidr},
		{name: "aks reserved serviceCIDR", createRequest: aksAzureCniCreateFullReservedServiceCidr, expectedModel: nil, expectedError: pkgErrors.ErrorAksNetworkServiceCidrIsReserved},
		{name: "aks Dns Service Ip is not valid", createRequest: aksAzureCniCreateFullDnsServiceIpNotValid, expectedModel: nil, expectedError: pkgErrors.ErrorAksNetworkDnsServiceIpNotValid},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aksCluster, _ := pipCluster.CreateAKSClusterFromRequest(tc.createRequest, organizationId, userId)
			err := aksCluster.ValidateNetwork(tc.createRequest, containerservice.NetworkPlugin(tc.createRequest.Properties.CreateClusterAKS.Network.NetworkPlugin))

			if tc.expectedError != nil {

				if err != nil {
					if !reflect.DeepEqual(tc.expectedError, errors.Cause(err)) {
						t.Errorf("Expected model: %v, got: %v", tc.expectedError, err)
					}
				} else {
					t.Errorf("Expected error: %s, but not got error!", tc.expectedError.Error())
					t.FailNow()
				}

			} else {
				if err != nil {
					t.Errorf("Error during CreateCommonClusterFromRequest: %s", err.Error())
					t.FailNow()
				}
			}

		})

	}
}

var (
	dockerBridgeCidr      = clusterRequestDockerBridgeCidr
	serviceCidr           = clusterRequestServiceCidr
	dnsServiceIp          = clusterRequestDnsServiceIp
	podCidr               = clusterRequestPodCidr
	networkPluginKubenet  = clusterRequestNetworkPluginKubenet
	networkPluginAzureCni = clusterRequestNetworkPluginAzureCNI

	aksAzureCniCreateFull = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      serviceCidr,
			DnsServiceIp:     dnsServiceIp,
			NetworkPlugin:    networkPluginAzureCni,
		})

	aksAzureCniCreateFullReservedServiceCidr = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      "169.254.0.0/16",
			DnsServiceIp:     dnsServiceIp,
			NetworkPlugin:    networkPluginAzureCni,
		})

	aksAzureCniCreateFullDnsServiceIpNotValid = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      serviceCidr,
			DnsServiceIp:     "15.12.12.2/16",
			NetworkPlugin:    networkPluginAzureCni,
		})


	aksAzureCniCreateFullDnsServiceIpNotInsideServiceCidr = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      serviceCidr,
			DnsServiceIp:     "11.0.0.1",
			NetworkPlugin:    networkPluginAzureCni,
		})

	aksAzureCniCreateFullDnsServiceIpShouldNotEndWithOne = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      serviceCidr,
			DnsServiceIp:     "10.0.0.1",
			NetworkPlugin:    networkPluginAzureCni,
		})


	aksKubenetCreateFull = CreateDefaultAksClusterRequestWithNetwork(
		&aks.NetworkCreate{
			DockerBridgeCidr: dockerBridgeCidr,
			ServiceCidr:      serviceCidr,
			DnsServiceIp:     dnsServiceIp,
			PodCidr:          &podCidr,
			NetworkPlugin:    networkPluginKubenet,
		})

	aksEmptyLocationCreate = &pkgCluster.CreateClusterRequest{
		Name:     clusterRequestName,
		Location: "",
		Cloud:    pkgCluster.Azure,
		SecretId: clusterRequestSecretId,
		Properties: &pkgCluster.CreateClusterProperties{
			CreateClusterAKS: &aks.CreateClusterAKS{
				ResourceGroup:     clusterRequestRG,
				KubernetesVersion: clusterRequestKubernetes,
				NodePools: map[string]*aks.NodePoolCreate{
					clusterRequestAgentName: {
						Count:            clusterRequestNodeCount,
						NodeInstanceType: clusterRequestNodeInstance,
					},
				},
				Network: &aks.NetworkCreate{
					DockerBridgeCidr: dockerBridgeCidr,
					ServiceCidr:      serviceCidr,
					DnsServiceIp:     dnsServiceIp,
					PodCidr:          &podCidr,
					NetworkPlugin:    networkPluginKubenet,
				},
			},
		},
	}
)

var (
	aksModelFullKubenet = CreateDefaultAksModelWithNetwork(&networkPluginKubenet,
		&podCidr, &serviceCidr, &dnsServiceIp, &dockerBridgeCidr)

	aksModelFullAzureCni = CreateDefaultAksModelWithNetwork(&networkPluginAzureCni,
		nil, &serviceCidr, &dnsServiceIp, &dockerBridgeCidr)
)

func CreateDefaultAksClusterRequestWithNetwork(network *aks.NetworkCreate) *pkgCluster.CreateClusterRequest {
	return &pkgCluster.CreateClusterRequest{
		Name:     clusterRequestName,
		Location: clusterRequestLocation,
		Cloud:    pkgCluster.Azure,
		SecretId: clusterRequestSecretId,
		Properties: &pkgCluster.CreateClusterProperties{
			CreateClusterAKS: &aks.CreateClusterAKS{
				ResourceGroup:     clusterRequestRG,
				KubernetesVersion: clusterRequestKubernetes,
				NodePools: map[string]*aks.NodePoolCreate{
					clusterRequestAgentName: {
						Autoscaling:      true,
						MinCount:         clusterRequestNodeCount,
						MaxCount:         clusterRequestNodeMaxCount,
						Count:            clusterRequestNodeCount,
						NodeInstanceType: clusterRequestNodeInstance,
						Labels:           clusterRequestNodeLabels,
					},
				},
				Network: network,
			},
		},
	}
}

func CreateDefaultAksModelWithNetwork(networkPlugin *string, podCidr *string, serviceCidr *string,
	dnsServiceIp *string, dockerBridgeCidr *string) *model.ClusterModel {
	return &model.ClusterModel{
		CreatedBy:      userId,
		Name:           clusterRequestName,
		Location:       clusterRequestLocation,
		SecretId:       clusterRequestSecretId,
		Cloud:          pkgCluster.Azure,
		Distribution:   pkgCluster.AKS,
		OrganizationId: organizationId,
		AKS: model.AKSClusterModel{
			ResourceGroup:     clusterRequestRG,
			KubernetesVersion: clusterRequestKubernetes,
			NetworkPlugin:     networkPlugin,
			PodCidr:           podCidr,
			DnsServiceIp:      dnsServiceIp,
			ServiceCidr:       serviceCidr,
			DockerBridgeCidr:  dockerBridgeCidr,
			NodePools: []*model.AKSNodePoolModel{
				{
					CreatedBy:        userId,
					Autoscaling:      true,
					NodeMinCount:     clusterRequestNodeCount,
					NodeMaxCount:     clusterRequestNodeMaxCount,
					Count:            clusterRequestNodeCount,
					NodeInstanceType: clusterRequestNodeInstance,
					Name:             clusterRequestAgentName,
					Labels:           clusterRequestNodeLabels,
				},
			},
		},
	}
}
