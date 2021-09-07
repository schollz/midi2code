package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"runtime"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/micmonay/keybd_event"
)

type Press struct {
	key   int
	shift bool
}

var keyMap = map[string]Press{
	"a":  Press{key: keybd_event.VK_A},
	"b":  Press{key: keybd_event.VK_B},
	"c":  Press{key: keybd_event.VK_C},
	"d":  Press{key: keybd_event.VK_D},
	"e":  Press{key: keybd_event.VK_E},
	"f":  Press{key: keybd_event.VK_F},
	"g":  Press{key: keybd_event.VK_G},
	"h":  Press{key: keybd_event.VK_H},
	"i":  Press{key: keybd_event.VK_I},
	"j":  Press{key: keybd_event.VK_J},
	"k":  Press{key: keybd_event.VK_K},
	"l":  Press{key: keybd_event.VK_L},
	"m":  Press{key: keybd_event.VK_M},
	"n":  Press{key: keybd_event.VK_N},
	"o":  Press{key: keybd_event.VK_O},
	"p":  Press{key: keybd_event.VK_P},
	"q":  Press{key: keybd_event.VK_Q},
	"r":  Press{key: keybd_event.VK_R},
	"s":  Press{key: keybd_event.VK_S},
	"t":  Press{key: keybd_event.VK_T},
	"u":  Press{key: keybd_event.VK_U},
	"v":  Press{key: keybd_event.VK_V},
	"w":  Press{key: keybd_event.VK_W},
	"x":  Press{key: keybd_event.VK_X},
	"y":  Press{key: keybd_event.VK_Y},
	"z":  Press{key: keybd_event.VK_Z},
	".":  Press{key: keybd_event.VK_DOT},
	"0":  Press{key: keybd_event.VK_0},
	"1":  Press{key: keybd_event.VK_1},
	"2":  Press{key: keybd_event.VK_2},
	"3":  Press{key: keybd_event.VK_3},
	"4":  Press{key: keybd_event.VK_4},
	"5":  Press{key: keybd_event.VK_5},
	"6":  Press{key: keybd_event.VK_6},
	"7":  Press{key: keybd_event.VK_7},
	"8":  Press{key: keybd_event.VK_8},
	"9":  Press{key: keybd_event.VK_9},
	";":  Press{key: keybd_event.VK_SEMICOLON},
	"(":  Press{keybd_event.VK_9, true},
	")":  Press{keybd_event.VK_0, true},
	"-":  Press{key: keybd_event.VK_MINUS},
	"_":  Press{keybd_event.VK_MINUS, true},
	"\n": Press{key: keybd_event.VK_ENTER}, // TODO: for live coding make this shift
}

func softcut_random_loop() string {
	start := float64(int((rand.Float64()*100)*100)) / 100
	duration := float64(int((rand.Float64()*1+0.05)*100)) / 100
	voice := rand.Intn(6) + 1
	t := template.Must(template.New("my").Parse("clock.run(function() softcut.loop_start({{.Voice}},{{.Start}}); softcut.loop_end({{.Voice}},{{.Start}}+10); softcut.rec_level({{.Voice}},0.5); softcut.position({{.Voice}},{{.Start}}); clock.sleep({{.Duration}}+0.2); softcut.loop_end({{.Voice}},{{.Start}}+{{.Duration}}); softcut.pos({{.Voice}},{{.Start}}) end)"))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, struct {
		Start    float64
		Duration float64
		Voice    int
	}{start, duration, voice})
	if err != nil {
		panic(err)
	}
	return tpl.String()
}

func main() {
	fmt.Println(softcut_random_loop())
	fmt.Println("running")
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	// For linux, it is very important to wait 2 seconds
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	line := "softcut.loop_start(1.0)\n"
	for i := 0; i < 5; i++ {
		line = line + line
	}
	fmt.Println(line)
	// for _, char := range line {
	// 	time.Sleep(time.Duration(100+rand.Intn(100)) * time.Millisecond)
	// 	s := string(char)
	// 	fmt.Print(s)
	// 	kb.SetKeys(keyMap[s].key)
	// 	kb.HasSHIFT(keyMap[s].shift)
	// 	err = kb.Launching()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	chars := []rune(line)
	i := 0

	// test with
	// for i in {1..10}; do; echo -n "/key" | nc -w 0 -u localhost 8765; sleep 0.2; done
	addr := "127.0.0.1:8765"
	d := osc.NewStandardDispatcher()
	d.AddMsgHandler("/key", func(msg *osc.Message) {
		s := string(chars[i])
		fmt.Print(s)
		kb.SetKeys(keyMap[s].key)
		kb.HasSHIFT(keyMap[s].shift)
		err = kb.Launching()
		if err != nil {
			panic(err)
		}
		i++
	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	fmt.Println("server listening on 8765")
	server.ListenAndServe()
}
