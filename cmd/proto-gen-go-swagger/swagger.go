package main

import (
    "encoding/json"
    "fmt"
    "google.golang.org/genproto/googleapis/api/annotations"
    "google.golang.org/protobuf/compiler/protogen"
    "google.golang.org/protobuf/proto"
    "os"
    "strings"
)

func generateFile(gen *protogen.Plugin, file *protogen.File, omitempty bool) *protogen.GeneratedFile {
    if len(file.Services) == 0 || (omitempty && !hasHTTPRule(file.Services)) {
        return nil
    }
    filename := file.GeneratedFilenamePrefix + ".swagger.json"
    g := gen.NewGeneratedFile(filename, file.GoImportPath)
    generateFileContent(gen, file, g, omitempty)
    return g
}

// generateFileContent generates the kratos errors definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, omitempty bool) {
    if len(file.Services) == 0 {
        return
    }
    bytes, _ := json.Marshal(*file)
    fmt.Fprintln(os.Stderr, string(bytes))
    //fmt.Fprintln(os.Stderr, file)
    s := &specification{
        OpenAPI: "3.0.0",
        Title: fmt.Sprintf("%s",file.Desc.Package()),
        Paths: map[string]map[string]path{},
    }
    //file.Messages[0].Fields[0].Comments
    for _, service := range file.Services {
        genService(gen, file, g, service, omitempty, s)
    }
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool, spec *specification) {
    //if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
    //    g.P("//")
    //    fmt.Fprintln(os.Stderr, service.Desc.Options().(*descriptorpb.ServiceOptions).String())
    //}
    fmt.Fprintln(os.Stderr, "generate swagger", service.Comments.Leading.String(), service.Comments.Trailing.String() )
    for _, comments := range service.Comments.LeadingDetached {
        fmt.Fprintln(os.Stderr,comments.String())
    }
    for _, method := range service.Methods {

        rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
        if rule != nil && ok {
            p, m := generatePath(rule)
            if n, ok := spec.Paths[p]; ok {
                n[strings.ToLower(m)] = path{Tags: []string{}}
            } else {
                spec.Paths[p] = map[string]path{strings.ToLower(m): {Tags: []string{}}}
            }
            //method.Location.SourceFile
            //method.Input.Comments.Trailing
            //fmt.Fprintln(os.Stderr, method.Input.Comments.Leading )

            for _, comments := range method.Input.Comments.LeadingDetached {
                fmt.Fprintln(os.Stderr,comments.String())
            }
        }
    }
    bytes, _ := json.Marshal(spec)
    g.Write(bytes)
}

func generatePath(rule *annotations.HttpRule) (string,string) {
    var (
        path         string
        method       string
    )
    switch pattern := rule.Pattern.(type) {
    case *annotations.HttpRule_Get:
        path = pattern.Get
        method = "GET"
    case *annotations.HttpRule_Put:
        path = pattern.Put
        method = "PUT"
    case *annotations.HttpRule_Post:
        path = pattern.Post
        method = "POST"
    case *annotations.HttpRule_Delete:
        path = pattern.Delete
        method = "DELETE"
    case *annotations.HttpRule_Patch:
        path = pattern.Patch
        method = "PATCH"
    case *annotations.HttpRule_Custom:
        path = pattern.Custom.Path
        method = pattern.Custom.Kind
    }
    return path, method
}

func parseComment([]string)  {

}
func hasHTTPRule(services []*protogen.Service) bool {
    for _, service := range services {
        for _, method := range service.Methods {
            if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
                continue
            }
            rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
            if rule != nil && ok {
                return true
            }
        }
    }
    return false
}