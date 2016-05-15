package main

import (
	"./Lists"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

var numbuckets = 4
var numthreads = 8

//Taken from https://stackoverflow.com/questions/5884154/golang-read-text-file-into-string-array-and-write
// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

type section struct {
	startln int
	endln   int //Should stop at this line (not read it)
	done    bool
	lock    *sync.Mutex
}

func initSections(lines []string, numthreads int) []*section {
	var numsections = (int)((float64)(numthreads) * math.Log2((float64)(numthreads)))
	if numsections == 0 {
		numsections = 1
	}
	sections := make([]*section, numsections)
	startln := 0
	sectionln := numsections / len(lines)
	for i := 0; i < numsections-1; i++ {
		thissection := new(section)
		thissection.startln = startln
		thissection.endln = startln + sectionln
		thissection.done = false
		thissection.lock = &sync.Mutex{}
		startln = startln + sectionln
		sections[i] = thissection
	}
	//Last section
	thissection := new(section)
	thissection.startln = startln
	thissection.endln = len(lines)
	thissection.done = false
	thissection.lock = &sync.Mutex{}
	startln = startln + sectionln
	sections[len(sections)-1] = thissection

	return sections
}

func wcParallel(lines []string, hmap *Lists.HashMap, numthreads int) {
	sections := initSections(lines, numthreads)
	var wg sync.WaitGroup
	wg.Add(numthreads)
	for i := 0; i < numthreads; i++ {
		go countlines(hmap, &wg, sections, lines)
	}
	wg.Wait()
}

func countlines(hmap *Lists.HashMap, wg *sync.WaitGroup, sections []*section, lines []string) {
	validsec := true
	var chosensection *section
	for { //Until no valid sections left
		validsec = false
		for i := 0; i < len(sections); i++ {
			if sections[i].done == false {
				sections[i].lock.Lock()
				if sections[i].done == false {
					//Found valid section to work on
					sections[i].done = true
					validsec = true
					chosensection = sections[i]
					break
				} else {
					sections[i].lock.Unlock()
					continue
				}
			}
		}
		if validsec == false {
			break
		}
		dosection(hmap, chosensection, lines)
		chosensection.lock.Unlock()
	}
	wg.Done()
}

func dosection(hmap *Lists.HashMap, chosensection *section, lines []string) {
	for i := chosensection.startln; i < chosensection.endln; i++ {
		words := strings.Fields(lines[i])
		for _, word := range words {
			word = strings.ToLower(strings.Trim(word, ".,?!\\ {}()/"))
			//fmt.Println(word)
			val, there := hmap.Get(word)
			if there == false {
				zero := new(int32)
				*zero = 0
				val = zero
				hmap.Insert(word, val)
			}
			count := val.(*int32)
			//CAS to increment pointer
			for {
				currentval := *count
				if atomic.CompareAndSwapInt32(count, currentval, currentval+1) {
					break
				}
			}
		}
	}

}

func main() {
	lines, err := readLines("test.txt")
	if err != nil {
		fmt.Println("Error reading file!")
		os.Exit(1)
	}
	hMap := new(Lists.HashMap)
	hMap.Init(numbuckets, Lists.LLListType)

	wcParallel(lines, hMap, numthreads)
	hMap.PrintMap()
}
