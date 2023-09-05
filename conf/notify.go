package conf

func GetNotifyEndPoint() string {
	return configuration.Notify.EndPoint
}

func GetAccessToken() string {
	return configuration.Notify.AccessToken
}
