package paths

import (
	"slices"
	"sort"

	"github.com/gin-gonic/gin"

	"me/xboxbedrock/minecraft/imageserver/util"

	"strings"

	"strconv"
)

// This the GIN Route to get the block image
func GetBlockImage(c *gin.Context) {

	//Pull the config injected by the middleware
	configUntyped, _ := c.Get("myConfig")

	//Cast the type of the config
	config := *configUntyped.(*util.Config)

	//Pull blocks data injected by middleware
	dataArrayUntyped, _ := c.Get("dataArray")

	//Cast the type of the data array
	dataArray := *dataArrayUntyped.(*[]util.DataEntry)

	//Pull the version query parameter
	version := c.Query("version")

	//Pull the comma seperated RGB query parameter
	rgbString := c.Query("rgb")

	//Get the count as a string from query
	countString := c.Query("count")

	//Get the noText parameter as a string
	noTextString := c.Query("noText")

	//Get the height parameter as a string
	heightString := c.Query("height")

	//Get the page parameter as a string
	pageString := c.Query("page")

	//If the version is empty, that means the parameter is empty, tell them input malformed
	if version == "" {
		c.JSON(400, gin.H{
			"message": "Missing Parameter: version",
		})

		return
	}

	//If version is outside of allowedVersions, input malformed
	if !slices.Contains(config.AllowedVersions, version) {
		c.JSON(400, gin.H{
			"message": "Invalid Parameter: version",
		})

		return
	}

	//If RGB string not provided, input malformed
	if rgbString == "" {
		c.JSON(400, gin.H{
			"message": "Missing Parameter: rgb",
		})

		return
	}

	//Split the RGB string by comma
	rgbSlice := strings.Split(rgbString, ",")

	//If the RGB string length is not three, the input is malformed and the user is told so
	if len(rgbSlice) != 3 {
		c.JSON(400, gin.H{
			"message": "Invalid Format: rgb must be an integer comma seperated list of 3 8 bit integers",
		})

		return
	}

	//Create a slice for RGB values, uint8 is all we need
	rgb := make([]uint8, 3)

	//Loop through the values to validate
	for i := 0; i < 3; i++ {
		colTemp, err := strconv.Atoi(rgbSlice[i])

		//This error will call if the provided values are not integers
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid Format: rgb must be an integer comma seperated list of 3 8 bit integers",
			})

			return
		}

		//Throw error if outside valid range for RGB color
		if colTemp < 0 || colTemp > 255 {
			c.JSON(400, gin.H{
				"message": "Invalid Format: rgb must be an integer comma seperated list of 3 8 bit integers",
			})

			return
		}

		//Put into out uint8 slice
		rgb[i] = uint8(colTemp)
	}

	//Check if count string is present
	if countString == "" {
		c.JSON(400, gin.H{
			"message": "Missing Parameter: count",
		})

		return
	}

	//Get count as an integer
	count, countError := strconv.Atoi(countString)

	//Check if count is a valid integer
	if countError != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Format: count must be a valid integer in range 1-9",
		})

		return
	}

	//Check if count in range 1-9
	if count < 1 || count > 9 {
		c.JSON(400, gin.H{
			"message": "Invalid Format: count must be a valid integer in range 1-9",
		})

		return
	}

	//Assume we want text by default
	noText := false

	//Check if noText param is given, if it is, then we give no text
	if noTextString != "" {
		noText = true
	}

	//Check if the height parameter is present
	if heightString == "" {
		c.JSON(400, gin.H{
			"message": "Missing Parameter: height",
		})

		return
	}

	//Get height as an integer
	height, heightError := strconv.Atoi(heightString)

	//Check if count is a valid integer
	if heightError != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Format: height must be a valid integer in the range 16-4096 that is a multiple of 16",
		})

		return
	}

	//Check if height is outside of range
	if height < 16 || height > 4096 {
		c.JSON(400, gin.H{
			"message": "Invalid Format: height must be a valid integer in the range 16-4096 that is a multiple of 16",
		})

		return
	}

	//Check if height is a multiple of 16
	if height%16 != 0 {
		c.JSON(400, gin.H{
			"message": "Invalid Format: height must be a valid integer in the range 16-4096 that is a multiple of 16",
		})

		return
	}

	if pageString == "" {
		c.JSON(400, gin.H{
			"message": "Missing Parameter: page",
		})

		return
	}

	//Get page as an integer
	page, pageError := strconv.Atoi(pageString)

	//Check if page is a valid integer
	if pageError != nil {
		c.JSON(400, gin.H{
			"message": "Invalid Format: page must be a valid integer in the range 1-100",
		})

		return
	}

	//Check if page is outside of range
	if page < 1 || page > 100 {
		c.JSON(400, gin.H{
			"message": "Invalid Format: page must be a valid integer in the range 1-100",
		})

		return
	}

	//Convert the RGB color into the CIELAB Colorspace for more accurate color comparision
	l, a, b := util.RgbToLab(rgb[0], rgb[1], rgb[2])

	//Create a slice for the filtered data
	filtered := []util.DataEntry{}

	//Loop through slice to filter blocks by version
	for _, value := range dataArray {

		//Remove all blocks intended to be shown in 3d, and remove all decorations, we are not capable of rendering them
		//Make sure the block is present in the version we want, validated by the data file
		if (!value.Show3D && !value.IsDecoration) && slices.Contains(value.Versions, version) {
			//Append to filtered if valid
			filtered = append(filtered, value)
		}
	}

	//Using sort.SliceStable, we will sort the blocks by how close they are to our target color
	sort.SliceStable(filtered, func(i, j int) bool {

		//Convert the first block into LAB Colorspace
		l1, a1, b1 := util.HexToLAB(filtered[i].Hex)

		//Get the distance between target and block 1
		deltaI := util.DeltaE(l1, a1, b1, l, a, b)

		//Convert second block to LAB Colorspace
		l2, a2, b2 := util.HexToLAB(filtered[j].Hex)

		//Get the distance between target and block 2
		deltaJ := util.DeltaE(l2, a2, b2, l, a, b)

		//Return stating wheter J is less than I
		return deltaJ > deltaI
	})

	//Get the n most accurate blocks, as specified by count
	threeSlice := filtered[(count * (page - 1)):(count * (page))]

	//Gets the block images as a JPEG byte array
	byteArr, err := util.AppendBlockImages(&threeSlice, height, noText)

	//Check if anything stupid happend
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server Error: Image processing failed",
		})
	}

	//Return data to requester
	c.Data(200, "image/jpeg", *byteArr)
}
