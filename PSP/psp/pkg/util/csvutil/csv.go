package csvutil

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yuansuan/ticp/PSP/psp/internal/common"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/strutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/util/timeutil"
	"github.com/yuansuan/ticp/PSP/psp/pkg/xtype"
)

type CsvHeaderEntity struct {
	// 表头名
	Name string
	// 字段名
	Column string
	// 转换器,可用于枚举等字段的转化
	Converter func(interface{}) string
}

func CSVContentWithTab(content string) string {
	return fmt.Sprintf("%v%v%v", common.Tab, content, common.Tab)
}

func CSVFormatTime(t time.Time, format, defaultReturn string) string {
	timeStr := timeutil.FormatTime(t, format)
	if strutil.IsEmpty(timeStr) {
		return CSVContentWithTab(defaultReturn)
	}

	return CSVContentWithTab(timeStr)
}

type ExportCSVFileInfo struct {
	CSVFileName string
	CSVHeaders  []string
	FillCSVData func(w *csv.Writer) error
}

func ExportCSVFilesToZip(ctx *gin.Context, zipFileName string, infos []*ExportCSVFileInfo) error {
	exportFileName := fmt.Sprintf("%v-%v", url.QueryEscape(zipFileName), time.Now().Format(common.DateOnly))
	disposition := fmt.Sprintf("attachment; filename=%s.zip", exportFileName)
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", disposition)

	buf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, v := range infos {
		csvFile, err := zipWriter.Create(v.CSVFileName + ".csv.gz")
		if err != nil {
			return err
		}
		err = internalExportCSVFile(ctx, csvFile, v)
		if err != nil {
			return err
		}
	}

	err := zipWriter.Close()
	if err != nil {
		return err
	}
	_, err = ctx.Writer.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func ExportCSVFile(ctx *gin.Context, info *ExportCSVFileInfo) error {
	return internalExportCSVFile(ctx, nil, info)
}

func internalExportCSVFile(ctx *gin.Context, writer io.Writer, info *ExportCSVFileInfo) error {
	if writer == nil {
		exportFileName := fmt.Sprintf("%v-%v", url.QueryEscape(info.CSVFileName), time.Now().Format(common.DateOnly))
		disposition := fmt.Sprintf("attachment; filename=%s.csv.gz", exportFileName)
		ctx.Header("Content-Type", "application/x-gzip")
		ctx.Header("Content-Disposition", disposition)
		writer = ctx.Writer
	}

	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	// Write BOM directly to the gzip writer
	bom := []byte{0xEF, 0xBB, 0xBF}
	_, err := gzWriter.Write(bom)
	if err != nil {
		return err
	}

	w := csv.NewWriter(gzWriter)
	defer w.Flush()

	_ = w.Write(info.CSVHeaders)
	err = info.FillCSVData(w)
	if err != nil {
		return err
	}

	return nil
}

// SimpleExportCsv 导出csv
func SimpleExportCsv(ctx *gin.Context, header []CsvHeaderEntity, data []interface{}, csvFleName string) error {
	setHttpHeader(ctx, csvFleName)

	gzWriter := gzip.NewWriter(ctx.Writer)
	defer gzWriter.Close()
	// Write BOM directly to the gzip writer
	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := gzWriter.Write(bom); err != nil {
		return err
	}

	w := csv.NewWriter(gzWriter)
	defer w.Flush()

	// 插入表头
	if err := insertHeader(w, header); err != nil {
		return err
	}

	// 插入数据
	if err := insertData(w, header, data); err != nil {
		return err
	}

	return nil
}

func insertData(w *csv.Writer, header []CsvHeaderEntity, data []interface{}) error {
	for _, d := range data {
		v := reflect.ValueOf(d)
		rowData := make([]string, 0, len(header))

		for _, head := range header {
			val := v.FieldByName(head.Column).Interface()
			//// 零值处理
			//if reflect.DeepEqual(val, reflect.Zero(reflect.TypeOf(val)).Interface()) {
			//	rowData = append(rowData, common.Bar)
			//	continue
			//}
			if head.Converter != nil {
				rowData = append(rowData, head.Converter(val))
			} else {
				rowData = append(rowData, fmt.Sprintf("%v", val))
			}
		}

		err := w.Write(rowData)
		if err != nil {
			return err
		}

	}
	w.Flush()
	return nil
}

// LargeDataExportCsv 导出csv
func LargeDataExportCsv(ctx *gin.Context, header []CsvHeaderEntity, getDate func(page *xtype.Page) ([]interface{}, error), csvFleName string) error {
	setHttpHeader(ctx, csvFleName)

	gzWriter := gzip.NewWriter(ctx.Writer)
	defer gzWriter.Close()
	// Write BOM directly to the gzip writer
	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := gzWriter.Write(bom); err != nil {
		return err
	}

	w := csv.NewWriter(gzWriter)
	defer w.Flush()

	// 插入表头
	if err := insertHeader(w, header); err != nil {
		return err
	}

	// 分批插入数据
	pageIndex, pageSize := common.DefaultPageIndex, common.CSVExportNumber
	for {
		data, err := getDate(&xtype.Page{
			Index: int64(pageIndex),
			Size:  int64(pageSize),
		})
		if err != nil {
			return err
		}
		if len(data) == 0 {
			break
		}

		// 塞数据
		if err := insertData(w, header, data); err != nil {
			return err
		}

		if len(data) < common.CSVExportNumber {
			break
		}
		// 继续下一页
		pageIndex++
	}

	return nil
}

func insertHeader(w *csv.Writer, header []CsvHeaderEntity) error {
	headerData := make([]string, 0, len(header))
	for _, head := range header {
		headerData = append(headerData, head.Name)
	}
	err := w.Write(headerData)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func setHttpHeader(ctx *gin.Context, csvFleName string) {
	disposition := fmt.Sprintf("attachment; filename=%s.csv.gz", url.QueryEscape(csvFleName))
	ctx.Header("Content-Type", "application/x-gzip")
	ctx.Header("Content-Disposition", disposition)
}
