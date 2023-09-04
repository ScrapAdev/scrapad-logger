package logger

func MethodA() {
	println("A")
	MethodB()
}

func MethodB() {
	println("B")
	println(MethodC())
}

func MethodC() string {
	return "C"
}

func Method() {
	MethodA()
}
