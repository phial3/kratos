package main

import (
    "encoding/json"
    "fmt"
    "google.golang.org/genproto/googleapis/api/annotations"
    "google.golang.org/protobuf/compiler/protogen"
    "google.golang.org/protobuf/proto"
    "google.golang.org/protobuf/types/descriptorpb"
    "os"
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
    fmt.Fprintln(os.Stderr, file.Desc.Package())

    for _, service := range file.Services {
        genService(gen, file, g, service, omitempty)
    }
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, omitempty bool) {
    if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
        g.P("//")
        fmt.Fprintln(os.Stderr, service.Desc.Options().(*descriptorpb.ServiceOptions).String())
    }
    fmt.Fprintln(os.Stderr, "generate swagger", service.Comments.Leading.String(), service.Comments.Trailing.String() )
    for _, comments := range service.Comments.LeadingDetached {
        fmt.Fprintln(os.Stderr,comments.String())
    }
    for _, method := range service.Methods {
        rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
        if rule != nil && ok {
            fmt.Fprintln(os.Stderr, rule)
        }
    }
    //service.Desc.Options().
    //// HTTP Server.
    //sd := &serviceDesc{
    //    ServiceType: service.GoName,
    //    ServiceName: string(service.Desc.FullName()),
    //    Metadata:    file.Desc.Path(),
    //}
    //for _, method := range service.Methods {
    //    if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
    //        continue
    //    }
    //    rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
    //    if rule != nil && ok {
    //        for _, bind := range rule.AdditionalBindings {
    //            sd.Methods = append(sd.Methods, buildHTTPRule(g, method, bind))
    //        }
    //        sd.Methods = append(sd.Methods, buildHTTPRule(g, method, rule))
    //    } else if !omitempty {
    //        path := fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
    //        sd.Methods = append(sd.Methods, buildMethodDesc(g, method, "POST", path))
    //    }
    //}
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