package dbgen

import (
	"database/sql"
	"fmt"
	"strings"
)

const (
	queryColumnsSQL = `SELECT TABLE_CATALOG, TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, ORDINAL_POSITION, COLUMN_DEFAULT, IS_NULLABLE,
DATA_TYPE, CHARACTER_MAXIMUM_LENGTH, CHARACTER_OCTET_LENGTH, NUMERIC_PRECISION, NUMERIC_SCALE, DATETIME_PRECISION, CHARACTER_SET_NAME,
COLLATION_NAME, COLUMN_TYPE, COLUMN_KEY, EXTRA, PRIVILEGES, COLUMN_COMMENT FROM COLUMNS WHERE TABLE_SCHEMA = ? ORDER BY TABLE_NAME, ORDINAL_POSITION`
)

type Column struct {
	TableCatalog           string
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int64
	ColumnDefault          sql.NullString
	IsNullAble             string
	DataType               string
	CharacterMaximumLength []uint8
	CharacterOctetLength   []uint8
	NumericPrecision       []uint8
	NumericScale           []uint8
	DatetimePrecision      []uint8
	CharacterSetName       sql.NullString
	CollationName          sql.NullString
	ColumnType             string
	ColumnKey              string
	Extra                  string
	Privileges             string
	ColumnComment          string
}

type ColumnSchema struct {
	Column
	GoFieldName  string
	GoFieldType  string
	GoTags       string
	IsNullAble   bool
	IsPrimaryKey bool
}

func (cs ColumnSchema) goFieldType() (string, string, error) {
	var (
		goType        string
		importPackage string
	)
	switch cs.DataType {
	case "char", "varchar", "text", "longtext", "json":
		if cs.IsNullAble {
			goType = "sql.NullString"
			importPackage = "database/sql"
		} else {
			goType = "string"
		}
	case "date", "time", "datetime", "timestamp":
		if cs.Column.IsNullAble == "YES" {
			goType = "*time.Time"
		} else {
			goType = "time.Time"
		}
		importPackage = "time"
	case "bit", "tinyint", "smallint", "int", "mediumint", "int64", "bigint":
		if strings.HasSuffix(cs.ColumnType, "(1)") {
			if cs.IsNullAble {
				goType = "sql.NullBool"
				importPackage = "database/sql"
			} else {
				goType = "bool"
			}
		} else {
			if cs.IsNullAble {
				goType = "sql.NullInt64"
				importPackage = "database/sql"
			} else {
				goType = "int64"
			}
		}
	case "float", "decimal", "double":
		if cs.IsNullAble {
			goType = "sql.NullFloat64"
			importPackage = "database/sql"
		} else {
			goType = "float64"
		}
	case "blob", "longblob", "binary":
		goType = "[]byte"
	default:
		return "", "", fmt.Errorf("dbgen: unsupported datatype: %s", cs.DataType)
	}

	return goType, importPackage, nil
}
