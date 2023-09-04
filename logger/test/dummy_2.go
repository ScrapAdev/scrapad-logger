package logger

import "errors"

func MethA() error {
	println("A")
	return MethB()
}

func MethB() error {
	println("B")
	println(MethC())
	return errors.New("ERROR GRAVE")
}

func MethC() string {
	return "C"
}

func Meth() error {
	return MethA()
}
