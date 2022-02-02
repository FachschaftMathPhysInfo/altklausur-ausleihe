package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/utils"
	"github.com/adjust/rmq/v3"
	render "github.com/brunsgaard/go-pdfium-render"
	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
	"github.com/minio/minio-go/v7"
	pdfcpu_api "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

const (
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 5

	reportBatchSize = 10000
	consumeDuration = time.Millisecond
	shouldLog       = false
)

type RMQConsumer struct {
	name        string
	count       int
	before      time.Time
	MinIOClient *minio.Client
}

func NewRMQConsumer(minioClient *minio.Client, tag int) *RMQConsumer {
	return &RMQConsumer{
		name:        fmt.Sprintf("consumer%d", tag),
		count:       0,
		before:      time.Now(),
		MinIOClient: minioClient,
	}
}

func (consumer *RMQConsumer) Consume(delivery rmq.Delivery) {
	// perform task
	var task utils.RMQMarkerTask
	if err := json.Unmarshal([]byte(delivery.Payload()), &task); err != nil {
		// is this the correct error handling?
		delivery.Reject()
	}
	log.Printf("%s working on task %q", consumer.name, task.ExamUUID)
	executeMarkerTask(consumer.MinIOClient, task)

	if err := delivery.Ack(); err != nil {
		log.Println(err)
	}
}

func applyWatermark(input io.ReadSeeker, output io.Writer, textLeft string, textDiagonal string) error {
	onTop := true
	update := false

	var watermarks []*pdfcpu.Watermark
	// Stamp all odd pages of the pdf in red at the right border of the document
	watermark1, err := pdfcpu_api.TextWatermark(textLeft, "font:Courier, points:20, col: 1 0 0, rot:-90, sc: 0.8 rel, opacity:0.4, off: -260 0", onTop, update, pdfcpu.POINTS)
	if err != nil {
		return err
	}
	watermarks = append(watermarks, watermark1)

	// Stamp all odd pages of the pdf in red at the right border of the document
	watermark2, err := pdfcpu_api.TextWatermark(textDiagonal, "font:Helvetica, points:40, col: 1 0 0, diagonal:1, sc:1 abs, opacity:0.2, pos: c", onTop, update, pdfcpu.POINTS)
	if err != nil {
		return err
	}
	watermarks = append(watermarks, watermark2)

	// add the mathphys logo to the top-left corner? :P
	// wm, err = pdfcpu_api.PDFWatermark("MathPhysLogoInfo.pdf", "pos:tr, rot:0, sc:0.5 abs, offset: -10 -10, opacity:0.5", onTop, update, pdfcpu.POINTS)
	// pdfcpu_api.AddWatermarks(input, output, nil, wm, nil)
	// if err != nil {
	//	return err
	// }

	if len(watermarks) == 0 {
		return fmt.Errorf("no watermarks in array")
	}

	var tempout, tempin bytes.Buffer
	// handle the first case seperately to save one copy
	if err = pdfcpu_api.AddWatermarks(input, &tempout, nil, watermarks[0], nil); err != nil {
		return err
	}
	for _, watermark := range watermarks[1:] {
		// swap the two buffers
		tempin, tempout = tempout, tempin
		tempout.Reset()
		tmpReader := bytes.NewReader(tempin.Bytes())

		if err = pdfcpu_api.AddWatermarks(tmpReader, &tempout, nil, watermark, nil); err != nil {
			return err
		}
	}
	bits := tempout.Bytes()
	doc, err := render.NewDocument(&bits)
	if err != nil {
		return err
	}
	pagesbuf := make([]io.Reader, 0)
	for i := 0; i < doc.GetPageCount(); i++ {
		img := doc.RenderPage(i, 150)
		buf := new(bytes.Buffer)
		png.Encode(buf, img)
		pagesbuf = append(pagesbuf, bytes.NewReader(buf.Bytes()))
	}
	conf := pdfcpu.NewDefaultConfiguration()
	imp := pdfcpu.DefaultImportConfig()
	imp.Gray = false
	buf2 := new(bytes.Buffer)
	pdfcpu_api.ImportImages(nil, buf2, pagesbuf, imp, conf)
	ctx, err := pdfcpu.Read(bytes.NewReader(buf2.Bytes()), pdfcpu.NewDefaultConfiguration())
	keySec, err := nacl.Load("b538ff93d9b028a767c2f8afc05d586936b2bc0ba5c04eddf0b58f381de2a433")
	if err != nil {
		log.Fatalln(err)
	}
	pkey, err := nacl.Load("bbf05a8f323315477201cd51176b86ee5267f459d1792b743a792be265c678a2")
	if err != nil {
		log.Fatalln(err)
	}
	for i, r := range ctx.XRefTable.Table {
		if cast, ok := r.Object.(pdfcpu.StreamDict); ok {

			encrypted := box.EasySeal([]byte(textLeft), keySec, pkey)
			encrypted2 := box.EasySeal([]byte(textDiagonal), keySec, pkey)
			cast.Dict.InsertString("ref1", base64.StdEncoding.EncodeToString(encrypted))
			cast.Dict.InsertString("ref2", base64.StdEncoding.EncodeToString(encrypted2))
			ctx.XRefTable.Table[i].Object = cast
		}
	}
	ctx.EnsureVersionForWriting()
	pdfcpu_api.WriteContext(ctx, output)
	return nil
}

func executeMarkerTask(minioClient *minio.Client, task utils.RMQMarkerTask) {

	context := context.Background()
	examBucket := os.Getenv("MINIO_EXAM_BUCKET")
	cacheBucket := os.Getenv("MINIO_CACHE_BUCKET")

	obj, err := minioClient.GetObject(
		context,
		examBucket,
		task.ExamUUID.String(),
		minio.GetObjectOptions{})

	if err != nil {
		log.Println(err)
	}

	objInfo, err := obj.Stat()
	if err != nil {
		log.Println(err)
	}

	// create a new buffer to write the watermarked pdf into
	var b bytes.Buffer
	bufWriter := bufio.NewWriter(&b)

	// read the object from the bucket storage
	// (dunno why this has to be done, doesnt work otherwise :P)
	exam, err := ioutil.ReadAll(obj)
	if err != nil {
		log.Println(err)
	}
	examReader := bytes.NewReader(exam)

	// apply the watermark to the PDF
	wmErr := applyWatermark(examReader, bufWriter, task.TextLeft, task.TextDiagonal)
	if wmErr != nil {
		log.Println(wmErr)
	}

	// write back the changed pdf to the bucket storage
	_, putErr := minioClient.PutObject(
		context,
		cacheBucket,
		utils.GetExamCachePath(task.UserID, task.ExamUUID),
		&b,
		int64(b.Len()),
		minio.PutObjectOptions{
			ContentType: objInfo.ContentType,
		},
	)
	if putErr != nil {
		log.Println(err)
	}
}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}

func main() {
	render.InitLibrary()
	minioClient := utils.InitMinIO()

	// get job from queue
	rmqClient := utils.InitRmq()

	tagQueue, err := rmqClient.OpenQueue("tag-queue")
	if err != nil {
		log.Fatalln("Error while opening RMQ Queue: ", err)
	}

	if err = tagQueue.StartConsuming(10, time.Second); err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("Started consumer #%d!", i)
		log.Println(name)
		if _, err := tagQueue.AddConsumer(name, NewRMQConsumer(minioClient, i)); err != nil {
			log.Fatalln(err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()

	<-rmqClient.StopAllConsuming() // wait for all Consume() calls to finish
}
