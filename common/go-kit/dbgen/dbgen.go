package dbgen

import (
	"bytes"
	"database/sql"
	"fmt"
)

type Table struct {
	Name            string
	GoStructureName string
	Columns         []*ColumnSchema
	ExtraFields     []*ExtraField
}

type Generator struct {
	db     *sql.DB
	dbName string
	opts   Options
}

func New(db *sql.DB, options ...Options) *Generator {
	g := &Generator{
		db:   db,
		opts: defaultOptions,
	}
	if len(options) > 0 {
		g.opts.merge(options[0])
	}
	return g
}

func (g *Generator) parse() (map[string]*Table, error) {
	if err := g.db.QueryRow(`SELECT DATABASE()`).Scan(&g.dbName); err != nil {
		return nil, err
	}

	if _, err := g.db.Exec(`USE information_schema`); err != nil {
		return nil, err
	}

	rows, err := g.db.Query(queryColumnsSQL, g.dbName)
	if err != nil {
		return nil, err
	}

	tableMap := map[string]*Table{}
	for rows.Next() {
		cs := new(ColumnSchema)
		column := &cs.Column
		if err := rows.Scan(&column.TableCatalog, &column.TableSchema, &column.TableName, &column.ColumnName, &column.OrdinalPosition, &column.ColumnDefault, &column.IsNullAble, &column.DataType, &column.CharacterMaximumLength,
			&column.CharacterOctetLength, &column.NumericPrecision, &column.NumericScale, &column.DatetimePrecision, &column.CharacterSetName, &column.CollationName, &column.ColumnType, &column.ColumnKey, &column.Extra,
			&column.Privileges, &column.ColumnComment,
		); err != nil {
			return nil, err
		}
		if g.opts.IgnoreColumnFunc(*cs) {
			continue
		}

		cs.IsNullAble, cs.IsPrimaryKey = !g.opts.DisableNull && column.IsNullAble == "YES", column.ColumnKey == "PRI"
		cs.GoFieldName = g.opts.FieldNameFunc(*cs)

		if g.opts.TagsFunc != nil {
			cs.GoTags = g.opts.TagsFunc(*cs).String()
		}

		table := tableMap[cs.TableName]
		if table == nil {
			table = &Table{
				Name: cs.TableName,
			}
			table.GoStructureName = g.opts.StructureNameFunc(*table)
			tableMap[table.Name] = table
		}
		table.Columns = append(table.Columns, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableMap, nil
}

func (g *Generator) Generate() (BufferMap, error) {
	tableMap, err := g.parse()
	if err != nil {
		return nil, err
	}

	bufferMap := make(BufferMap)
	writeBuffer := func(table *Table) error {
		outputPath := ""
		if g.opts.OutputPathFunc != nil {
			outputPath = g.opts.OutputPathFunc(*table)
		}
		buf := bufferMap.get(outputPath, g.opts)
		for _, cs := range table.Columns {
			goType, importPackage, err := g.opts.FieldTypeFunc(*cs)
			if err != nil {
				return err
			}
			cs.GoFieldType = goType
			if importPackage != "" {
				buf.importPackageMap[importPackage] = true
			}
		}

		for _, f := range table.ExtraFields {
			if f.ImportPackage != "" {
				buf.importPackageMap[f.ImportPackage] = true
			}
		}

		tplParam, requiredPackages := g.opts.TemplatePreparationFunc(*table)
		for _, p := range requiredPackages {
			buf.importPackageMap[p] = true
		}
		b := &bytes.Buffer{}
		if err := g.opts.Template.Execute(b, tplParam); err != nil {
			return err
		}

		buf.structures = append(buf.structures, b.String())
		return nil
	}

	if g.opts.Tables == nil || len(g.opts.Tables) == 0 {
		// all tables
		for _, table := range tableMap {
			if g.opts.ExtraFieldFunc != nil {
				table.ExtraFields = g.opts.ExtraFieldFunc(*table)
			}
			if err := writeBuffer(table); err != nil {
				return nil, err
			}
		}
	} else {
		for _, tableName := range g.opts.Tables {
			table, ok := tableMap[tableName]
			if !ok {
				return nil, fmt.Errorf("dbgen: table not found, %s", tableName)
			}
			if g.opts.ExtraFieldFunc != nil {
				table.ExtraFields = g.opts.ExtraFieldFunc(*table)
			}

			if err := writeBuffer(table); err != nil {
				return nil, err
			}
		}
	}

	return bufferMap, nil
}
