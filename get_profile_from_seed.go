// GoProfileFromSeed generates user profiles based on a seed.
// A Profile is a structure with:
//   - First Name
//   - Last Name
//     -Username
//     -Email
//     -Region
//     -Address
//     -seed
package GoProfileFromSeed

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"unicode"
)

type Profile struct {
	FirstName string
	LastName  string
	Username  string
	Email     string
	Region    string
	Address   string
	Seed      string
}

// A seed is a 5 character string. Each character can a number or a lower case or upper case English letter. This function returns a Profile struct.
// For more info https://github.com/Tharusha-dev/GoProfileFromSeed?tab=readme-ov-file#how-it-works-%EF%B8%8F
func GetProfileFromSeed(seed string) Profile {

	asciiPoints := make(map[string]int)
	asciiPoints["asciiLowerCaseLetterStarts"] = 97
	asciiPoints["asciiLowerCaseLetterEnds"] = 122
	asciiPoints["asciiUpperCaseLetterStarts"] = 65
	asciiPoints["asciiUpperCaseLetterEnds"] = 90
	asciiPoints["asciiNumberStarts"] = 48
	asciiPoints["asciiNumberEnds"] = 57
	asciiPoints["asciiNumberToUpperCaseLetterOffset"] = 17

	regions := [4]string{"US", "UK", "LK", "AUS"}

	// var fnameFile string

	seed = getInput(seed)

	dataFileLocation, err := getDataFileLocation()

	if err != nil {
		fmt.Println(err)
	}

	userNameFile := dataFileLocation + "/common_templates/usernames.txt"
	emailFormatsFile := dataFileLocation + "/common_templates/emails.txt"

	region := determineRegion(int(seed[0]), regions[:], asciiPoints)

	//gender letter is set to the SECOND letter
	genderDeterminingLetter := int(seed[1])

	//fname letter is the first non numerical letter and lname letter is the second non numerical letter
	firstNameDeterminingLetter, lastNameDeterminingLetter := findFirstTwoNonNumericalCharacters(seed, asciiPoints)

	// 0 => Male & 1 => Female
	gender := determineGender(genderDeterminingLetter)

	firstNameFile := getFirstNameFile(gender, firstNameDeterminingLetter, region, dataFileLocation)

	firstNameOffset, lastNameOffset, addressOffset, emailOffset, usernameOffset := getOffsets(seed, asciiPoints)

	firstName := getLineAtIndex(firstNameOffset, firstNameFile)

	lastNameFile := dataFileLocation + "/by_region/" + region + "/names/last_names/lname_" + strings.ToUpper(lastNameDeterminingLetter) + ".txt"

	lastName := getLineAtIndex(lastNameOffset, lastNameFile)

	email := getFormattedString(emailOffset, firstName, lastName, seed, emailFormatsFile)

	addressFile := dataFileLocation + "/by_region/" + region + "/addresses/addresses.txt"

	address := getLineAtIndex(addressOffset, addressFile)

	username := getFormattedString(usernameOffset, firstName, lastName, seed, userNameFile)

	generatedProfile := Profile{FirstName: firstName, LastName: lastName, Username: username, Email: email, Region: region, Address: address, Seed: seed}

	return generatedProfile

}

func getDataFileLocation() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("Can't locate data file")
	}

	return strings.Replace(filename, "get_profile_from_seed.go", "data", 1), nil

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

func getFirstNameFile(gender int, firstNameDeterminingLetter string, region string, dataFileLocation string) string {

	var firstNameFile string

	if gender == 1 {

		firstNameFile = dataFileLocation + "/by_region/" + region + "/names/first_names/female_" + strings.ToUpper(firstNameDeterminingLetter) + ".txt"

	} else {
		firstNameFile = dataFileLocation + "/by_region/" + region + "/names/first_names/male_" + strings.ToUpper(firstNameDeterminingLetter) + ".txt"

	}

	return firstNameFile

}

func getInput(seed string) string {

	seed = strings.TrimSpace(seed)

	inputError := validateInput(seed)

	if inputError != nil {
		log.Println(inputError)
		getInput(seed)

	} else {
		return seed
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
