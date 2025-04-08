/*
	Problem: Go reflection does not support enumerating types, variables and functions of packages.

pkgreflect generates a file named pkgreflect.go in every parsed package directory.
This file contains the following maps of exported names to reflection types/values:

// Types Types

	var Types = map[string]reflect.Type{ ... }
	var Functions = map[string]reflect.Value{ ... }
	var Variables = map[string]reflect.Value{ ... }

Command line usage:

	pkgreflect --help
	pkgreflect [-notypes][-nofuncs][-novars][-unexported][-norecurs][-gofile=filename.go] [DIR_NAME]

If -norecurs is not set, then pkgreflect traverses recursively into sub-directories.
If no DIR_NAME is given, then the current directory is used as root.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

func init() {
	router = regexp.MustCompile(`@(GET|POST|PUT|DELETE|PATCH|OPTIONS|TRACE|HEAD)\s+(\S+)`)
}

const routerName = "router"

var (
	stdout        bool
	gofile        string
	notests       bool
	routerPkgPath string
	rootPkgPath   string
	router        *regexp.Regexp
)

func main() {
	flag.StringVar(&gofile, "gofile", "generated_router.go", "Name of the generated .go file")
	flag.BoolVar(&stdout, "stdout", false, "Write to stdout.")
	flag.BoolVar(&notests, "notests", true, "Don't list test related code")
	flag.Parse()

	Start(flag.Args())
	if len(flag.Args()) > 0 {
		for _, dir := range flag.Args() {
			parseDir(dir)
		}
	} else {
		parseDir(".")
	}
}

func parseDir(dir string) {
	dirFile, err := os.Open(dir)
	if er, ok := err.(*os.PathError); ok && er.Err == syscall.ENOENT {
		fmt.Println("directory ", dir, " not exists, ignore.")
		return
	}
	if err != nil {
		panic(err)
	}
	defer dirFile.Close()
	info, err := dirFile.Stat()
	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		panic("Path is not a directory: " + dir)
	}
	//fmt.Println("scan dir:", dir)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filter, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, pkg := range pkgs {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, WriteHeadFile(pkg.Name))
		for _, f := range pkg.Files {
			//ast.Print(fset, f)
			for _, d := range f.Decls {
				switch x := d.(type) {
				case *ast.FuncDecl:
					class := ""
					if x.Recv != nil {
						for _, f := range x.Recv.List {
							switch exp := f.Type.(type) {
							case *ast.StarExpr:
								switch expType := exp.X.(type) {
								case *ast.Ident:
									class = expType.Name
									break
								}
							}
						}

						if x.Doc != nil && class != "" {
							for _, d := range x.Doc.List {
								code := HandleMethodAndComment(pkg.Name, class, x.Name.String(), d.Text)
								if code != "" {
									fmt.Fprint(&buf, code)
								}
							}
						}
					}
				}
			}
		}
		fmt.Fprint(&buf, WriteTailFile())

		if stdout {
			io.Copy(os.Stdout, &buf)
		} else {
			//fmt.Println(rootPkgPath)
			routerPath, _ := filepath.Abs(rootPkgPath + string(os.PathSeparator) + routerName)
			if _, err := os.Stat(routerPath); err != nil {
				if err := os.Mkdir(routerPath, os.ModePerm); err != nil {
					panic(err)
				}
			}

			filename, _ := filepath.Abs(routerPath + string(os.PathSeparator) + gofile)
			//fmt.Println(filename)
			//filename := gofile
			newFileData := buf.Bytes()
			oldFileData, _ := ioutil.ReadFile(filename)
			if !bytes.Equal(newFileData, oldFileData) {
				err = ioutil.WriteFile(filename, newFileData, 0660)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func filter(info os.FileInfo) bool {
	name := info.Name()
	if info.IsDir() {
		return false
	}
	if name == gofile {
		return false
	}
	if filepath.Ext(name) != ".go" {
		return false
	}
	if strings.HasSuffix(name, "_test.go") && notests {
		return false
	}
	return true
}

// Start Start
func Start(args []string) string {
	/*
		root, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		path2, _ := filepath.Abs(root)
		var handlerPath string
		if strings.Contains(path1, path2) {
			handlerPath = path1
		} else {
			handlerPath, _ = filepath.Abs(path2 + string(os.PathSeparator) + path1)
		}
	*/
	handlerPath, _ := filepath.Abs(args[0])
	//fmt.Println(handlerPath)
	rootPkgPath = getParentDirectory(handlerPath)
	//fmt.Println(rootPkgPath)
	rp, _ := filepath.Abs(handlerPath)
	firstGoPath := strings.Split(os.Getenv("ROOTPATH"), string(os.PathListSeparator))[0]
	tmp0, _ := filepath.Abs(filepath.Join(firstGoPath, "/src"))
	routerPkgPath = rp[len(tmp0)+1:]
	//fmt.Println(routerPkgPath)
	routerPkgPath = strings.Replace(routerPkgPath, "\\", "/", -1)
	fmt.Println(firstGoPath, routerPkgPath)

	return ``
}

// WriteHeadFile WriteHeadFile
func WriteHeadFile(name string) string {
	return `` +
		`// This code is generated, DO NOT EDIT.

package router

import "github.com/gin-gonic/gin"
import "` + routerPkgPath + `"

// UseRoutersGenerated UseRoutersGenerated
func UseRoutersGenerated(server *gin.Engine) {
`
}

// WriteTailFile WriteTailFile
func WriteTailFile() string {
	return `
}
	
// This code is generated, DO NOT EDIT.	
`
}

// map[pkg+"~"+class]
var createdHandler = map[string]bool{}

// HandleMethodAndComment HandleMethodAndComment
func HandleMethodAndComment(pkg, class, funcName, comment string) string {
	if len(funcName) == 0 || len(comment) == 0 {
		return ""
	}
	method, rule := parseComment(&comment)
	if len(method) == 0 || len(rule) == 0 || len(class) == 0 {
		return ""
	}
	result := ""
	if key := fmt.Sprintf("%s~%s", pkg, class); !createdHandler[key] {
		result = fmt.Sprintf("%s\t%s := handler.Create%s()\n", result, class, class)
		createdHandler[key] = true
	}
	return fmt.Sprintf("%s\tserver.%s(\"%s\", %s.%s)\n", result, method, rule, class, funcName)
}

func parseComment(commnet *string) ([]byte, []byte) {
	var method, rule []byte = nil, nil
	if route := router.FindSubmatch([]byte(*commnet)); len(route) > 0 {
		method, rule = route[1], route[2]
	}
	return method, rule
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, string(os.PathSeparator)))
}
