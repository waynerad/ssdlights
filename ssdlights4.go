package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

//    0x05: SetWhiteChannelGroup - 5 G Ch I
//    0x06: SetWhiteChannelUnit - 6 U Ch I
func danLightsWhite(zone int, unitNum int, groupNum int, intensities []int) {
	panic("Point 24")
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.50.131:9000")
	// serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	checkError(err)
	conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err)
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
func danLightsColor(zone int, unitNum int, groupNum int, hue int, saturation int, brightness int, fade int) {
	if (zone == 202) && (unitNum == 217) {
		if brightness > 128 {
			brightness = 128
		}
	}
	serverAddr, err := net.ResolveUDPAddr("udp", "192.168.50.131:9000")
	// serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9000")
	checkError(err)
	conn, err := net.DialUDP("udp", nil, serverAddr)
	checkError(err)
	defer conn.Close()
	buf := make([]byte, 9)
	// turn wall light red
	if saturation == 255 {
		saturation = 254 // bug in Dan's lights
	}
	fadeLow := fade & 255
	fadeHigh := fade >> 8
	if groupNum != 0 {
		buf[0] = byte(3)          // 3 == SetHSVGroup
		buf[1] = byte(zone)       // Z -- zone address
		buf[2] = byte(groupNum)   // G -- unit address
		buf[3] = byte(14)         // FL -- HSV Flags
		buf[4] = byte(hue)        // H -- hue
		buf[5] = byte(saturation) // S -- saturation
		buf[6] = byte(brightness) // V -- value
		buf[7] = byte(fadeHigh)   // FTh -- fade time, high bits
		buf[8] = byte(fadeLow)    // FTl -- fade time, low bits
	} else {
		buf[0] = byte(4)          // 2 == SetRGBunit, 4 == SetHSVunit
		buf[1] = byte(zone)       // Z -- zone address
		buf[2] = byte(unitNum)    // U -- unit address
		buf[3] = byte(14)         // FL -- HSV Flags
		buf[4] = byte(hue)        // H -- hue
		buf[5] = byte(saturation) // S -- saturation
		buf[6] = byte(brightness) // V -- value
		buf[7] = byte(fadeHigh)   // FTh -- fade time, high bits
		buf[8] = byte(fadeLow)    // FTl -- fade time, low bits
	}
	fmt.Println("Point 88 Sending: ", buf)
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println(buf, err)
		panic("Write to connection")
	}
	return
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
	case "lobbywallandlanterns":
		zone = 101
		offsetStart = 101
		count = 9
		sequences = [][]int{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8}}
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

func getStandardMessageDelay() time.Duration {
	return time.Duration(200000000)
}

// for unitNum := offsetStart; unitNum < (offsetStart + count); unitNum++ {

func doWhiteLights(whiteBrightness int) {
	panic("Point 216")
	bank := "baywhite"
	zone, offsetStart, _, sequences := getBankInfo(bank)
	delay := getStandardMessageDelay()
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
	panic("Point 250")
	bank := "baywhite"
	zone, offsetStart, count, _ := getBankInfo(bank)
	delay := getStandardMessageDelay()

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
	delay := getStandardMessageDelay()
	for unitNum := offsetStart; unitNum < (offsetStart + count); unitNum++ {
		hue := rand.Intn(255)
		saturation := 254
		brightness := finalBrightness
		fade := 20
		groupNum := 0
		danLightsColor(zone, unitNum, groupNum, hue, saturation, brightness, fade)
		time.Sleep(delay)
	}
}

func getGroupsForBank(bankName string) (int, []int) {
	switch bankName {
	case "baycolorcolumns":
		return 201, []int{20, 21, 22}
	case "baycolorrows":
		return 201, []int{10, 11, 12, 13}
	default:
		panic("Bank name is missing or invalid")
	}
}

func getUnitsForBank(bankName string) (int, []int) {
	if bankName == "lobbywallandlanterns" {
		return 101, []int{101, 102, 103, 109, 104, 108, 105, 106, 107}
		return 101, []int{101, 102, 103, 104, 105, 106, 107, 108, 109}
		return 101, []int{101, 102, 107, 103, 108, 104, 105, 106, 109}
	} else {
		zone, offsetStart, count, _ := getBankInfo(bankName)
		unitOrGroupNumbers := make([]int, count)
		for ii := 0; ii < count; ii++ {
			unitOrGroupNumbers[ii] = offsetStart + ii
		}
		return zone, unitOrGroupNumbers
	}
}

type lightingEffect struct {
	zone           int
	unitNum        int
	groupNum       int
	fromHue        int
	fromSaturation int
	fromBrightness int
	toHue          int
	toSaturation   int
	toBrightness   int
	fadeTime       int64
	pauseTime      int64
	endOfStream bool
}

