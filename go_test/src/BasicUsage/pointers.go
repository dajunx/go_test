package basicuse

import "fmt"

//TestPointers 指针的使用
func TestPointers() {
	i, j := 42, 2701

	//Go 没有指针运算
	p := &i         // point to i
	fmt.Println(*p) // read i through the pointer
	*p = 21         // set i through the pointer
	fmt.Println(i)  // see the new value of i

	p = &j         // point to j
	*p = *p / 37   // divide j through the pointer
	fmt.Println(j) // see the new value of j
}
