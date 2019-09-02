package forms

import "html/template"

var formHeader = template.Must(template.New("form/header").Parse(`
<script>
	if (!String.prototype.replaceAll) {
		String.prototype.replaceAll = function(search, replacement) {
			var target = this;
			return target.split(search).join(replacement);
		};
	}
	if (!window.goWebForms) {
		window.goWebForms = {
			indexes: {},
			converters: {
				int: function(value) {
					return parseInt(value);
				},
				float64: function(value) {
					return parseFloat(value);
				},
				bool: function(value) {
					return value == "true" ? true : false;
				}
			},
			togglePtrField: function(id, setCaption, unsetCaption) {
				var cnt = document.getElementById('template-'+id).innerHTML;
				var el = document.getElementById('ptr-'+id);
				var elBtn = document.getElementById('ptr-btn-'+id);
				if (el.childElementCount === 0) {
					el.style.display = '';
					el.innerHTML = cnt;
					elBtn.innerHTML = unsetCaption;
					elBtn.classList.add('btn-danger');
					elBtn.classList.remove('btn-outline-secondary');
				} else {
					el.style.display = 'none';
					el.innerHTML = '';
					elBtn.innerHTML = setCaption;
					elBtn.classList.remove('btn-danger');
					elBtn.classList.add('btn-outline-secondary');
				}
			},
			toggleExpand: function(el, query) {
				var toggle = el.querySelector(query);
				if (toggle.style.display === 'none') {
					toggle.style.display = '';
				} else {
					toggle.style.display = 'none';
				}
			},
			addArrayItem: function(id, indexMax) {
				if (!window.goWebForms.indexes[id]) {
					window.goWebForms.indexes[id] = indexMax
				}
				var index = window.goWebForms.indexes[id];

				var templateId = id.replace(new RegExp(/\[\d+\]/, 'g'), '');

				var cnt = document.getElementById('template-'+templateId)
									.innerHTML
									.replaceAll(templateId, id+'['+window.goWebForms.indexes[id]+']');

				var e = document.createElement('div');
				e.innerHTML = cnt;
				document.getElementById(id).appendChild(e);
				window.goWebForms.indexes[id] = index + 1;
			},
			submit: function(id) {
				var data = window.goWebForms.JSON(document.getElementById(id));
				fetch(id, {
					method: "POST",
					headers: {
						"Content-type":"application/json"
					},
					body: JSON.stringify(data)
				})
			},
			JSON: function(form) {
				var data = {}, form_arr = [];

				if(typeof HTMLFormElement === "function" && form instanceof HTMLFormElement) {
					for(var i in form.elements) {
						if(form.elements[i] instanceof HTMLInputElement ||
							form.elements[i] instanceof HTMLSelectElement ||
							form.elements[i] instanceof HTMLTextAreaElement) {
							var converter = window.goWebForms.converters[form.elements[i].getAttribute('data-value-type')];
							var name = form.elements[i].name;
							form_arr.push({name:name, value: converter ? converter(form.elements[i].value) : form.elements[i].value });
						}
					}
				}
				else if(Array.isArray(form)) {
					form_arr = form;
				}

				data = form_arr.reduce(function (r, o) {
					var s = r, arr = o.name.split('.');
					arr.forEach((n, k) => {
						var key = n.replace(/\[[0-9]*\]$/, "");
						if (!s.hasOwnProperty(key))
							s[key] = (new RegExp("\[[0-9]*\]$").test(n)) ? [] : {};

						if (s[key] instanceof Array) {
							if(k === arr.length - 1) {
								return s[key].push(o.value);
							} else {
								var i = parseInt((n.match(new RegExp("([0-9]+)\]$")) || []).pop(), 10);
								i = isNaN(i) ? s[key].length : i;
								s[key][i] = s[key][i] || {};
								return s = s[key][i];
							}
						} else {
							return (k === arr.length - 1) ? s[key] = o.value : s = s[key];
						}
					});
					return r;
				}, {});
				return data;
			}
		}
	}

	if (!Array.prototype.reduce) {
		Object.defineProperty(Array.prototype, 'reduce', {
			value: function(callback /*, initialValue*/) {
			if (this === null) {
				throw new TypeError( 'Array.prototype.reduce ' +
				'called on null or undefined' );
			}
			if (typeof callback !== 'function') {
				throw new TypeError( callback +
				' is not a function');
			}

			// 1. Let O be ? ToObject(this value).
			var o = Object(this);

			// 2. Let len be ? ToLength(? Get(O, "length")).
			var len = o.length >>> 0;

			// Steps 3, 4, 5, 6, 7
			var k = 0;
			var value;

			if (arguments.length >= 2) {
				value = arguments[1];
			} else {
				while (k < len && !(k in o)) {
				k++;
				}

				// 3. If len is 0 and initialValue is not present,
				//    throw a TypeError exception.
				if (k >= len) {
				throw new TypeError( 'Reduce of empty array ' +
					'with no initial value' );
				}
				value = o[k++];
			}

			// 8. Repeat, while k < len
			while (k < len) {
				// a. Let Pk be ! ToString(k).
				// b. Let kPresent be ? HasProperty(O, Pk).
				// c. If kPresent is true, then
				//    i.  Let kValue be ? Get(O, Pk).
				//    ii. Let accumulator be ? Call(
				//          callbackfn, undefined,
				//          « accumulator, kValue, k, O »).
				if (k in o) {
				value = callback(value, o[k], k, o);
				}

				// d. Increase k by 1.
				k++;
			}

			// 9. Return accumulator.
			return value;
			}
		});
	}
</script>
<form id="{{.ID}}" action="javascript:window.goWebForms.submit('{{.ID}}')">
`))

