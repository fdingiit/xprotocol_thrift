package eureka

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

type Publisher struct {
	instanceInfo *InstanceInfo
	published    bool
}

var (
	publishers       = NewConcurrentMap()
	heartbeatStarted = false
)

func init() {

}

func (c *Client) RegisterInstance(dataId string, instanceInfo *InstanceInfo) error {
	publisher := &Publisher{
		instanceInfo: instanceInfo,
		published:    false,
	}

	publishers.Set(dataId, publisher)

	if len(instanceInfo.InstanceID) == 0 {
		instanceInfo.InstanceID = uuid.New().String()
	}

	instanceInfo.HostName = instanceInfo.App

	values := []string{"apps", dataId}
	path := strings.Join(values, "/")
	instance := &Instance{
		Instance: instanceInfo,
	}
	body, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	pubSuccess := false
	for i := 0; i < 3; i++ {
		_, err = c.Post(path, body)
		if err == nil {
			pubSuccess = true
			break
		}
	}

	publisher.published = pubSuccess

	return err
}

func (c *Client) UnregisterInstance(dataId string) error {
	pubItem, ok := publishers.Get(dataId)
	if !ok {
		return nil
	}

	publisher := pubItem.(*Publisher)

	instanceId := publisher.instanceInfo.InstanceID
	values := []string{"apps", dataId, instanceId}
	path := strings.Join(values, "/")
	_, err := c.Delete(path)

	if err == nil {
		publishers.Remove(dataId)
	}

	return err
}

func (c *Client) SendHeartbeat(dataId, instanceId string) error {
	values := []string{"apps", dataId, instanceId}
	path := strings.Join(values, "/")
	resp, err := c.Put(path, nil)
	if err != nil {
		return err
	}
	switch resp.StatusCode {
	case http.StatusNotFound:
		return newError(ErrCodeInstanceNotFound,
			"Instance resource not found when sending heartbeat", 0)
	}
	return nil
}

func (c *Client) StartHeartbeat(d time.Duration) {
	if heartbeatStarted {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("[error] refresh subscriber panic.")
			}
		}()
		for true {
			time.Sleep(d)

			for dataId, pub := range publishers.Items() {
				if pub != nil {
					publisher := pub.(*Publisher)
					if publisher.published {
						err := c.SendHeartbeat(dataId, publisher.instanceInfo.InstanceID)
						if err != nil {
							log.Printf("[WAIN] eureka heartbeat failed. dataId = %s", dataId)
						}
					}
				}
			}
		}
	}()

	heartbeatStarted = true

	log.Printf("[INFO] eureka heartbeat job started.")
}
