/*
Copyright 2019 kubeflow.org.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
)

func TestCreateModelServingContainer(t *testing.T) {

	var requestedResource = v1.ResourceRequirements{
		Limits: v1.ResourceList{
			"cpu": resource.Quantity{
				Format: "100",
			},
		},
		Requests: v1.ResourceList{
			"cpu": resource.Quantity{
				Format: "90",
			},
		},
	}
	var config = FrameworksConfig{
		TensorRT: FrameworkConfig{
			ContainerImage: "someOtherImage",
		},
	}
	var spec = TensorRTSpec{
		StorageURI:     "gs://someUri",
		Resources:      requestedResource,
		RuntimeVersion: "19.05-py3",
	}
	g := gomega.NewGomegaWithT(t)

	expectedContainer := &v1.Container{
		Image:     "someOtherImage:19.05-py3",
		Resources: requestedResource,
		Args: []string{
			"trtserver",
			"--model-store=/mnt/models",
			"--allow-poll-model-repository=false",
			"--allow-grpc=true",
			"--allow-http=true",
			"--grpc-port=9000",
			"--http-port=8080",
		},
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				ContainerPort: 8080,
			},
		},
	}

	// Test Create without config
	container := spec.CreateModelServingContainer("someName", &config)
	g.Expect(container).To(gomega.Equal(expectedContainer))

	// Test Create with config
	expectedContainer.Image = "nvcr.io/nvidia/tensorrtserver:19.05-py3"
	emptyConfig := FrameworksConfig{TensorRT: FrameworkConfig{}}
	container = spec.CreateModelServingContainer("someName", &emptyConfig)
	g.Expect(container).To(gomega.Equal(expectedContainer))
}
