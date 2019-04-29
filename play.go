package main

import (
	"fmt"
	"strings"
)

func main(){
	my_array:=[]string{"a","b","v"}
	stringSlices := strings.Join(my_array, "\n")
	fmt.Print(stringSlices)
}

