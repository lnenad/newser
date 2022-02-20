package pkg

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

func SetupPdf() fpdf.Pdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.AddUTF8Font("robotoserif", "", filepath.Join("fonts", "RobotoSerif-Regular.ttf"))
	pdf.AddUTF8Font("robotocondensed", "", filepath.Join("fonts", "RobotoCondensed-Regular.ttf"))

	pdf.SetFont("robotoserif", "", 4)
	pdf.CellFormat(
		0,
		2,
		strings.Join(
			[]string{"Generated at", time.Now().Format("2006-01-02T15:04:05Z07:00")},
			" ",
		),
		"",
		2,
		"R",
		false,
		0,
		"",
	)
	pdf.SetFont("robotoserif", "", 16)
	pdf.CellFormat(
		0,
		25,
		strings.Join(
			[]string{"Newspaper date", time.Now().Format("2006-02-01")},
			" ",
		),
		"",
		2,
		"C",
		false,
		0,
		"",
	)
	return pdf
}

func WriteHeader(pdf fpdf.Pdf, header string) {
	pdf.SetFont("robotoserif", "", 11)
	pdf.CellFormat(0, 20, header, "", 2, "L", false, 0, "")
}

func WriteArticle(config Config, pdf fpdf.Pdf, article Article) {
	pdf.SetFont("robotoserif", "", float64(config.Font.Title))
	pdf.MultiCell(0, 6, article.Title, "", "", false)
	pdf.CellFormat(0, 5, "", "", 2, "L", false, 0, "")
	pdf.SetFont("robotocondensed", "", float64(config.Font.Content))
	pdf.MultiCell(0, 4.5, article.Content, "", "", false)
	pdf.CellFormat(0, 10, "", "", 2, "L", false, 0, "")
}

func RegisterImage(pdf fpdf.Pdf, urlStr string) {
	const (
		margin   = 10
		ht       = 30
		fontSize = 15
	)

	var (
		rsp *http.Response
		err error
		tp  string
	)

	ln := pdf.PointConvert(fontSize)
	rsp, err = http.Get(urlStr)
	if err == nil {
		tp = pdf.ImageTypeFromMime(rsp.Header["Content-Type"][0])
		infoPtr := pdf.RegisterImageReader(urlStr, tp, rsp.Body)
		if pdf.Ok() {
			imgWd, imgHt := infoPtr.Extent()
			pdf.Image(urlStr, pdf.GetX()+ln, pdf.GetY()+ln,
				imgWd, imgHt, false, tp, 0, "")
		}
	} else {
		pdf.SetError(err)
	}
}
