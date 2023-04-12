package sti2023

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
)

const HASHSIZE = sha256.Size

func ReadJsonFile(dir, name string, data any) bool {
	h := Hash(name)
	filename := fmt.Sprintf("%x.json", h)
	byteValue, err := os.ReadFile(dir + "/" + filename)
	if os.IsNotExist(err) {
		return false
	}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func WriteJsonFile(dir, name string, data any) bool {
	h := Hash(name)
	filename := fmt.Sprintf("%x.json", h)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			fmt.Printf("error %s", err)
			return false
		}
	}
	byteValue, _ := json.Marshal(&data)
	if err := os.WriteFile(dir+"/"+filename, byteValue, 0644); err != nil {
		fmt.Printf("error %s", err)
		return false
	}
	return true
}

func Hash(signature string) [HASHSIZE]byte {
	return sha256.Sum256([]byte(signature))
}
