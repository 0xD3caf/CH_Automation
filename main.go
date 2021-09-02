package main

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/thedevsaddam/gojsonq/v2"
)

/* TODO
-Add reset for progression mode
	//Check if the current zone is the same as HZE - 1, yes, press button to turn on
	//or check if progress mode != on, true, turn on
-Add logic for combining skills, try to maximize times they all sync

-Add DLL code injection to allow better control
//collect handle to game, inject dll

*/
type Skill struct {
	name     string
	cd       int //as seconds
	duration int //as seconds
}

type Ancient struct {
	name  string
	level int64
	power float32
	uid   int
}

var hero_ready robotgo.CHex = 0xC0D431
var hero_unavail robotgo.CHex = 0x46687F

//each string is a set of info uid:levelbonus, split then use

//gamesave location should be set here
var gamesave = "clickerHeroSave.txt"
var hash = "7a990d4252c6fb53aacfbb0ec1a3b23"

//ORDER MATTERS
var Cooldowns = [9]Skill{
	{"Clickstorm", 600, 10},
	{"Powersurge", 600, 10},
	{"Lucky Strikes", 1800, 30},
	{"Metal Detector", 1800, 30},
	{"Golden Clicks", 3200, 30},
	{"The Dark Ritual", 28800, -1},
	{"Super Clicks", 1800, 60},
	{"Energize", 3200, -1},
	{"Reload", 3200, -1},
}

var my_ancients = [26]Ancient{
	{"Libertas", 0, 0.25, 4},
	{"Siyalatas", 0, 0.25, 5},
	{"Mammon", 0, 0.05, 8},
	{"Mimzee", 0, 0.5, 9},
	{"Pluto", 0, .30, 10},
	{"Dogcog", 0, 0, 11},
	{"Fortuna", 0, 0, 12},
	{"Atman", 0, 0, 13},
	{"Dora", 0, 0, 14},
	{"Bhaal", 0, 0.15, 15},
	{"Morgulis", 0, 0.11, 16},
	{"Chronos", 0, 0, 17},
	{"Bubos", 0, 0, 18},
	{"Fragsworth", 0, 0.20, 19},
	{"Vaagur", 0, 0, 20},
	{"Kumawakamaru", 0, 0, 21},
	{"Chawedo", 0, 2.0, 22},
	{"Hetoncheir", 0, 2.0, 23},
	{"Beserker", 0, 2.0, 24},
	{"Sniperino", 0, 2.0, 25},
	{"Kleptos", 0, 2.0, 26},
	{"Energon", 0, 2.0, 27},
	{"Argaiv", 0, 2.0, 28},
	{"Juggernaut", 0, 0.01, 29},
	{"Revloc", 0, 0, 31},
	{"Nogardnit", 0, 0.10, 32}}

func init() {
	//starts by adding extra time to length and computing reduced cooldowns\
	//sets the window to be Clicker heroes
	//robotgo.SetActive()
	parse_save(gamesave)
	update_data("save_JSON.txt")
	for i := 0; i < len(Cooldowns); i++ {
		//Cooldowns[i].cd := Cooldowns[i].cd / modval
		//Cooldowns[i].duration := Cooldowns[i].duration / modval
	}
}
func main() {

	fmt.Println("Clicker Heroes Automation")
	main_hero_upgrade()
	auto_abilities()
	/*
		var wait_group sync.WaitGroup
		wait_group.Add(1)
		//go autoClick_Polling(&wait_group)
		go bird_collector()
		wait_group.Wait()
	*/
	fmt.Println("Done")

	/*count := 0
	for {
		bird_collector()
		if count > 5 {
			break
		} else {
			time.Sleep(10 * time.Second)
			count += 1
		}
	}*/
}

func bird_collector() {
	//!TODO
	// change format to use robotgo.AddMouse() and not use the physical mouse
	tolerance := 0.02
	var color robotgo.CHex = 0xC0D431
	var color2 robotgo.CHex = 0xA6C22F

	for {
		bitmap := robotgo.CaptureScreen()
		// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
		a, b := robotgo.FindColor(color, bitmap, tolerance)
		c, d := robotgo.FindColor(color2, bitmap, tolerance)
		if a != -1 || b != -1 {
			robotgo.MoveClick(a, b, "LEFT_BUTTON", false)
			fmt.Println("Got One")
		} else if c != -1 || d != -1 {
			robotgo.MoveClick(c, d, "LEFT_BUTTON", false)
		}
		robotgo.FreeBitmap(bitmap)
		time.Sleep(5 * time.Second)
	}
}

func main_hero_upgrade() {
	//keeps the current max hero at max levels
	robotgo.MoveClick(0, 0, "LEFT_BUTTON", false)

}

func prev_hero_upgrade() {
	//scroll to bottom
	//check if bottom is purchaseable and has 0 levels
	//or has non-zero levels
	//if no levels and no purchase, use prev hero
}

