package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type treeNode struct {
	fileInfo  os.FileInfo
	childrens []treeNode
}

func (n treeNode) Name() string {
	if n.fileInfo.IsDir() {
		return n.fileInfo.Name()
	} else {
		return fmt.Sprintf("%s (%s)", n.fileInfo.Name(), n.Size())
	}
}

func (n treeNode) Size() string {
	if n.fileInfo.Size() > 0 {
		return fmt.Sprintf("%db", n.fileInfo.Size())
	} else {
		return "empty"
	}
}

func getTreeNodes(path string, withFiles bool) ([]treeNode, error) {
	files, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, err
	}

	var nodes []treeNode

	for _, file := range files {
		if !file.IsDir() && !withFiles {
			continue
		}

		node := treeNode{
			fileInfo: file,
		}

		if node.fileInfo.IsDir() {
			childrens, err := getTreeNodes(path+string(os.PathSeparator)+file.Name(), withFiles)

			if err != nil {
				return nil, err
			}

			node.childrens = childrens
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func printTree(out io.Writer, nodes []treeNode, parentPrefix string) {
	var prefix string = "├───"
	var childPrefix string = "│\t"

	for i, node := range nodes {
		if i == len(nodes)-1 {
			prefix = "└───"
			childPrefix = "\t"
		}

		fmt.Fprint(out, parentPrefix, prefix, node.Name(), "\n")

		if node.fileInfo.IsDir() {
			printTree(out, node.childrens, parentPrefix+childPrefix)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	nodes, err := getTreeNodes(path, printFiles)
	if err != nil {
		return err
	}

	printTree(out, nodes, "")

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
