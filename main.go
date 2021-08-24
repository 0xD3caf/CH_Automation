package main

import (
	"bytes"
	"compress/zlib"
	b64 "encoding/base64"
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/thedevsaddam/gojsonq/v2"
	"io"
	"math"
	"os"
	"strconv"
)

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

var gamesave = "clickerHeroSave.txt"

var hash = "7a990d4252c6fb53aacfbb0ec1a3b23"
var alphabet = "abcdefghijklmnopqrstuvwxyz"
var salt = "af0ik392jrmt0nsfdghy0"
var delimiter = "Fe12NAfA3R6z4k0z"
var Cooldowns = [9]Skill{
	// -1 indicates NA for that skill
	//ORDER CURRENTLY MATTERS
	Skill{"Clickstorm", 600, 30},
	Skill{"Powersurge", 600, 30},
	Skill{"Lucky Strikes", 1800, 30},
	Skill{"Metal Detector", 1800, 30},
	Skill{"Golden Clicks", 3200, 30},
	Skill{"The Dark Ritual", 28800, -1},
	Skill{"Super Clicks", 1800, 60},
	Skill{"Energize", 3200, -1},
	Skill{"Reload", 3200, -1}}

var my_ancients = [26]Ancient{
	//ORDER CURRENTLY MATTERS
	Ancient{"Libertas", 0, 0.25, 4},
	Ancient{"Siyalatas", 0, 0.25, 5},
	Ancient{"Mammon", 0, 0.05, 8},
	Ancient{"Mimzee", 0, 0.5, 9},
	Ancient{"Pluto", 0, .30, 10},
	Ancient{"Dogcog", 0, 0, 11},
	Ancient{"Fortuna", 0, 0, 12},
	Ancient{"Atman", 0, 0, 13},
	Ancient{"Dora", 0, 0, 14},
	Ancient{"Bhaal", 0, 0.15, 15},
	Ancient{"Morgulis", 0, 0.11, 16},
	Ancient{"Chronos", 0, 0, 17},
	Ancient{"Bubos", 0, 0, 18},
	Ancient{"Fragsworth", 0, 0.20, 19},
	Ancient{"Vaagur", 0, 0, 20},
	Ancient{"Kumawakamaru", 0, 0, 21},
	Ancient{"Chawedo", 0, 2.0, 22},
	Ancient{"Hetoncheir", 0, 2.0, 23},
	Ancient{"Beserker", 0, 2.0, 24},
	Ancient{"Sniperino", 0, 2.0, 25},
	Ancient{"Kleptos", 0, 2.0, 26},
	Ancient{"Energon", 0, 2.0, 27},
	Ancient{"Argaiv", 0, 2.0, 28},
	Ancient{"Juggernaut", 0, 0.01, 29},
	Ancient{"Revloc", 0, 0, 31},
	Ancient{"Nogardnit", 0, 0.10, 32},
}

func init() {
	//starts by adding extra time to length and computing reduced cooldowns\
	parse_save(gamesave)
	update_ancients("save_JSON.txt")
	update_cooldowns()
	for i := 0; i < len(Cooldowns); i++ {
		//Cooldowns[i].cd := Cooldowns[i].cd / modval
		//Cooldowns[i].duration := Cooldowns[i].duration / modval
	}

}
func main() {
	fmt.Println("Clicker Heroes Automation")
	main_hero_upgrade()
	auto_abilities()
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
	bitmap := robotgo.CaptureScreen()
	// use `defer robotgo.FreeBitmap(bit)` to free the bitmap
	tolerance := 0.01
	var color robotgo.CHex = 0xC0D431
	x, y := robotgo.FindColor(color, bitmap, tolerance)
	if x == -1 && y == -1 {
		fmt.Println("The color could not be found")
	} else {
		robotgo.MoveClick(x, y, "LEFT_BUTTON", false)
	}
	robotgo.FreeBitmap(bitmap)
}