type iEffectProducer interface {
	init(hue int, delay int64, colorCycleSpeed int, alwaysCycleColors bool, bankName string, bankIsGroups bool, direction int, timeToRun int64) 
	nextEffect() lightingEffect
}

type waveEffect struct {
	delay              int64
	colorCycleSpeed    int
	leadingEdge        int
	trailingEdge       int
	currentHue         int
	lastHues           []int
	alwaysCycleColors bool
	zone               int
	unitsOrGroups      int
	unitOrGroupNumbers []int
	useGroups          bool
	falling            bool
	direction          int
	timeLeft int64
}

func (self *waveEffect) init(hue int, delay int64, colorCycleSpeed int, alwaysCycleColors bool, bankName string, bankIsGroups bool, direction int, timeToRun int64) {
	self.delay = delay
	self.colorCycleSpeed = colorCycleSpeed
	self.leadingEdge = -1
	self.trailingEdge = -2
	self.currentHue = hue
	self.falling = false
	self.alwaysCycleColors = alwaysCycleColors
	self.direction = direction
	if bankIsGroups {
		self.useGroups = true
		self.zone, self.unitOrGroupNumbers = getGroupsForBank(bankName)
	} else {
		self.zone, self.unitOrGroupNumbers = getUnitsForBank(bankName)

	}
	self.unitsOrGroups = len(self.unitOrGroupNumbers)
	if self.direction == -1 {
		self.leadingEdge = self.unitsOrGroups
		self.trailingEdge = self.unitsOrGroups + 1
	}
	self.lastHues = make([]int, self.unitsOrGroups)
	self.timeLeft = timeToRun
}

func (self *waveEffect) nextEffect() lightingEffect {
	var result lightingEffect
	result.zone = self.zone
	result.unitNum = 0
	result.groupNum = 0
	result.fromSaturation = 255
	result.toSaturation = 255

	if self.falling {
		// falling edge
		self.trailingEdge += self.direction
		if self.trailingEdge == self.unitsOrGroups {
			self.trailingEdge = 0
		}
		if self.trailingEdge < 0 {
			self.trailingEdge = self.unitsOrGroups - 1
		}
		if (self.trailingEdge >= 0) && (self.trailingEdge < self.unitsOrGroups) {
			if self.useGroups {
				result.groupNum = self.unitOrGroupNumbers[self.trailingEdge]
			} else {
				result.unitNum = self.unitOrGroupNumbers[self.trailingEdge]
			}
			result.fromHue = self.lastHues[self.trailingEdge]
			result.fromBrightness = 255
			result.toHue = self.lastHues[self.trailingEdge]
			result.toBrightness = 0
			result.fadeTime = self.delay
			result.pauseTime = 0
		} else {
			self.falling = false // flip it and fall through to execute rising code
		}
	}
	if !self.falling {
		// rising edge
		self.leadingEdge += self.direction
		cyclecolors := false
		if self.leadingEdge == self.unitsOrGroups {
			self.leadingEdge = 0
			cyclecolors = true
		}
		if self.leadingEdge < 0 {
			self.leadingEdge = self.unitsOrGroups - 1
			cyclecolors = true
		}
		if self.alwaysCycleColors {
			cyclecolors = true // always cycle colors? // Andrew's suggestion
		}
		if cyclecolors {
			self.currentHue += self.colorCycleSpeed
			if self.currentHue > 255 {
				self.currentHue -= 256
			}
		}
		fmt.Println("point 348 self.leadingEdge", self.leadingEdge)
		self.lastHues[self.leadingEdge] = self.currentHue
		// rising
		if self.useGroups {
			result.groupNum = self.unitOrGroupNumbers[self.leadingEdge]
		} else {
			result.unitNum = self.unitOrGroupNumbers[self.leadingEdge]
		}
		result.fromHue = self.currentHue
		result.fromBrightness = 0
		result.toHue = self.currentHue
		result.toBrightness = 255
		result.fadeTime = self.delay
		result.pauseTime = self.delay
	}
	self.falling = !self.falling
	self.timeLeft -= result.pauseTime
	result.endOfStream = false
	if self.timeLeft < 0 {
	result.endOfStream = true
	}
	return result
}

type effectElement struct {
	producer       iEffectProducer
	nextEffectTime int64
}

