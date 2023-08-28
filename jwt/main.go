package main

import (
	"fmt"
	"jwt/util"
)

func main() {

	str, _ := util.GenerateToken(1, "zhangsan")
	token, _ := util.ParseToken(str)

	fmt.Println(fmt.Sprintf("GenerateToken: %v", str))
	fmt.Println(fmt.Sprintf("ParseToken: %v", token))
}
