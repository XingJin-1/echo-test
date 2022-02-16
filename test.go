package main

import "strings"

func Testmain() {
	/*
		arg := "--footer=test1="

		i_minus := strings.Index(arg, "=")
		name := arg[2:i_minus]
		value := arg[i_minus+1:]

		println(name, " ", value)
	*/
	//s1 := "true"
	//b1, _ := strconv.ParseBool(s1)
	//fmt.Printf("%T, %v\n", b1, b1)
	//s := "google.com,infineon.com"

	//l := strings.Split(s, ",")

	//println(l[0], ":", l[1])
	//s := "--footer=test1"
	//res := split(s, "=")
	//println(res[0], " ", res[1])
	//println(res)

	userID := "XingJin"

	token, err := CreateToken(userID)

	if err == nil {
		println("token: ", token)
	} else {
		println("error while generating tokens")
	}

}

func Split(s, sep string) []string {
	var result []string
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):]
		i = strings.Index(s, sep)
	}
	return append(result, s)
}
