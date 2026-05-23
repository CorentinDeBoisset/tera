package cfg

import (
	"maps"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type ConfigFile struct {
	BasePath    string                   `yaml:"-"`
	LogFilePath string                   `yaml:"log_file"`
	Jobs        map[string]JobConfig     `yaml:"jobs,omitempty"`
	Services    map[string]ServiceConfig `yaml:"services,omitempty"`
}

type DependencyConfig struct {
	Target            string `yaml:"target"`
	RestartWithTarget bool   `yaml:"restart_with_target"`
	WaitTargetStarted bool   `yaml:"wait_target_restarted"`
}

type HealthcheckConfig struct {
	Port int `yaml:"port"`
	// File       string `yaml:"file"`
	// OutputText string `yaml:"output_text"`
}

type ServiceConfig struct {
	TaskConfig `yaml:"task_config,inline"`

	Healthcheck  HealthcheckConfig  `yaml:"healthcheck"`
	AutoRestart  bool               `yaml:"auto_restart"`
	Dependencies []DependencyConfig `yaml:"dependencies"`

	OpenTarget string `yaml:"open_target"`
}

type JobConfig struct {
	Steps    []StepConfig `yaml:"steps"`
	RunAfter []CmdConfig  `yaml:"run_after,omitempty"`
}

type StepConfig struct {
	Name      string       `yaml:"name"`
	Tasks     []TaskConfig `yaml:"tasks"`
	RunBefore []CmdConfig  `yaml:"run_before,omitempty"`
	RunAfter  []CmdConfig  `yaml:"run_after,omitempty"`
}

type TaskConfig struct {
	CmdConfig `yaml:"cmd_config,inline"`

	Name      string      `yaml:"name"`
	RunBefore []CmdConfig `yaml:"run_before,omitempty"`
	RunAfter  []CmdConfig `yaml:"run_after,omitempty"`
}

type CmdConfig struct {
	Cmd  string `yaml:"cmd"`
	Path string `yaml:"path,omitempty"`
	// TODO: add a FailedWhen: a template calculated with the exit code, the stdout and stderr
}

// findConfig tries to find a configuration file.
// If no path is given in argument, it tries to find a .tera.yml file in the parent directories of the current working directory.
// If a .tera.yml file is found, it checks if it has a sibling .tera.extra.yml file to also import.
func findConfig(givenPath string) (ret []string, err error) {
	configs := make([]string, 0, 1)
	if len(givenPath) > 0 {
		var finalPath string
		if !filepath.IsAbs(givenPath) {
			finalPath = givenPath
		} else {
			finalPath, err = filepath.Abs(givenPath)
			if err != nil {
				return nil, newConfigError("An error occured calculating an absolute path: %s", err)
			}
		}

		stat, err := os.Stat(finalPath)
		if err != nil {
			return nil, newConfigError("An error occured when checking the path \"%s\":\n%s", finalPath, err)
		}
		if stat.IsDir() {
			return nil, newConfigError("The path \"%s\" is a directory", finalPath)
		}

		configs = append(configs, finalPath)

		return configs, nil
	}

	curDir, err := os.Getwd()
	if err != nil {
		return nil, newConfigError("Failed to read the current working directory: %s", err)
	}
	for {
		testedPath := filepath.Join(curDir, ".tera.yml")
		stat, err := os.Stat(testedPath)
		if err == nil && !stat.IsDir() {
			configs = append(configs, testedPath)

			// Check if the config file has an "extra" sibling
			siblingPath := filepath.Join(curDir, ".tera.extra.yml")
			if stat, err := os.Stat(siblingPath); err == nil && !stat.IsDir() {
				configs = append(configs, siblingPath)
			}

			return configs, nil
		}

		if curDir == filepath.Dir(curDir) {
			// We have reached the root directory
			break
		}

		// Go to the parent directory
		curDir = filepath.Dir(curDir)
	}

	return nil, newConfigError("No configuration file could be found")
}

func mergeConfigs(base, extra *ConfigFile) *ConfigFile {
	resultingConfig := &ConfigFile{}

	// This should not be set from the yaml alone, but just in case we copy it anyway
	if len(base.BasePath) > 0 {
		resultingConfig.BasePath = base.BasePath
	}

	// for logFilePath, extra has precedence over base
	if len(extra.LogFilePath) > 0 {
		resultingConfig.LogFilePath = extra.LogFilePath
	} else if len(base.LogFilePath) > 0 {
		resultingConfig.LogFilePath = base.LogFilePath
	}

	resultingConfig.Jobs = make(map[string]JobConfig)
	maps.Copy(resultingConfig.Jobs, base.Jobs)
	maps.Copy(resultingConfig.Jobs, extra.Jobs)

	resultingConfig.Services = make(map[string]ServiceConfig)
	maps.Copy(resultingConfig.Services, base.Services)
	maps.Copy(resultingConfig.Services, extra.Services)

	return resultingConfig
}

func validateCommands(configs []CmdConfig) error {
	for cmdIdx, cmd := range configs {
		if len(cmd.Cmd) == 0 {
			return newConfigError("The task #%d has no command declared", cmdIdx)
		}
	}

	return nil
}

func validateJobConfig(jobId string, job *JobConfig) error {
	if len(job.Steps) == 0 {
		return newConfigError("No step is declared in the job \"%s\"", jobId)
	}

	for stepIdx, step := range job.Steps {
		if len(step.Name) == 0 {
			return newConfigError("The step #%d in the job \"%s\" has no name declared", stepIdx, jobId)
		}

		// Check the hooks
		if err := validateCommands(step.RunBefore); err != nil {
			return newConfigError("The step \"%s\" in the job \"%s\" has invalid run_before hooks: %s", step.Name, jobId, err)
		}
		if err := validateCommands(step.RunAfter); err != nil {
			return newConfigError("The step \"%s\" in the job \"%s\" has invalid run_after hooks: %s", step.Name, jobId, err)
		}

		// Check all the tasks.
		if len(step.Tasks) == 0 {
			return newConfigError("The step \"%s\" in the job \"%s\" has no task declared", step.Name, jobId)
		}
		for taskIdx, task := range step.Tasks {
			if err := validateTaskConfig(task); err != nil {
				if len(task.Name) > 0 {
					return newConfigError("The task \"%s\" in the step \"%s\" in the job \"%s\" is invalid: %s", task.Name, step.Name, jobId, err)
				}
				return newConfigError("The task #%d in the step \"%s\" in the job \"%s\" is invalid: %s", taskIdx, step.Name, jobId, err)
			}
		}
	}

	return nil
}

func validateTaskConfig(task TaskConfig) error {
	if len(task.Name) == 0 {
		return newConfigError("No name is declared")
	}

	if err := validateCommands(task.RunBefore); err != nil {
		return newConfigError("The run_before hooks are invalid: %s", err)
	}
	if err := validateCommands(task.RunAfter); err != nil {
		return newConfigError("The run_after hooks are invalid: %s", err)
	}
	if len(task.Cmd) == 0 {
		return newConfigError("No command is declared")
	}

	return nil
}

func validateConfig(cfg *ConfigFile) error {
	if len(cfg.Jobs) == 0 && len(cfg.Services) == 0 {
		return newConfigError("No job and no service is declared in the configuration")
	}

	for jobId, job := range cfg.Jobs {
		if err := validateJobConfig(jobId, &job); err != nil {
			return err
		}
	}

	for serviceId, service := range cfg.Services {
		if len(serviceId) == 0 {
			return newConfigError("The key of a service is not defined")
		}

		if err := validateTaskConfig(service.TaskConfig); err != nil {
			return newConfigError("The service \"%s\" is invalid: %s", serviceId, err)
		}
	}

	return nil
}

func ParseConfigs(fileContents [][]byte) (*ConfigFile, error) {
	var resultingConfig *ConfigFile
	for _, fileContent := range fileContents {
		rawConfig := &ConfigFile{}
		if err := yaml.Unmarshal(fileContent, rawConfig); err != nil {
			return nil, newConfigError("The file could not be parsed from YAML: %s", err.Error())
		}

		if resultingConfig == nil {
			resultingConfig = rawConfig
		} else {
			resultingConfig = mergeConfigs(resultingConfig, rawConfig)
		}
	}

	if err := validateConfig(resultingConfig); err != nil {
		return nil, newConfigError("The configuration is invalid: %s", err)
	}

	return resultingConfig, nil
}

func FindAndParseConfig(givenPath string) (*ConfigFile, error) {
	configPaths, err := findConfig(givenPath)
	if err != nil {
		return nil, err
	}
	if len(configPaths) < 1 {
		return nil, newConfigError("No valid configuration could be found")
	}

	configContents := make([][]byte, 0, len(configPaths))
	for _, configPath := range configPaths {
		fileContent, err := os.ReadFile(configPath)
		if err != nil {
			return nil, newConfigError("The contents of the file \"%s\" could not be read: %s", configPath, err)
		}
		configContents = append(configContents, fileContent)
	}

	config, err := ParseConfigs(configContents)
	if err != nil {
		return nil, err
	}

	config.BasePath = filepath.Dir(configPaths[0])

	return config, nil
}
