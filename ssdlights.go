package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

// "rand"

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(err)
		fmt.Println(msg)
		panic(msg)
	}
}

//    0x05: SetWhiteChannelGroup - 5 G Ch I
//    0x06: SetWhiteChannelUnit - 6 U Ch I
func danLightsWhite(zone int, unitNum int, groupNum int, intensities []int) {
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.100.148:9000")
	// serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	checkError(err, "point 21")
	conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err, "point 23")
	defer conn.Close()
	buf := make([]byte, 4)
	// turn wall light red
	for chNum := 0; chNum < 4; chNum++ {
		if groupNum != 0 {
			buf[0] = byte(5)        // 5 == SetWhiteChannelGroup
			buf[1] = byte(groupNum) // U -- unit address
			buf[2] = byte(chNum)    // channel
			buf[3] = byte(intensities[chNum])
		} else {
			buf[0] = byte(6)       // 6 == SetWhiteChannelUnit
			buf[1] = byte(unitNum) // U -- unit address
			buf[2] = byte(chNum)   // channel
			buf[3] = byte(intensities[chNum])
		}
		fmt.Println("Point 43 Sending: ", buf)
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println(buf, err)
			panic("Write to connection")
		}
	}
}

//    0x03: SetHSVgroup - 3 Z G Fl H S V Ft Ft
//    0x04: SetHSVunit -  4 Z U Fl H S V Ft Ft
func danLightsColor(zone int, unitNum int, groupNum int, hueStart int, hueStop int, saturation int, brightness int, fade int) {
	if (zone == 202) && (unitNum == 217) {
		if brightness > 128 {
			brightness = 128
		}
	}
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.100.148:9000")
	// serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	checkError(err, "point 21")
	conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err, "point 23")
	defer conn.Close()
	buf := make([]byte, 9)
	// turn wall light red
	if saturation == 255 {
		saturation = 254 // bug in Dan's lights
	}
	if groupNum != 0 {
		buf[0] = byte(3)          // 3 == SetHSVGroup
		buf[1] = byte(zone)       // Z -- zone address
		buf[2] = byte(groupNum)   // G -- unit address
		buf[3] = byte(14)         // FL -- HSV Flags
		buf[4] = byte(hueStart)   // H -- hue
		buf[5] = byte(saturation) // S -- saturation
		buf[6] = byte(brightness) // V -- value
		buf[7] = byte(0)          // FTh -- fade time, high bits
		buf[8] = byte(0)          // FTl -- fade time, low bits
	} else {
		buf[0] = byte(4)          // 2 == SetRGBunit, 4 == SetHSVunit
		buf[1] = byte(zone)       // Z -- zone address
		buf[2] = byte(unitNum)    // U -- unit address
		buf[3] = byte(14)         // FL -- HSV Flags
		buf[4] = byte(hueStart)   // H -- hue
		buf[5] = byte(saturation) // S -- saturation
		buf[6] = byte(brightness) // V -- value
		buf[7] = byte(0)          // FTh -- fade time, high bits
		buf[8] = byte(0)          // FTl -- fade time, low bits
	}
	fmt.Println("Point 88 Sending: ", buf)
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println(buf, err)
		panic("Write to connection")
	}
	if fade > 0 {
		brightness = brightness >> 2 // Hack to make pulses
		FTl := fade & 255
		FTh := (fade - FTl) >> 8
		buf[0] = byte(4)          // 2 == SetRGBunit
		buf[1] = byte(zone)       // Z -- zone address
		buf[2] = byte(unitNum)    // U -- unit address
		buf[3] = byte(14)         // U -- unit address
		buf[4] = byte(hueStop)    // H -- hue
		buf[5] = byte(saturation) // S -- saturation
		buf[6] = byte(brightness) // V -- value
		buf[7] = byte(FTh)        // FTh -- fade time, high bits
		buf[8] = byte(FTl)        // FTl -- fade time, low bits
		_, err = conn.Write(buf)
		if err != nil {
			fmt.Println(buf, err)
			panic("Write to connection")
		}
	}
}

