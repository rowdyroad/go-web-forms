package forms

import "html/template"

var formHeader = template.Must(template.New("form/header").Parse(`
	<script>
		function goWebFormsUUID() {
			return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
				var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
				return v.toString(16);
			});
		}

		var goWebFormsIndexes = {};
		function goWebFormsAddArrayItem(id, indexMax) {
			if (!goWebFormsIndexes[id]) {
				goWebFormsIndexes[id] = indexMax
			}
			var index = goWebFormsIndexes[id];
			var cnt = document.getElementById('template-' + id)
								.innerHTML
								.replace(new RegExp('template-'+id+'-id', 'g'), goWebFormsUUID())
								.replace(new RegExp('template-'+id+'-name', 'g'), id+'['+goWebFormsIndexes[id]+']')
								.replace(new RegExp('display:none;', 'g'), '');

			var e = document.createElement('div');
			e.innerHTML = cnt;

			document.getElementById('array-'+id).appendChild(e);
			goWebFormsIndexes[id] = index + 1;
		}
	</script>
	<form {{if .}}{{if .URL}}action="{{.URL}}"{{end}}{{end}} method="POST">
`))

var formFooter = template.Must(template.New("form/footer").Parse(`
		<button type="submit" class="btn btn-primary">{{if .}}{{if .SubmitBtnCaption}}{{.SubmitBtnCaption}}{{else}}Submit{{end}}{{else}}Submit{{end}}</button>
	</form>
`))

var formArrayHeader = template.Must(template.New("form/arrayHeader").Parse(`
	{{if .Label}}
		<h4>
			{{.Label}}
			{{if .Description}}<small><div>{{.Description}}</div></small>{{end}}
		</h4>
	{{end}}
	<div id="array-{{.Name}}">
`))

var formArrayFooter = template.Must(template.New("form/arrayFooter").Parse(`
	</div>
	{{if not .Readonly}}
		<div style="margin:0.4em;margin-bottom:1em">
			<input type="button" class="btn btn-secondary" value="{{if .AddBtnCaption}}{{.AddBtnCaption}}{{else}}Add{{end}}" onclick="goWebFormsAddArrayItem('{{.Name}}', {{.Length}})"/>
		</div>
	{{end}}
`))
var formArrayItemWrapperHeader = template.Must(template.New("form/arrayItemWrapperHeader").Parse(`
	<div id="item-{{.Name}}" class="form-group" style="padding-left: {{.Indent}}em">
`))

var formStructHeader = template.Must(template.New("form/structHeader").Parse(`
	{{if .IsArrayItem}}
		{{if .ItemLabel}}
			<h4>
				{{.ItemLabel}}
				{{if .Description}}<small>{{.Description}}</small>{{end}}
			</h4>
		{{end}}
	{{else}}
		{{if .Label}}
			<h4>
				{{.Label}}
				{{if .Description}}<small>{{.Description}}</small>{{end}}
			</h4>
		{{end}}
	{{end}}
`))

var formArrayItemWrapperFooter = template.Must(template.New("form/arrayItemWrapperFooter").Parse(`
	{{if not .Readonly}}
		<div style="text-align:right; padding:0.4em 0">
			<input type="button" class="btn btn-danger" onclick="javascript:document.getElementById('item-{{.Name}}').remove()" value="{{if .DeleteBtnCaption}}{{.DeleteBtnCaption}}{{else}}Delete{{end}}"/>
		</div>
	{{end}}
	</div>
`))

var formInput = template.Must(template.New("form/input").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<input type="{{.Type}}"  name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} class="form-control" id="{{.ID}}" value="{{.Value}}" placeholder="{{.Placeholder}}"/>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))

var formTextarea = template.Must(template.New("form/textarea").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<textarea class="form-control" name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} placeholder="{{.Placeholder}}" id="{{.ID}}" rows="{{.Rows}}">{{.Value}}</textarea>
	{{if not .IsArrayItem }}
		</div>
	{{end}}


	`))

var formCheckbox = template.Must(template.New("form/checkbox").Parse(`
	<div class="form-group" style="padding-left: {{.Indent}}em">
		<div class="form-check">
			<input class="form-check-input" type="{{.Type}}" name="{{.Name}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} {{if .Value}}checked{{end}}  class="form-control" id="{{.ID}}" value="true" placeholder="{{.Placeholder}}"/>
			<input type="hidden" name="{{.Name}}" value="false"/>
			{{if not .IsArrayItem }} {{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}} {{end}}
		</div>
	</div>
`))

var formSelect = template.Must(template.New("form/select").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
	{{$self := .}}
	<select class="form-control" name="{{.Name}}" id="{{.ID}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} >
		{{range $v, $t := .Options}}
			<option {{if eq $self.Value $v}}selected{{end}} value={{$v}}>{{$t}}</option>
		{{end}}
	</select>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))
