//
// original : https://qiita.com/souhub/items/9572bd260b71d17c7128
// changes: to check upload file extension, to respond file save result using ajax, to prepare css for visibility improvement
//
// compile option for Raspberry PI zero
// $ GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w" -o "image_name"
//
// using port: WEB server(HTTP/8000), image data sending port(TCP/5000)
//

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	proc "imgproc/proc"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	combinedMessage := fmt.Sprintf("%s (Status: %d)", message, statusCode)
	json.NewEncoder(w).Encode(Response{Success: false, Message: combinedMessage})
}

func sendSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(Response{Success: true, Message: message})
}

// respond to file upload button
func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendErrorResponse(w, "許可されていないメソッド", http.StatusMethodNotAllowed)
		return
	}
	// Parse the multipart form
	err := r.ParseMultipartForm(2 << 20)
	if err != nil {
		sendErrorResponse(w, "ファイルサイズが大き過ぎます", http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		sendErrorResponse(w, "ファイルのアップロード失敗", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".png" {
		sendErrorResponse(w, "許可されていないファイルタイプ", http.StatusBadRequest)
		return
	}

	// to clear img directory files
	if err := clearImageDirectory(); err != nil {
		sendErrorResponse(w, "既存の画像ファイルの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	fixedFilename := "uImage" + ext
	imagePath := filepath.Join("img", fixedFilename)

	saveImage, err := os.Create(imagePath)
	if err != nil {
		sendErrorResponse(w, "ファイル領域の確保失敗", http.StatusInternalServerError)
		return
	}
	defer saveImage.Close()

	// Copy the file
	_, err = io.Copy(saveImage, file)
	if err != nil {
		sendErrorResponse(w, "アップロードしたファイルの書き込み失敗", http.StatusInternalServerError)
		return
	}

	proc.ImgProc()

	sendSuccessResponse(w, "ファイルが正常にラズピコに送信されました。")
}

// respond to text send button : create text to image & send out to rapberry pi pico
func uploadText(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		sendErrorResponse(w, "許可されていないメソッド", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		sendErrorResponse(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		sendErrorResponse(w, "テキストが空です", http.StatusBadRequest)
		return
	}

	// to check text length
	if !isValidTextLength(text) {
		sendErrorResponse(w, "テキストの長さが制限を超えています（半角8文字または全角4文字以内）", http.StatusBadRequest)
		return
	}

	// to clear img directory files
	if err := clearImageDirectory(); err != nil {
		sendErrorResponse(w, "既存の画像ファイルの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	_, err = proc.CreateImageFromText(text)
	if err != nil {
		sendErrorResponse(w, "テキストからイメージの作成に失敗しました", http.StatusInternalServerError)
		return
	}

	proc.ImgProc()

	sendSuccessResponse(w, "テキストからイメージが作成され、ラズピコに送信されました。")
}

// CORS対応ラッパー
/*
func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 開発中は "*"、本番ではオリジンを制限
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}*/

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	// 許可したいオリジンのホワイトリスト（複数も可）
	allowedOrigins := map[string]struct{}{
		"http://rasp-b.local": {},
		"http://localhost": {},
		"http://127.0.0.1": {},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			// キャッシュに影響するのでVaryヘッダーもセット推奨
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// プリフライトリクエストはここで終了
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h(w, r)
	}
}

// JS専用：uploadTextと同等だがCORS対応
func uploadTextJS(w http.ResponseWriter, r *http.Request) {
	log.Println("uploadTextJS called", r.Method, r.URL.Path)

	if r.Method != "POST" {
		sendErrorResponse(w, "許可されていないメソッド", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		sendErrorResponse(w, "フォームの解析に失敗しました", http.StatusBadRequest)
		return
	}

	text := r.FormValue("text")
	if text == "" {
		sendErrorResponse(w, "テキストが空です", http.StatusBadRequest)
		return
	}

	if !isValidTextLength(text) {
		sendErrorResponse(w, "テキストの長さが制限を超えています（半角8文字または全角4文字以内）", http.StatusBadRequest)
		return
	}

	if err := clearImageDirectory(); err != nil {
		sendErrorResponse(w, "既存の画像ファイルの削除に失敗しました", http.StatusInternalServerError)
		return
	}

	_, err = proc.CreateImageFromText(text)
	if err != nil {
		sendErrorResponse(w, "テキストからイメージの作成に失敗しました", http.StatusInternalServerError)
		return
	}

	proc.ImgProc()

	sendSuccessResponse(w, "JS経由: テキストからイメージが作成され、ラズピコに送信されました。")
}

// respond document root(index.html)
func index(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles("index.html"))
	if err := tmp.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

func main() {

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/upload", upload)
	http.HandleFunc("/uploadText", uploadText)
	http.HandleFunc("/uploadText_js", withCORS(uploadTextJS))
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}

// clear the files in the img directory
func clearImageDirectory() error {
	dir := "img"
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	// to clear all image files
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}

	return nil
}

func isValidTextLength(text string) bool {
    runeCount := utf8.RuneCountInString(text)
    byteCount := len(text)

    // full-width length check <= 4
    fullWidthCount := (byteCount - runeCount) / 2

    // 
    return fullWidthCount <= 4 && runeCount <= 8
}