var formFooter = template.Must(template.New("form/footer").Parse(`
		<button type="submit" class="btn btn-primary">{{if .}}{{if .SubmitBtnCaption}}{{.SubmitBtnCaption}}{{else}}Submit{{end}}{{else}}Submit{{end}}</button>
	</form>
`))

var formPtrHeader = template.Must(template.New("form/ptrHeader").Parse(`
	<div class="mb-2"
			data-template="form/ptr"
			data-name="{{.Name}}"
			data-id="{{.ID}}">
		<div class="mb-2">
			{{.Label}}
			{{if .Description}}<small><div>{{.Description}}</div></small>{{end}}
			{{if not .Readonly}}
				<button type="button" id="ptr-btn-{{.Name}}" onclick="window.goWebForms.togglePtrField('{{.Name}}', '{{if .SetBtnCaption}}{{.SetBtnCaption}}{{else}}Set{{end}}','{{if .UnsetBtnCaption}}{{.UnsetBtnCaption}}{{else}}Unset{{end}}')"
				{{if .IsNil}}
					class="btn btn-outline-secondary btn-sm"> {{if .SetBtnCaption}}{{.SetBtnCaption}}{{else}}Set{{end}}
				{{else}}
					class="btn btn-danger btn-sm "> {{if .UnsetBtnCaption}}{{.UnsetBtnCaption}}{{else}}Unset{{end}}
				{{end}}
				</button>
			{{end}}
		</div>
		<div id="ptr-{{.Name}}" {{if .IsNil}}style="display:none"{{end}}>
`))

var formPtrFooter = template.Must(template.New("form/ptrFooter").Parse(`
		</div>
	</div>
`))

var formArrayHeader = template.Must(template.New("form/arrayHeader").Parse(`
		<div class="card mb-2 w-100"
			data-template="form/array"
			data-name="{{.Name}}"
			data-id="{{.ID}}">
			<div class="card-header" style="cursor:pointer" onclick="goWebForms.toggleExpand(document, '#array-body-{{.Name}}')">
				{{.Label}}
				{{if .Description}}<small><div>{{.Description}}</div></small>{{end}}
			</div>
			<div class="card-body" id="array-body-{{.Name}}" style="display:{{if .Expanded}}block{{else}}none{{end}}">
				<div id="{{.Name}}">
`))

