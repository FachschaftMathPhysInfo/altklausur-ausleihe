package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/utils"
	"github.com/adjust/rmq/v3"
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
	log.Printf("working on task %q", task.ExamUUID)
	executeMarkerTask(consumer.MinIOClient, task)

	if err := delivery.Ack(); err != nil {
		log.Println(err)
	}
}

func applyWatermark(input io.ReadSeeker, output io.Writer, text string) error {
	onTop := true
	update := false

	var watermarks []*pdfcpu.Watermark
	// Stamp all odd pages of the pdf in red at the right border of the document
	watermark1, err := pdfcpu_api.TextWatermark(text, "font:Courier, points:40, col: 1 0 0, rot:-90, sc:1 abs, opacity:0.4, pos: l, offset: -190 0", onTop, update, pdfcpu.POINTS)
	if err != nil {
		return err
	}
	watermarks = append(watermarks, watermark1)

	// Stamp all odd pages of the pdf in red at the right border of the document
	watermark2, err := pdfcpu_api.TextWatermark("eq192 - eq192 - eq192", "font:Helvetica, points:40, col: 1 0 0, diagonal:1, sc:1 abs, opacity:0.2, pos: c", onTop, update, pdfcpu.POINTS)
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

	io.Copy(output, &tempout)

	return nil
}

func executeMarkerTask(minioClient *minio.Client, task utils.RMQMarkerTask) {

	context := context.Background()
	examBucket := os.Getenv("MINIO_EXAM_BUCKET")
	cacheBucket := os.Getenv("MINIO_CACHE_BUCKET")

	obj, err := minioClient.GetObject(context, examBucket, task.ExamUUID.String(), minio.GetObjectOptions{})
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
	wmErr := applyWatermark(examReader, bufWriter, task.Text+"1. Mai 2021")
	if wmErr != nil {
		log.Println(wmErr)
	}

	// write back the changed pdf to the bucket storage
	_, putErr := minioClient.PutObject(
		context,
		cacheBucket,
		objInfo.Key,
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
	minioClient := utils.InitMinIO()

	// get job from queue
	errChan := make(chan error, 10)
	go logErrors(errChan)
	rmqClient, err := rmq.OpenConnection(
		os.Getenv("RMQ_QUEUE_NAME"),
		"tcp",
		os.Getenv("REDIS_CONNECTION_STRING"),
		1,
		errChan,
	)
	if err != nil {
		log.Fatalln(err)
	}

	tagQueue, err := rmqClient.OpenQueue("tag-queue")

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
