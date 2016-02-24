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
	Name     string
	Text     string
	Children []*Elem
	isTextNode bool
}

func parse(html string) *Elem {
	var tree *Elem
	var cursor *Elem
	var stack []*Elem
	decoder := xml.NewDecoder(strings.NewReader(html))
	decoder.Strict = false
	decoder.AutoClose = xml.HTMLAutoClose
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}

		switch entity := token.(type) {
		case xml.StartElement:
			cursor = new(Elem)
			cursor.Name = entity.Name.Local
			if (tree != nil) {
				tree.Children = append(tree.Children, cursor)
			}
			tree = cursor
			stack = append(stack, cursor)
		case xml.EndElement:
			if tree.Name != entity.Name.Local {
				panic("Tag names mismatch: got " + entity.Name.Local + ", expected " + tree.Name)
			}
			stack = stack[:len(stack) - 1]
			if len(stack) > 0 {
				tree = stack[len(stack) - 1]
			}
		case xml.CharData:
			if (tree != nil) {
				cursor = new(Elem)
				cursor.isTextNode = true
				cursor.Text = fmt.Sprintf("%s", entity)
				tree.Children = append(tree.Children, cursor)
			}
		default:
			fmt.Printf("%#v\n", token)
		}
	}
	return tree
}

func renderMarkdown(tree *Elem) string {

	return renderRecursive(tree, false)
}

func renderRecursive(tree *Elem, inBody bool) string {

	if tree.isTextNode {
		if inBody {
			return tree.Text
		}
		return ""
	}


	template := "%s"
	switch tree.Name {
	case "body":
		inBody = true
	case "p":
		template = "%s\n\n";
	case "b":
		template = "*%s*";
	case "strong":
		template = "*%s*";
	case "h1":
		//content = " #
	}

	content := ""
	for _, elem := range tree.Children {
		content += renderRecursive(elem, inBody)
	}

	return fmt.Sprintf(template, content)
}

func printStack(stack []*Elem) {
	for i, elem := range stack {
		if i > 0 {
			fmt.Printf(" > ")
		}
		fmt.Printf("%s", elem.Name)
	}
	fmt.Println()
}

func main() {
	if len(os.Args) <= 1 {
		os.Exit(2)
	}
	html := readFile(os.Args[1])
	//fmt.Println(html)
	tree := parse(html)
	fmt.Printf("%s\n", renderMarkdown(tree))
}