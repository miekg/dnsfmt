package main

import (
	"fmt"
	"strconv"

	"github.com/miekg/dnsfmt/zonefile"
)

const (
	Second = 1
	Minute = Second * 60
	Hour   = Minute * 60
	Day    = Hour * 24
	Week   = Day * 7
)

func TimeToHumanByte(ttl []byte) []byte {
	i, err := strconv.ParseUint(string(ttl), 10, 64)
	if err != nil {
		j, ok := zonefile.StringToTTL(string(ttl))
		if !ok {
			return ttl
		}
		k := int(j)
		return []byte(TimeToHuman(&k))
	}
	j := int(i)
	return []byte(TimeToHuman(&j))
}

func TimeToHuman(ttl *int) string {
	// not for these smaller ones
	if *ttl <= 600 {
		return fmt.Sprintf("%d", *ttl)
	}

	// round to nearest minute?
	t := *ttl
	week := t / Week
	t -= week * Week

	day := t / Day
	t -= day * Day

	hour := t / Hour
	t -= hour * Hour

	min := t / Minute
	t -= min * Minute

	//	sec := t / Second
	//	t -= sec * Second

	s := ""
	if week > 0 {
		s += fmt.Sprintf("%dW", week)
	}
	if day > 0 {
		s += fmt.Sprintf("%dD", day)
	}
	if hour > 0 {
		s += fmt.Sprintf("%dH", hour)
	}
	if min > 0 {
		s += fmt.Sprintf("%dM", min)
	}
	// discard sec
	return s
}
