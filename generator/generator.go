package generator

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const MaxID int = 100

type Generator struct {
	GenRootDir string
}

func NewGenerator() *Generator {
	return &Generator{
		GenRootDir: "",
	}
}

func (g *Generator) Generate(idlPath string) bool {
	//.go 파일 파싱
	idlParser := NewIDLParser()

	if idlParser.Parse(idlPath) == false {
		return false
	}

	g.GenRootDir = path.Join(idlPath, "gen")

	if g.MakePacketID(idlParser) == false {
		return false
	}

	//protocol 세팅 인터페이스 생성
	if g.MakePacketConstructor(idlParser) == false {
		return false
	}

	//task 생성
	if g.MakeTask(idlParser) == false {
		return false
	}

	return true
}

func (g *Generator) MakePacketID(parser *IDLParser) bool {
	var source strings.Builder

	source.WriteString("//generated...\n\n")
	source.WriteString("package gen_PacketID\n\n")

	index := 0
	for _, packets := range parser.IdlMap {
		packetID := index * MaxID

		source.WriteString("const (\n")

		for packetName := range packets {
			source.WriteString(fmt.Sprintf("\t%s_ID = %d\n", packetName, packetID))
			packetID++
		}

		source.WriteString(")\n")
		source.WriteString("\n")
		index++
	}

	genFilePath := path.Join(g.GenRootDir, "gen_PacketID.go")
	if g.makeGenFile(genFilePath, source.String()) == false {
		return false
	}

	return true
}

func (g *Generator) MakePacketConstructor(idlParser *IDLParser) bool {

	return true
}

func (g *Generator) MakeTask(idlParser *IDLParser) bool {
	return true
}

func (g *Generator) makeGenFile(filePath string, data string) bool {
	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, os.ModePerm)

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}

	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false
	}

	return true
}
