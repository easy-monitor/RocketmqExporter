package wrapper

import (
	"RocketmqExporter/model"
	"RocketmqExporter/utils"

	//"bytes"
	"encoding/json"
	"fmt"
)

func GetTopicNameList(rocketmqConsoleIPAndPort string, sessionId string) []string {
	var url = "http://" + rocketmqConsoleIPAndPort + "/topic/list.query"
	var content = utils.HttpUrl("GET", url, sessionId)

	var jsonData model.TopicList
	err := json.Unmarshal([]byte(content), &jsonData)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return jsonData.Data.TopicList
}

func GetConsumerListByTopic(rocketmqConsoleIPAndPort string, topicName string, sessionId string) *model.ConsumerList_By_Topic {

	var url = "http://" + rocketmqConsoleIPAndPort + "/topic/queryConsumerByTopic.query?topic=" + topicName
	var content = utils.HttpUrl("GET", url, sessionId)

	var jsonData *model.ConsumerList_By_Topic
	err := json.Unmarshal([]byte(content), &jsonData)

	if err != nil {
		return nil
	}

	return jsonData
}
