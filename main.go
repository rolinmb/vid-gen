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
    angle := math.Pi * SCALE * float64(x*x+y*y)
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

func generatePngOne(fnameInc,fnameDec string, complexity float64, colorfactor float64) {
    newPng := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
    for x := 0; x < WIDTH; x++ {
        for y := 0; y < HEIGHT; y++ {
            dx, dy := distort(x, y)
            r, g, b := getPixelColorOne(dx, dy, complexity, colorfactor)
            newPng.Set(x, y, color.RGBA{r, g, b, 255})
        }
    }
    savePng(fnameInc, newPng)
    savePng(fnameDec, newPng)
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
        generatePngOne(fnameInc, fnameDec, multiplier*float64(i), float64(4*i-1))
    }
    runFfmpegOne()
}

func main() {
    routineOne()
}