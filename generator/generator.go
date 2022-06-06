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
	parser     *IDLParser
	GenRootDir string
}

func NewGenerator() *Generator {
	return &Generator{
		parser:     NewIDLParser(),
		GenRootDir: "",
	}
}

func (g *Generator) Generate(idlPath string) bool {
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

func (g *Generator) makePacketID() bool {
	if g.parser.IsParsed == false {
		return false
	}

	var source strings.Builder

	index := 0
	for fileName, packets := range g.parser.IdlMap {
		source.WriteString("//generated...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", g.parser.IdlPkgMap[fileName]))

		packetID := index * MaxID
		source.WriteString("const (\n")

		for packetName := range packets {
			source.WriteString(fmt.Sprintf("\t%s_ID = %d\n", packetName, packetID))
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

	for fileName, packets := range g.parser.IdlMap {
		source.WriteString("//generated...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", g.parser.IdlPkgMap[fileName]))

		for packetName := range packets {
			source.WriteString(fmt.Sprintf("func NEW_%s() (uint16, %s) {\n", packetName, packetName))
			source.WriteString(fmt.Sprintf("\treturn %s_ID ,%s {\n", packetName, packetName))
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

	//make task
	for fileName, packets := range g.parser.IdlMap {
		source.WriteString("//generator...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", g.parser.IdlPkgMap[fileName]))

		for packetName := range packets {
			taskName := fmt.Sprintf("%s_TASK", packetName)
			source.WriteString(fmt.Sprintf("type %s struct {\n", taskName))
			source.WriteString(fmt.Sprintf("\t%s\n", packetName))
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
	for fileName, packets := range g.parser.IdlMap {
		source.WriteString("//generator...\n\n")

		source.WriteString(fmt.Sprintf("package %s\n\n", g.parser.IdlPkgMap[fileName]))

		source.WriteString("import (\n")
		source.WriteString("\t\"fmt\"\n")
		source.WriteString("\t\"gonetlib/task\"\n")
		source.WriteString("\t\"gonetlib/message\"\n")
		source.WriteString(")\n\n")

		for packetName := range packets {
			//make register
			registerName := fmt.Sprintf("%s_TASK_REGISTER", packetName)

			source.WriteString(fmt.Sprintf("type %s struct {\n", registerName))
			source.WriteString("}\n\n")

			source.WriteString(fmt.Sprintf("func(t *%s) CreateTask(packet *message.Message) task.ITask {\n", registerName))
			source.WriteString("\tif packet == nil {\n")
			source.WriteString("\t\tfmt.Println(\"packet is nullpt\")\n")
			source.WriteString("\t\treturn nil\n")
			source.WriteString("\t}\n\n")

			source.WriteString(fmt.Sprintf("\tvar newTask %s_TASK\n\n", packetName))
			source.WriteString("\tpacket.Pop(&newTask)\n\n")
			source.WriteString("\treturn &newTask\n")
			source.WriteString("}\n\n")

			source.WriteString(fmt.Sprintf("func Add_%s() {\n", registerName))
			source.WriteString(fmt.Sprintf("\ttask.AddTaskRegister(%s_ID, &%s{})\n", packetName, registerName))
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
