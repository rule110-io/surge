package models

type Topic struct {
	Name        string
	NameEncoded string
}

type TopicInfo struct {
	Name              string
	Subscribers       int
	FileCount         int
	Permissions       TopicPermissions
	SubscriptionState int
}
