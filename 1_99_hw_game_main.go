package main

import (
	"fmt"
	"strings"
)

type Player struct {
	position  BaseRoom
	backpack  bool
	inventory []Object
}

func NewPlayer() *Player {
	return &Player{}
}
func (p *Player) GetPosition() BaseRoom {
	return p.position
}
func (p *Player) SetPosition(value BaseRoom) {
	p.position = value
}
func (p *Player) HasBackpack() bool {
	return p.backpack
}
func (p *Player) SetBackpack(value bool) {
	p.backpack = value
}
func (p *Player) GetInventory() []Object {
	return p.inventory
}
func (p *Player) SetInventory(value []Object) {
	p.inventory = value
}
func (p *Player) AddInventory(value Object) {
	p.inventory = append(p.inventory, value)
}
func (p *Player) GoRoom(placeName string) string {
	temp := p.GetPosition().GetName() == placeName
	if temp {
		return "находишься в том же месте"
	}
	roomCheck := p.GetPosition().GetConnectedRoom(placeName)
	if roomCheck == nil {
		return "нет пути в " + placeName
	}
	room := *p.GetPosition().GetConnectedRoom(placeName)
	if room.GetAccessStatus() {
		player.SetPosition(room)
		response := player.GetPosition().GetDescription()
		response += ". можно пройти - "
		place := ""
		for _, s := range player.GetPosition().GetAllConnectedRooms() {
			place += s.GetName() + ", "
		}
		place = strings.TrimSuffix(strings.TrimSpace(place), ",")
		response += place
		return response
	}
	return room.GetAccessCondition()
}
func (p *Player) LookAround() string {
	flag := false
	objects := ""
	for k, v := range p.GetPosition().GetStock() {
		object := ""
		if len(p.GetPosition().GetStock()[k]) > 0 {
			flag = true
			for _, s := range v {
				object += s.name + ", "
			}
			objects += k + ": " + object
		}
	}
	if !flag {
		objects = "пустая комната"
	} else {
		objects = strings.TrimSuffix(strings.TrimSpace(objects), ",")
	}

	extra := ""
	cond1, _ := ContainsObject(player.GetInventory(), "конспекты")
	cond2, _ := ContainsObject(player.GetInventory(), "ключи")
	if !cond1 && !cond2 && player.GetPosition().GetName() == "кухня" {
		extra += "собрать рюкзак и "
	}

	places := ""
	for _, s := range p.GetPosition().GetAllConnectedRooms() {
		places += s.GetName() + ", "
	}
	places = strings.TrimSuffix(strings.TrimSpace(places), ",")

	return fmt.Sprintf(p.GetPosition().GetLookAroundFormat(), objects, extra, places)
}
func (p *Player) Take(objectName string) string {
	if p.HasBackpack() {
		flag, subStock, object := ContainsInRoom(p.GetPosition().GetStock(), objectName)
		if flag {
			p.AddInventory(*object)
			p.GetPosition().GetStock()[subStock] = RemoveFromSlice(p.GetPosition().GetStock()[subStock], objectName)
			return "предмет добавлен в инвентарь: " + objectName
		}
		return "нет такого"
	}
	return "некуда класть"
}
func (p *Player) PutOn(objectName string) string {
	flag, subStock, _ := ContainsInRoom(p.GetPosition().GetStock(), objectName)
	if flag {
		p.SetBackpack(true)
		p.GetPosition().GetStock()[subStock] = RemoveFromSlice(p.GetPosition().GetStock()[objectName], objectName)
		return "вы надели: " + objectName
	}
	return "нет такого"
}
func (p *Player) UseOn(object, activeObject string) string {
	flag, _ := ContainsObject(player.GetInventory(), object)
	if flag && p.GetPosition().GetActiveObjects()[object] != nil && p.GetPosition().GetActiveObjects()[object].GetName() == activeObject {
		p.GetPosition().GetActiveObjects()[object].Action()
		return p.GetPosition().GetActiveObjects()[object].GetResponse()
	} else if !flag {
		return "нет предмета в инвентаре - " + object
	} else {
		return "не к чему применить"
	}
}

type Object struct {
	name        string
	description string
}

func (o *Object) GetName() string {
	return o.name
}
func (o *Object) GetDescription() string {
	return o.description
}

type ActiveObject interface {
	GetName() string
	GetResponse() string
	Action()
}
type Door struct {
	name     string
	response string
	accessTo BaseRoom
}

func (d *Door) GetName() string {
	return d.name
}
func (d *Door) GetResponse() string {
	return d.response
}
func (d *Door) Action() {
	d.accessTo.SetAccessStatus(true)
}

type BaseRoom interface {
	GetName() string
	GetDescription() string
	GetLookAroundFormat() string
	GetActiveObjects() map[string]ActiveObject
	GetStock() map[string][]Object
	GetAllConnectedRooms() []BaseRoom
	GetConnectedRoom(string) *BaseRoom
	GetAccessStatus() bool
	SetAccessStatus(bool)
	GetAccessCondition() string
}
type Room struct {
	name              string
	description       string
	lookAroundFormat  string
	stock             map[string][]Object
	activeObjects     map[string]ActiveObject
	allConnectedRooms []BaseRoom
	accessStatus      bool
	accessCondition   string
}

