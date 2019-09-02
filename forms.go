package forms

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
)

var durationType = reflect.TypeOf(time.Duration(0))

//MakeHTML main function to create html of form
func MakeHTML(id string, data interface{}, out io.Writer) string {

	templates := map[string]bytes.Buffer{}
	options := map[string]interface{}{
		"ID": id,
	}
	formHeader.Execute(out, options)
	processField(out, reflect.ValueOf(data), nil, templates)
	for _, xx := range templates {
		out.Write(xx.Bytes())
	}
	formFooter.Execute(out, options)
	return id
}

func processField(f io.Writer, value reflect.Value, field *formField, xTemplates map[string]bytes.Buffer) {

	switch value.Type().Kind() {
	case reflect.Ptr:
		field.IsNil = value.IsNil()
		field.ValueType = value.Type().Elem().String()
		formPtrHeader.Execute(f, field)
		field.Label = ""
		tx := bytes.Buffer{}
		tx.WriteString(fmt.Sprintf(`<script type="x-template" id="template-%s">`, field.Name))
		if field.IsNil {
			field.Value = reflect.Zero(value.Type().Elem())
			processField(&tx, reflect.Zero(value.Type().Elem()), field, xTemplates)
		} else {
			processField(&tx, value.Elem(), field, xTemplates)
		}
		tx.WriteString("</script>")
		xTemplates[field.Name] = tx
		formPtrFooter.Execute(f, field)
		return
	case reflect.Array, reflect.Slice:
		var fieldTemplate *formField
		if field != nil {
			field.Length = value.Len()
			ft := field.Copy()
			fieldTemplate = &ft
			fieldTemplate.IsArrayItem = true
			formArrayHeader.Execute(f, field)

			tx := bytes.Buffer{}
			tx.WriteString(fmt.Sprintf(`<script type="x-template" id="template-%s">`, field.Name))
			formArrayItemWrapperHeader.Execute(&tx, fieldTemplate)
			processField(&tx, reflect.Zero(reflect.TypeOf(value.Interface()).Elem()), fieldTemplate, xTemplates)
			formArrayItemWrapperFooter.Execute(&tx, fieldTemplate)
			tx.WriteString("</script>")
			xTemplates[fieldTemplate.Name] = tx
		}

		for i := 0; i < value.Len(); i++ {
			subField := field.Copy()
			subField.Expanded = field.ItemsExpanded
			subField.ID = uuid.New().String()
			subField.Name = fmt.Sprintf("%s[%d]", field.Name, i)
			subField.Value = value.Index(i).Interface()
			subField.Index = i + 1
			subField.IsArrayItem = true

			formArrayItemWrapperHeader.Execute(f, subField)
			processField(f, value.Index(i), &subField, xTemplates)
			formArrayItemWrapperFooter.Execute(f, subField)
		}

		if field != nil {
			formArrayFooter.Execute(f, field)
		}

		return
	case reflect.Bool:
		if field.Type == "" {
			field.Type = "checkbox"
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Type == "" {
			if value.Type() == durationType {
				field.Type = "string"
			} else {
				field.Type = "number"
			}
		}
	case reflect.Struct:

		if field != nil {
			formStructHeader.Execute(f, field)
		}

		for i := 0; i < value.NumField(); i++ {
			fd := parseTags(value.Type().Field(i))
			if fd.Skip {
				continue
			}

			if field != nil {
				fd.Name = field.Name + "." + fd.Name
				fd.Readonly = field.Readonly
				fd.Disabled = field.Disabled
			}
			fd.Value = value.Field(i).Interface()

			processField(f, value.Field(i), &fd, xTemplates)
		}
		if field != nil {
			formStructFooter.Execute(f, field)
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
		panic(err)
	}
	ff.ID = uuid.New().String()
	if ff.Name == "" {
		ff.Name = sf.Name
	}
	if ff.Label == "" {
		ff.Label = sf.Name
	}
	ff.ValueType = sf.Type.Kind().String()
	if sf.Type.Kind() == reflect.Slice || sf.Type.Kind() == reflect.Array {
		ff.ValueType = sf.Type.Elem().Kind().String()
	}

	return ff
}
