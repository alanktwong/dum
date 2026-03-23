package playbook

import (
	e "awong/dotfiles/pkg/enums"
	u "awong/dotfiles/pkg/external"
	l "awong/dotfiles/pkg/logging"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Factory constructs a Input and its PlayBook from YAML.
type Factory struct {
	Log   l.Logger
	Utils u.Ext
}

// NewFactory constructs a Factory.
func NewFactory() *Factory {
	return &Factory{
		Log:   l.Log(),
		Utils: u.NewExt(),
	}
}

// InputOptions provides the parameterized struct for constructing a Input.
type InputOptions struct {
	File   string
	Group  string
	DryRun bool
}

// Provide constructs a Input given Input Options.
func (f *Factory) Provide(options InputOptions) (*Input, error) {
	yml, err := f.getYaml(options.File)
	if err != nil {
		return nil, err
	}
	playbook, err := f.ProvidePlayBook(yml)
	if err != nil {
		return nil, fmt.Errorf("failed to provide playbook: %w", err)
	}
	return NewInput(options.DryRun, options.Group, playbook), nil
}

func (f *Factory) getYaml(file string) (map[string]interface{}, error) {
	f.Log.Debug("Loading playbook from file", "file", file)
	absPath, err := f.Utils.ToAbsolutePath(file)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for file %s: %w", file, err)
	}
	byteArray, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get read file %s: %w", absPath, err)
	}
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(byteArray, &yamlMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal playbook from file %s: %w", file, err)
	}
	return yamlMap, nil
}

// ProvidePlayBook constructs a PlayBook given a map that expresses the YML.
func (f *Factory) ProvidePlayBook(yml map[string]interface{}) (*PlayBook, error) {
	pbData, ok := yml["playbook"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("playbook key not found or invalid format in YAML")
	}

	plays, err := f.providePlays(pbData)
	if err != nil {
		return nil, fmt.Errorf("failed to provide plays: %w", err)
	}
	id := f.getString(pbData, "id", "")
	description := f.getString(pbData, "description", id)
	apps := f.provideJetBrainsApps(pbData)
	attributes, err := NewAttributes(
		id,
		description,
		f.getBool(pbData, "enabled", true),
		f.getBool(pbData, "sudo", false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}
	return NewPlayBook(attributes, plays, apps)
}

func (f *Factory) providePlays(yml map[string]interface{}) ([]*Play, error) {
	var plays []*Play
	if arr, ok := yml["plays"].([]interface{}); ok {
		for _, t := range arr {
			if m, ok := t.(map[string]interface{}); ok {
				play, err := f.providePlay(m)
				if err != nil {
					return nil, fmt.Errorf("failed to provide play: %w", err)
				}
				plays = append(plays, play)
			} else {
				return nil, fmt.Errorf("play is not a map[string]interface{}")
			}
		}
	}
	return plays, nil
}

func (f *Factory) providePlay(yml map[string]interface{}) (*Play, error) {
	tasks, err := f.provideTasks(yml)
	if err != nil {
		return nil, fmt.Errorf("failed to provide tasks: %w", err)
	}
	id := f.getString(yml, "id", "")
	description := f.getString(yml, "description", id)
	if len(tasks) == 0 {
		return nil, fmt.Errorf("play %v has no tasks", id)
	}
	attributes, err := NewAttributes(
		id,
		description,
		f.getBool(yml, "enabled", true),
		f.getBool(yml, "sudo", false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}
	return NewPlay(attributes, tasks)
}

func (f *Factory) provideTasks(yml map[string]interface{}) ([]Task, error) {
	var tasks []Task
	if arr, ok := yml["tasks"].([]interface{}); ok {
		for _, t := range arr {
			if m, ok := t.(map[string]interface{}); ok {
				task, err := f.provideTask(m)
				if err != nil {
					return nil, fmt.Errorf("failed to provide task: %w", err)
				}
				tasks = append(tasks, task)
			} else {
				return nil, fmt.Errorf("task is not a map[string]interface{}")
			}
		}
	}
	return tasks, nil
}

func (f *Factory) provideTask(yml map[string]interface{}) (Task, error) {
	id := f.getString(yml, "id", "")
	description := f.getString(yml, "description", id)
	attributes, err := NewAttributes(
		id,
		description,
		f.getBool(yml, "enabled", true),
		f.getBool(yml, "sudo", false),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attributes: %w", err)
	}

	typeStr := f.getString(yml, "type", "")
	taskType, err := e.ParseTaskType(typeStr)
	if err != nil {
		return nil, fmt.Errorf("unknown task type %s for task %s", typeStr, id)
	}
	switch taskType {
	case e.TaskTypeDir:
		return NewDirTask(attributes)
	case e.TaskTypeMas:
		return NewMasTask(attributes)
	case e.TaskTypeVscode:
		return NewVsCodePluginTask(attributes)
	case e.TaskTypeFunction:
		return NewFunctionTask(attributes)
	case e.TaskTypeLink:
		return NewLinkTask(attributes,
			f.getString(yml, "root", ""),
			f.getString(yml, "target", ""))
	case e.TaskTypeGit:
		return NewGitTask(attributes,
			f.getString(yml, "root", ""),
			f.getString(yml, "name", ""))
	case e.TaskTypeJetbrains:
		apps := f.getStrings(yml, "apps", make([]string, 0))
		return NewJetBrainsPluginTask(attributes, apps)
	case e.TaskTypeBrew:
		return NewBrewTask(attributes,
			f.getString(yml, "tap", ""))
	case e.TaskTypeCask:
		return NewBrewCaskTask(attributes,
			f.getString(yml, "tap", ""))
	case e.TaskTypeCellar:
		return NewBrewCellarTask(attributes,
			f.getString(yml, "tap", ""))
	case e.TaskTypeBash:
		attributes.Command = f.getString(yml, "command", "")
		attributes.Script = f.getString(yml, "script", "")
		return NewBashTask(attributes)
	default:
		return nil, fmt.Errorf("unknown task type %s for task %s", taskType, id)
	}
}

func (f *Factory) provideJetBrainsApps(yml map[string]interface{}) map[string]string {
	apps := make(map[string]string)
	appkeys := e.JetBrainsTypeNames()

	if jetbrainsArray, ok := yml["jetbrains"].([]interface{}); ok {
		for _, each := range appkeys {
			for _, it := range jetbrainsArray {
				if m, ok2 := it.(map[string]interface{}); ok2 {
					version := f.getString(m, each, "")
					if version != "" {
						apps[each] = version
					}
				}
			}
		}
	}
	return apps
}

func (f *Factory) getString(data map[string]interface{}, key string, def string) string {
	if value, ok := data[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return def
}

func (f *Factory) getStrings(data map[string]interface{}, key string, def []string) []string {
	if value, ok := data[key]; ok {
		if generic, ok := value.([]interface{}); ok {
			stringValues := make([]string, 0)
			for _, v := range generic {
				if strValue, ok := v.(string); ok {
					stringValues = append(stringValues, strValue)
				}
			}
			return stringValues
		}
	}
	return def
}

func (f *Factory) getBool(data map[string]interface{}, key string, def bool) bool {
	if value, ok := data[key]; ok {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return def
}
