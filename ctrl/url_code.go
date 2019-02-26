package ctrl

import (
	"math"
)

func GetNextString(s string) string {
	if s == "" {
		return "0"
	}
	dictionary := GetDictionary()

	char := int(s[len(s)-1])
	next := 0
	if char == 57 {
		next = 65
	} else {
		next = char + 1
	}

	if isInSlice(dictionary, next) {
		return s[:len(s)-1] + string(next)
	} else {
		ns := ""
		if s[:len(s)-1] != "" {
			ns = GetNextString(s[:len(s)-1])
		} else {
			ns = string(dictionary[0])
		}
		return ns + string(dictionary[0])
	}
}

func isInSlice(s []int, val int) bool  {
	m := make(map[int]bool)
	for i := 0; i < len(s); i++ {
		m[s[i]] = true
	}

	if _, ok := m[val]; ok {
		return true;
	} else {
		return false;
	}
}

func GetDictionary() []int {
	dictionary := []int{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89 , 90}
	return dictionary
}

func GetDictionaryMap() map[int]int  {

	dictionaryMap := make(map[int]int)
	dictionary := GetDictionary()
	for i, d := range dictionary {
		dictionaryMap[d] = i + 1
	}
	return dictionaryMap
}

func GetNumberFromCode(code string) float64 {
	code = Reverse(code)
	dictLength := float64(len(GetDictionary()))
	dictMap := GetDictionaryMap()

	value := 0.0
	for i, c := range code {
		value = value + math.Pow(dictLength, float64(i)) * float64(dictMap[int(c)]);
	}
	return value
}

func GetNthCode(n int) string {
	code := "0"
	count := 1

	for count < n {
		code = GetNextString(code)
		count += 1
	}
	return code
}