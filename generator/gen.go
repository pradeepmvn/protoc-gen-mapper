package generator

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

//ParentStructParam
const (
	ParentStructParam = "parent"
	OutFileName       = "mapper"
	OutDir            = "proto"
	ToMapFuncName     = "ToMap"
	FromMapFuncName   = "FromMap"
	charset           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//MapperGen is main struct that constains requ and responses for generation
type MapperGen struct {
	Request  *plugin.CodeGeneratorRequest
	Response *plugin.CodeGeneratorResponse
	// parent name
	parentStruct string
	// proto message definitions
	messages map[string]*ProtoMessageDetails
	//pkg imports to be added
	pkgImports map[string]struct{}
}

// ProtoMessageDetails with package name and message details
type ProtoMessageDetails struct {
	mtype *descriptorpb.DescriptorProto
	p     string
}

// New creates a new mapper instance
func New() *MapperGen {
	g := new(MapperGen)
	g.Request = new(plugin.CodeGeneratorRequest)
	g.Response = new(plugin.CodeGeneratorResponse)
	g.pkgImports = make(map[string]struct{})
	return g
}

// Generate generates the output for all the files we're outputting.
func (m *MapperGen) Generate() {
	params := m.Request.GetParameter()
	//process command line param and get parent struct name. It is case sensitive. Exact match
	for _, p := range strings.Split(params, ",") {
		// ignore other params for now
		if i := strings.Index(p, "="); i > 0 {
			if p[0:i] == ParentStructParam {
				m.parentStruct = p[i+1:]
			}
		}
	}
	if len(m.parentStruct) < 1 {
		log.Print("No  Parent Param found in the input. Usage : --map_out=\"parent=Product:.\" ")
		os.Exit(1)
	}

	// file: package| messageslmap:msgname|message
	// Load all structs into map by reading all proto files.
	m.messages = make(map[string]*ProtoMessageDetails)
	for _, f := range m.Request.ProtoFile {
		// open up each file and load messages into the map.
		for _, mtype := range f.MessageType {
			m.messages[mtype.GetName()] = &ProtoMessageDetails{mtype: mtype, p: f.GetPackage()}
		}
	}

	//  start the generation process
	resp := []string{}
	msg := m.messages[m.parentStruct]
	resp = append(resp, "//////////////////////////////////////////////////////////////////////////////////////////////")
	resp = append(resp, "// Code generated by protoc-gen-mapper. .")
	resp = append(resp, "// DO NOT EDIT.")
	resp = append(resp, "//////////////////////////////////////////////////////////////////////////////////////////////")
	resp = append(resp, "")
	resp = append(resp, "package "+msg.p)
	//add imports
	resp = append(resp, "import \"strconv\"")
	resp = append(resp, "import \"fmt\"")
	resp = append(resp, "")

	//storage for all feilds, non-nil
	revMapping := []string{}
	mapping := []string{}

	//process recursively
	processMessageTypes(msg, &resp, &mapping, &revMapping, m.parentStruct, m.messages, "")

	//create a map function
	resp = append(resp, "// "+ToMapFuncName+" Convert a struct into a Map")
	resp = append(resp, "func (p *"+m.parentStruct+") "+ToMapFuncName+"()  map[string]string {")
	resp = append(resp, "m := make(map[string]string)")
	resp = append(resp, strings.Join(mapping, "\n"))
	resp = append(resp, "return m")
	resp = append(resp, "}")
	// create frommap function
	resp = append(resp, "")
	resp = append(resp, "")
	resp = append(resp, "// "+FromMapFuncName+" Convert a Map into a Struct")
	resp = append(resp, "func "+FromMapFuncName+"(m map[string]string) *"+m.parentStruct+"{")
	resp = append(resp, "var p=new("+m.parentStruct+") ")
	resp = append(resp, strings.Join(revMapping, "\n"))
	resp = append(resp, "return p")
	resp = append(resp, "}")

	m.Response.File = append(m.Response.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(OutDir + "/" + OutFileName + ".gen.go"),
		Content: proto.String(strings.Join(resp, "\n")),
	})
}

