package config

type Conf struct {
	Port   string    `json:"port" yaml:"port"`
	IgnoredTopics	string `json:"ignoredTopics" yaml:"ignoredTopics"`
}
