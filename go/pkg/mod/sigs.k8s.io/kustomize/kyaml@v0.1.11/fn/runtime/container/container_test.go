// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package container

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/kyaml/fn/runtime/runtimeutil"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestFilter_setupExec(t *testing.T) {
	var tests = []struct {
		name           string
		functionConfig string
		expectedArgs   []string
		instance       Filter
	}{
		{
			name: "command",
			functionConfig: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
`,
			expectedArgs: []string{
				"run",
				"--rm",
				"-i", "-a", "STDIN", "-a", "STDOUT", "-a", "STDERR",
				"--network", "none",
				"--user", "nobody",
				"--security-opt=no-new-privileges",
			},
			instance: Filter{Image: "example.com:version"},
		},

		{
			name: "network",
			functionConfig: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
`,
			expectedArgs: []string{
				"run",
				"--rm",
				"-i", "-a", "STDIN", "-a", "STDOUT", "-a", "STDERR",
				"--network", "test-1",
				"--user", "nobody",
				"--security-opt=no-new-privileges",
			},
			instance: Filter{Image: "example.com:version", Network: "test-1"},
		},

		{
			name: "storage_mounts",
			functionConfig: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
`,
			expectedArgs: []string{
				"run",
				"--rm",
				"-i", "-a", "STDIN", "-a", "STDOUT", "-a", "STDERR",
				"--network", "none",
				"--user", "nobody",
				"--security-opt=no-new-privileges",
				"--mount", fmt.Sprintf("type=%s,src=%s,dst=%s:ro", "bind", "/mount/path", "/local/"),
				"--mount", fmt.Sprintf("type=%s,src=%s,dst=%s:ro", "volume", "myvol", "/local/"),
				"--mount", fmt.Sprintf("type=%s,src=%s,dst=%s:ro", "tmpfs", "", "/local/"),
			},
			instance: Filter{
				Image: "example.com:version",
				StorageMounts: []runtimeutil.StorageMount{
					{MountType: "bind", Src: "/mount/path", DstPath: "/local/"},
					{MountType: "volume", Src: "myvol", DstPath: "/local/"},
					{MountType: "tmpfs", Src: "", DstPath: "/local/"},
				},
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := yaml.Parse(tt.functionConfig)
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			tt.instance.Exec.FunctionConfig = cfg

			os.Setenv("KYAML_TEST", "FOO")
			tt.instance.setupExec()

			// configure expected env
			for _, e := range os.Environ() {
				// the process env
				tt.expectedArgs = append(tt.expectedArgs, "-e", strings.Split(e, "=")[0])
			}
			tt.expectedArgs = append(tt.expectedArgs, tt.instance.Image)

			if !assert.Equal(t, "docker", tt.instance.Exec.Path) {
				t.FailNow()
			}
			if !assert.Equal(t, tt.expectedArgs, tt.instance.Exec.Args) {
				t.FailNow()
			}
		})
	}
}

func TestFilter_Filter(t *testing.T) {
	cfg, err := yaml.Parse(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: foo
`)
	if !assert.NoError(t, err) {
		return
	}

	input, err := (&kio.ByteReader{Reader: bytes.NewBufferString(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-foo
---
apiVersion: v1
kind: Service
metadata:
  name: service-foo
`)}).Read()
	if !assert.NoError(t, err) {
		return
	}

	instance := Filter{}
	instance.Exec.FunctionConfig = cfg
	instance.Exec.Path = "sed"
	instance.Exec.Args = []string{"s/Deployment/StatefulSet/g"}
	output, err := instance.Filter(input)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	b := &bytes.Buffer{}
	err = kio.ByteWriter{Writer: b, KeepReaderAnnotations: true}.Write(output)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, `apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: deployment-foo
  annotations:
    config.kubernetes.io/index: '0'
    config.kubernetes.io/path: 'statefulset_deployment-foo.yaml'
---
apiVersion: v1
kind: Service
metadata:
  name: service-foo
  annotations:
    config.kubernetes.io/index: '1'
    config.kubernetes.io/path: 'service_service-foo.yaml'
`, b.String()) {
		t.FailNow()
	}
}
func TestFilter_String(t *testing.T) {
	instance := Filter{Image: "foo"}
	if !assert.Equal(t, "foo", instance.String()) {
		t.FailNow()
	}

	instance.Exec.DeferFailure = true
	if !assert.Equal(t, "foo deferFailure: true", instance.String()) {
		t.FailNow()
	}
}

func TestFilter_ExitCode(t *testing.T) {
	instance := Filter{}
	instance.Exec.Path = "/not/real/command"
	instance.Exec.DeferFailure = true
	_, err := instance.Filter(nil)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.EqualError(t, instance.GetExit(), "fork/exec /not/real/command: no such file or directory") {
		t.FailNow()
	}
}
