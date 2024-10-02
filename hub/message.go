package hub

type Message struct {
	topic   string
	content []byte
}

func ExtractTopic(message []byte) string {
	return string(message)
}
