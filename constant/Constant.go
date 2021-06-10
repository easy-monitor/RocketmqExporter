package constant

import (
	"strings"
)
var IgnoredTopicArray []string

//rockemq有一些默认topic，这些不需要纳入监控，需要支持可配置
//var ignoredTopicList = []string{"RMQ_SYS_TRANS_HALF_TOPIC", "BenchmarkTest", "OFFSET_MOVED_EVENT", "TBW102", "SELF_TEST_TOPIC", "DefaultCluster", "broker-b", "broker-a"}
func SetIgnoredTopicArray(ignoredTopicsInEnv string) {
	IgnoredTopicArray = strings.Split(ignoredTopicsInEnv, ",")
}

//定义metrics接口path
func GetMetricsPath() string {
	return "/metrics"
}



//定义metrics数据名称的前缀：比如定义前缀是rocketmq，则metrics为:rocketmq_msg_diff_detail等
func GetMetricsPrefix() string {
	return "rocketmq"
}
