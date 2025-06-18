package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/beevik/ntp"
	"log"
	"math"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"
)

func NTPCurrentTime() string {
	t, err := ntp.Time("0.pool.ntp.org")
	layout := "Mon Jan _2 15:04:05 2006"
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("Time: %v\n", t.Format(layout))
}

func StringShortener() (string, error) {
	scanner := bufio.NewReader(os.Stdin)
	curr, _ := scanner.ReadString('\n')
	curr = strings.TrimSpace(curr)
	curr = strings.Replace(curr, " ", "", -1)
	runes := []rune(curr)
	currDigit, digit := 0, 0
	builder := strings.Builder{}
	for i := len(runes) - 1; i >= 0; i-- {
		for i >= 0 && unicode.IsDigit(runes[i]) {
			currDigit += int(runes[i]-'0') * int(math.Pow(10, float64(digit)))
			i--
			digit++
		}
		if i == -1 {
			return "", errors.New("incorrect input")
		}
		if currDigit == 0 {
			builder.WriteRune(runes[i])
		}
		builder.WriteString(strings.Repeat(string(runes[i]), currDigit))
		currDigit = 0
		digit = 0
	}
	b := builder.String()
	runesB := []rune(b)
	for l, r := 0, len(b)-1; l < r; {
		runesB[l], runesB[r] = runesB[r], runesB[l]
		l++
		r--
	}
	if curr != "" && string(runesB) == "" {
		return "", errors.New("invalid input")
	}
	return string(runesB), nil
}

func MapAnagrams(arr *[]string) *map[string]*[]string {
	m := make(map[string]*[]string)
	supported := make(map[string]string)
	alreadyWas := make(map[string]struct{})
	for _, word := range *arr {
		word = strings.ToLower(word)
		splitted := strings.Split(word, "")
		sort.Strings(splitted)
		wordSplitted := strings.Join(splitted, "")
		if _, ok := supported[wordSplitted]; !ok {
			supported[wordSplitted] = word
		}
		if m[supported[wordSplitted]] == nil {
			empty := []string{}
			m[supported[wordSplitted]] = &empty
		}
		if _, ok := alreadyWas[word]; !ok {
			alreadyWas[word] = struct{}{}
			*m[supported[wordSplitted]] = append(*m[supported[wordSplitted]], word)
		}
	}
	for key, value := range m {
		sort.Strings(*value)
		if len(*value) == 1 {
			delete(m, key)
		}
	}
	return &m
}

func OrChannel(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		return channels[0]
	}
	if len(channels) == 1 {
		return channels[0]
	}

	done := make(chan interface{})
	go func() {
		defer close(done)
		select {
		case <-channels[0]:
		case <-channels[1]:
		case <-OrChannel(channels[2:]...):
		}
	}()
	return done
}

func practicalOrChannel() *time.Duration {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-OrChannel(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	ans := time.Since(start)
	return &ans
}

func main() {
	if val := practicalOrChannel(); val == nil {
		log.Fatal("Amount of input channels is zero")
	} else {
		fmt.Printf("Time duration is: %.6f seconds\n", val.Seconds())
	}
}
