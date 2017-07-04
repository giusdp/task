package task

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

// ReadTaskfile parses Taskfile from the disk
func (e *Executor) ReadTaskfile() error {
	path := filepath.Join(e.Dir, TaskFilePath)

	var err error
	e.Tasks, err = e.readTaskfileData(path)
	if err != nil {
		return err
	}

	osTasks, err := e.readTaskfileData(fmt.Sprintf("%s_%s", path, runtime.GOOS))
	if err != nil {
		switch err.(type) {
		case taskFileNotFound:
		default:
			return err
		}
	}
	if err := mergo.MapWithOverwrite(&e.Tasks, osTasks); err != nil {
		return err
	}
	if err := e.readTaskvarsFile(); err != nil {
		return err
	}
	return nil
}

func (e *Executor) readTaskfileData(path string) (tasks map[string]*Task, err error) {
	if b, err := ioutil.ReadFile(path + ".yml"); err == nil {
		return tasks, yaml.Unmarshal(b, &tasks)
	}
	if b, err := ioutil.ReadFile(path + ".json"); err == nil {
		return tasks, json.Unmarshal(b, &tasks)
	}
	if b, err := ioutil.ReadFile(path + ".toml"); err == nil {
		return tasks, toml.Unmarshal(b, &tasks)
	}
	return nil, taskFileNotFound{path}
}

func (e *Executor) readTaskvarsFile() error {
	file := filepath.Join(e.Dir, TaskvarsFilePath)

	if b, err := ioutil.ReadFile(file + ".yml"); err == nil {
		if err := yaml.Unmarshal(b, &e.taskvars); err != nil {
			return err
		}
		return nil
	}
	if b, err := ioutil.ReadFile(file + ".json"); err == nil {
		if err := json.Unmarshal(b, &e.taskvars); err != nil {
			return err
		}
		return nil
	}
	if b, err := ioutil.ReadFile(file + ".toml"); err == nil {
		if err := toml.Unmarshal(b, &e.taskvars); err != nil {
			return err
		}
		return nil
	}
	return nil
}
