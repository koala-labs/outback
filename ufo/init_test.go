package main

import (
	"testing"
)
//
//func cleanTempDir(t *testing.T, path string) {
//	if _, err := fs.Stat(path); os.IsNotExist(err) {
//		err := os.Remove(path)
//
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//}

// @todo it will probably be easier to test if broken down into smaller funcs
func TestItCreatesUFODirectoryIfItDoesNotExist(t *testing.T) {
	//temp := os.TempDir() + "/gotest"
	//defer cleanTempDir(t, temp)

}