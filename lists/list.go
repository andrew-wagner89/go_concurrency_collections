package main

type List interface {
	printlist()
	insert(int) bool
	contains(int) bool
	remove(int) bool
}