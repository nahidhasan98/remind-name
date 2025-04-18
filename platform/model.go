package platform

type Platform struct {
	Name        string `bson:"name"`
	BotName     string `bson:"botName"`
	BotUsername string `bson:"botUsername"`
	BotLink     string `bson:"botLink"`
}
