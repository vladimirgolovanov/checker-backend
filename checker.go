package main

type Checker interface {
	Check(name string) bool
	GetId() int
}
