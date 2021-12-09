package util

//Checks palindromity by comparing letters to the letter with same pos starting from end
func CheckPalindrome(text string) bool {
	runes := []rune(text)
	length := len(runes)
	for i := range runes {
		if i < length/2 && runes[i] != runes[length-1-i] {
			return false
		}
	}
	return true
}
