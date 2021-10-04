package main

func IntToRoman(num int) string {
	values := []int{
		1000, 900, 500, 400,
		100, 90, 50, 40,
		10, 9, 5, 4, 1,
	}

	romans := []string{
		"M", "CM", "D", "CD",
		"C", "XC", "L", "XL",
		"X", "IX", "V", "IV",
		"I"}
	roman := ""
	i := 0
	for num > 0 {
		k := num / values[i]
		for j := 0; j < k; j++ {
			roman += romans[i]
			num -= values[i]
		}
		i++
	}
	return roman
}