func getBankInfo(bank string) (int, int, int, [][]int) {
	zone := 200
	offsetStart := 200
	count := 1
	var sequences [][]int
	switch bank {
	case "everything":
		zone = 0
		offsetStart = 0
		count = 1
	case "lobbywhite":
		zone = 100
		offsetStart = 1
		count = 1
	case "lobbywall":
		zone = 101
		offsetStart = 101
		count = 6
		sequences = [][]int{
			[]int{0, 1, 2, 3, 4, 5},
			{4, 5, 3, 2, 0, 1}}
	case "lobbylanterns":
		zone = 101
		offsetStart = 107
		count = 3
		sequences = [][]int{
			[]int{0, 1, 2},
			{0, 2, 1},
			{2, 1, 0},
			{1, 0, 2},
			{1, 2, 0},
			{2, 0, 1}}
	case "baywhite":
		zone = 200
		offsetStart = 201
		count = 12
		sequences = [][]int{
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			[]int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
			[]int{2, 1, 0, 5, 4, 3, 8, 7, 6, 11, 10, 9},
			[]int{9, 10, 11, 6, 7, 8, 3, 4, 5, 0, 1, 2},
			[]int{0, 5, 6, 11, 10, 7, 4, 1, 2, 3, 8, 9},
			[]int{9, 8, 3, 2, 1, 4, 7, 10, 11, 6, 5, 0},
			[]int{0, 5, 1, 2, 4, 6, 11, 7, 3, 8, 10, 9},
			[]int{9, 8, 10, 11, 7, 3, 2, 4, 6, 5, 1, 0}}
	case "baycolor":
		zone = 201
		offsetStart = 201
		count = 12

		sequences = [][]int{
			[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			[]int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0},
			[]int{2, 1, 0, 5, 4, 3, 8, 7, 6, 11, 10, 9},
			[]int{9, 10, 11, 6, 7, 8, 3, 4, 5, 0, 1, 2},
			[]int{0, 5, 6, 11, 10, 7, 4, 1, 2, 3, 8, 9},
			[]int{9, 8, 3, 2, 1, 4, 7, 10, 11, 6, 5, 0},
			[]int{0, 5, 1, 2, 4, 6, 11, 7, 3, 8, 10, 9},
			[]int{9, 8, 10, 11, 7, 3, 2, 4, 6, 5, 1, 0}}

	case "bayaccents":
		zone = 202
		offsetStart = 213
		count = 8
		sequences = [][]int{
			[]int{0, 1, 2, 3, 4, 5, 6, 7}}

	case "bayaccentsminusparty":
		zone = 202
		offsetStart = 213
		count = 4
		sequences = [][]int{
			[]int{0, 1, 2, 3}}

	case "partylights":
		zone = 202
		offsetStart = 217
		count = 4
		sequences = [][]int{
			[]int{0, 1, 2, 3}}

	default:
		fmt.Println("bank unrecognized: ", bank)
		panic("bank unrecognized")
	}
	return zone, offsetStart, count, sequences
}

// for unitNum := offsetStart; unitNum < (offsetStart + count); unitNum++ {

func doWhiteLights(whiteBrightness int) {
	bank := "baywhite"
	zone, offsetStart, _, sequences := getBankInfo(bank)
	delay := time.Duration(200000000)
	count := 0
	warmBright := whiteBrightness
	coldBright := int(float64(whiteBrightness) * 0.85)
	for _, seqce := range sequences {
		for _, seqVal := range seqce {
			fmt.Println("seqVal", seqVal)
			unitNum := offsetStart + seqVal
			var intensities []int
			switch count {
			case 0:
				intensities = []int{warmBright, 0, 0, 0}
			case 1:
				intensities = []int{0, warmBright, 0, 0}
			case 2:
				intensities = []int{0, 0, 0, coldBright}
			case 3:
				intensities = []int{0, 0, coldBright, 0}
			}
			groupNum := 0
			danLightsWhite(zone, unitNum, groupNum, intensities)
			time.Sleep(delay)
		}
		count++
		if count == 4 {
			count = 0
		}
	}
}

func blankWhiteLights(finalBrightness int) {
	bank := "baywhite"
	zone, offsetStart, count, _ := getBankInfo(bank)
	delay := time.Duration(200000000)
	for unitNum := offsetStart; unitNum < (offsetStart + count); unitNum++ {
		intensities := []int{finalBrightness, finalBrightness, finalBrightness, finalBrightness}
		groupNum := 0
		danLightsWhite(zone, unitNum, groupNum, intensities)
		time.Sleep(delay)
	}
}

func blankLobbyWhite(finalBrightness int) {
	bank := "lobbywhite"
	zone, offsetStart, count, _ := getBankInfo(bank)
	delay := time.Duration(200000000)
	for unitNum := offsetStart; unitNum < (offsetStart + count); unitNum++ {
		hueStart := rand.Intn(255)
		saturation := 254
		brightness := finalBrightness
		hueStop := rand.Intn(255)
		fade := 20
		groupNum := 0
		danLightsColor(zone, unitNum, groupNum, hueStart, hueStop, saturation, brightness, fade)
		time.Sleep(delay)
	}
}

type bankInfo struct {
	name            string
	zone            int
	start           int
	count           int
	sequences       [][]int
	currentSequence int
	currentPosition int
}

