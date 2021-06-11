package hubitat

import (
	"regexp"
	"strings"
)

type Device struct {
	api *Api
	Name string
	Label string
	Type string
	Id string `storm:"id"`
	Date string
	Model string
	Manufacturer string
	Capabilities []string
	Attributes DeviceAttribute
	Commands []DeviceCommand
	State map[string]string
}

type DeviceAttribute struct{
	Level string
	Name string
	CurrentValue string
	Battery string
	DataType string
	Values []string
	Switch string
	Motion string
}

type DeviceCommand struct{
	Command string
}

func (d Device) GetState(key string) string{


	if(d.State==nil){

		d.State = make(map[string]string)
	}


	if  val,ok := d.State[key]; ok {

		return val

	}

   return ""

}

func (d *Device) Cmd(cmd string){

	go d.api.DeviceCmd(d.Id,cmd)

}

func (d Device) ToHandle() string{

	re := regexp.MustCompile(`_+`)


	label := strings.Replace(d.Label,"-","_",-1)
	label = strings.Replace(label," ","_",-1)
	label = strings.ToLower(label)
	label = re.ReplaceAllString(label,"_")
	return label


}