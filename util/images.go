package util

import (
	"github.com/davidbyttow/govips/v2/vips"
)

// Gets an image for the block from the data folder given the texture name, returns a libvips ImageRef pointer
func GetBlockImage(textureID string) (*vips.ImageRef, error) {

	//Tell vips to load the image
	image1, err := vips.NewImageFromFile("./blocks/images/" + textureID + ".png")

	//Return error if present
	if err != nil {
		return nil, err
	}

	//Return the ImageRef
	return image1, nil
}

// This function does all the heavy lifting, it generates an image for all the blocks in the list of dataentries provided
// The provided height is the pixel height of the image, the horizantal is scaled respectively
// This returns the image as JPEG bytes
func AppendBlockImages(data *[]DataEntry, height int, noText bool) (*[]byte, error) {

	//Create a big image that will be the size of the final image, to paste everything on to
	bigImage, err := vips.Black(height*len(*data), height)

	//Error handling
	if err != nil {
		return nil, err
	}

	//Add an alpha channel to support transparency and make sure that transparent blocks overlay properly
	bigImage.AddAlpha()

	//Make sure we are in the SRGB color space
	bigImage.ToColorSpace(vips.InterpretationSRGB)

	//Draw A Rectangle of background color, this does nothing if the blocks have no transparency, but it provides a backdrop for transparent textures
	bigImage.DrawRect(vips.ColorRGBA{R: 110, G: 101, B: 91, A: 255}, 0, 0, height*len(*data), height, true)

	//Loop through the entire slice of block data
	for idx, entry := range *data {

		//Grab the block image for the given block
		imageRef, err := GetBlockImage(entry.TextureName)

		//error handle
		if err != nil {
			return nil, err
		}

		//Make sure color spaces match throughout
		imageRef.ToColorSpace(vips.InterpretationSRGB)

		//Resize the image to make the height match the height supplied by the user
		//This uses nearest neighbor to retain sharpness
		imageRef.Resize(float64(height/imageRef.Height()), vips.KernelNearest)

		//Overlay this image on top of the background image, offset is based on its list index to make sure two images do not overlap
		bigImage.Composite(imageRef, vips.BlendModeOver, height*idx, 0)

		//Close the image to prevent funny memory leaks that take 6 hours to debug (definetly not from experience)
		imageRef.Close()

	}

	//Do all the text stuff if we even want text

	if !noText {

		//Create an image for the white bar at the bottom
		transp, err := vips.Black(height*len(*data), height/16)

		//Error handling
		if err != nil {
			return nil, err
		}

		//Add alpha channel to image
		transp.AddAlpha()

		//Convert to SRGB color space to match
		transp.ToColorSpace(vips.InterpretationSRGB)

		//Draw a white rectangle that is opaque to the image, we will overlay this later into the main image
		transp.DrawRect(vips.ColorRGBA{R: 255, G: 255, B: 255, A: 150}, 0, 0, height*len(*data), height/16, true)

		//Composite both the images to include the white bar
		bigImage.Composite(transp, vips.BlendModeAdd, 0, height-(height/16))

		//Close our temporary image
		transp.Close()

		//Flatten the image, provide bg (although irrelevent at this point), this allows us to add Labels without error later
		bigImage.Flatten(&vips.Color{R: 110, G: 101, B: 91})

		//Loop through data to add labe;s
		for idx, entry := range *data {

			//The label parameters for text
			lParams := vips.LabelParams{
				//Use the display name
				Text: entry.DisplayName,
				//A minecraft font, the docker image already installs this for you to make life easier
				Font: "Minecraftia",
				//Have some padding
				Width: vips.Scalar{Value: float64(height - (height / 16))},
				//Have some padding
				Height: vips.Scalar{Value: float64((height / 16) - (height / 64))},
				//Center a bit
				OffsetX: vips.Scalar{Value: float64((height * idx) + height/32)},
				//Center a bit
				OffsetY: vips.Scalar{Value: float64(height - (height / 16) + (height / 128))},
				//No opacity
				Opacity: 255,
				//Text should be black
				Color: vips.Color{R: 0, G: 0, B: 0},
				//Put text to center
				Alignment: vips.AlignCenter,
			}

			//Add the label to the image
			bigImage.Label(&lParams)
		}
	}

	//Export parameters to export our result as a JPEG
	ep := vips.NewDefaultJPEGExportParams()
	//Export the image as JPEG bytes
	image1bytes, _, _ := bigImage.Export(ep)

	//Close the image to prevent memory leak
	bigImage.Close()

	//Return the JPEG bytes
	return &image1bytes, nil
}
