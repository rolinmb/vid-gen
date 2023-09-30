package main

import (
    "fmt"
    "log"
    "math"
    "strconv"
    "image"
    "image/color"
    "image/png"
    "os"
    "os/exec"
)

const (
    WIDTH = 1600
    HEIGHT = 1600
    SCALE = 0.001
    // COMPLEXITY = 10
    // COLORFACTOR = 40
    dist_amp = 50.0
	dist_freq = 0.01
	dist_phase = 0.0
)

func savePng(fname string, newPng *image.RGBA) {
    out, err := os.Create(fname)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()
    err = png.Encode(out, newPng)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Successfully created/rewritten", fname)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func distort(x, y int) (int, int) {
    dx := y + int(dist_amp * math.Sin(dist_freq * float64(x) + dist_phase))
	//dx := x + int(dist_amp * math.Sin(dist_freq * float64(x) + dist_phase))
    dy := x + int(dist_amp * math.Sin(dist_freq * float64(y) + dist_phase))
	//dy := y + int(dist_amp * math.Sin(dist_freq * float64(y) + dist_phase))
	dx = clamp(dx, 0, WIDTH-1)
	dy = clamp(dy, 0, HEIGHT-1)
    return dx, dy
}

func getPixelColorOne(x,y int, complexity,colorfactor float64) (uint8, uint8, uint8) {
    // angle := (math.Pi * SCALE * math.Sin(float64(x*x*y/5))) // Change to +/-/* or divide/modulus by x+1 or y+1
    // angle := math.Pi * SCALE * (math.Cos(float64(x+y))*0.2)
    // angle := math.Pi * SCALE * math.Tan(float64(x+(y/((x*y)+1))) - math.Sin(float64(x*y)*0.1))
    // angle := math.Pi * 2.0 * SCALE *  math.Sqrt(float64(x*x+y*y)) // circular gradient
    angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) // sine wave ripple
    // angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))  // hyperbolic spiral
    // angle := math.Pi * SCALE * float64(x*x+y*y)  // square
    // angle := math.Pi * SCALE * math.Exp(-0.01 * math.Sqrt(float64(x*x+y*y)))  // exponential decay
    // angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4  // offset sine wave
    // angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y))  // lissajous curve
    // angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
    // angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2) // diamond
    // angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))  // hypotochoid
    // angle := math.Pi * SCALE * (math.Sin(float64(x)*0.05) * math.Exp(-float64(y)*0.1))
    // Remix Set 1 Generators
    //  centerX := float64(WIDTH) / 2
    //  centerY := float64(HEIGHT) / 2
    /*angle := math.Atan2(float64(y)-centerY, float64(x)-centerX)
    angle = angle + complexity*math.Sin(angle)
    r := uint8((math.Sin(angle) + 1) * 127.5 * colorfactor)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * colorfactor)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * colorfactor)*/
    /*distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    distance = distance + complexity*math.Sin(distance*0.1)
    r := uint8((math.Sin(distance*0.1) + 1) * 127.5 * colorfactor)
    g := uint8((math.Sin(distance*0.1+2*math.Pi/3) + 1) * 127.5 * colorfactor)
    b := uint8((math.Sin(distance*0.1+4*math.Pi/3) + 1) * 127.5 * colorfactor)*/
    /*angleX := math.Pi * complexity * (float64(x) - centerX) / centerX
    angleY := math.Pi * complexity * (float64(y) - centerY) / centerY
    angle := math.Sin(angleX) + math.Cos(angleY)
    angle = angle + complexity*math.Sin(angle*10)
    r := uint8((math.Sin(angle) + 1) * 127.5 * colorfactor)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * colorfactor)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * colorfactor)*/
    /*distance := math.Sqrt(math.Pow(float64(x)-centerX, 2) + math.Pow(float64(y)-centerY, 2))
    angleX := math.Pi * complexity * (float64(x) - centerX) / centerX
    angleY := math.Pi * complexity * (float64(y) - centerY) / centerY
    combinedAngle := angleX + angleY + distance
    trippyAngle := combinedAngle + complexity*math.Sin(combinedAngle)
    r := uint8((math.Sin(trippyAngle) + 1) * 127.5 * colorfactor)
    g := uint8((math.Sin(trippyAngle+2*math.Pi/3) + 1) * 127.5 * colorfactor)
    b := uint8((math.Sin(trippyAngle+4*math.Pi/3) + 1) * 127.5 * colorfactor)*/
    // Main RGB functions
    distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * complexity + frequency) * colorfactor + 128)
    // r := uint8(0)
    g := uint8(math.Sin(angle * complexity + frequency + 2*math.Pi/3) * colorfactor + 128)
	// g := uint8(0)
    b := uint8(math.Sin(angle * complexity + frequency + 4*math.Pi/3) * colorfactor + 128)
    // b := uint8(0)
    return r, g, b
}

func generatePngOne(fname string, complexity float64, colorfactor float64) {
    newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx, dy := distort(x, y)
            r, g, b := getPixelColorOne(dx, dy, complexity, colorfactor)
            newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    savePng(fname, newPng)
}

