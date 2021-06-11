package hubitat

import (
	"context"
	"encoding/json"
	"github.com/recws-org/recws"
	"log"
	"net/url"
	"time"
)


var cfg Cfg

var Events chan Event

type Cfg struct {
    Addr string
    Enabled bool
    Devices map[string][]string
}


type Event struct{
	Source string
    Name string
	DisplayName string
	Value string
	Unit string
	DeviceId int
	HubId string  `json:"-"`
	InstalledAppId string `json:"-"`
	DescriptionText string
}

func Config(addr string,enabled bool,dev map[string][]string){
    
     cfg = Cfg{Addr: addr,Enabled: enabled,Devices: dev}
    
    
    if(enabled==false){
        
        return
    }
    
    go setupListener()
    
    Events = make(chan Event)
    
}



func setupListener(){
    
	u := url.URL{Scheme: "ws", Host: cfg.Addr, Path: "/eventsocket"}
	
	log.Printf("connecting to %s", u.String())
    

    ctx,_ := context.WithCancel(context.Background())
    
	ws := recws.RecConn{
		KeepAliveTimeout: 30 * time.Second,
	}
	
	ws.Dial(u.String(), nil)

 
    for {
    		select {
    		case <-ctx.Done():
    			go ws.Close()
    			log.Printf("Websocket closed %s", ws.GetURL())
    			return
    		default:
    			if !ws.IsConnected() {
    				log.Printf("Websocket disconnected %s", ws.GetURL())
    				continue
    			}
    
    			_, message, err := ws.ReadMessage()
    			if err != nil {
    				log.Printf("Error: ReadMessage %s", ws.GetURL())
    				return
    			}
    
    			parseMessage(message)
    		}
    	}
    
}


func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}		


//{"source":"DEVICE","name":"zone_4","displayName":"Sprinkler - Zone 1","value":"off","unit":null,"deviceId":2061,"hubId":null,"locationId":null,"installedAppId":null,"descriptionText":null}
///we look for events from the system on the device
func parseMessage(message []byte){
    



          var dev Event

          if err := json.Unmarshal([]byte(message), &dev); err != nil {
			    log.Println(string(message))
				log.Printf("JSON Error: %s",err)
          }

          Events <- dev

}
