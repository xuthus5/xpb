package common

func getInt(a uint8) int {
	if a-'0' > 0 && a <= '9' {
		return int(a - '0')
	}

	if a-'a' > 0 && a <= 'z' {
		return int(a-'a') + 10
	}

	if a-'A' > 0 && a <= 'Z' {
		return int(a-'A') + 20
	}

	return int('@')
}

func Base62Encoder(str1 string, str2 string) string {
	List62 := []string{"f", "x", "H", "M", "u", "a", "R", "5", "i", "V", "G", "s", "T", "c", "2", "1", "y", "4", "t", "U", "n", "o", "q", "v", "F", "b", "J", "7", "j", "3", "K", "p", "P", "g", "N", "l", "Y", "6", "W", "r", "e", "S", "I", "z", "O", "L", "m", "d", "A", "h", "D", "w", "Q", "E", "Z", "0", "8", "C", "B", "9", "k", "X"}
	i := len(str1) - 1
	j := len(str2) - 1
	var sum string
	var tem int //进位
	for i >= 0 && j >= 0 {
		s := getInt(str1[i]) + getInt(str2[j]) + tem
		if s >= 62 {
			tem = 1
			sum = List62[s%62] + sum
		} else {
			tem = 0
			sum = List62[s] + sum
		}
		i--
		j--
	}
	for i >= 0 {
		s := getInt(str1[i]) + tem
		if s >= 62 {
			tem = 1
			sum = List62[s%62] + sum
		} else {
			tem = 0
			sum = List62[s] + sum
		}
		i--
	}
	for j >= 0 {
		s := getInt(str2[i]) + tem
		if s >= 62 {
			tem = 1
			sum = List62[s%62] + sum
		} else {
			tem = 0
			sum = List62[s] + sum
		}
		j--
	}
	if tem != 0 {
		sum = "1" + sum
	}
	return sum
}
