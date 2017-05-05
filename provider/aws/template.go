package aws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"html/template"
	"strings"
)

func formationHelpers() template.FuncMap {
	return template.FuncMap{
		"lower": func(s string) string {
			return strings.ToLower(s)
		},
		"priority": func(app, service string) uint32 {
			return crc32.ChecksumIEEE([]byte(fmt.Sprintf("%s-%s", app, service))) % 50000
		},
		"resource": func(s string) string {
			return upperName(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"upper": func(s string) string {
			return strings.ToUpper(s)
		},
	}
}

func formationTemplate(name string, data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	tn := fmt.Sprintf("%s.json.tmpl", name)
	tf := fmt.Sprintf("provider/aws/formation/%s", tn)

	t, err := template.New(tn).Funcs(formationHelpers()).ParseFiles(tf)
	if err != nil {
		return nil, err
	}

	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	var v interface{}

	fmt.Printf("buf.String() = %+v\n", buf.String())
	if err := json.Unmarshal(buf.Bytes(), &v); err != nil {
		switch t := err.(type) {
		case *json.SyntaxError:
			return nil, fmt.Errorf("%s \n\n\n\n\n %s", jsonSyntaxError(t, buf.Bytes()), buf.String())
		}
		return nil, err
	}

	return json.MarshalIndent(v, "", "  ")
}

func jsonSyntaxError(err *json.SyntaxError, data []byte) error {
	start := bytes.LastIndex(data[:err.Offset], []byte("\n")) + 1
	line := bytes.Count(data[:start], []byte("\n"))
	pos := int(err.Offset) - start - 1
	ltext := strings.Split(string(data), "\n")[line]

	return fmt.Errorf("json syntax error: line %d pos %d: %s: %s", line, pos, err.Error(), ltext)
}
