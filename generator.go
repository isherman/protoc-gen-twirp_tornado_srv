package main

import (
	"bytes"
	"path"
	"strings"
	"text/template"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pkg/errors"
)

type generator struct{}

func (g *generator) Generate(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	resp := &plugin.CodeGeneratorResponse{}

	for _, name := range req.FileToGenerate {
		file, err := getFileDescriptor(req, name)
		if err != nil {
			return nil, err
		}

		genFile, err := g.generateFile(file)
		if err != nil {
			return nil, errors.Wrapf(err, "generating %q", name)
		}

		// Add the generated file to the response
		resp.File = append(resp.File, genFile)
	}

	return resp, nil
}

func (g *generator) generateFile(file *descriptor.FileDescriptorProto) (*plugin.CodeGeneratorResponse_File, error) {
	var err error

	buffer := &bytes.Buffer{}

	tmpl := template.New("python_file")
	tmpl, err = tmpl.Parse(pythonSrvTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing template")
	}

	presenter := &pythonSrvTemplatePresenter{
		proto:   file,
		srcInfo: file.GetSourceCodeInfo(),

		Version:        "1.0.0",
		SourceFilename: file.GetName(),
		Package:        file.GetPackage(),
	}

	err = tmpl.Execute(buffer, presenter)
	if err != nil {
		return nil, errors.Wrapf(err, "rendering template for %q", file.GetName())
	}

	resp := &plugin.CodeGeneratorResponse_File{}
	resp.Name = proto.String(getOutputFilename(file))
	resp.Content = proto.String(buffer.String())

	return resp, nil
}

// getFileDescriptor finds the FileDescriptorProto for the given filename.
// Returns a error if the descriptor could not be found.
func getFileDescriptor(req *plugin.CodeGeneratorRequest, name string) (*descriptor.FileDescriptorProto, error) {
	for _, descriptor := range req.ProtoFile {
		if descriptor.GetName() == name {
			return descriptor, nil
		}
	}

	return nil, errors.Errorf("could not find descriptor for %q", name)
}

// getOutputFilename determines what the filename should be for the generated
// code.
func getOutputFilename(file *descriptor.FileDescriptorProto) string {
	name := file.GetName()
	name = strings.TrimSuffix(name, path.Ext(name))

	return name + "_twirp_srv.py"
}
