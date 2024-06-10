package main

import (
	"me/xboxbedrock/minecraft/imageserver/paths"
	"strconv"

	"me/xboxbedrock/minecraft/imageserver/util"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gin-gonic/gin"
)

// This is some hacky GIN middleware that injects the config into the GIN request
func InjectConfig(conf *util.Config, data *[]util.DataEntry) gin.HandlerFunc {
	//middleware
	return func(c *gin.Context) {
		//Inject config
		c.Set("myConfig", conf)
		//Inject block data
		c.Set("dataArray", data)
		//Allow request to continue
		c.Next()
	}

}

func main() {

	//Startup VIPS for image processing
	vips.Startup(nil)
	//Murder VIPS when the program terminates
	defer vips.Shutdown()

	//Load the config file
	config := util.LoadJson[util.Config]("config.json")

	//Load the block data.json
	blockDataEarly := util.LoadJson[util.DataItems]("blocks/data.json")

	//Create a slice to store this block data without keys
	blockData := make([]util.DataEntry, len(blockDataEarly))

	//This loop adds all the values in the json file to the blockData array
	idx := 0
	for _, v := range blockDataEarly {
		blockData[idx] = v
		idx++
	}

	//Create a router with GIN
	router := gin.Default()
	//Add the middleware
	router.Use(InjectConfig(&config, &blockData))
	//Register the path
	router.GET("/getBlockImage", paths.GetBlockImage)

	//Bind to port specified in config
	router.Run("0.0.0.0:" + strconv.Itoa(int(config.Port)))

}
