package handler

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImage_rectangles(t *testing.T) {
	i := NewImage("testData/haarcascade_frontalface_alt.xml")
	src, err := ioutil.ReadFile("testData/dest.jpg")
	buffer := bytes.NewBuffer(nil)
	if assert.NoError(t, err) {
		err = i.rectangles(src, buffer)
		if assert.NoError(t, err) {
			assert.NotNil(t, buffer.Bytes())
			log.Println(buffer.String())
		}
	}
}