var formArrayFooter = template.Must(template.New("form/arrayFooter").Parse(`
				</div>
				{{if not .Readonly}}
					<div class="mb-3">
						<button type="button" class="btn btn-secondary" onclick="window.goWebForms.addArrayItem('{{.Name}}', {{.Length}})">{{if .AddBtnCaption}}{{.AddBtnCaption}}{{else}}Add{{end}}</button>
					</div>
				{{end}}
			</div>
		</div>
	`))
var formArrayItemWrapperHeader = template.Must(template.New("form/arrayItemWrapperHeader").Parse(`
		<div id="{{.Name}}" class="input-group mb-2" data-template="form/arrayItem" data-name="{{.Name}}" data-id="{{.ID}}">
`))

var formStructHeader = template.Must(template.New("form/structHeader").Parse(`
		<div class="card mb-2 w-100">
			<div class="card-header" style="cursor:pointer" onclick="goWebForms.toggleExpand(document, '#struct-body-{{.Name}}')">
			{{if .IsArrayItem}}
				{{if .ItemLabel}}
						{{.ItemLabel}}
						{{if .Description}}<small>{{.Description}}</small>{{end}}
				{{end}}
			{{else}}
				{{if .Label}}
					{{.Label}}
					{{if .Description}}<small>{{.Description}}</small>{{end}}
				{{end}}
			{{end}}
			</div>
		<div class="card-body" id="struct-body-{{.Name}}" style="display:{{if .Expanded}}block{{else}}none{{end}}">
`))

var formStructFooter = template.Must(template.New("form/structFooter").Parse(`
		</div>
	</div>
`))

var formArrayItemWrapperFooter = template.Must(template.New("form/arrayItemWrapperFooter").Parse(`
	{{if not .Readonly}}
		<div class="input-group-append" data-template="form/arrayItemFooter">
			<button type="button" class="btn btn-danger" onclick="javascript:document.getElementById('{{.Name}}').remove()">{{if .DeleteBtnCaption}}{{.DeleteBtnCaption}}{{else}}Delete{{end}}</button>
		</div>
	{{end}}
	</div>
`))

var formInput = template.Must(template.New("form/input").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<input type="{{.Type}}" name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} class="form-control" id="{{.ID}}" value="{{.Value}}" placeholder="{{.Placeholder}}"/>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))

var formTextarea = template.Must(template.New("form/textarea").Parse(`
		{{if not .IsArrayItem }}
		<div class="form-group">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<textarea class="form-control" name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} placeholder="{{.Placeholder}}" id="{{.ID}}" rows="{{.Rows}}">{{.Value}}</textarea>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))

var formCheckbox = template.Must(template.New("form/checkbox").Parse(`
		<div class="form-group">
		<div class="form-check">
			<input class="form-check-input" type="{{.Type}}" name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} {{if .Value}}checked{{end}}  class="form-control" id="{{.ID}}" value="true" placeholder="{{.Placeholder}}"/>
			<input type="hidden" name="{{.Name}}" data-value-type="{{.ValueType}}" value="false"/>
			{{if not .IsArrayItem }} {{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}} {{end}}
		</div>
	</div>
`))

var formSelect = template.Must(template.New("form/select").Parse(`
		{{if not .IsArrayItem }}
		<div class="form-group">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
	{{$self := .}}
	<select class="form-control" name="{{.Name}}" data-value-type="{{.ValueType}}" id="{{.ID}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} >
		{{range $v, $t := .Options}}
			<option {{if eq $self.Value $v}}selected{{end}} value={{$v}}>{{$t}}</option>
		{{end}}
	</select>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))
