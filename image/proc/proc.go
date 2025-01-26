package proc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/image/draw"
	"image"
	"image/png"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"image/color"
	//"log"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	_ "image/jpeg"
)

// to get the running server IP address
func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return " "
	}

	var ip = string("")
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				fmt.Printf("Server running on: http://%s\n", ipnet.IP.String())
			}
		}
	}
	parts := strings.Split(ip, ".")

	if len(parts) == 4 {
		parts[3] = "200:5000"
		ip = strings.Join(parts, ".")
		fmt.Printf("Server running on: http://%s\n", ip)
	}
	return ip
}

// to prepare image path & file name, and call image compression
func ImgProc() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	imgDir := filepath.Join(currentDir, "img")
	files, err := os.ReadDir(imgDir)
	if err != nil {
		fmt.Println("Error reading img directory:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No files found in img directory")
		return
	}

	var inputPath, outputPath string
	maxWidth, maxHeight := 64, 32

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
			inputPath = filepath.Join(imgDir, file.Name())
			outputPath = filepath.Join(imgDir, "compress_"+strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))+".png")
			break
		}
	}

	if inputPath == "" {
		fmt.Println("No supported image file found in img directory")
		return
	}

	// compresss, extruct g/b/r elements and send to raspberry pi pico via TCP/IP protcol
	//
	compressImage(inputPath, outputPath, maxWidth, maxHeight)
	fmt.Println("Image compressed and saved successfully as PNG")

	bgr(outputPath)
	fmt.Println("RGB values extracted and saved to proc/output.txt and img/output.bin")

	sendTcpIp()
}

// compress the recceived image file to 64*32 png imagefile and save it to the same directory
func compressImage(inputPath, outputPath string, maxWidth, maxHeight int) {
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	dst := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer out.Close()

	// Always encode as PNG
	err = png.Encode(out, dst)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		return
	}

}

// extract r/g/b from compressed image file and write it to a binary file
func bgr(outputPath string) {
	// open the image file
	file, err := os.Open(outputPath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// decode the original file
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// check the image file size
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	if width != 64 || height != 32 {
		fmt.Println("Image size is not 64x32")
		return
	}

	// three dimensions array for R/G/B storing
	gbrArray := make([][][]uint8, 3) // R, G, B
	for i := range gbrArray {
		gbrArray[i] = make([][]uint8, height)
		for y := range gbrArray[i] {
			gbrArray[i][y] = make([]uint8, width)
		}
	}

	// to get BGR value
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			gbrArray[2][y][x] = uint8(math.Round(float64(r>>8) / 16))
			gbrArray[1][y][x] = uint8(math.Round(float64(g>>8) / 16))
			gbrArray[0][y][x] = uint8(math.Round(float64(b>>8) / 16))
		}
	}

	// Write to the binary file
	binaryOutputPath := filepath.Join(filepath.Dir(outputPath), "output.bin")
	binaryFile, err := os.Create(binaryOutputPath)
	if err != nil {
		fmt.Println("Error creating binary output file:", err)
		return
	}
	defer binaryFile.Close()

	// Write BGR values
	for i := 0; i < 3; i++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				_, err = binaryFile.Write([]byte{gbrArray[i][y][x]})
				if err != nil {
					fmt.Println("Error writing BGR value to binary file:", err)
					return
				}
			}
		}
	}
}

func sendTcpIp() {
	// raspberry pi PICO address & port number
	serverAddr := getIp()

	// Set a timeout for the entire operation
	timeout := 30 * time.Second
	conn, err := net.DialTimeout("tcp", serverAddr, timeout)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Set deadlines for write and read operations
	deadline := time.Now().Add(timeout)
	err = conn.SetDeadline(deadline)
	if err != nil {
		fmt.Println("Error setting deadline:", err)
		return
	}

	// Read the binary file
	binaryFilePath := "img/output.bin"
	fileData, err := os.ReadFile(binaryFilePath)
	if err != nil {
		fmt.Println("Error reading binary file:", err)
		return
	}

	// Prepare the data to send
	dataLength := uint16(len(fileData))

	// Create a buffer to hold all the data
	buffer := new(bytes.Buffer)

	// Write length (2 bytes)
	binary.Write(buffer, binary.BigEndian, dataLength)

	// Write file data
	buffer.Write(fileData)

	// Send all data
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Data sent successfully")

	// After sending data, reset the deadline for reading the response
	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		fmt.Println("Error setting read deadline:", err)
		return
	}

	// Receive the response
	respBuffer := make([]byte, 1024)
	n, err := conn.Read(respBuffer)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Read timeout")
		} else {
			fmt.Println("Error receiving data:", err)
		}
		return
	}

	fmt.Printf("Received response: %s\n", string(respBuffer[:n]))
}

func CreateImageFromText(text string) (string, error) {
    width, height := 400, 200

    img := image.NewRGBA(image.Rect(0, 0, width, height))

    // 背景を黒にします
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            img.Set(x, y, color.Black)
        }
    }

	// in raspberry pi case
	// sudo apt install fonts-noto-cjk
	// change the locale
	//
	fontPath := "/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc"
	// in mac case
    //fontPath := "/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc"
    fontBytes, err := os.ReadFile(fontPath)
    if err != nil {
        return "", fmt.Errorf("フォントファイルの読み込みに失敗しました: %v", err)
    }

    fontCollection, err := opentype.ParseCollection(fontBytes)
    if err != nil {
        return "", fmt.Errorf("フォントコレクションのパースに失敗しました: %v", err)
    }

    f, err := fontCollection.Font(0)
    if err != nil {
        return "", fmt.Errorf("フォントの取得に失敗しました: %v", err)
    }

    fontSize := 30.0*3 // フォントサイズを調整
    face, err := opentype.NewFace(f, &opentype.FaceOptions{
        Size: fontSize,
        DPI:  72,
    })
    if err != nil {
        return "", fmt.Errorf("フォントフェイスの作成に失敗しました: %v", err)
    }

    textColor := color.White

    // テキストの幅と高さを計算
    bounds, _ := font.BoundString(face, text)
    textWidth := float64(bounds.Max.X - bounds.Min.X) / 64
    textHeight := float64(bounds.Max.Y - bounds.Min.Y) / 64

    // 描画位置を計算（中央揃え）
    x := (float64(width) - textWidth) / 2
    y := (float64(height) + textHeight) / 2

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(textColor),
        Face: face,
        Dot:  fixed.Point26_6{
            X: fixed.Int26_6(x * 64),
            Y: fixed.Int26_6(y * 64),
        },
    }
    d.DrawString(text)

    // 画像を保存
    imagePath := filepath.Join("img", "uImage.png")
    outputFile, err := os.Create(imagePath)
    if err != nil {
        return "", fmt.Errorf("出力ファイルの作成に失敗しました: %v", err)
    }
    defer outputFile.Close()

    err = png.Encode(outputFile, img)
    if err != nil {
        return "", fmt.Errorf("画像のエンコードに失敗しました: %v", err)
    }

    return imagePath, nil
}
