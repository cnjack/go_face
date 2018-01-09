package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInvalidImage   = errors.New("invalid image")
	ErrPeopleNotFound = errors.New("people not found in image")
)

type Image struct {
	c gocv.CascadeClassifier
}

func NewImage(xmlPath string) *Image {
	classifier := gocv.NewCascadeClassifier()
	classifier.Load(xmlPath)
	return &Image{
		c: classifier,
	}
}

func (i *Image) Html(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (i *Image) Rectangles(c *gin.Context) {
	src, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || src == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewCommonRespWithError(http.StatusBadRequest, ErrBadRequest))
		return
	}

	if err := i.rectangles(src, c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewCommonRespWithError(http.StatusBadRequest, err))
		return
	}
}

type RectangleResp struct {
	Min Point `json:"min"`
	Max Point `json:"max"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (i *Image) rectangles(src []byte, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	mat := gocv.IMDecode(src, gocv.IMReadUnchanged)
	if mat.Empty() {
		return ErrInvalidImage
	}
	defer mat.Close()

	rectangles := i.c.DetectMultiScale(mat)
	if len(rectangles) == 0 {
		return ErrPeopleNotFound
	}

	resp := make([]*RectangleResp, 0)
	for keys := range rectangles {
		resp = append(resp, &RectangleResp{
			Point(rectangles[keys].Min),
			Point(rectangles[keys].Max),
		})
	}

	e := json.NewEncoder(w)
	return e.Encode(NewCommonRespWithData(resp))
}
