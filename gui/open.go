package gui

import (
	"fmt"

	"github.com/mechiko/utility"
)

// открываем в браузере ссылку
func Open(url string) error {
	if url == "" {
		return fmt.Errorf("пустой url")
	}
	if err := utility.OpenHttpLinkInShell(url); err != nil {
		return err
	}
	return nil
}

// открываем в эксплорере текущую папку программы
func OpenDir(dir string) (err error) {
	if utility.PathOrFileExists(dir) {
		if err := utility.OpenFileInShell(dir); err != nil {
			return err
		}
	}
	return nil
}

// открываем в эксплорере файл по имени
func OpenFile(file string) (err error) {
	if utility.PathOrFileExists(file) {
		if err := utility.OpenFileInShell(file); err != nil {
			return err
		}
	}
	return nil
}