func executeEffect(effect lightingEffect) {
	if effect.fadeTime == 0 {
		// no fade, use only "from" hue, saturation, and brightness
		danLightsColor(effect.zone, effect.unitNum, effect.groupNum, effect.fromHue, effect.fromSaturation, effect.fromBrightness, 0)
	} else {
		// fade: use "from" and "to" hue, saturation, and brightness, and convert fade time units
		danLightsColor(effect.zone, effect.unitNum, effect.groupNum, effect.fromHue, effect.fromSaturation, effect.fromBrightness, 0)
		convertedFadeTime := effect.fadeTime / 100000000
		danLightsColor(effect.zone, effect.unitNum, effect.groupNum, effect.toHue, effect.toSaturation, effect.toBrightness, int(convertedFadeTime))
	}
	delay := getStandardMessageDelay()
	time.Sleep(delay)
}

func instantiateNewBayWaveEffect(rnd *rand.Rand) effectElement {
	var waver1 waveEffect
	hue := rnd.Intn(256)
	var fastestDelay int64
	var slowestDelay int64
	var delay int64
	fastestDelay = 300000000
	slowestDelay = 950000000
	delay = int64(rnd.Intn(int(slowestDelay-fastestDelay+1))) + fastestDelay
	fastestCycle := 5
	slowestCycle := 23
	colorCycleSpeed := rnd.Intn(slowestCycle-fastestCycle+1) + fastestCycle
	var alwaysCycleColors bool
	alwaysCycleColors = (rnd.Intn(2) == 0)
	var bankName string
	if rnd.Intn(2) == 0 {
		bankName = "baycolorrows"
	} else {
		bankName ="baycolorcolumns"
	}
	bankIsGroups := true
	direction := (2 * rnd.Intn(2)) - 1 // +1 or -1
	var minTime int64
	var maxTime int64
	var timeToRun int64
	minTime = 120000000000 // 120 seconds = 2 minutes
	maxTime = 960000000000 // 960 seconds = 16 minutes
	timeToRun = int64(rnd.Intn(int(maxTime-minTime+1))) + minTime
	waver1.init(hue, delay, colorCycleSpeed, alwaysCycleColors, bankName, bankIsGroups, direction, timeToRun)
	var wavEle1 effectElement
	wavEle1.producer = &waver1
	wavEle1.nextEffectTime = time.Now().UnixNano()
	return wavEle1
}

func instantiateNewLobbyWaveEffect(rnd *rand.Rand) effectElement {
	var waver2 waveEffect
	hue := rnd.Intn(256)
	var fastestDelay int64
	var slowestDelay int64
	var delay int64
	fastestDelay = 300000000
	slowestDelay = 950000000
	delay = int64(rnd.Intn(int(slowestDelay-fastestDelay+1))) + fastestDelay
	fastestCycle := 5
	slowestCycle := 23
	colorCycleSpeed := rnd.Intn(slowestCycle-fastestCycle+1) + fastestCycle
	var alwaysCycleColors bool
	alwaysCycleColors = (rnd.Intn(2) == 0)
	bankName := "lobbywallandlanterns"
	bankIsGroups := false
	direction := (2 * rnd.Intn(2)) - 1 // +1 or -1
	var minTime int64
	var maxTime int64
	var timeToRun int64
	minTime = 120000000000 // 120 seconds = 2 minutes
	maxTime = 960000000000 // 960 seconds = 16 minutes
	timeToRun = int64(rnd.Intn(int(maxTime-minTime+1))) + minTime
	waver2.init(hue, delay, colorCycleSpeed, alwaysCycleColors, bankName, bankIsGroups, direction, timeToRun)
	var wavEle2 effectElement
	wavEle2.producer = &waver2
	wavEle2.nextEffectTime = time.Now().UnixNano()
	return wavEle2
}

func main() {
	rnd := rand.New(rand.NewSource(0))
	var effectBank []effectElement
	effectBank = make([]effectElement, 2)
	//
	// Set up bay wave effect
	effectBank[0] = instantiateNewBayWaveEffect(rnd)
 
	//
	// Set up lobby wave effect
	effectBank[1] = instantiateNewLobbyWaveEffect(rnd)

	// loop forever to multiplex signals on to the same UDP message stream
	for {
		currentTime := time.Now().UnixNano()
		// fmt.Println("currentTime", currentTime)
		for ii := 0; ii < len(effectBank); ii++ {
			if currentTime >= effectBank[ii].nextEffectTime {
				nextEffect := effectBank[ii].producer.nextEffect()
				fmt.Println("nextEffect!", nextEffect)
				executeEffect(nextEffect)
				if nextEffect.endOfStream {
					if ii == 0 {
						effectBank[0] = instantiateNewBayWaveEffect(rnd ) 
					}
					if ii == 1 {
						effectBank[1] = instantiateNewLobbyWaveEffect(rnd)
					}					
				} else {
					effectBank[ii].nextEffectTime = currentTime + nextEffect.pauseTime
				}
			}
		}
	}
}
