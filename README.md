# Block Image Server

This is a simple microservice that returns block images for the BuildThEarth project, provided a color

It uses the CIELAB color space to see which blocks closely match the color, and returns a well-styled JPEG describing them

# API

There is only one route

```
/getBlockImage?version={version}&rgb={rgb}&count={count}&height={height}&noText={noText}&page={page}
```

`version` - One of 1.12 or 1.20, describes the minecraft version for the blocks in the response

`rgb` - RGB value to match blocks to, no-spaces comma seperated list, eg: 255,255,255

`count` - The amount of blocks you want the image to have, preferably 3. Must be in the range 1-9 (Suggested Value: 3)

`height` - The height of the output image, must be in the range 16-4096 and a multiple of 16 (Suggested Value: 512)

`noText` - If there should be text at the bottom of the image or not, if you do not want text, supply this parameter with any value

`page` - This should be a number between 1 and 100, this is the page of the results you want to get

# Config

The config is in `config.json`

This should contain `port`, the port the service should bind to and `allowedVersions`, which should not be changed as the data file does not support versions besides 1.12 and 1.20