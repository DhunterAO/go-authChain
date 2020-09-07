package fileutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func LoadJson(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func DumpJson(filename string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(data))
	err = ioutil.WriteFile(filename, data, 0755)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}