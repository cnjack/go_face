package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"image"
	"image/color"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInvalidImage   = errors.New("invalid image")
	ErrPeopleNotFound = errors.New("people not found in image")
	ErrEncode         = errors.New("encode image error")
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

func (i *Image) Draw(c *gin.Context) {
	src, err := ioutil.ReadAll(c.Request.Body)
	if err != nil || src == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewCommonRespWithError(http.StatusBadRequest, ErrBadRequest))
		return
	}

	if err := i.draw(src, c.Writer); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, NewCommonRespWithError(http.StatusBadRequest, err))
		return
	}
}

func (i *Image) draw(src []byte, w http.ResponseWriter) error {
	mat := gocv.IMDecode(src, gocv.IMReadUnchanged)
	if mat.Empty() {
		return ErrInvalidImage
	}
	defer mat.Close()
	rectangles := i.c.DetectMultiScale(mat)
	if len(rectangles) == 0 {
		return ErrPeopleNotFound
	}
	for keys := range rectangles {
		color.Black.RGBA()
		gocv.Line(
			mat,
			rectangles[keys].Min,
			image.Point{X: rectangles[keys].Min.X, Y: rectangles[keys].Max.Y},
			color.RGBA{0, 0, 0, 1},
			2,
		)
		gocv.Line(
			mat,
			rectangles[keys].Min,
			image.Point{X: rectangles[keys].Max.X, Y: rectangles[keys].Min.Y},
			color.RGBA{0, 0, 0, 1},
			2,
		)
		gocv.Line(
			mat,
			rectangles[keys].Max,
			image.Point{X: rectangles[keys].Min.X, Y: rectangles[keys].Max.Y},
			color.RGBA{0, 0, 0, 1},
			2,
		)
		gocv.Line(
			mat,
			rectangles[keys].Max,
			image.Point{X: rectangles[keys].Max.X, Y: rectangles[keys].Min.Y},
			color.RGBA{0, 0, 0, 1},
			2,
		)
	}
	if mat.Empty() {
		return ErrEncode
	}
	img, err := gocv.IMEncode(".png", mat)
	if err != nil {
		return ErrEncode
	}
	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(img)
	return err
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
