package generator

import (
	"gonetlib/netlogger"
	"os"
	"path"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(path string) bool {
	//.go 파일 파싱
	idlParser := NewIDLParser()

	if idlParser.Parse(path) == false {
		return false
	}

	//id 생성
	g.GenerateID(idlParser)

	//protocol 세팅 인터페이스 생성

	//task 생성

	return true
}

func (g *Generator) GenerateID(parser *IDLParser) bool {
	for pkgName, structure := range parser.IdlMap {
		pkgPath := path.Join(".", pkgName)
		err := os.Mkdir(pkgPath, os.ModeDir)
		if err != nil && os.IsExist(err) == false {
			netlogger.GetLogger().Error("Failed to mkdir for generate id | path[%s]", pkgPath)
			return false
		}

		for structureName, _ := range structure {

		}
	}

	return true
}
