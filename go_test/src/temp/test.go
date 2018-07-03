/*
演示使用syscall的3种方法调用dll文件函数
Lib.dll里面有一个add函数,是通过_declspec(dllexport)extern "C"和导出的
DllTestDef.dll里面有一个add函数,是通过.def文件导出的
c++code里面是2个库文件的源代码,基于c++6
*/

package main

import (
	"fmt"
	"syscall"
)

func IntPtr(n int) uintptr {
	return uintptr(n)
}

//func Lib_add(a, b int) {
//	lib := syscall.NewLazyDLL("lib.dll")
//	fmt.Println("dll:", lib.Name)
//	add := lib.NewProc("add")
//	fmt.Println("+++++++NewProc:", add, "+++++++")
//
//	ret, _, err := add.Call(IntPtr(a), IntPtr(b))
//	if err != nil {
//		fmt.Println("lib.dll运算结果为:", ret)
//	}
//
//}

func DllTestDef_add(a, b int) {
	DllTestDef, _ := syscall.LoadLibrary("DllTestDef.dll")
	fmt.Println("+++++++syscall.LoadLibrary:", DllTestDef, "+++++++")
	defer syscall.FreeLibrary(DllTestDef)
	add, err := syscall.GetProcAddress(DllTestDef, "add")
	fmt.Println("GetProcAddress", add)

	ret, _, err := syscall.Syscall(
		uintptr(add),
		2,
		uintptr(a),
		uintptr(b),
		0)
	if err != nil {
		fmt.Println("DllTestDef.dll运算结果为:", ret)
	}
}

func testDll() {
	lin, _ := syscall.LoadLibrary(`E:\vs2010\libui-alpha3.5\build\out\Debug\libui.dll`)
	fmt.Println("+++++++syscall.LoadLibrary:", lin, "+++++++")
}

func DllTestDef_add2(a, b int) {
	DllTestDef := syscall.MustLoadDLL("DllTestDef.dll")
	add := DllTestDef.MustFindProc("add")

	fmt.Println("+++++++MustFindProc：", add, "+++++++")
	ret, _, err := add.Call(IntPtr(a), IntPtr(b))
	if err != nil {
		fmt.Println("DllTestDef的运算结果为:", ret)
	}
}

func main() {
	//Lib_add(4, 5)
	testDll()
	DllTestDef_add(4, 5)
	DllTestDef_add2(4, 5)
}
