package hubitat

import (
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/recws-org/recws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"fmt"
	"log"
	//"reflect"
	"github.com/asdine/storm/v3"
	"context"
)

var db *storm.DB

func init(){


	var err error
	db, err = storm.Open("cacheddevices.db")
    if(err!=nil){
    	panic(err)
	}
}

type Api struct{
	access_token string
	base_url string
	Devices []*Device
	usecache bool
}

func NewApi(token string,url string,usecache bool) (*Api){


	api :=&Api{
		access_token: token,
		base_url: url,
		usecache: usecache,
	}

   return api
}

func (a *Api) url(req string) (string){

	out :=  fmt.Sprintf("%s/%s?access_token=%s",a.base_url,req,a.access_token)

    log.Println(out)

	return out

}

func (a *Api) Request(url string) ([]byte,error){

	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	resp,err := client.Get(url)

    if(err!=nil){

    	return []byte(""),err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if(err!=nil){
		return body,err
	}

	return body,nil

}


func (a *Api) Device(id string) (*Device){

	  var em *Device

     for _,d := range a.Devices {

     	if(id==d.Id){

     		return d
		}

	 }

	 return em

}


func (a *Api) DeviceCmd(id string,cmd string) (Device,error){
	var dev Device

	url := a.url(fmt.Sprintf("devices/%s/%s",id,cmd))

	body,err := a.Request(url)

	if(err!=nil){

		return dev,err

	}


	err = json.Unmarshal(body,&dev)

	if(err!=nil){

		return dev,err
	}
	return dev,nil

}

func (a *Api) CachedDevices() ([]*Device,error){

	var all []*Device
	if(a.usecache==true) {

		err := db.All(&all)

		if (err != nil) {

			return all, err
		}
	}

	a.Devices = all

	for _,d := range a.Devices {

		d.api =  a

	}

	return a.Devices,nil

}

func (a *Api) saveToCache(){

	 if(!a.usecache){
	 	return
	 }

	 for _,i := range a.Devices {
	 	db.Save(&i)
	 }

}

func (a *Api) AllDevices() (error){

	url := a.url("devices/all")

	var all []*Device

    body,err := a.Request(url)

    if(err!=nil){

    	return err
	}

	err = json.Unmarshal(body,&all)

	if(err!=nil){

		return err
	}

	a.Devices = all


	for _,d := range a.Devices {

		d.api =  a

	}


	a.saveToCache()

	//save to cache

    return nil

}

func (a *Api) Attributes(name string) ([]*Device,error){

	url := a.url(fmt.Sprintf("attribute/%s",name))

	var all []*Device

	body,err := a.Request(url)

	if(err!=nil){

		return all,err
	}

	err = json.Unmarshal(body,&all)

	if(err!=nil){

		return all,err
	}


	return all,nil
}

func (a *Api) EstablishWatch() {

    up,perr := url.Parse(a.base_url)

	if perr != nil {
		panic(perr)
	}


	u := url.URL{Scheme: "ws", Host: up.Host, Path: "/eventsocket"}

	log.Printf("connecting to %s", u.String())


	ctx,_ := context.WithCancel(context.Background())

	ws := recws.RecConn{
		KeepAliveTimeout: 30 * time.Second,
	}

	ws.Dial(u.String(), nil)

	eventChan := make(chan string)

	go debounceEvent(300*time.Millisecond, eventChan, func(id string) {

		log.Println("Got Debounce Event for device id",id)

		//at this point the state is up to date, we could also emit any event flows we want here

		d:=a.Device(id)
		log.Printf("%#v",d)

	})


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

			id := a.parseMessage(message)
            eventChan <- id

		}
	}


}

func (a *Api) parseMessage(message []byte) string {

	id,_,_,_ := jsonparser.Get(message, "deviceId")

	for _,d := range a.Devices {

		if( string(id) != d.Id ){
			continue
		}

		source,_,_,_ := jsonparser.Get(message, "source")
		name,_,_,_ := jsonparser.Get(message, "name")
		display,_,_,_ := jsonparser.Get(message, "displayName")

		if(string(source)=="DEVICE"){

			value,_,_,_ := jsonparser.Get(message, "value")

			d:=a.Device(string(id))


			if(d.State==nil){

				d.State = make(map[string]string)
			}

			d.State[strings.Title(string(name))] = string(value)

		    log.Printf("Hubitat: %s->%s %s:%s (%s)\n",string(display),string(value),"device",string(id),d.ToHandle())
			//log.Printf("%#v",d)
			//Events <- Device{Name: string(display),Id: string(id),Type: key,Action: string(value)}

			return string(id)
		}

	}

	return string(id)

	//is the deviceId in our list?


}


func debounceEvent(interval time.Duration, input chan string, cb func(arg string)) {

	var itemhold string
	var lister = make(map[string]string)
	timer := time.NewTimer(interval)

	for {


		select {
		case itemhold = <-input:

			lister[itemhold]=itemhold

		case <-timer.C:

			for key,item := range lister{

				cb(item)

				delete(lister,key)
			}
			timer.Reset(interval)
		}
	}
}
