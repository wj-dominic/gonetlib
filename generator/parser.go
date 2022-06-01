package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

type Field struct {
	Name string
	Type string
}

type IDLParser struct {
	IdlMap map[string]map[string][]Field
}

func NewIDLParser() *IDLParser {
	return &IDLParser{
		IdlMap: make(map[string]map[string][]Field),
	}
}

func (p *IDLParser) Parse(path string) bool {
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

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
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
								if _, exist := p.IdlMap[file.Name.Name]; exist == false {
									p.IdlMap[file.Name.Name] = make(map[string][]Field)
								}

								structType := typeSpec.Type.(*ast.StructType)
								if _, exist := p.IdlMap[file.Name.Name][typeSpec.Name.Name]; exist == true {
									continue
								}

								p.IdlMap[file.Name.Name][typeSpec.Name.Name] = make([]Field, len(structType.Fields.List))

								for index, field := range structType.Fields.List {
									fieldType := field.Type
									typeString := types.ExprString(fieldType)
									for _, fieldName := range field.Names {
										tmpField := Field{Name: fieldName.Name, Type: typeString}
										p.IdlMap[file.Name.Name][typeSpec.Name.Name][index] = tmpField
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return true
}
