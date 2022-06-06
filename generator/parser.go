package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

type Field struct {
	Name string
	Type string
}

type IDLParser struct {
	IdlMap    map[string]map[string][]Field
	IdlPkgMap map[string]string
	IsParsed  bool
}

func NewIDLParser() *IDLParser {
	return &IDLParser{
		IdlMap:    make(map[string]map[string][]Field),
		IdlPkgMap: make(map[string]string),
		IsParsed:  false,
	}
}

func (p *IDLParser) Parse(path string) bool {
	p.IsParsed = false

	p.clearMap()

	if len(path) <= 0 {
		return false
	}

	fset := token.NewFileSet()
	if fset == nil {
		fmt.Println("Failed to make a new file set from generator")
		return false
	}

	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Failed to parse dir | path[%s] err[%s]\n", path, err)
		return false
	}

	for pkgName, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			fileName = filepath.Base(fileName)
			fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

			for _, node := range file.Decls {
				switch node.(type) {
				case *ast.GenDecl:
					genDecl := node.(*ast.GenDecl)
					for _, spec := range genDecl.Specs {
						switch spec.(type) {
						case *ast.TypeSpec:
							typeSpec := spec.(*ast.TypeSpec)
							switch typeSpec.Type.(type) {
							case *ast.StructType:
								if _, exist := p.IdlMap[fileName]; exist == false {
									p.IdlMap[fileName] = make(map[string][]Field)
								}

								structType := typeSpec.Type.(*ast.StructType)
								if _, exist := p.IdlMap[fileName][typeSpec.Name.Name]; exist == true {
									continue
								}

								p.IdlMap[fileName][typeSpec.Name.Name] = make([]Field, len(structType.Fields.List))
								p.IdlPkgMap[fileName] = pkgName

								for index, field := range structType.Fields.List {
									fieldType := field.Type
									typeString := types.ExprString(fieldType)
									for _, fieldName := range field.Names {
										tmpField := Field{Name: fieldName.Name, Type: typeString}
										p.IdlMap[fileName][typeSpec.Name.Name][index] = tmpField
									}
								}
							}
						}
					}
				}
			}
		}
	}

	p.IsParsed = true

	return true
}

func (p *IDLParser) clearMap() {
	for k := range p.IdlMap {
		for j := range p.IdlMap[k] {
			delete(p.IdlMap[k], j)
		}

		delete(p.IdlMap, k)
	}
}
