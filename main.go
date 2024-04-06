package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Profile struct {
	firstName string
	lastName  string
	username  string
	email     string
	address   string
}

func main() {

	asciiPoints := make(map[string]int)
	asciiPoints["asciiLowerCaseLetterStarts"] = 97
	asciiPoints["asciiLowerCaseLetterEnds"] = 122
	asciiPoints["asciiUpperCaseLetterStarts"] = 65
	asciiPoints["asciiUpperCaseLetterEnds"] = 90
	asciiPoints["asciiNumberStarts"] = 48
	asciiPoints["asciiNumberEnds"] = 57
	asciiPoints["asciiNumberToUpperCaseLetterOffset"] = 17

	fmt.Println(asciiPoints)

	const userNameFile = "data/common_templates/usernames.txt"
	const emailFormatsFile = "data/common_templates/emails.txt"
	const usAddressesFile = "data/by_region/US/addresses/addresses.txt"

	regions := [4]string{"US", "UK", "LK", "AUS"}

	// var fnameFile string

	seed := getInput()

	region := determineRegion(int(seed[0]), regions[:], asciiPoints)

	//gender letter is set to the SECOND letter
	genderDeterminingLetter := int(seed[1])

	//fname letter is the first non numerical letter and lname letter is the second non numerical letter
	firstNameDeterminingLetter, lastNameDeterminingLetter := findFirstTwoNonNumericalCharacters(seed, asciiPoints)

	// 0 => Male & 1 => Female
	gender := determineGender(genderDeterminingLetter)

	firstNameFile := getFirstNameFile(gender, firstNameDeterminingLetter, region)

	firstNameOffset, lastNameOffset, addressOffset, emailOffset, usernameOffset := getOffsets(seed, asciiPoints)

	firstName := getLineAtIndex(firstNameOffset, firstNameFile)

	lastNameFile := "data/by_region/" + region + "/names/last_names/lname_" + strings.ToUpper(lastNameDeterminingLetter) + ".txt"

	lastName := getLineAtIndex(lastNameOffset, lastNameFile)

	email := getFormattedString(emailOffset, firstName, lastName, seed, emailFormatsFile)

	addressFile := "data/by_region/" + region + "/addresses/addresses.txt"

	address := getLineAtIndex(addressOffset, addressFile)

	username := getFormattedString(usernameOffset, firstName, lastName, seed, userNameFile)

	generatedProfile := Profile{firstName: firstName, lastName: lastName, username: username, email: email, address: address}

	UpdateSVG(seed, firstName)

	generatedProfile.printDetails()
	fmt.Println(region)

}

func (profile Profile) printDetails() {

	log.Println(profile.firstName)
	log.Println(profile.lastName)
	log.Println(profile.email)
	log.Println(profile.address)
	log.Println(profile.username)

}

func determineRegion(letter int, regions []string, asciiPoints map[string]int) string {

	asciiNumbersChunk := 10 / len(regions)
	asciiLettersChunk := 26 / len(regions)

	n := asciiPoints["asciiNumberStarts"]
	l := asciiPoints["asciiLowerCaseLetterStarts"]
	u := asciiPoints["asciiUpperCaseLetterStarts"]

	for i := 0; i < len(regions); i++ {

		if i == len(regions)-1 {
			if letter >= n && letter <= asciiPoints["asciiNumberEnds"] || letter >= l && letter <= asciiPoints["asciiLowerCaseLetterEnds"] || letter >= u && letter <= asciiPoints["asciiUpperCaseLetterEnds"] {
				return regions[i]
			}
		} else {
			if letter >= n && letter < n+asciiNumbersChunk || letter >= l && letter < l+asciiLettersChunk || letter >= u && letter < u+asciiLettersChunk {

				return regions[i]
			}
		}

		n += asciiNumbersChunk
		l += asciiLettersChunk
		u += asciiLettersChunk

	}
	return ""
}

func getOffsets(seed string, asciiPoints map[string]int) (int, int, int, int, int) {

	// first name offset is calculated by THIRD letter * FOURTH letter
	firstNameOffset := findIndexFromLetter(seed[2], asciiPoints) * findIndexFromLetter(seed[3], asciiPoints)

	// last name offset is calculated by FOUTH letter * FIFTH (LAST) letter
	lastNameOffset := findIndexFromLetter(seed[3], asciiPoints) * findIndexFromLetter(seed[4], asciiPoints)

	// address offset is calculated by SECOND letter * FOURTH letter
	addressOffset := findIndexFromLetter(seed[1], asciiPoints) * findIndexFromLetter(seed[3], asciiPoints)

	// email offset set to the ascii of FORTH letter
	emailOffset := int(seed[3])

	// username offset set to the ascii of FIFTH (LAST) letter
	usernameOffset := int(seed[4])

	return firstNameOffset, lastNameOffset, addressOffset, emailOffset, usernameOffset

}

