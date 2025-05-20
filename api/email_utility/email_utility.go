package email_utility

import (
	"finala/api/config"
	"fmt"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Email sender interface for sending emails
type EmailSender interface {
	Send(to, subject, body string) error
}

// SMTPSender is the implementation of gmail.v2
type SMTPSender struct {
	Dialer *gomail.Dialer
}

const (
	marginH = 10.0
	lineHt  = 3.0
	cellGap = 0.5
)

type cellType struct {
	str  string
	list [][]byte
	ht   float64
}

var cell cellType

func NewSMTPSender(host string, port int, username, password string) *SMTPSender {
	return &SMTPSender{
		Dialer: gomail.NewDialer(host, port, username, password),
	}
}

func (s *SMTPSender) Send(to, subject, body, attachment string) error {
	mailArray := strings.Split(to, ",")
	m := gomail.NewMessage()
	m.SetHeader("From", s.Dialer.Username)
	m.SetHeader("To", mailArray...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attachment != "" {
		m.Attach(attachment)
	}
	d := s.Dialer
	return d.DialAndSend(m)
}

// CreatePDF generates a PDF document from the given data.
func CreatePDF(pdfFileName, description string, data []map[string]interface{}, sendEmailInfo config.SendEmailInfo) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 10)

	// Title
	pdf.Cell(0, 10, "Finala Report")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 8)
	pdf.Cell(0, 10, fmt.Sprintf("Resource Type: %s", sendEmailInfo.ResourceType))
	pdf.Ln(5)

	// Title
	pdf.Cell(0, 10, description)
	pdf.Ln(10)

	// Set font
	pdf.SetFont("Arial", "", 6)
	alignList := []string{"L", "C", "R"}
	orderKeys, columnWidths := configureHeaderColumn(sendEmailInfo, data)
	// set header
	for i, colNames := range orderKeys {
		pdf.CellFormat(columnWidths[i], 10, colNames, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	// set table cell values
	y := pdf.GetY()
	for _, tableRecords := range data {
		for _, rows := range tableRecords {
			if colValues, ok := rows.(map[string]interface{}); ok {
				maxHt := lineHt
				// calculate the height of the cell
				var cellList []cellType
				for colJ, orderKey := range orderKeys {
					cell.str = fmt.Sprintf("%v", colValues[orderKey])
					cell.list = pdf.SplitLines([]byte(cell.str), columnWidths[colJ]-cellGap-cellGap)
					cell.ht = (float64(len(cell.list)) * lineHt) + 3
					if cell.ht > maxHt {
						maxHt = cell.ht
					}
					cellList = append(cellList, cell)
				}
				// bing values to the cell
				x := marginH
				for colJ := range orderKeys {
					cell = cellList[colJ]
					cellY := 3 + y + cellGap + (maxHt-cell.ht)/2
					if cellY > 240 {
						//log.WithFields(log.Fields{"events": len("test")}).Info("new pdf page", "cellY-", cellY, "maxHt-cell.ht- ", maxHt, cell.ht, x, y)
						pdf.AddPage()
						// Reset positions
						x = 10.0
						y = 20.0
						cellY = y + cellGap + (maxHt-cell.ht)/2
					}
					pdf.Rect(x, y, columnWidths[colJ], maxHt+cellGap+cellGap, "D")

					for splitJ := 0; splitJ < len(cell.list); splitJ++ {
						pdf.SetXY(x+cellGap, cellY)
						pdf.CellFormat(columnWidths[colJ]-cellGap-cellGap, 0, string(cell.list[splitJ]), "", 0,
							alignList[1], false, 0, "")
						cellY += lineHt
					}
					x += columnWidths[colJ]
				}
				y += maxHt + cellGap + cellGap
			}
		}
	}

	// Save to file
	err := pdf.OutputFileAndClose(pdfFileName)
	if err != nil {
		fmt.Println("Error creating PDF:", err)
		log.WithFields(log.Fields{"events": len("test")}).Info("Error creating PDF:", err)
	}
	log.WithFields(log.Fields{"events": len("test")}).Info("Pdf created", pdfFileName)
}

func configureHeaderColumn(sendEmailInfo config.SendEmailInfo, data []map[string]interface{}) ([]string, []float64) {
	var orderKeys []string
	var columnWidths []float64
	if len(sendEmailInfo.Columns) == 0 {
		for _, headerCol := range data[0] {
			if innerMap, ok := headerCol.(map[string]interface{}); ok {
				for colKeys := range innerMap {
					if colKeys != "LaunchTime" {
						orderKeys = append(orderKeys, colKeys)
					}
				}
			}
		}
	} else {
		for _, cols := range sendEmailInfo.Columns {
			if cols != "LaunchTime" {
				orderKeys = append(orderKeys, cols)
			}
		}
	}

	sort.Strings(orderKeys)

	var numCols = len(orderKeys)

	for range orderKeys {

		columnWidths = append(columnWidths, float64(195/numCols))
	}

	return orderKeys, columnWidths
}

/* // Function filterExecutnData removed as it was unused
func filterExecutnData(sendEmailInfo config.SendEmailInfo, data []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if len(sendEmailInfo.Filters) != 0 {
		for _, tableRecords := range data {
			for _, rows := range tableRecords {
				if colValues, ok := rows.(map[string]interface{}); ok {
					for key, filter := range sendEmailInfo.Filters {
						if filter == colValues[key] {
							result = append(result, colValues)
						}
					}
				}
			}
		}
		return result
	}
	return data
}
*/
