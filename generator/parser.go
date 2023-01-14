package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"gonetlib/netlogger"
	"gonetlib/util"
)

type Field struct {
	Name string
	Type string
}

type IDLStruct struct {
	Name   string
	Fields []*Field
}

type IDLFile struct {
	Name       string
	PakageName string
	IDLs       []*IDLStruct
}

type IDLParser struct {
	cache map[string]*IDLFile
}

func NewIDLParser() *IDLParser {
	return &IDLParser{
		cache: make(map[string]*IDLFile),
	}
}

func (p *IDLParser) Parse(srcPath string) bool {
	if p.isValid(srcPath) == false {
		return false
	}

	fset := token.NewFileSet()
	if fset == nil {
		netlogger.Error("Failed to make a new file set from generator\n")
		return false
	}

	pkgs, err := parser.ParseDir(fset, srcPath, nil, parser.ParseComments)
	if err != nil {
		netlogger.Error("Failed to parse dir | path[%s] err[%s]\n", srcPath, err)
		return false
	}

	p.reset()

	for pkgName, pkg := range pkgs {
		for filePath, file := range pkg.Files {
			idlFile := p.parse(file)
			if idlFile == nil {
				netlogger.Error("Failed to parse to idl file [path:%s]\n", filePath)
				continue
			}

			idlFile.Name = util.GetFileNameWithoutExt(filePath)
			idlFile.PakageName = pkgName

			p.cache[idlFile.Name] = idlFile
		}
	}

	return true
}

func (p *IDLParser) Get() map[string]*IDLFile {
	return p.cache
}

func (p *IDLParser) parse(file *ast.File) *IDLFile {
	var idlFile IDLFile
	idlFile.IDLs = make([]*IDLStruct, 0)

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok == false {
			continue
		}

		for _, spec := range genDecl.Specs {
			typespec, ok := spec.(*ast.TypeSpec)
			if ok == false {
				continue
			}

			structType, ok := typespec.Type.(*ast.StructType)
			if ok == false {
				continue
			}

			var idlStruct IDLStruct
			idlStruct.Name = typespec.Name.Name
			idlStruct.Fields = make([]*Field, 0)

			for _, field := range structType.Fields.List {
				fieldType := field.Type
				typeString := types.ExprString(fieldType)

				for _, fieldIndent := range field.Names {
					tmpField := Field{Name: fieldIndent.Name, Type: typeString}
					idlStruct.Fields = append(idlStruct.Fields, &tmpField)
				}
			}

			idlFile.IDLs = append(idlFile.IDLs, &idlStruct)
		}
	}

	return &idlFile
}

func (p *IDLParser) reset() {
	for k := range p.cache {
		delete(p.cache, k)
	}
}

func (p *IDLParser) isValid(filePath string) bool {
	return util.IsExistPath(filePath)
}