func auto_abilities() {
	fmt.Println(Cooldowns)
	/*for i := 0; i < len(Cooldowns); i++ {
		fmt.Println("Skill : ", Cooldowns[i].name, "has cooldown: ", Cooldowns[i].cd, " and duration: ", Cooldowns[i].duration)
	}*/
	//auto clicks on any upgrade where the duration is longer than cooldown, ie it can always be active
}

func update_save() {
	//saves the game and then recopies and parses save
}
func parse_save(input_file string) {
	//var final_str string
	file, err := os.Open(input_file)
	if err != nil {
		panic(err)
	}
	buff, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()
	save := buff[32:]
	save_string := string(save)
	//base64 decode
	decoded_save_string, err := b64.StdEncoding.DecodeString(save_string)
	if err != nil {
		panic(err)
	}
	//zlib decompress
	save_bytes := bytes.NewReader(decoded_save_string)
	output, err := zlib.NewReader(save_bytes)
	if err != nil {
		panic(err)
	}
	defer output.Close()
	final_string, err := io.ReadAll(output)
	if err != nil {
		panic(err)
	}
	file, err = os.OpenFile("save_JSON.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(final_string)
	//first loop goes through the set of chars pulled from file
}

func update_data(save_JSON string) {
	//takes cd reducation and duration increase from file
	//updates ability table

	//opens file
	file, err := os.Open(save_JSON)
	if err != nil {
		panic(err)
	}
	//collect data and close file
	save_data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	file.Close()

	save_string := string(save_data) //convert bytes to string
	var ancient_search_string string
	var ancient_output_int int64
	for i := 0; i < len(my_ancients); i++ {
		ancient_id := my_ancients[i].uid
		ancient_search_string = "ancients.ancients." + strconv.Itoa(ancient_id) + ".level"
		ancient_JSON_resp := gojsonq.New().FromString(save_string).Find(ancient_search_string)
		ancient_output := ancient_JSON_resp.(string)
		big_val := false
		for i := 0; i < len(ancient_output); i++ {
			if string(ancient_output[i]) == "e" {
				big_val = true
			}
		}
		if !big_val {
			ancient_output_float, _ := strconv.ParseFloat(ancient_output, 64)
			ancient_output_int = int64(math.Round(ancient_output_float))
		} else {
			ancient_output_int = int64(-1)
		}
		my_ancients[i].level = int64(ancient_output_int)
	}
	count := 1
	//updates any ancients that have item bonus stats

	//collect list of item slots
	//items.items.uid
	list_string := "items.slots"
	item_list_JSON_resp := gojsonq.New().FromString(save_string).From(list_string).Get()
	fmt.Println(item_list_JSON_resp)
	for {
		//iterate over items in slots and get info for item
		item_string := "items"
		item_JSON_resp := gojsonq.New().FromString(save_string).Find(item_string)
		//fmt.Println(item_JSON_resp)
		if item_JSON_resp != nil {
			//needs to be singe item to use string method
			//item_output_string := item_JSON_resp.(string)
			//fmt.Println(item_output_string)
			count = count + 1
			fmt.Println("Looped")
			if count > 25 {
				fmt.Println("break")
				break
			}
		}
	}
	//fmt.Println(item_output)
	var cooldown_level int64
	for i := 0; i < len(my_ancients); i++ {
		if my_ancients[i].uid == 20 {
			cooldown_level = my_ancients[i].level
			break
		}
	}
	//updates length of all abilities
	//not complete
	//!TODO
	Cooldowns[0].duration = Cooldowns[0].duration + int(my_ancients[16].level*2) //Clickstorm
	Cooldowns[1].duration = Cooldowns[1].duration + int(my_ancients[18].level*2) //Powersurge
	Cooldowns[2].duration = Cooldowns[2].duration + int(my_ancients[19].level*2) //Lucky Strikes
	Cooldowns[3].duration = Cooldowns[3].duration + int(my_ancients[21].level*2) //Metal Detector
	Cooldowns[4].duration = Cooldowns[4].duration + int(my_ancients[20].level*2) //GOlden Clicks
	Cooldowns[6].duration = Cooldowns[6].duration + int(my_ancients[17].level*2) //Super Clicks

	var cooldown_reduction float32
	if cooldown_level >= 280 {
		//if the level is greater than 280 we can simply set to 75%, no calcs needed
		cooldown_reduction = 73.56
	} else {
		//if under 280, we need to find the % reduction
		var mod_val float64 = -0.026
		test := mod_val * float64(cooldown_level)
		effect_mod := math.Pow(1.92, test)
		cooldown_reduction = float32(75 * (1 - effect_mod))
	}
	final_cooldown_mod := 1 - (cooldown_reduction / 100)
	for i := 0; i < len(Cooldowns); i++ { //works with static value, should be correct
		rounded_cooldown := float32(Cooldowns[i].cd) + .5
		Cooldowns[i].cd = int(rounded_cooldown * final_cooldown_mod)
	}
}

func update_ancients(save_JSON string) {
}

func autoClick_Polling(wait_group *sync.WaitGroup) {
	defer wait_group.Done()
}

func ascend() {
	//triggers ascension at best possible level
}

func transcend() {
	//triggers transcension at best possbile level
}