func main_hero_upgrade() {
	//keeps the current max hero at max levels
	return
}

func prev_hero_upgrade() {
	return
}

func auto_abilities() {
	/*for i := 0; i < len(Cooldowns); i++ {
		fmt.Println("Skill : ", Cooldowns[i].name, "has cooldown: ", Cooldowns[i].cd, " and duration: ", Cooldowns[i].duration)
	}*/
	//auto clicks on any upgrade where the duration is longer than cooldown, ie it can always be active
	return
}

func update_save() {
	//saves the game and then recopies and parses save
	return
}
func parse_save(input_file string) {
	//var final_str string
	file, err := os.Open(input_file)
	if err != nil {
		fmt.Println("The file was invalid")
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
		fmt.Println("We failed to decode the string")
	}
	//zlib decompress
	save_bytes := bytes.NewReader(decoded_save_string)
	output, err := zlib.NewReader(save_bytes)
	defer output.Close()
	final_string, err := io.ReadAll(output)
	file, err = os.OpenFile("save_JSON.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(final_string)
	//first loop goes through the set of chars pulled from file
	return
}

func replace_save(curent_save string, new_save string) {
	//replaces old save with new version including new edits
	//start by zlib compressing

	//then b64 encode

	//add hash to start of file

	//restart game ???
	return
}

func update_cooldowns() {
	//takes cd reducation and duration increase from file
	//updates ability table
	cooldown_level := my_ancients[14].level
	Cooldowns[0].duration = Cooldowns[0].duration + int(my_ancients[16].level*2) //Clickstorm
	Cooldowns[1].duration = Cooldowns[1].duration + int(my_ancients[18].level*2) //Powersurge
	Cooldowns[2].duration = Cooldowns[2].duration + int(my_ancients[19].level*2) //Lucky Strikes
	Cooldowns[3].duration = Cooldowns[3].duration + int(my_ancients[21].level*2) //Metal Detector
	Cooldowns[4].duration = Cooldowns[4].duration + int(my_ancients[20].level*2) //GOlden Clicks
	Cooldowns[6].duration = Cooldowns[6].duration + int(my_ancients[17].level*2) //Super Clicks
	var cooldown_reduction float32
	if cooldown_level >= 280 {
		//if the level is greater than 280 we can simply set to 75%, no calcs needed
		cooldown_reduction = 0.75
	} else {
		//if under 280, we need to find the % reduction
		var mod_val float64 = -0.026
		test := mod_val * float64(cooldown_level)
		effect_mod := math.Pow(1.92, test)
		fmt.Println(effect_mod)
		cooldown_reduction = float32(75 * (1 - effect_mod))
	}
	final_cooldown_reduction := 1 - (cooldown_reduction / 100)
	for i := 0; i < len(Cooldowns); i++ {
		Cooldowns[i].cd = int(float32(Cooldowns[i].cd) * final_cooldown_reduction)
	}
	return
}

func update_ancients(save_JSON string) {
	//opens file
	file, err := os.Open(save_JSON)
	if err != nil {
		panic(err)
	}
	//collect data and close file
	save_data, err2 := io.ReadAll(file)
	if err2 != nil {
		panic(err)
	}
	file.Close()

	save_string := string(save_data) //convert bytes to string
	var search_string string
	var output_int int64
	for i := 0; i < len(my_ancients); i++ {
		ancient_id := my_ancients[i]
		search_string = "ancients.ancients." + strconv.Itoa(ancient_id.uid) + ".level"
		JSON_resp := gojsonq.New().FromString(save_string).Find(search_string)
		output := JSON_resp.(string)
		big_val := false
		for i := 0; i < len(output); i++ {
			if string(output[i]) == "e" {
				big_val = true
			}
		}
		if big_val == false {
			output_float, _ := strconv.ParseFloat(output, 64)
			output_int = int64(math.Round(output_float))
		} else {
			output_int = int64(-1)
		}
		my_ancients[i].level = int64(output_int)
	}
}
