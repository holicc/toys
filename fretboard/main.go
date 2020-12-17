package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	STRINGS        = 6
	POSITION_RANGE = 5
	FRET_GAP       = 11
)

var (
	Black  = Color("\033[1;30m%s\033[0m")
	Red    = Color("\033[1;31m%s\033[0m")
	Green  = Color("\033[1;32m%s\033[0m")
	Yellow = Color("\033[1;33m%s\033[0m")

	BaseNote = []string{"E", "B", "G", "D", "A", "E"}
	Notes    = []string{"C", "#C/bD", "D", "#D/bE", "E", "F", "#F/bG", "G", "#G/bA", "A", "#A/bB", "B"}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var line strings.Builder

	fretboard := make([]string, STRINGS)

	c := RandInt(STRINGS)
	r := RandInt(POSITION_RANGE)
	//
	line.WriteString(Yellow("1=C	"))
	for i := 1; i <= POSITION_RANGE; i++ {
		g := FRET_GAP
		line.WriteString(strconv.Itoa(i))
		line.WriteString(strings.Repeat(" ", g))
	}
	fmt.Println(line.String())
	line.Reset()
	for i := range fretboard {
		line.WriteString(Black("│"))
		for j := 1; j <= POSITION_RANGE; j++ {
			if i+1 == c && j == r {
				line.WriteString(strings.Repeat("-", POSITION_RANGE))
				line.WriteString(Red("?"))
				line.WriteString(strings.Repeat("-", POSITION_RANGE))
			} else {
				line.WriteString(strings.Repeat("-", FRET_GAP))
			}
			line.WriteString("|")
		}
		fretboard[i] = BaseNote[i] + ":" + line.String()
		line.Reset()
	}
	for _, s := range fretboard {
		fmt.Println(s)
	}
	//
	Quiz(r, c)
}

func Quiz(r int, c int) {
	var (
		correct = Answer(r, c)
		ans     string
		start   = time.Now()
	)
	fmt.Println("What is the note ? ")
	fmt.Print("> ")
	_, _ = fmt.Scan(&ans)
	fmt.Println("⌛️", math.Floor(time.Since(start).Seconds()), "s")
	if ans == correct || strings.Contains(correct, ans) {
		fmt.Println(Green("Correct !"))
	} else {
		fmt.Println(Red("Wrong,") + "answer is: " + Green(correct))
	}
}

func Answer(p, s int) string {
	r := BaseNote[s-1]
	index := 0
	for i := range Notes {
		if Notes[i] == r {
			index = i
			break
		}
	}
	if index+p >= len(Notes) {
		return Notes[(index+p)%len(Notes)]
	} else {
		return Notes[index+p]
	}
}

func RandInt(n int) (r int) {
	for {
		r = rand.Intn(n)
		if r != 0 {
			return
		}
	}
}

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}
