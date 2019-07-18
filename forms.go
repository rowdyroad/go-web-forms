package forms

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

type iPage interface {
	Title() string
}

func MakeHTML(data interface{}, out io.Writer) bool {
	templates := map[string]bytes.Buffer{}
	out.Write([]byte(`<form method="POST">`))
	processField(out, reflect.ValueOf(data), nil, templates)
	for _, xx := range templates {
		out.Write(xx.Bytes())
	}
	out.Write([]byte(`<button type="submit" class="btn btn-primary">Submit</button></form>`))
	return true
}

func ParseForm(form url.Values, data interface{}) {
	src := reflect.ValueOf(data)
	parseField(reflect.Indirect(reflect.ValueOf(data)), src, "", &form)
}

func processField(f io.Writer, value reflect.Value, field *formField, xTemplates map[string]bytes.Buffer) {
	switch value.Type().Kind() {
	case reflect.Array, reflect.Slice:
		var fieldTemplate *formField
		if field != nil {
			ft := field.Copy()
			fieldTemplate = &ft
			fieldTemplate.ID = fmt.Sprintf("id-%s-new", field.Name)
			fieldTemplate.Name = fmt.Sprintf("name-%s-new", field.Name)
			fieldTemplate.IsArrayItem = true
			fieldTemplate.Length = value.Len()
			formArrayHeader.Execute(f, fieldTemplate)

			tx := bytes.Buffer{}
			tx.Write([]byte(fmt.Sprintf(`<script type="x-template" id="new-%s">`, fieldTemplate.Name)))
			formArrayItemWrapperHeader.Execute(&tx, fieldTemplate)
			processField(&tx, reflect.Zero(reflect.TypeOf(value.Interface()).Elem()), fieldTemplate, xTemplates)
			formArrayItemWrapperFooter.Execute(&tx, fieldTemplate)
			tx.Write([]byte(`</script>`))
			xTemplates[fieldTemplate.Name] = tx
		}

		for i := 0; i < value.Len(); i++ {
			subField := field.Copy()
			subField.ID = uuid.New().String()
			subField.Name = fmt.Sprintf("%s[%d]", field.Name, i)
			subField.Value = value.Index(i).Interface()
			subField.IsArrayItem = true

			formArrayItemWrapperHeader.Execute(f, subField)
			processField(f, value.Index(i), &subField, xTemplates)
			formArrayItemWrapperFooter.Execute(f, subField)
		}

		if field != nil {
			formArrayFooter.Execute(f, fieldTemplate)
		}

		return
	case reflect.Bool:
		if field.Type == "" {
			field.Type = "checkbox"
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Type == "" {
			field.Type = "number"
		}
	case reflect.Struct:
		indent := 0

		if field != nil {
			formStructHeader.Execute(f, field)
			indent += field.Indent + 1
		}

		for i := 0; i < value.NumField(); i++ {
			fd := parseTags(value.Type().Field(i))
			if fd.Skip {
				continue
			}

			if field != nil {
				fd.Name = field.Name + "[" + fd.Name + "]"
				fd.Readonly = field.Readonly
				fd.Disabled = field.Disabled
			}
			fd.Value = value.Field(i).Interface()
			fd.Indent = indent

			processField(f, value.Field(i), &fd, xTemplates)
		}
		return
	default:
		if field.Type == "" {
			field.Type = "text"
		}
	}
	tmpl := formInput

	switch field.Type {
	case "checkbox":
		tmpl = formCheckbox
	case "textarea":
		tmpl = formTextarea
	case "select":
		tmpl = formSelect
	}
	tmpl.Execute(f, field)
}

var reLineDelim = regexp.MustCompile(`([^\\]);`)
var reHeadSpaces = regexp.MustCompile(`^\s*`)
var reHeadSpacesML = regexp.MustCompile(`\n\s*`)

func parseTags(sf reflect.StructField) formField {
	var ff formField
	tag := sf.Tag.Get("htmlForm")
	if tag == "-" {
		ff.Skip = true
		return ff
	}
	tag = reLineDelim.ReplaceAllString(tag, "$1\n")
	tag = strings.Replace(tag, "\\;", ";", -1)
	tag = reHeadSpaces.ReplaceAllString(tag, "")
	tag = reHeadSpacesML.ReplaceAllString(tag, "\n")
	if err := yaml.Unmarshal([]byte(tag), &ff); err != nil {
		log.Println("Error on", sf.Name, sf.Tag)
		panic(err)
	}
	ff.ID = uuid.New().String()
	if ff.Name == "" {
		ff.Name = sf.Name
	}
	if ff.Label == "" {
		ff.Label = sf.Name
	}

	log.Println(sf.Name, ff)
	return ff
}

func parseField(value reflect.Value, src reflect.Value, name string, form *url.Values) {
	switch value.Type().Kind() {
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			ff := parseTags(value.Type().Field(i))
			if src.IsValid() && (ff.Skip || ff.Disabled || ff.Readonly) {
				value.Field(i).Set(src.Field(i))
				continue
			}
			if name != "" {
				ff.Name = name + "[" + ff.Name + "]"
			}
			if src.IsValid() {
				parseField(value.Field(i), src.Field(i), ff.Name, form)
			} else {
				parseField(value.Field(i), reflect.Value{}, ff.Name, form)
			}
		}
	case reflect.Array, reflect.Slice:
		i := 0
		re := regexp.MustCompile(fmt.Sprintf("%s\\[(\\d+)\\]", name))
		for key := range *form {
			if regs := re.FindStringSubmatch(key); len(regs) == 2 {
				id, _ := strconv.ParseInt(regs[1], 10, 64)
				value.Set(reflect.Append(value, reflect.Zero(value.Type().Elem())))
				if !src.IsValid() || int(id) >= src.Len() {
					parseField(reflect.Indirect(value.Index(i)), reflect.Value{}, regs[0], form)
				} else {
					parseField(reflect.Indirect(value.Index(i)), src.Index(int(id)), regs[0], form)
				}
				i++
				for kk := range *form {
					if strings.Index(kk, fmt.Sprintf("%s[%d]", name, id)) != -1 {
						form.Del(kk)
					}
				}
			}
		}

	case reflect.Bool:
		v, _ := strconv.ParseBool(form.Get(name))
		value.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(form.Get(name), 10, 64)
		value.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(form.Get(name), 10, 64)
		value.SetUint(v)
	case reflect.String:
		value.SetString(form.Get(name))
	}
}
