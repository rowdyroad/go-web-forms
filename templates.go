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
		if (!window.goWebFormsUUID) {
			window.goWebFormsUUID = function() {
				return 'xxxxxxxxxxxxxxxxyxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
					var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
					return v.toString(16);
				});
			}
		}
		if (!window.goWebFormsIndexes) {
			window.goWebFormsIndexes = {};
		}

		if (!window.goWebFormsIgnores) {
			window.goWebFormsIgnores = [];
		}

		if (!window.goWebFormsTogglePtrField) {
			window.goWebFormsTogglePtrField = function(id, setCaption, unsetCaption) {
				var css = document.getElementById('ptr-'+id).style;
				if (css.display == 'none') {
					css.display = '';
					document.getElementById('ptr-btn-'+id).innerHTML = unsetCaption;
					var index = window.goWebFormsIgnores.indexOf(id);
					if (index !== -1) {
						window.goWebFormsIgnores.splice(index, 1);
					}
				} else {
					css.display = 'none';
					document.getElementById('ptr-btn-'+id).innerHTML = setCaption;
					window.goWebFormsIgnores.push(id);
				}
			}
		}
		if (!window.goWebFormsAddArrayItem) {
			window.goWebFormsAddArrayItem = function(id, indexMax) {
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

		  if (!window.goForm2JSON) {
			window.goFormConverters = {
				int: function(value) {
					return parseInt(value);
				},
				float64: function(value) {
					return parseFloat(value);
				},
				bool: function(value) {
					return value == "true" ? true : false;
				}
			}

			window.goForm2JSON = function(form) {
				var data = {}, form_arr = [];
				if(typeof HTMLFormElement === "function" && form instanceof HTMLFormElement) {
					for(var i in form.elements) {
						if(form.elements[i] instanceof HTMLInputElement ||
							form.elements[i] instanceof HTMLSelectElement ||
							form.elements[i] instanceof HTMLTextAreaElement) {
							var converter = window.goFormConverters[form.elements[i].getAttribute('data-value-type')];
							var name = form.elements[i].name;
							if (!window.goWebFormsIgnores.some(function(n) {
								console.log(name, n, name.indexOf(n));
								return name.indexOf(n) === 0;
							})) {
								form_arr.push({name:name, value: converter ? converter(form.elements[i].value) : form.elements[i].value });
							}
						}
					}
				}
				else if(Array.isArray(form)) {
					form_arr = form;
				}
				data = form_arr.reduce(function (r, o) {
					var s = r, arr = o.name.split('.');
					arr.forEach((n, k) => {
						var ck = n.replace(/\[[0-9]*\]$/, "");
						if (!s.hasOwnProperty(ck))
							s[ck] = (new RegExp("\[[0-9]*\]$").test(n)) ? [] : {};
						if (s[ck] instanceof Array) {
							var i = parseInt((n.match(new RegExp("([0-9]+)\]$")) || []).pop(), 10);
							i = isNaN(i) ? s[ck].length : i;
							s[ck][i] = s[ck][i] || {};
							if(k === arr.length - 1) {
								return s[ck][i] = o.value;
							}
							else {
								return s = s[ck][i];
							}
						}
						else {
							return (k === arr.length - 1) ? s[ck] = o.value : s = s[ck];
						}
					});
					return r;
				}, {});
				return data;
			}
		}

		if (!window.goWebFormSubmit) {
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

var formPtrHeader = template.Must(template.New("form/ptrHeader").Parse(`
	{{if .Label}}
		<label>
			{{.Label}}
			{{if .Description}}<small><div>{{.Description}}</div></small>{{end}}
		</label>
	{{end}}
	{{if .IsNil}}
	<script>
		window.goWebFormsIgnores.push('{{.Name}}');
	</script>
	{{end}}
	</script>
	<div id="ptr-{{.Name}}" {{if .IsNil}}style="display:none"{{end}}>
`))

var formPtrFooter = template.Must(template.New("form/ptrFooter").Parse(`
	</div>
	{{if not .Readonly}}
		<div style="margin:0.4em;margin-bottom:1em">
			<button type="button" class="btn btn-secondary" id="ptr-btn-{{.Name}}" onclick="goWebFormsTogglePtrField('{{.Name}}', '{{if .SetBtnCaption}}{{.SetBtnCaption}}{{else}}Set{{end}}','{{if .UnsetBtnCaption}}{{.UnsetBtnCaption}}{{else}}Unset{{end}}')">
			{{if .IsNil}}
				{{if .SetBtnCaption}}{{.SetBtnCaption}}{{else}}Set{{end}}
			{{else}}
				{{if .UnsetBtnCaption}}{{.UnsetBtnCaption}}{{else}}Unset{{end}}
			{{end}}
			</button>
		</div>
	{{end}}
`))

var formArrayHeader = template.Must(template.New("form/arrayHeader").Parse(`
	<div class="shadow-sm p-3 mb-5 bg-white rounded">
		{{if .Label}}
			<label>
				{{.Label}}
				{{if .Description}}<small><div>{{.Description}}</div></small>{{end}}
			</label>
		{{end}}
		<div id="array-{{.Name}}">
`))

var formArrayFooter = template.Must(template.New("form/arrayFooter").Parse(`
		</div>
		{{if not .Readonly}}
			<div style="margin:0.4em;margin-bottom:1em">
				<button type="button" class="btn btn-secondary" onclick="goWebFormsAddArrayItem('{{.Name}}', {{.Length}})">{{if .AddBtnCaption}}{{.AddBtnCaption}}{{else}}Add{{end}}</button>
			</div>
		{{end}}
	</div>
`))
var formArrayItemWrapperHeader = template.Must(template.New("form/arrayItemWrapperHeader").Parse(`
	<div id="item-{{.Name}}" class="form-group" style="padding-left: {{.Indent}}em">
`))

var formStructHeader = template.Must(template.New("form/structHeader").Parse(`
	<div class="shadow-sm p-3 mb-5 bg-white rounded">
		{{if .IsArrayItem}}
			{{if .ItemLabel}}
				<label>
					{{.ItemLabel}}
					{{if .Description}}<small>{{.Description}}</small>{{end}}
				</label>
			{{end}}
		{{else}}
			{{if .Label}}
				<label>
					{{.Label}}
					{{if .Description}}<small>{{.Description}}</small>{{end}}
				</label>
			{{end}}
		{{end}}
`))

var formStructFooter = template.Must(template.New("form/structFooter").Parse(`</div>`))

var formArrayItemWrapperFooter = template.Must(template.New("form/arrayItemWrapperFooter").Parse(`
	{{if not .Readonly}}
		<div style="text-align:right; padding:0.4em 0">
			<button type="button" class="btn btn-danger" onclick="javascript:document.getElementById('item-{{.Name}}').remove()">{{if .DeleteBtnCaption}}{{.DeleteBtnCaption}}{{else}}Delete{{end}}</button>
		</div>
	{{end}}
	</div>
`))

var formInput = template.Must(template.New("form/input").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<input type="{{.Type}}"  name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} class="form-control" id="{{.ID}}" value="{{.Value}}" placeholder="{{.Placeholder}}"/>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))

var formTextarea = template.Must(template.New("form/textarea").Parse(`
	{{if not .IsArrayItem }}
		<div class="form-group" style="padding-left: {{.Indent}}em">
		{{if .Label}} <label class="form-check-label" for="{{.ID}}">{{.Label}}</label>{{end}}
	{{end}}
		<textarea class="form-control" name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} placeholder="{{.Placeholder}}" id="{{.ID}}" rows="{{.Rows}}">{{.Value}}</textarea>
	{{if not .IsArrayItem }}
		</div>
	{{end}}


	`))

var formCheckbox = template.Must(template.New("form/checkbox").Parse(`
	<div class="form-group" style="padding-left: {{.Indent}}em">
		<div class="form-check">
			<input class="form-check-input" type="{{.Type}}" name="{{.Name}}" data-value-type="{{.ValueType}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} {{if .Value}}checked{{end}}  class="form-control" id="{{.ID}}" value="true" placeholder="{{.Placeholder}}"/>
			<input type="hidden" name="{{.Name}}" data-value-type="{{.ValueType}}" value="false"/>
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
	<select class="form-control" name="{{.Name}}" data-value-type="{{.ValueType}}" id="{{.ID}}" {{if .Disabled}}disabled{{end}} {{if .Readonly}}readonly{{end}} >
		{{range $v, $t := .Options}}
			<option {{if eq $self.Value $v}}selected{{end}} value={{$v}}>{{$t}}</option>
		{{end}}
	</select>
	{{if not .IsArrayItem }}
		</div>
	{{end}}
`))
