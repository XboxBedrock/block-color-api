package util

import (
	"fmt"

	"os"

	"io"

	"encoding/json"
)

// This is the config file that lets us see what are valid versions for the `version` query parameter
type Config struct {
	//Port to run webserver on
	Port int32 `json:"port"`
	//Allowed Minecraft versions
	AllowedVersions []string `json:"allowedVersions"`
}

// This is used to represent every entry in the data.json file
type DataEntry struct {
	//Display name of the block
	DisplayName string `json:"display_name"`
	//Average Hex color of the block
	Hex string `json:"hex"`
	//Average color of the block in the LAB format
	Lab []float64 `json:"lab"`
	//Name of the texture, can be used to find /blocks/images/{texture_name}.png
	TextureName string `json:"texture_name"`
	//Is this a decoration or a block?
	IsDecoration bool `json:"is_decoration"`
	//Should this block be rendered in 3d or 2d
	Show3D bool `json:"show_3d"`
	//Minecraft versions that the block is present in, should be a subset of AllowedVersions
	Versions []string `json:"versions"`
}

// To represent the JSON file, we create a map with the keys being the JSON keys, however the keys are never really used
type DataItems map[string]DataEntry

// Given the file name, this generic function will load the JSON file into an interface and return an instance
func LoadJson[T interface{}](fileName string) T {

	//Open the file
	configFile, errorConf := os.Open(fileName)

	//Completely nuke the system if an error is detected
	if errorConf != nil {
		fmt.Println(errorConf)
		os.Exit(1)
	}

	fmt.Println("Loaded " + fileName)

	//Get the bytes from the file
	configBytes, _ := io.ReadAll(configFile)

	//create an empty variable to store the file
	var config T

	//Unmarshall the JSON into the config
	json.Unmarshal(configBytes, &config)

	return config
}
