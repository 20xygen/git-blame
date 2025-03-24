package internal

import "fmt"

func TestDir() {
	path := "gitfame"

	d, err := getDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d.asList())

	fmt.Println()

	files := d.files()
	for _, fl := range files {
		fmt.Printf("%s [ %s ; %s ]\n", fl.fl.name, fl.ext, fl.lang)
	}

	fmt.Println()

	fl := d.kids["cmd"].(*dir).kids["gitfame"].(*dir).kids["main.go"].(*file)

	bo, err := parseBlame(fl.path())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(bo.String())
}
