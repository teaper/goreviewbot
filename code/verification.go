package code

import (
	"embed"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"os"
	"strconv"
)

//go:embed imgs/*
var imgNums embed.FS

//9：生成验证码
func CreateVerificationCode(codes [4] int32) {
	file1, _ := imgNums.Open("imgs/"+strconv.FormatInt(int64(codes[0]), 10)+".gif")
	file2, _ := imgNums.Open("imgs/"+strconv.FormatInt(int64(codes[1]), 10)+".gif")
	file3, _ := imgNums.Open("imgs/"+strconv.FormatInt(int64(codes[2]), 10)+".gif")
	file4, _ := imgNums.Open("imgs/"+strconv.FormatInt(int64(codes[3]), 10)+".gif")
	defer file1.Close()
	defer file2.Close()
	defer file3.Close()
	defer file4.Close()


	// 加入合并的图片
	var (
		img1, img2, img3, img4 image.Image
		err        error
	)
	if img1, _, err = image.Decode(file1); err != nil {
		log.Println(err)
	}
	if img2, _, err = image.Decode(file2); err != nil {
		log.Println(err)
	}
	if img3, _, err = image.Decode(file3); err != nil {
		log.Println(err)
	}
	if img4, _, err = image.Decode(file4); err != nil {
		log.Println(err)
	}

	// 将四张图片合成一张
	newWidth := (img1.Bounds().Max.X)*4 //新图片宽度
	newHeight := img2.Bounds().Max.Y //新图片高度
	newImg := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight)) //创建一个新RGBA图像
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), img1, img1.Bounds().Min, draw.Over) //画上第一张缩放后的图片
	draw.Draw(newImg, newImg.Bounds(), img2, img2.Bounds().Min.Sub(image.Pt(img1.Bounds().Max.X, 0)), draw.Over) //画上第二张缩放后的图片（注意X值的起始位置）
	draw.Draw(newImg, newImg.Bounds(), img3, img3.Bounds().Min.Sub(image.Pt(img1.Bounds().Max.X, 0)).Sub(image.Pt(img2.Bounds().Max.X, 0)), draw.Over)
	draw.Draw(newImg, newImg.Bounds(), img4, img4.Bounds().Min.Sub(image.Pt(img1.Bounds().Max.X, 0)).Sub(image.Pt(img2.Bounds().Max.X, 0)).
		Sub(image.Pt(img3.Bounds().Max.X, 0)), draw.Over)

	// 保存成图片
	imgfile, _ := os.Create("code.gif")
	defer imgfile.Close()
	gif.Encode(imgfile, newImg, &gif.Options{NumColors: 5000})
	log.Println("随机4位数 ==> ",codes[0],codes[1],codes[2],codes[3])
}
