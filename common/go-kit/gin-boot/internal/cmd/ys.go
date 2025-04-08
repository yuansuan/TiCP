package cmd

import (
	"fmt"
	"github.com/alecthomas/kong"
	"go/importer"
	"go/types"
)

// Context ...
type Context struct {
	AppName string
}

// SyncCmd ...
type ArtisanCmd struct {
	Sync SyncCmd `cmd:"" help:"Sync model to database."`
}

type SyncCmd struct {
	//Project string `arg:"" name:"project" help:"Project name."`
	Model string `arg:"" optional:"" name:"model" help:"Model name."`
}

// Run 计划添加Model到数据库的xorm同步命令，然并卵，似乎无法实现
// 目前仅实现了动态展示Model列表
// 占位符，``` go run . artisan sync model ```，TODO
func (r *SyncCmd) Run(ctx Context) error {
	modelPath := fmt.Sprintf("%s/dao/model", ctx.AppName)
	pkg, err := importer.For("source", nil).Import(modelPath)
	if err != nil {
		return err
	}
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)
		if tn, ok := obj.Type().(*types.Named); ok {
			fmt.Printf("%v,%#v\n", tn, tn.NumMethods())
		}
	}
	//engine, err := xorm.NewEngine("mysql", "lambdacal:1234yskj@(0.0.0.0:3306)/cw?charset=utf8")
	//if err != nil {
	//	panic(err)
	//}
	//err = engine.Sync2(new(model.License))
	//if err != nil {
	//	panic(err)
	//}
	return nil
}

// Cli ...
var Cli struct {
	Artisan ArtisanCmd `cmd:"" help:"Artisan command line."`
}

func main() {
	ctx := kong.Parse(&Cli)
	fmt.Println(ctx)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{AppName: "test"})
	ctx.FatalIfErrorf(err)
}