func getPixelColorTwo(x,y int, complexity,colorfactor float64) (uint8, uint8, uint8) {
	// angle := math.Pi * SCALE * (float64(x)/10 + float64(y)/5) + math.Pi * SCALE * math.Sin(2*float64(x))  // hypotochoid
	// angle := math.Pi * SCALE * math.Sin(3*float64(x)) + math.Pi * SCALE * math.Sin(4*float64(y)) // lissajous curve
	// angle := math.Pi * SCALE * math.Sin(0.1*float64(x)+0.1*float64(y)) + math.Pi/4  // offset sine wave
	// angle := math.Pi * SCALE * math.Sin(float64(x)*0.1) + math.Pi * SCALE * math.Sin(float64(y)*0.1) // sine wave ripple
	// angle := math.Pi * SCALE * math.Abs(float64(x)-WIDTH/2) + math.Pi * SCALE * math.Abs(float64(y)-HEIGHT/2) // diamond
	// angle := math.Pi * SCALE * (1 / (float64(x)*float64(x) + float64(y)*float64(y) + 1))  // hyperbolic spiral
	// angle := math.Pi * SCALE * math.Cos(3*float64(x*y)) + math.Pi * SCALE * math.Sin(4*float64(y))
	// Remix Set 1 Generators
    centerX := float64(WIDTH) / 2
    centerY := float64(HEIGHT) / 2
    angle := math.Atan2(float64(x*y)-centerY, float64(x*x*y*y)-centerX)
    // angle = angle + complexity*math.Sin(angle)
    r := uint8((math.Sin(angle) + 1) * 127.5 * colorfactor)
    g := uint8((math.Sin(angle+2*math.Pi/3) + 1) * 127.5 * colorfactor)
    b := uint8((math.Sin(angle+4*math.Pi/3) + 1) * 127.5 * colorfactor)
	/*distance := math.Sqrt(math.Pow(float64(x-WIDTH/2), 2) + math.Pow(float64(y-HEIGHT/2), 2))
	frequency := distance * SCALE
    r := uint8(math.Sin(angle * complexity + frequency) * colorfactor + 128)
    // r := uint8(0)
    g := uint8(math.Sin(angle * complexity + frequency + 2*math.Pi/3) * colorfactor + 128)
	// g := uint8(0)
    b := uint8(math.Sin(angle * complexity + frequency + 4*math.Pi/3) * colorfactor + 128)
    // b := uint8(0)*/
    return r, g, b
}

func generatePngTwo(fname string, complexity float64, colorfactor float64) {
	newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx, dy := distort(x, y)
            r, g, b := getPixelColorTwo(dx, dy, complexity, colorfactor)
            newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    savePng(fname, newPng)
}

func runFfmpegOne() {
    ffmpegCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/trial1/trial1_09302023_%d.png",
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "vid_out/trial1_09302023.mp4",
    )
    cmdOutput, err := ffmpegCmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Output:")
    fmt.Println(string(cmdOutput))
}

func routineOne() {
	multiplier := 0.08333
    for i := 1; i < 31; i++ {
        fnameInc := "png_out/trial1/trial1_09302023_"+strconv.FormatInt(int64(i-1), 10)+".png"
        fnameDec := "png_out/trial1/trial1_09302023_"+strconv.FormatInt(int64(59-(i-1)), 10)+".png"
        generatePngOne(fnameInc, multiplier*float64(i), float64(4*i-1))
        generatePngOne(fnameDec, multiplier*float64(i), float64(4*i-1))
    }
    runFfmpegOne()
}

func interpolatePngs(fIn1 string, fIn2 string, fOut string, factor float64) {
	fileIn1, err := os.Open(fIn1)
	if err != nil {
		fmt.Printf("Failed to decode input .png file: %v", err)
		return
	}
	defer fileIn1.Close()
	pngIn1, err := png.Decode(fileIn1)
	if err != nil {
		fmt.Printf("Failed to decode input .png file: %v", err)
		return
	}
	fileIn2, err := os.Open(fIn2)
	if err != nil {
		fmt.Printf("Failed to decode input .png file: %v", err)
		return
	}
	defer fileIn2.Close()
	pngIn2, err := png.Decode(fileIn2)
	if err != nil {
		fmt.Printf("Failed to decode input .png file: %v", err)
		return
	}
	bounds := pngIn1.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	newPng := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c1 := pngIn1.At(x, y)
			c2 := pngIn2.At(x, y)
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			r := uint8(float64(r1)*(float64(1)-factor) + float64(r2)*factor)
			g := uint8(float64(g1)*(float64(1)-factor) + float64(g2)*factor)
			b := uint8(float64(b1)*(float64(1)-factor) + float64(b2)*factor)
			a := uint8(float64(a1)*(float64(1)-factor) + float64(a2)*factor)
			newPng.SetRGBA(x, y, color.RGBA{r, g, b, a})
		}
	}
	savePng(fOut, newPng)
	fileIn1.Close()
	fileIn2.Close()
}

func runFfmpegTwo() {
	ffmpegCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "30",
        "-i", "png_out/trial10/trial10_09302023i_%d.png",
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "vid_out/trial10_09302023i.mp4",
    )
	cmdOutput, err := ffmpegCmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Output:")
    fmt.Println(string(cmdOutput))
}

func routineTwo() {
	multiplier := 4.0
	interpFactor := 0.5
	for i := 1; i < 31; i++ {
		fnameOne := "png_out/trial10/trial10_09302023_"+strconv.FormatInt(int64(i-1), 10)+".png"
		fnameTwo := "png_out/trial10/trial10_09302023_"+strconv.FormatInt(int64(59-(i-1)), 10)+".png"
		fnameInterp := "png_out/trial10/trial10_09302023i_"+strconv.FormatInt(int64(i-1), 10)+".png"
		generatePngOne(fnameOne, multiplier*float64(i), float64(10*i-1))
        generatePngTwo(fnameTwo, multiplier*float64(i), float64(10*i-1))
		interpolatePngs(fnameOne, fnameTwo, fnameInterp, interpFactor)
	}
	runFfmpegTwo()
}

func main() {
    //routineOne()
	routineTwo()
}