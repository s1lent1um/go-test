package main

import (
	"fmt"
	"io"
	"os"
	"encoding/xml"
	"strings"
)

// this is a comment

func readFile(filename string) string {
	//fmt.Println(filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	text := ""
	for {
		buf := make([]byte, 1024)
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		text += fmt.Sprintf("%s", buf)
	}
	return text
}

type Elem struct {
	Name string
	Text string
	Children []*Elem
}

func autoClose(value string) bool {
	value = strings.ToLower(value)
	for _, s := range xml.HTMLAutoClose {
		if strings.ToLower(s) == value {
			return true
		}
	}
	return false
}

func parse(html string) {
	var tree *Elem
	var cursor *Elem
	var stack []*Elem
	var token xml.Token
	decoder := xml.NewDecoder(strings.NewReader(html))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	for {
		if token == nil {
			token, _ = decoder.Token()
			if token == nil {
				break
			}
		}

		switch entity := token.(type) {
		case xml.StartElement:
			//fmt.Printf("<%s>\n", entity.Name.Local)
			cursor = new(Elem)
			cursor.Name = entity.Name.Local
			if (tree != nil) {
				tree.Children = append(tree.Children, cursor)
			}
			tree = cursor
			stack = append(stack, cursor)
			printStack(stack)
		case xml.EndElement:
			fmt.Printf("</%s>\n", entity.Name.Local)
		case xml.CharData:
			fmt.Printf("%s\n", entity)
		default:
			fmt.Printf("%#v\n", token)
		}
		token = nil
	}
	fmt.Println(tree)
}

func printStack (stack []*Elem) {
	for i, elem := range stack {
		if i > 0 {
			fmt.Printf(" > ")
		}
		fmt.Printf("%s", elem.Name)
	}
	fmt.Println()
}
func isClosing(token xml.Token, name string) (result bool) {
	defer func() {
		if err := recover(); err != nil {
			result = false
		}
	}()
	entity := token.(xml.EndElement)
	return strings.ToLower(name) == entity.Name.Local
}

func main() {
	if len(os.Args) <= 1 {
		os.Exit(2)
	}
	html := readFile(os.Args[1])
	//fmt.Println(html)
	parse(html)
}