func processMessageTypes(msg *ProtoMessageDetails, resp *[]string, mapping *[]string, revMapping *[]string, parent string, messagesMap map[string]*ProtoMessageDetails, val string) {
	log.Println("generating for : " + msg.mtype.GetName())
	for _, f := range msg.mtype.Field {
		switch f.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
			//recursion for message
			v := f.GetTypeName()[strings.LastIndex(f.GetTypeName(), ".")+1:]
			newVal := ""
			if len(val) > 0 {
				newVal = val + "."
			}
			// add a null check before recusion
			*mapping = append(*mapping, "if p."+strings.Title(newVal)+strings.Title(f.GetJsonName())+" != nil {")
			*revMapping = append(*revMapping, "if p."+strings.Title(newVal)+strings.Title(f.GetJsonName())+" != nil {")
			processMessageTypes(messagesMap[v], resp, mapping, revMapping, parent+"."+f.GetJsonName(), messagesMap, newVal+f.GetJsonName())
			*mapping = append(*mapping, "}")
			*revMapping = append(*revMapping, "}")

		case descriptorpb.FieldDescriptorProto_TYPE_INT32:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = strconv.Itoa(int(p."+lval+strings.Title(f.GetJsonName())+"))")
			fn := "i" + randString(1)
			//convert back to str
			*revMapping = append(*revMapping, fn+", _ := strconv.Atoi(m["+cKey+"])")
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"=int32("+fn+")")
		case descriptorpb.FieldDescriptorProto_TYPE_INT64:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = strconv.Itoa(int(p."+lval+strings.Title(f.GetJsonName())+"))")
			fn := "i" + randString(1)
			//convert back to str
			*revMapping = append(*revMapping, fn+", _ :=  strconv.ParseInt(m["+cKey+"], 10, 64)")
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"="+fn)
		case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = strconv.FormatBool(p."+lval+strings.Title(f.GetJsonName())+")")
			fn := "b" + randString(1)
			//convert back to str
			*revMapping = append(*revMapping, fn+", _ := strconv.ParseBool(m["+cKey+"])")
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"="+fn)
		case descriptorpb.FieldDescriptorProto_TYPE_FLOAT, descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = fmt.Sprintf(\"%f\", p."+lval+strings.Title(f.GetJsonName())+")")
			fn := "f" + randString(1)
			//convert back to str
			*revMapping = append(*revMapping, fn+", _ := strconv.ParseFloat(m["+cKey+"],64)")
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"="+fn)
		case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = p."+lval+strings.Title(f.GetJsonName())+".String()")
			// p.Code= Product_StatusCode(Product_StatusCode_value[m[ProductCode]])
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"= "+replaceTypeNameForEnums(f.GetTypeName())+"("+replaceTypeNameForEnums(f.GetTypeName())+"_value[m["+cKey+"]])")

		default:
			cKey := replaceDotWithCamelCase(parent) + strings.Title(f.GetJsonName())
			*resp = append(*resp, "const "+cKey+"= \""+strings.ToLower(parent)+"."+f.GetJsonName()+"\"")
			lval := ""
			if len(val) > 0 {
				lval = strings.Title(val) + "."
			}
			*mapping = append(*mapping, "m["+cKey+"] = p."+lval+strings.Title(f.GetJsonName()))
			*revMapping = append(*revMapping, "p."+lval+strings.Title(f.GetJsonName())+"= m["+cKey+"]")
		}
	}
}

// replaceDotWithCamelCase Converts the constant names to CameCase for public use
func replaceDotWithCamelCase(input string) string {
	s := strings.Split(input, ".")
	if len(s) == 0 {
		return input
	}
	out := []string{}
	for i := range s {
		out = append(out, strings.Title(s[i]))
	}
	return strings.Join(out, "")
}

// replaceTypeNameForEnums
func replaceTypeNameForEnums(input string) string {
	s := strings.Split(input, ".")
	if len(s) == 0 {
		return input
	}
	out := []string{}
	for i := 2; i < len(s); i++ {
		out = append(out, strings.Title(s[i]))
	}
	return strings.Join(out, "_")
}

func randString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
