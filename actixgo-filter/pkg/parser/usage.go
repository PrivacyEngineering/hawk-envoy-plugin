package parser

import (
	"bytes"
	"fmt"
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/format"
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/parser/body"
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/parser/headers"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"strings"
)

type Field struct {
	Format    string
	Namespace string
	Path      string
	Count     int
}

const templateUsage = `
{
  "id": "{{id}}",
  "metadata": {
    "side": "{{side}}",
    "phase": "{{phase}}",
    "timestamp": "{{timestamp}}"
  },
  "endpoint": {
    "id": "string",
    "host": "{{host}}",
    "protocol": "{{protocol}}",
    "method": "{{method}}",
    "path": "{{path}}"
  },
  "initiator": {
    "host": "{{initiator}}"
  },
  "fields": [ {{fields}} ],
  "tags": { {{tags}} }
}
`

const templateField = `
{ "format": "{{format}}", "namespace": "{{namespace}}", "path": "{{path}}", "count": {{count}} }`

func Transform(reqBody, resBody []byte, reqHeaders, resHeaders [][2]string) []byte {
	id := string(getProperty([]string{"request", "id"}))

	reqUsage := build(reqHeaders, reqBody, id, "REQUEST")
	resUsage := build(resHeaders, resBody, id, "RESPONSE")

	usages := "[" + reqUsage + "," + resUsage + "]"
	proxywasm.LogWarnf("GO-transform: payload: %v", usages)
	return []byte(usages)
}

func build(header [][2]string, payload []byte, id string, phase string) string {
	var fs []Field
	paths := headers.ParseHeader(header)
	ff := buildFields(paths, "properties", "header")
	fs = append(fs, ff...)

	if paths, err := body.ParseBody(payload); err == nil {
		ff := buildFields(paths, "json", "body")
		fs = append(fs, ff...)
	}

	fields := buildManyFields(fs)
	return buildUsage(id, phase, fields)
}

func buildUsage(id string, phase string, fields string) string {
	// time library is not supported by tinygo neither wasm proxy
	//timestamp := time.Now().Format(isoFormat)
	//timestamp := "2020-01-02T00:00:00Z"
	timestamp := fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"request", "time"})))

	usage := strings.Replace(templateUsage, "{{id}}", id, 1)
	usage = strings.Replace(usage, "{{side}}", "CLIENT", 1)
	usage = strings.Replace(usage, "{{phase}}", phase, 1)
	usage = strings.Replace(usage, "{{timestamp}}", timestamp, 1)
	usage = strings.Replace(usage, "{{tags}}", buildTags(), 1)
	usage = strings.Replace(usage, "{{initiator}}", string(getProperty([]string{"source", "address"})), 1)
	usage = strings.Replace(usage, "{{host}}", string(getProperty([]string{"request", "host"})), 1)
	usage = strings.Replace(usage, "{{protocol}}", string(getProperty([]string{"request", "protocol"})), 1)
	usage = strings.Replace(usage, "{{method}}", string(getProperty([]string{"request", "method"})), 1)
	usage = strings.Replace(usage, "{{path}}", string(getProperty([]string{"request", "path"})), 1)

	return strings.Replace(usage, "{{fields}}", fields, 1)
}

func buildTags() string {
	// parameters:
	buff := bytes.NewBufferString("")

	// request attributes
	buff.WriteString(`"request.time": "`)
	buff.WriteString(fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"request", "time"}))))
	buff.WriteString(`",`)

	buff.WriteString(`"request.path": "`)
	buff.Write(getProperty([]string{"request", "path"}))
	buff.WriteString(`",`)

	buff.WriteString(`"request.url_path": "`)
	buff.Write(getProperty([]string{"request", "url_path"}))
	buff.WriteString(`",`)

	buff.WriteString(`"request.host": "`)
	buff.Write(getProperty([]string{"request", "host"}))
	buff.WriteString(`",`)

	buff.WriteString(`"request.method": "`)
	buff.Write(getProperty([]string{"request", "method"}))
	buff.WriteString(`",`)

	buff.WriteString(`"request.id": "`)
	buff.Write(getProperty([]string{"request", "id"}))
	buff.WriteString(`",`)

	buff.WriteString(`"request.protocol": "`)
	buff.Write(getProperty([]string{"request", "protocol"}))
	buff.WriteString(`",`)

	//response attributes
	buff.WriteString(`"response.code": "`)
	buff.WriteString(fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"response", "code"}))))
	buff.WriteString(`",`)

	buff.WriteString(`"response.code_details": "`)
	buff.Write(getProperty([]string{"response", "code_details"}))
	buff.WriteString(`",`)

	//connection attributes
	buff.WriteString(`"source.address": "`)
	buff.Write(getProperty([]string{"source", "address"}))
	buff.WriteString(`",`)

	buff.WriteString(`"source.port": "`)
	buff.WriteString(fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"source", "port"}))))
	buff.WriteString(`",`)

	buff.WriteString(`"destination.address": "`)
	buff.Write(getProperty([]string{"destination", "address"}))
	buff.WriteString(`",`)

	buff.WriteString(`"destination.port": "`)
	buff.WriteString(fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"destination", "port"}))))
	buff.WriteString(`",`)

	//upstream attributes
	buff.WriteString(`"upstream.address": "`)
	buff.Write(getProperty([]string{"upstream", "address"}))
	buff.WriteString(`",`)

	buff.WriteString(`"upstream.port": "`)
	buff.WriteString(fmt.Sprintf("%f", format.FmtNumber(getProperty([]string{"upstream", "port"}))))
	buff.WriteString(`",`)

	buff.WriteString(`"upstream.local_address": "`)
	buff.Write(getProperty([]string{"upstream", "local_address"}))
	buff.WriteString(`",`)

	//wasm attributes
	buff.WriteString(`"plugin_name": "`)
	buff.Write(getProperty([]string{"plugin_name"}))
	buff.WriteString(`",`)

	buff.WriteString(`"plugin_root_id": "`)
	buff.Write(getProperty([]string{"plugin_root_id"}))
	buff.WriteString(`",`)

	buff.WriteString(`"plugin_vm_id": "`)
	buff.Write(getProperty([]string{"plugin_vm_id"}))
	buff.WriteString(`",`)

	buff.WriteString(`"cluster_name": "`)
	buff.Write(getProperty([]string{"cluster_name"}))
	buff.WriteString(`"`)

	//buff.WriteString(`"listener_direction": "`)
	//buff.Write(getProperty([]string{"listener_direction"}))
	//buff.WriteString(`"`)

	return buff.String()
}

func buildManyFields(fs []Field) string {
	buff := bytes.NewBufferString("")
	for i, f := range fs {
		buff.WriteString(buildField(f))
		if i != len(fs)-1 {
			buff.WriteString(",")
		}
	}
	return buff.String()
}

func buildField(f Field) string {
	template := strings.Replace(templateField, "{{format}}", f.Format, 1)
	template = strings.Replace(template, "{{namespace}}", f.Namespace, 1)
	template = strings.Replace(template, "{{path}}", f.Path, 1)
	return strings.Replace(template, "{{count}}", fmt.Sprintf("%d", f.Count), 1)
}

func buildFields(paths map[string]int, format, ns string) []Field {
	if len(paths) == 0 {
		return nil
	}
	var fields []Field
	for key, value := range paths {
		field := Field{
			Format:    format,
			Namespace: ns,
			Path:      key,
			Count:     value,
		}
		fields = append(fields, field)
	}
	return fields
}
