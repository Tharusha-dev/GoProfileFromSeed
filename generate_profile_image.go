package main

import (
	"fmt"
	"os"
	"strings"
)

func UpdateSVG(seed string, firstName string) string {

	const svgHeaders = `<svg width="250" height="250" xmlns="http://www.w3.org/2000/svg"><rect width="250" height="250" fill="white"/><g fill="white">`
	const svgFooter = `</g></svg>`
	const startX = 50
	const startY = 50

	var builder strings.Builder
	builder.WriteString(svgHeaders)

	// Open the file for writing
	file, errs := os.Create(fmt.Sprintf(`%s.svg`, firstName))
	if errs != nil {
		fmt.Println("Failed to create file:", errs)

	}
	defer file.Close()

	x, y := startX, startY

	for i := 0; i < 5; i++ {

		builder.WriteString(fmt.Sprintf(`<rect x="%s" y="%s" width="50" height="50" fill="%s"/>`, fmt.Sprint(x), fmt.Sprint(y), "#FF0000"))

		x_shift, y_shift := determineShift(int(seed[i]))

		if x+x_shift < 0 {
			x += 100
		} else {
			x += x_shift
		}

		if y+y_shift < 0 {
			y += 100

		} else {
			y += y_shift

		}

	}
	builder.WriteString(svgFooter)

	_, errs = file.WriteString(builder.String())
	if errs != nil {
		fmt.Println("Failed to write to file:", errs)

	}
	fmt.Println("generated profile image")
	return "line"
}

func determineShift(determiningLetter int) (int, int) {

	if determiningLetter >= 97 && determiningLetter <= 105 {
		return 0, 50 //up

	} else if determiningLetter >= 106 && determiningLetter <= 114 {
		return 50, 50 //right-up

	} else if determiningLetter >= 115 && determiningLetter <= 122 {
		return 50, 0 //right

	} else if determiningLetter >= 65 && determiningLetter <= 73 {
		return 50, -50 //right-down

	} else if determiningLetter >= 74 && determiningLetter <= 82 {
		return 0, -50 //down

	} else if determiningLetter >= 83 && determiningLetter <= 90 {
		return -50, -50 //left-down

	} else if determiningLetter >= 48 && determiningLetter <= 52 {
		return -50, 0 //left

	} else {
		return -50, 50 //left-up

	}

}
