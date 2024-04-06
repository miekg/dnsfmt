package main

import "fmt"

/*

   #s = seconds = # x 1 seconds (really!)
   #m = minutes = # x 60 seconds
   #h = hours = # x 3600 seconds
   #d = day = # x 86400 seconds
   #w = week = # x 604800 seconds
*/

const (
	Second = 1
	Minute = Second * 60
	Hour   = Minute * 60
	Day    = Hour * 24
	Week   = Day * 7
)

func TimeToHuman(ttl *int) string {
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

	sec := t / Second
	t -= sec * Second

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
	if sec > 0 {
		s += fmt.Sprintf("%dS", sec)
	}
	return s
}
