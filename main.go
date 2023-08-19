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
    WIDTH = 1500
    HEIGHT = 1500
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
    fmt.Println("\nSuccessfully created/rewritten", fname)
}

func getPixelColor(x,y int, complexity,colorfactor float64) (uint8, uint8, uint8) {
    angle := math.Pi * SCALE * (float64(y-x)-math.Sin(float64(x*y))) // Change to +/-/* or divide/modulus by x+1 or y+1
    // angle := math.Pi * SCALE * math.Tan(x+y)*0.2)
    // angle := math.Pi * SCALE * math.Tan(float64(x+(y/((x*y)+1))) - math.Sin(float64(x*y)*0.1))
    // angle := math.Pi * SCALE * math.Sin(float64(x-y)*0.1)/math.Exp(float64(x-y))
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
    // dx := x + int(dist_amp * math.Sin(dist_freq * float64(x) + dist_phase))
	dx := x + int(dist_amp * math.Sin(dist_freq * float64(x) + dist_phase))
    // dy := x + int(dist_amp * math.Sin(dist_freq * float64(y) + dist_phase))
	dy := y + int(dist_amp * math.Sin(dist_freq * float64(y) + dist_phase))
	dx = clamp(dx, 0, WIDTH-1)
	dy = clamp(dy, 0, HEIGHT-1)
    return dx, dy
}

func generatePng(fname string, complexity float64, colorfactor float64) {
    newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx, dy := distort(x, y)
            r, g, b := getPixelColor(dx, dy, complexity, colorfactor)
            newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    savePng(fname, newPng)
}

func runFfmpeg() {
    ffmpegCmd := exec.Command(
        "ffmpeg", "-y",
        "-framerate", "80",
        "-i", "png_out/trial18/trial18_%d.png",
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "vid_out/trial18.mp4",
    )
    cmdOutput, err := ffmpegCmd.CombinedOutput()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Output:")
    fmt.Println(string(cmdOutput))
}

func main() {
    multiplier := 2.0
    for i := 1; i < 81; i++ {
        fnameInc := "png_out/trial18/trial18_"+strconv.FormatInt(int64(i-1), 10)+".png"
        fnameDec := "png_out/trial18/trial18_"+strconv.FormatInt(int64(159-(i-1)), 10)+".png"
        generatePng(fnameInc, float64(multiplier*float64(i)), float64(2*i-1))
        generatePng(fnameDec, float64(multiplier*float64(i)), float64(2*i-1))
    }
    runFfmpeg()
}