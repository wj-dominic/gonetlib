package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

type Generator struct {
	IdlMap map[string]map[string]ast.StructType
}

func NewGenerator() *Generator {
	return &Generator{
		IdlMap: make(map[string]map[string]ast.StructType),
	}
}

func (g *Generator) Generate(path string) bool {
	//.go 파일 파싱
	if g.parse(path) == false {
		return false
	}

	//id 생성

	//protocol 세팅 인터페이스 생성

	//task 생성

	return true
}

func (g *Generator) parse(path string) bool {
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
								if _, exist := g.IdlMap[file.Name.Name]; exist == false {
									g.IdlMap[file.Name.Name] = make(map[string]ast.StructType)
								}

								structType := typeSpec.Type.(*ast.StructType)
								if _, exist := g.IdlMap[file.Name.Name][typeSpec.Name.Name]; exist == true {
									continue
								}

								g.IdlMap[file.Name.Name][typeSpec.Name.Name] = *structType

								//g.IdlMap[typeSpec.Name.Name] = *structType
								// for _, field := range structType.Fields.List {
								// 	var fieldType string

								// 	switch field.Type.(type) {
								// 	case *ast.ArrayType:
								// 		fieldArray := field.Type.(*ast.ArrayType)
								// 		fieldType = fieldArray.Elt.(*ast.Ident).Name
								// 		break

								// 	case *ast.Ident:
								// 		fieldIdent := field.Type.(*ast.Ident)
								// 		fieldType = fieldIdent.Name
								// 		break

								// 	case *ast.MapType:
								// 		fieldMap := field.Type.(*ast.MapType)
								// 		fieldType = fieldMap.Key.(*ast.Ident).Name
								// 	}

								// 	for _, name := range field.Names {
								// 		fmt.Printf("\tField: name=%s type=%s\n", name.Name, fieldType)
								// 	}

								// }

							}
						}
					}
				}
			}
		}
	}

	return true
}
