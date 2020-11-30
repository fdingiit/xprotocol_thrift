package sofaregistry

type Handler interface {
	// Received data will be in the form of a map (zoneName->dataList).
	// LocalZone is pushed by server for priority routing.
	OnRegistryPush(dataID string, data map[string][]string, localZone string)
}

type HandlerFunc func(dataID string, data map[string][]string, localZone string)

func (f HandlerFunc) OnRegistryPush(dataID string, data map[string][]string, localZone string) {
	f(dataID, data, localZone)
}
