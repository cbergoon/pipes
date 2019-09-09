package pipeline

import (
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

func findPluginPathsInDir(pluginDir string) ([]string, error) {
	pluginPaths := []string{}
	err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".so") {
			pluginPaths = append(pluginPaths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return pluginPaths, err
}

func openPlugins(pluginPaths []string) ([]*plugin.Plugin, error) {
	plugins := []*plugin.Plugin{}
	for _, pp := range pluginPaths {
		plug, err := plugin.Open(pp)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, plug)
	}
	return plugins, nil
}

// LoadProcessPluginNameMap reads a specified plugin directory and returns a map of the opened plugins
func LoadProcessPluginNameMap(pluginDir string) (map[string]*plugin.Plugin, error) {
	pluginPaths, err := findPluginPathsInDir(pluginDir)
	if err != nil {
		return nil, err
	}

	pluginList, err := openPlugins(pluginPaths)
	if err != nil {
		return nil, err
	}

	pluginMap := make(map[string]*plugin.Plugin)
	for _, p := range pluginList {
		pluginTypeNameSym, err := p.Lookup("TypeName")
		if err != nil {
			return nil, err
		}
		pluginTypeName, ok := pluginTypeNameSym.(func() string)
		if ok {
			pluginMap[pluginTypeName()] = p
		}
	}
	return pluginMap, nil
}
