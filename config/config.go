package config

type Conf struct {
	Port          string          `json:"port" yaml:"port"`
	IgnoredTopics string          `json:"ignoredTopics" yaml:"ignoredTopics"`
	Module        []*ModuleStruct `json:"module" yaml:"module"`
}
type ModuleStruct struct {
	Name     string `json:"name" yaml:"name"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}
type LoginResponse struct {
	Data DataStruct `json:"data" yaml:"data"`
}
type DataStruct struct {
	SessionId string `json:"sessionId" yaml:"sessionId"`
}
