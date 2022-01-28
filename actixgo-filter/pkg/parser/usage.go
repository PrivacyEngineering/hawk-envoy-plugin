package parser

import (
	"bytes"
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/logger"
	"github.com/TUB-CNPE-TB/rust-envoy-proxy/actixgo-filter/pkg/parser/body"
	"strings"
)

type Side string
type Phase string

const (
	CLIENT Side = "CLIENT"
	SERVER Side = "SERVER"

	REQUEST  Phase = "REQUEST"
	RESPONSE Phase = "RESPONSE"
)

type Field struct {
	Format    string
	Namespace string
	Path      string
	Count     int
}

const templateActixCollector = `
{
	"from": "actix-go-wasm-filter",
	"content": {{content}}
}
`

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
  "fields": [ {{fields}} ]
}
`

const templateField = `
{ "format": "{{format}}", "namespace": "{{namespace}}", "path": "{{path}}", "count": {{count}} }`

func Transform(reqBody, resBody []byte, reqHeaders, resHeaders [][2]string) []byte {
	var fs []Field
	//paths := headers.ParseHeader(reqHeaders)
	//ff := buildFields(paths, "properties", "header")
	//fs = append(fs, ff...)
	//
	//paths = headers.ParseHeader(resHeaders)
	//ff = buildFields(paths, "properties", "header")
	//fs = append(fs, ff...)

	if paths, err := body.ParseBody(reqBody); err == nil {
		ff := buildFields(paths, "json", "body")
		fs = append(fs, ff...)
	}
	if paths, err := body.ParseBody(resBody); err == nil {
		ff := buildFields(paths, "json", "body")
		fs = append(fs, ff...)
	}

	fields := buildManyFields(fs)
	content := buildUsage("RESPONSE", fields)
	last := buildActix(content)

	logger.WriteLog("To be sent. %s", last)

	return []byte(last)
}

func buildUsage(phase, fields string) string {
	usage := strings.Replace(templateUsage, "{{phase}}", phase, 1)
	return strings.Replace(usage, "{{fields}}", fields, 1)
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
	return strings.Replace(template, "{{count}}", "11", 1)
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

func buildActix(content string) string {
	return strings.Replace(templateActixCollector, "{{content}}", content, 1)
}
