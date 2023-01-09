package generator

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const MaxID int = 100

type Parser interface {
	Parse(path string) bool
	Get() map[string]*IDLFile
}

type Generator struct {
	parser     Parser
	GenRootDir string
}

func NewGenerator(parser Parser) *Generator {
	return &Generator{
		parser:     parser,
		GenRootDir: "",
	}
}

func (g *Generator) Generate(idlPath string) bool {
	if g.checkValid(idlPath) == false {
		return false
	}

	g.GenRootDir = idlPath
	if g.deleteGenFiles() == false {
		return false
	}

	if g.parser.Parse(idlPath) == false {
		return false
	}

	if g.makePacketID() == false {
		return false
	}

	if g.makePacketConstructor() == false {
		return false
	}

	if g.makeTask() == false {
		return false
	}

	return true
}

func (g *Generator) checkValid(idlPath string) bool {
	_, err := os.Stat(idlPath)
	if os.IsNotExist(err) {
		log.Fatalf("Path is not exist | path[%s]", idlPath)
		return false
	}

	/*
		캐시 파일 검증

		알고리즘
		- 캐시 파일이 없으면 그냥 리턴
		- 캐시 파일이 있으면
		ㄴ 캐시 파일과 프로토콜 파일과 비교 (추가, 제거된 프로토콜이 있는지 확인, 추가, 제거된 멤버 변수가 있는지 확인)
		ㄴ 추가된 프로토콜이 있으면 신규 Task 생성 대상으로 체크
		ㄴ 이후 외부에서는 생성 대상으로 체크된 놈들만 Task ID 생성할 수 있도록 수정

	*/

	files, err := os.ReadDir(idlPath)
	if err != nil {
		return false
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "gen") {
			fmt.Printf("file.Name(): %v\n", file.Name())
		}
	}

	return true
}

func (g *Generator) makePacketID() bool {
	var source strings.Builder

	index := 0

	IdlMap := g.parser.Get()

	for fileName, idlFile := range IdlMap {
		source.WriteString("//generated...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", idlFile.PakageName))

		packetID := index * MaxID
		source.WriteString("const (\n")

		for _, idl := range idlFile.IDLs {
			source.WriteString(fmt.Sprintf("\t%s_ID = %d\n", idl.Name, packetID))
			packetID++
		}

		source.WriteString(")\n")
		source.WriteString("\n")
		index++

		genFilePath := path.Join(g.GenRootDir, fmt.Sprintf("gen_%s_ID.go", fileName))
		if g.makeGenFile(genFilePath, source.String()) == false {
			return false
		}

		source.Reset()
	}

	return true
}

func (g *Generator) makePacketConstructor() bool {
	var source strings.Builder

	IdlMap := g.parser.Get()

	for fileName, idlFile := range IdlMap {
		source.WriteString("//generated...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", idlFile.PakageName))

		for _, idl := range idlFile.IDLs {
			source.WriteString(fmt.Sprintf("func NEW_%s() (uint16, %s) {\n", idl.Name, idl.Name))
			source.WriteString(fmt.Sprintf("\treturn %s_ID ,%s {\n", idl.Name, idl.Name))
			source.WriteString("\t}\n")
			source.WriteString("}\n")
			source.WriteString("\n")
		}

		genFilePath := path.Join(g.GenRootDir, fmt.Sprintf("gen_%s_Constructor.go", fileName))
		if g.makeGenFile(genFilePath, source.String()) == false {
			return false
		}

		source.Reset()
	}

	return true
}

func (g *Generator) makeTask() bool {
	var source strings.Builder

	IdlMap := g.parser.Get()

	for fileName, idlFile := range IdlMap {
		source.WriteString("//generator...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", idlFile.PakageName))

		for _, idl := range idlFile.IDLs {
			taskName := fmt.Sprintf("%s_TASK", idl.Name)
			source.WriteString(fmt.Sprintf("type %s struct {\n", taskName))
			source.WriteString(fmt.Sprintf("\t%s\n", idl.Name))
			source.WriteString("}\n\n")

			source.WriteString(fmt.Sprintf("func (t *%s) Run() bool {\n", taskName))
			source.WriteString("\t//do something what you want...\n")
			source.WriteString("\treturn true\n")
			source.WriteString("}\n\n")
		}

		genFilePath := path.Join(g.GenRootDir, fmt.Sprintf("gen_%s_Task.go", fileName))
		if g.makeGenFile(genFilePath, source.String()) == false {
			return false
		}

		source.Reset()
	}

	//make register
	for fileName, idlFile := range IdlMap {
		source.WriteString("//generator...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", idlFile.PakageName))

		source.WriteString("import (\n")
		source.WriteString("\t\"fmt\"\n")
		source.WriteString("\t\"gonetlib/task\"\n")
		source.WriteString("\t\"gonetlib/message\"\n")
		source.WriteString(")\n\n")

		for _, idl := range idlFile.IDLs {
			//make register
			registerName := fmt.Sprintf("%s_TASK_REGISTER", idl.Name)

			source.WriteString(fmt.Sprintf("type %s struct {\n", registerName))
			source.WriteString("}\n\n")

			source.WriteString(fmt.Sprintf("func(t *%s) CreateTask(packet *message.Message) task.ITask {\n", registerName))
			source.WriteString("\tif packet == nil {\n")
			source.WriteString("\t\tfmt.Println(\"packet is nullpt\")\n")
			source.WriteString("\t\treturn nil\n")
			source.WriteString("\t}\n\n")

			source.WriteString(fmt.Sprintf("\tvar newTask %s_TASK\n\n", idl.Name))
			source.WriteString("\tpacket.Pop(&newTask)\n\n")
			source.WriteString("\treturn &newTask\n")
			source.WriteString("}\n\n")

			source.WriteString(fmt.Sprintf("func Add_%s() {\n", registerName))
			source.WriteString(fmt.Sprintf("\ttask.AddTaskRegister(%s_ID, &%s{})\n", idl.Name, registerName))
			source.WriteString("}\n\n")
		}

		genFilePath := path.Join(g.GenRootDir, fmt.Sprintf("gen_%s_TaskRegister.go", fileName))
		if g.makeGenFile(genFilePath, source.String()) == false {
			return false
		}

		source.Reset()
	}

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

func (g *Generator) deleteGenFiles() bool {
	files, err := os.ReadDir(g.GenRootDir)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "gen_") == false {
			continue
		}

		filePath := path.Join(g.GenRootDir, file.Name())

		if err = os.Remove(filePath); err != nil {
			fmt.Printf("err: %v\n", err)
			return false
		}
	}

	return true
}
