package forms

import "html/template"

var formHeader = template.Must(template.New("form/header").Parse(`
	<script>
		if (!window.goInit) {
			window.goInit = true;
			if (!String.prototype.replaceAll) {
				String.prototype.replaceAll = function(search, replacement) {
					var target = this;
					return target.split(search).join(replacement);
				};
			}

			window.goWebFormsUUID = function() {
				return 'xxxxxxxxxxxxxxxxyxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
					var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
					return v.toString(16);
				});
			}

			window.goWebFormsIndexes = {};
			window.goWebFormsIndexes = function(id, indexMax) {
				if (!goWebFormsIndexes[id]) {
					goWebFormsIndexes[id] = indexMax
				}
				var index = goWebFormsIndexes[id];
				var cnt = document.getElementById('template-' + id)
									.innerHTML
									.replaceAll('template-'+id+'-id', goWebFormsUUID())
									.replaceAll('template-'+id+'-name', id+'['+goWebFormsIndexes[id]+']');

				var e = document.createElement('div');
				e.innerHTML = cnt;
				document.getElementById('array-'+id).appendChild(e);
				goWebFormsIndexes[id] = index + 1;
			}

			Array.prototype.reduce||Object.defineProperty(Array.prototype,"reduce",{value:function(e){if(null===this)throw new TypeError("Array.prototype.reduce called on null or undefined");if("function"!=typeof e)throw new TypeError(e+" is not a function");var n,r=Object(this),t=r.length>>>0,o=0;if(arguments.length>=2)n=arguments[1];else{for(;t>o&&!(o in r);)o++;if(o>=t)throw new TypeError("Reduce of empty array with no initial value");n=r[o++]}for(;t>o;)o in r&&(n=e(n,r[o],o,r)),o++;return n}}),window.goForm2JSON||(window.goForm2JSON=function(e){var n={},r=[];if("function"==typeof HTMLFormElement&&e instanceof HTMLFormElement)for(var t in e.elements)(e.elements[t]instanceof HTMLInputElement||e.elements[t]instanceof HTMLSelectElement||e.elements[t]instanceof HTMLTextAreaElement)&&r.push({name:e.elements[t].name,value:e.elements[t].value});else Array.isArray(e)&&(r=e);return n=r.reduce(function(e,n){var r=e,t=n.name.split(".");return t.forEach(function(e,o){var a=e.replace(/\[[0-9]*\]$/,"");if(r.hasOwnProperty(a)||(r[a]=new RegExp("[[0-9]*]$").test(e)?[]:{}),r[a]instanceof Array){var i=parseInt((e.match(new RegExp("([0-9]+)]$"))||[]).pop(),10);return i=isNaN(i)?r[a].length:i,r[a][i]=r[a][i]||{},o===t.length-1?r[a][i]=n.value:r=r[a][i]}return o===t.length-1?r[a]=n.value:r=r[a]}),e},{})});
			window.goWebFormSubmit = function(id) {
				var data = window.goForm2JSON(document.getElementById(id));
				fetch(id, {
					method: "POST",
					headers: {
						"Content-type":"application/json"
					},
					body: JSON.stringify(data)
				})
			}
		}
	</script>
	<form id="{{.ID}}" action="javascript:window.goWebFormSubmit('{{.ID}}')">
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
			<button class="btn btn-secondary" onclick="goWebFormsAddArrayItem('{{.Name}}', {{.Length}})">{{if .AddBtnCaption}}{{.AddBtnCaption}}{{else}}Add{{end}}</button>
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
			<button class="btn btn-danger" onclick="javascript:document.getElementById('item-{{.Name}}').remove()">{{if .DeleteBtnCaption}}{{.DeleteBtnCaption}}{{else}}Delete{{end}}"</button>
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
