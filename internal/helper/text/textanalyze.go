package TEXT_HELPER

import (
	"regexp"
)

type TextType int
type URLType int

const (
	TextNull TextType = iota
	TextEvent
	TextURL
)

const (
	URLNull URLType = iota
	URLYoutube
	URLYandex
)

var matchTime = regexp.MustCompile(`\d{1,2}:\d{2}`)
var matchDay = regexp.MustCompile(`(?i)(?:[^а-яё]|^)(понедельник|вторник|сред[ау]|четверг|пятниц[ау]|суббот[ау]|воскресенье|пн|вт|ср|чт|пт|сб|вс)(?:[^а-яё]|$)`)
var matchURLYoutybe = regexp.MustCompile(`\bhttps?://music\.youtube\.com[-a-zA-Z0-9()@:%_\+.~#?&//=]*`)
var matchURLYandex = regexp.MustCompile(`\bhttps?://music\.yandex\.(?:ru|by|kz|uz|com)[-a-zA-Z0-9()@:%_\+.~#?&//=]*`)

func AnalyzeTextType(text string) TextType {
	result := TextNull
	if isEvent(text) {
		result = TextEvent
	} else if isURL(text) {
		result = TextURL
	}
	return result
}

func URLAndType(text string) ([]string, URLType) {
	urlType := GetURLType(text)
	var url []string
	switch urlType {
	case URLYoutube:
		url = matchURLYoutybe.FindStringSubmatch(text)
	case URLYandex:
		url = matchURLYandex.FindStringSubmatch(text)
	}
	return url, urlType
}

func GetURLType(text string) URLType {
	if matchURLYoutybe.MatchString(text) {
		return URLYoutube
	} else if matchURLYandex.MatchString(text) {
		return URLYandex
	}
	return URLNull
}

func isEvent(text string) bool {
	return matchDay.MatchString(text) || matchTime.MatchString(text)
}

func isURL(text string) bool {
	return matchURLYoutybe.MatchString(text) || matchURLYandex.MatchString(text)
}
