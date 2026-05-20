package cfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseConfigString(input string) (*ConfigFile, error) {
	return ParseConfigs([][]byte{[]byte(input)})
}

func TestConfigParse(t *testing.T) {
	t.Parallel()

	sampleConfig := `
jobs:
  first-job:
    steps:
      - name: first_step
        run_before:
          - name: my command
            cmd: my-command
        tasks:
          - name: echoes
            cmd: echo 12 && echo 13

          - name: sleep
            cmd: sleep 5
        run_after:
          - name: other command
            cmd: other-command

      - name: second_step
        tasks:
          - name: echoes 1
            cmd: echo 12 && echo 13

services:
  serviceA:
    name: Service A
    cmd: echo 'A'
    path: .
    healthcheck:
      port: 4242
    open_target: "http://localhost:4242/something"
    dependencies:
      - target: ServiceB


  serviceB:
    name: Service B
    cmd: echo 'B'
    path: .
`

	config, err := parseConfigString(sampleConfig)
	assert.Nil(t, err)
	assert.Len(t, config.Jobs, 1)
	assert.Contains(t, config.Jobs, "first-job")

	assert.Len(t, config.Jobs["first-job"].Steps, 2)

	assert.Len(t, config.Jobs["first-job"].Steps[0].RunBefore, 1)
	assert.Equal(t, config.Jobs["first-job"].Steps[0].RunBefore[0].Cmd, "my-command")
	assert.Len(t, config.Jobs["first-job"].Steps[0].Tasks, 2)
	assert.Equal(t, config.Jobs["first-job"].Steps[0].Tasks[1].Cmd, "sleep 5")
	assert.Len(t, config.Jobs["first-job"].Steps[0].RunAfter, 1)

	assert.Len(t, config.Jobs["first-job"].Steps[1].RunBefore, 0)
	assert.Len(t, config.Jobs["first-job"].Steps[1].Tasks, 1)
	assert.Len(t, config.Jobs["first-job"].Steps[1].RunAfter, 0)

	assert.Len(t, config.Services, 2)
	assert.Equal(t, config.Services["serviceA"].Name, "Service A")
	assert.Equal(t, config.Services["serviceA"].Cmd, "echo 'A'")
	assert.Equal(t, config.Services["serviceA"].Healthcheck.Port, 4242)
	assert.Equal(t, config.Services["serviceA"].OpenTarget, "http://localhost:4242/something")
	assert.Equal(t, config.Services["serviceB"].Name, "Service B")
	assert.Len(t, config.Services["serviceA"].Dependencies, 1)
	assert.Len(t, config.Services["serviceB"].Dependencies, 0)
}

func TestGlobalErrors(t *testing.T) {
	t.Parallel()

	invalidYaml := `some plaintext`
	_, err := parseConfigString(invalidYaml)
	assert.ErrorContains(t, err, "The file could not be parsed from YAML")

	emptyConfig := `
jobs: {}
services:
`
	_, err = parseConfigString(emptyConfig)
	assert.ErrorContains(t, err, "No job and no service is declared in the configuration")
}

func TestJobErrors(t *testing.T) {
	t.Parallel()

	emptyJobConfig := `
jobs:
  first-job:
    steps: []
`
	_, err := parseConfigString(emptyJobConfig)
	assert.ErrorContains(t, err, "No step is declared in the job \"first-job\"")

	noCommand := `
jobs:
  first-job:
    steps:
        - name: First step
          tasks:
            - name: first_task
`
	_, err = parseConfigString(noCommand)
	assert.ErrorContains(t, err, "The task \"first_task\" in the step \"First step\" in the job \"first-job\" is invalid: No command is declared")
}

func TestMergeConfigs(t *testing.T) {
	t.Parallel()

	baseConfig := `
jobs:
  first-job:
    steps:
      - name: first_step
        run_before:
          - name: my command
            cmd: my-command
        tasks:
          - name: echoes
            cmd: echo 12 && echo 13

          - name: sleep
            cmd: sleep 5
        run_after:
          - name: other command
            cmd: other-command

      - name: second_step
        tasks:
          - name: echoes 1
            cmd: echo 12 && echo 13

  second-job:
    steps:
      - name: other_step
        tasks:
          - name: echoes 1234
            cmd: echo 1234

services:
  serviceA:
    name: Service A
    cmd: echo 'A'
    path: .
    healthcheck:
      port: 4242
    open_target: "http://localhost:4242/something"
    dependencies:
      - target: ServiceB


  serviceB:
    name: Service B
    cmd: echo 'B'
    path: .
`

	extraConfig := `
jobs:
  first-job:
    steps:
      - name: override_step
        tasks:
          - name: echoes 4567
            cmd: echo 4567

  extra-job:
    steps:
      - name: something-else
        tasks:
          - name: yes
            cmd: exit 0

services:
  serviceB:
    name: Service B override
    cmd: exit 123
    path: ..

  extraService:
    name: Extra service
    cmd: exit 123
    path: ..
    dependencies:
      - target: ServiceA
`

	mergedConfig, err := ParseConfigs([][]byte{[]byte(baseConfig), []byte(extraConfig)})
	assert.Nil(t, err)
	assert.Len(t, mergedConfig.Jobs, 3)
	assert.Contains(t, mergedConfig.Jobs, "first-job")
	assert.Contains(t, mergedConfig.Jobs, "second-job")
	assert.Equal(t, mergedConfig.Jobs["first-job"].Steps[0].Tasks[0].Name, "echoes 4567")
	assert.Contains(t, mergedConfig.Jobs, "extra-job")

	assert.Len(t, mergedConfig.Services, 3)
	assert.Contains(t, mergedConfig.Services, "serviceA")
	assert.Contains(t, mergedConfig.Services, "serviceB")
	assert.Equal(t, mergedConfig.Services["serviceB"].Name, "Service B override")
	assert.Contains(t, mergedConfig.Services, "extraService")
}
