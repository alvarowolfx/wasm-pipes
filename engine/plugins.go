package engine

var plugins = make(map[string]Plugin)

func RegisterPlugin(name string, plugin Plugin) {
	plugins[name] = plugin
}
