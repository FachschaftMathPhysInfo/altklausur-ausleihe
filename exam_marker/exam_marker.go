package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/exam_marker/v2/utils"
	"github.com/minio/minio-go/v7"
	api "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func applyWatermark(input io.ReadSeeker, output io.Writer, text string) error {
	onTop := true
	update := false

	// Stamp all odd pages of the pfd in red at the right border of the document
	watermark, err := api.TextWatermark(text, "font:Courier, points:40, col: 1 0 0, rot:-90, sc:1 abs, opacity:0.4, pos: l, offset: -190 0", onTop, update, pdfcpu.POINTS)
	if err != nil {
		return err
	}
	err = api.AddWatermarks(input, output, nil, watermark, nil)
	if err != nil {
		return err
	}

	// add the mathphys logo to the top-left corner? :P
	// wm, err = api.PDFWatermark("MathPhysLogoInfo.pdf", "pos:tr, rot:0, sc:0.5 abs, offset: -10 -10, opacity:0.5", onTop, update, pdfcpu.POINTS)
	// api.AddWatermarks(input, output, nil, wm, nil)
	// if err != nil {
	//	return err
	// }

	return nil
}

func watermarkFile(filename string) {

	context := context.Background()
	examBucket := os.Getenv("MINIO_EXAM_BUCKET")
	cacheBucket := os.Getenv("MINIO_CACHE_BUCKET")

	minioClient := utils.InitMinIO()
	obj, err := minioClient.GetObject(context, examBucket, filename, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	objInfo, err := obj.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	// create a new buffer to write the watermarked pdf into
	var b bytes.Buffer
	bufWriter := bufio.NewWriter(&b)

	// read the object from the bucket storage
	// (dunno why this has to be done, doesnt work otherwise :P)
	exam, err := ioutil.ReadAll(obj)
	if err != nil {
		log.Fatalln(err)
	}
	examReader := bytes.NewReader(exam)

	// apply the watermark to the PDF
	wmErr := applyWatermark(examReader, bufWriter, "test123 1. Mai 2021")
	if wmErr != nil {
		log.Fatalln(wmErr)
	}

	// write back the changed pdf to the bucket storage
	minioClient.PutObject(
		context,
		cacheBucket,
		objInfo.Key,
		&b,
		int64(b.Len()),
		minio.PutObjectOptions{ContentType: objInfo.ContentType},
	)
}

func main() {
	// get job from queue
	// TODO: implement

	// watermark file
	watermarkFile("d7b99c10-dd23-487f-af02-53e6a41b4b65")
	watermarkFile("044fb0f6-48de-4c27-b778-8cff8db7c67c")

	// flatten the PDF by going from PDF to IMG to PDF so the watermark cant be removed
	// TODO: implement
}
