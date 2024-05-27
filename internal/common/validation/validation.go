package validations

import (
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	LatitudePattern  string = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	LongitudePattern string = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
)

var ImgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_\-\:x\=\(\)]+)*(\.(jpg|jpeg|png))?$`, s)
	return match
}, "image url is not valid")

var LatitudeValidation = func(lat float64) bool {
	latitude := fmt.Sprintf("%f", lat)
	match, err := regexp.MatchString(LatitudePattern, latitude)
	if err != nil {
		return false
	}
	if !match {
		return false
	}
	return true
}

var LongitudeValidation = func(lat float64) bool {
	longitude := fmt.Sprintf("%f", lat)
	match, err := regexp.MatchString(LongitudePattern, longitude)
	if err != nil {
		return false
	}
	if !match {
		return false
	}
	return true
}
