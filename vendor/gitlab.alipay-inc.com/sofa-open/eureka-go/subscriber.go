package eureka

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"
	"time"
)

type Subscriber struct {
	dataId      string
	checkSum    string
	application *Application
}

var (
	subscribers = NewConcurrentMap()
	listeners   []*func(dataId string, app *Application, err error)
	loopStarted = false
)

func (c *Client) Subscribe(dataId string) error {
	subscriber := &Subscriber{dataId, "", nil}
	subscribers.Set(dataId, subscriber)

	//同步拉取一次
	getSuccess := false
	for i := 0; i < 3; i++ {
		app, err := c.GetApplication(dataId)
		if err == nil {
			getSuccess = true
			subscriber.application = app
			subscriber.checkSum = checksum(app)
			break
		}
	}

	if getSuccess {
		//同步触发一次回调事件
		c.triggerSubscriberChangedEvent(subscriber)
	}

	return nil
}

func (c *Client) UnSubscribe(dataId string) {
	subscribers.Remove(dataId)
}

//定时从 eureka 获取最新数据，如果数据md5有变化就触发回调事件
func (c *Client) StartLoopSubscriber(d time.Duration) {
	if loopStarted {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("[error] refresh subscriber panic.")
			}
		}()
		for true {
			time.Sleep(d)

			for dataId, item := range subscribers.Items() {
				if item != nil {
					Subscriber := item.(*Subscriber)
					newApp, err := c.GetApplication(dataId)
					if err == nil {
						oldCheckSum := Subscriber.checkSum
						newCheckSum := checksum(newApp)

						//有数据变更
						if !strings.EqualFold(newCheckSum, oldCheckSum) {
							Subscriber.checkSum = newCheckSum
							Subscriber.application = newApp
							//触发回调事件
							c.triggerSubscriberChangedEvent(Subscriber)
						}

					}
				}
			}
		}

	}()
	loopStarted = true

	log.Printf("[INFO] eureka loop subscriber job started.")
}

func checksum(app *Application) (result string) {
	//蚂蚁内部所有的数据都放在 metadata 里，所以通过 metadata 来判断 instance 唯一性
	if app == nil {
		return ""
	}
	pubDatas := make([]*MetaData, 0)
	instances := app.Instances
	if len(instances) > 0 {
		for _, instance := range instances {
			pubDatas = append(pubDatas, instance.Metadata)
		}
	}
	byte, _ := json.Marshal(pubDatas)
	res := md5.Sum(byte)
	result = hex.EncodeToString(res[:])
	return
}

func (c *Client) triggerSubscriberChangedEvent(subscriber *Subscriber) {
	if subscriber == nil || subscriber.application == nil || len(subscriber.application.Instances) == 0 {
		return
	}
	
	for _, funcItem := range listeners {
		(*funcItem)(subscriber.dataId, subscriber.application, nil)
	}
}

func (c *Client) AddChangedListener(listener func(dataId string, application *Application, err error)) {
	log.Printf("[INFO] add eureka listener")
	listeners = append(listeners, &listener)
}