func doRandomBank(banksRunning []bankInfo, colorBrightness int) {
	delay := time.Duration(200000000)
	randomBank := 2 // rand.Intn(len(banksRunning) - 1)
	numSeqs := len(banksRunning[randomBank].sequences)
	for seqNum := 0; seqNum < numSeqs; seqNum++ {
		numPoses := len(banksRunning[randomBank].sequences[seqNum])
		for posNum := 0; posNum < numPoses; posNum++ {
			unitNum := banksRunning[randomBank].start + banksRunning[randomBank].sequences[seqNum][posNum]

			zone := banksRunning[randomBank].zone
			hueStart := rand.Intn(255)
			saturation := 254
			brightness := colorBrightness
			hueStop := rand.Intn(255)
			fade := 30
			groupNum := 0
			danLightsColor(zone, unitNum, groupNum, hueStart, hueStop, saturation, brightness, fade)
			time.Sleep(delay)
		}
	}
}

func strobeWhiteLights(whiteBrightness int) {
	bank := "baywhite"
	zone, _, _, _ := getBankInfo(bank)
	delay := time.Duration(200000000)

	unitNum := 0
	hueStart := rand.Intn(255)
	saturation := 0
	brightness := whiteBrightness
	hueStop := rand.Intn(255)
	fade := 10
	groupNum := 0
	danLightsColor(zone, unitNum, groupNum, hueStart, hueStop, saturation, brightness, fade)
	time.Sleep(delay)

	hueStart = rand.Intn(255)
	saturation = 0
	brightness = 0
	hueStop = rand.Intn(255)
	fade = 10
	groupNum = 0
	danLightsColor(zone, unitNum, groupNum, hueStart, hueStop, saturation, brightness, fade)
	time.Sleep(delay)

}

func multiplexColorLights(bankList []string, colorBrightness int, whiteBrightness int) {
	bankCount := len(bankList)
	banksRunning := make([]bankInfo, bankCount)
	for idx, bankName := range bankList {
		banksRunning[idx].name = bankName
		banksRunning[idx].zone, banksRunning[idx].start, banksRunning[idx].count, banksRunning[idx].sequences = getBankInfo(bankName)
		banksRunning[idx].currentSequence = 0
		banksRunning[idx].currentPosition = 0
	}
	delay := time.Duration(200000000)
	currentBank := 0
	tick := 1
	for {
		fmt.Println("tick", tick)
		if (tick & 127) == 0 {
			if (tick & 1023) == 0 {
				strobeWhiteLights(whiteBrightness)
			} else {
				if (tick * 511) == 0 {
					doWhiteLights(whiteBrightness)
					blankWhiteLights(0)
					blankLobbyWhite(0)
				}
				doRandomBank(banksRunning, colorBrightness)
			}
		}

		zone := banksRunning[currentBank].zone
		unitNum := banksRunning[currentBank].start + banksRunning[currentBank].sequences[banksRunning[currentBank].currentSequence][banksRunning[currentBank].currentPosition]
		hueStart := rand.Intn(255)
		saturation := 254
		brightness := colorBrightness
		hueStop := rand.Intn(255)
		fade := 35
		groupNum := 0
		danLightsColor(zone, unitNum, groupNum, hueStart, hueStop, saturation, brightness, fade)
		time.Sleep(delay)

		banksRunning[currentBank].currentPosition++
		if banksRunning[currentBank].currentPosition >= len(banksRunning[currentBank].sequences[banksRunning[currentBank].currentSequence]) {
			banksRunning[currentBank].currentPosition = 0
			banksRunning[currentBank].currentSequence++
			if banksRunning[currentBank].currentSequence >= len(banksRunning[currentBank].sequences) {
				banksRunning[currentBank].currentSequence = 0
			}
		}

		currentBank++
		if currentBank >= bankCount {
			currentBank = 0
		}
		tick++

	}
}

func main() {
	// get flags
	var wflag = flag.Int("w", 128, "white brightness")
	var cflag = flag.Int("c", 255, "color brightness")
	var eflag = flag.Bool("e", false, "end (party's over)")
	flag.Parse()
	whiteBrightness := *wflag
	colorBrightness := *cflag
	endPartysOver := *eflag
	// if whiteBrightness > 254 {
	//	whiteBrightness = 254 // dan bug
	// }
	fmt.Println("whiteBrightness", whiteBrightness)
	fmt.Println("colorBrightness", colorBrightness)
	fmt.Println("endPartysOver", endPartysOver)

	current := time.Now()
	asint := current.UnixNano()
	rand.Seed(asint)

	if endPartysOver {
		doWhiteLights(whiteBrightness)
		blankWhiteLights(255)
		blankLobbyWhite(255)
	} else {
		doWhiteLights(whiteBrightness)
		blankWhiteLights(0)
		blankLobbyWhite(0)
		multiplexColorLights([]string{"lobbywall", "lobbylanterns", "baycolor", "bayaccentsminusparty"}, colorBrightness, whiteBrightness)

	}
}