func (r *Room) GetName() string {
	return r.name
}
func (r *Room) GetDescription() string {
	return r.description
}
func (r *Room) GetLookAroundFormat() string {
	return r.lookAroundFormat
}
func (r *Room) GetStock() map[string][]Object {
	return r.stock
}
func (r *Room) GetAllConnectedRooms() []BaseRoom {
	return r.allConnectedRooms
}
func (r *Room) GetConnectedRoom(roomName string) *BaseRoom {
	for room := range r.GetAllConnectedRooms() {
		if r.GetAllConnectedRooms()[room].GetName() == roomName {
			return &r.GetAllConnectedRooms()[room]
		}
	}
	return nil
}
func (r *Room) GetActiveObjects() map[string]ActiveObject {
	return r.activeObjects
}
func (r *Room) GetAccessStatus() bool {
	return r.accessStatus
}
func (r *Room) SetAccessStatus(value bool) {
	r.accessStatus = value
}
func (r *Room) GetAccessCondition() string {
	return r.accessCondition
}

func ContainsObject(slice []Object, name string) (bool, *Object) {
	for _, a := range slice {
		if a.GetName() == name {
			return true, &a
		}
	}
	return false, nil
}
func ContainsInRoom(map1 map[string][]Object, objectName string) (bool, string, *Object) {
	for subStock, slice := range map1 {
		flag, object := ContainsObject(slice, objectName)
		if flag {
			return true, subStock, object
		}
	}
	return false, "", nil
}
func RemoveFromSlice(list []Object, item string) []Object {
	length := len(list)
	for i, o := range list {
		if item == o.name {
			list[length-1], list[i] = list[i], list[length-1]
			return list[:length-1]
		}
	}
	return list
}

var player = NewPlayer()

func main() {
	initGame()
	/*
		fmt.Println(handleCommand("осмотреться"))
		fmt.Println(handleCommand("идти коридор"))
		fmt.Println(handleCommand("идти комната"))
		fmt.Println(handleCommand("осмотреться"))
		fmt.Println(handleCommand("надеть рюкзак"))
		fmt.Println(handleCommand("взять ключи"))
		fmt.Println(handleCommand("взять конспекты"))
		fmt.Println(handleCommand("идти коридор"))
		fmt.Println(handleCommand("применить ключи дверь"))
		fmt.Println(handleCommand("идти улица"))
	*/
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("завтракать"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("взять телефон"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("взять конспекты"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти кухня"))
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти улица"))
	fmt.Println(handleCommand("применить ключи дверь"))
	fmt.Println(handleCommand("применить телефон шкаф"))
	fmt.Println(handleCommand("применить ключи шкаф"))
	fmt.Println(handleCommand("идти улица"))
}

func initGame() {
	backpack := &Object{"рюкзак", "позволяет брать другие предметы"}
	keys := &Object{"ключи", "открывает дверь на улицу"}
	notes := &Object{"конспекты", "нужно взять в универ"}
	tea := &Object{"чай", "чашка чая"}

	hall := &Room{
		name:              "коридор",
		description:       "ничего интересного",
		lookAroundFormat:  "%s%s. можно пройти - %s",
		activeObjects:     nil,
		stock:             nil,
		allConnectedRooms: nil,
		accessStatus:      true,
		accessCondition:   "no",
	}
	home := &Room{
		name:              "домой",
		description:       "ты дома",
		lookAroundFormat:  "%s%s. можно пройти - %s",
		activeObjects:     nil,
		stock:             nil,
		allConnectedRooms: nil,
		accessStatus:      true,
		accessCondition:   "no",
	}
	kitchen := &Room{
		name:              "кухня",
		description:       "кухня, ничего интересного",
		lookAroundFormat:  "ты находишься на кухне, %s, надо %sидти в универ. можно пройти - %s",
		activeObjects:     nil,
		stock:             map[string][]Object{"на столе": {*tea}},
		allConnectedRooms: nil,
		accessStatus:      true,
		accessCondition:   "no",
	}
	privateRoom := &Room{
		name:             "комната",
		description:      "ты в своей комнате",
		lookAroundFormat: "%s%s. можно пройти - %s",
		activeObjects:    nil,
		stock: map[string][]Object{
			"на столе": {*keys, *notes},
			"на стуле": {*backpack},
		},
		allConnectedRooms: nil,
		accessStatus:      true,
		accessCondition:   "no",
	}
	street := &Room{
		name:              "улица",
		description:       "на улице весна",
		lookAroundFormat:  "%s%s. можно пройти - %s",
		activeObjects:     nil,
		stock:             nil,
		allConnectedRooms: nil,
		accessStatus:      false,
		accessCondition:   "дверь закрыта",
	}

	hall.allConnectedRooms = []BaseRoom{kitchen, privateRoom, street}
	home.allConnectedRooms = hall.allConnectedRooms
	kitchen.allConnectedRooms = []BaseRoom{hall}
	privateRoom.allConnectedRooms = []BaseRoom{hall}
	street.allConnectedRooms = []BaseRoom{home}

	hall.activeObjects = map[string]ActiveObject{"ключи": &Door{"дверь", "дверь открыта", street}}

	player.SetPosition(kitchen)
	player.SetBackpack(false)
	player.SetInventory(nil)
}

func handleCommand(command string) string {
	s := strings.Split(command, " ")

	switch s[0] {
	case "осмотреться":
		return player.LookAround()
	case "идти":
		return player.GoRoom(s[1])
	case "надеть":
		return player.PutOn(s[1])
	case "взять":
		return player.Take(s[1])
	case "применить":
		return player.UseOn(s[1], s[2])
	}
	return "неизвестная команда"
}
