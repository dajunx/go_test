package basicuse

import (
	"math"
	"fmt"
)

//接口类型是由一组方法定义的集合，接口类型的值可以存放实现这些方法的任何值
type Abser interface {
	Abs() float64
}

//TestUseInterface 接口的使用
func TestUseInterface() {
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f  // a MyFloat 实现了 Abser
	a.Abs() // 分别调用对应类型绑定的函数
	a = &v // a *Vertex 实现了 Abser
	a.Abs()

	// 下面一行，v 是一个 Vertex（而不是 *Vertex）
	// 所以没有实现 Abser。
	//a = v
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 匿名变量使用
type Human struct {
	name string
	age int
}
type student struct {
	hum Human
	pwd string
}

func TestAnonymousVariableUsage() {
	var nullInterface interface{}
	var i int = 5
	var str string
	str = "Hello world"

	Jim := student{Human{"Jim", 27}, "101"}

	nullInterface = i
	fmt.Println(nullInterface)
	nullInterface = str
	fmt.Println(nullInterface)
	nullInterface = Jim
	fmt.Println(nullInterface)
}