func getFirstNameFile(gender int, firstNameDeterminingLetter string, region string) string {

	var firstNameFile string

	if gender == 1 {
		log.Println(strings.ToUpper(firstNameDeterminingLetter) + ".txt" + "  Female")
		firstNameFile = "data/by_region/" + region + "/names/first_names/female_" + strings.ToUpper(firstNameDeterminingLetter) + ".txt"
		log.Println(firstNameFile)

	} else {
		log.Println(strings.ToUpper(firstNameDeterminingLetter) + ".txt" + "  Male")
		firstNameFile = "data/by_region/" + region + "/names/first_names/male_" + strings.ToUpper(firstNameDeterminingLetter) + ".txt"
		log.Println(firstNameFile)

	}

	return firstNameFile

}

func getInput() string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter seed:")

	read, err_read := reader.ReadString('\n')

	if err_read != nil {
		log.Printf("could not open the file: %v", err_read)
	}

	read = strings.TrimSpace(read)

	inputError := validateInput(read)

	if inputError != nil {
		log.Println(inputError)
		getInput()

	} else {
		return read
	}

	return ""

}

func validateInput(str string) error {
	for _, char := range str {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("only enter numbers or english letters")
		}

	}

	if len(str) > 5 {
		return errors.New("only enter 5 characters")
	}

	return nil
}

func determineGender(genderDeterminingLetter int) int {

	if genderDeterminingLetter >= 78 && genderDeterminingLetter <= 90 || genderDeterminingLetter >= 110 && genderDeterminingLetter <= 122 || genderDeterminingLetter >= 52 && genderDeterminingLetter <= 57 {

		return 0

	} else {

		return 1
	}

}

func findFirstTwoNonNumericalCharacters(seed string, asciiPoint map[string]int) (string, string) {
	var letters []int
	var numbers []int

	for _, l := range seed[2:] {
		kl := int(l)

		if kl >= asciiPoint["asciiLowerCaseLetterStarts"] && kl <= asciiPoint["asciiLowerCaseLetterEnds"] || kl >= asciiPoint["asciiUpperCaseLetterStarts"] && kl <= asciiPoint["asciiUpperCaseLetterEnds"] {

			letters = append(letters, kl)
		} else if kl >= asciiPoint["asciiNumberStarts"] && kl <= asciiPoint["asciiNumberEnds"] {
			numbers = append(numbers, kl)

		}

	}

	if len(letters) >= 2 {

		return string(letters[0]), string(letters[1])

	} else {

		return string(numbers[0] + asciiPoint["asciiNumberToUpperCaseLetterOffset"]), string(numbers[1] + asciiPoint["asciiNumberToUpperCaseLetterOffset"])

	}

}

func findIndexFromLetter(letter byte, asciiPoint map[string]int) int {

	l := int(letter)

	if l >= asciiPoint["asciiLowerCaseLetterStarts"] && l <= asciiPoint["asciiLowerCaseLetterEnds"] {
		return l - 96
	} else if l >= asciiPoint["asciiUpperCaseLetterStarts"] && l <= asciiPoint["asciiUpperCaseLetterEnds"] {
		return l - 38
	} else {
		return l + 5
	}

}

func getLineAtIndex(index int, fileToRead string) string {

	linesread := 0
	var line string

	for {

		if linesread >= index-1 {

			break
		}

		file, err_file := os.Open(fileToRead)

		if err_file != nil {
			log.Printf("could not open the file: %v", err_file)
		}

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {

			line = scanner.Text()

			if linesread >= index-1 {

				break
			}

			linesread++

		}

		file.Close()

	}

	return line
}

func getFormattedString(index int, fname string, lname string, seed string, fileToRead string) string {

	var formattedString string

	formattedString = getLineAtIndex(index, fileToRead)

	formattedString = strings.Replace(formattedString, "<fname>", strings.ToLower(fname), -1)
	formattedString = strings.Replace(formattedString, "<lname>", strings.ToLower(lname), -1)
	formattedString = strings.Replace(formattedString, "<int1>", fmt.Sprintf("%d", seed[0]), -1)
	formattedString = strings.Replace(formattedString, "<int2>", fmt.Sprintf("%d", seed[1]), -1)
	formattedString = strings.Replace(formattedString, "<int3>", fmt.Sprintf("%d", seed[2]), -1)
	formattedString = strings.Replace(formattedString, "<int4>", fmt.Sprintf("%d", seed[3]), -1)
	formattedString = strings.Replace(formattedString, "<int5>", fmt.Sprintf("%d", seed[4]), -1)

	return formattedString

}
