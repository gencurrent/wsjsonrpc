package wsjsonrpc

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func Run(url url.URL) error {
	fmt.Println("The client is starting")

	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not connect to the URL %v", url.String()))
		return err
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- done:
			return
		case t := <- ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil {
					log.Println("write: ", err)
					return err
				}
			case <-interrupt:
				log.Println("interrupt")
			
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close: ", err)
					return err 
				}
				select {
				case <- done: 
				case <-time.After(time.Second)
				}
				return 

		}
	}

}
