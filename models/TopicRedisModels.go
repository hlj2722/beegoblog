package models

import ()

///region TopicRedis
func AddTopicRedis(title, category, lable, content, attachment, author string) error {

	return nil
}

func GetTopicRedis(tid string) (*Topic, error) {

	return nil, nil
}

func ModifyTopicRedis(tid, title, category, lable, content, attachment string) error {

	return nil
}

func DeleteTopicRedis(tid string) error {

	return nil
}

func GetAllTopicsRedis(category, lable string, isDesc bool) (topics []*Topic, err error) {

	return nil, nil
}

///endRegion
