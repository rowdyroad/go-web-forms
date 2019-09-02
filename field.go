package forms

type formField struct {
	ID               string            `yaml:"-"`
	Name             string            `yaml:"name"`
	Description      string            `yaml:"description"`
	Type             string            `yaml:"type"`
	ValueType        string            `yaml:"-"`
	Label            string            `yaml:"label"`
	ItemLabel        string            `yaml:"itemLabel"`
	DeleteBtnCaption string            `yaml:"deleteBtnCaption"`
	AddBtnCaption    string            `yaml:"addBtnCaption"`
	SetBtnCaption    string            `yaml:"setBtnCaption"`
	UnsetBtnCaption  string            `yaml:"unsetBtnCaption"`
	IsNil            bool              `yaml:"-"`
	Value            interface{}       `yaml:"-"`
	Disabled         bool              `yaml:"disabled"`
	Rows             int               `yaml:"rows"`
	Placeholder      string            `yaml:"placeholder"`
	Index            int               `yaml:"-"`
	Readonly         bool              `yaml:"readonly"`
	Options          map[string]string `yaml:"options"`
	Skip             bool              `yaml:"-"`
	IsArrayItem      bool              `yaml:"-"`
	Length           int               `yaml:"-"`
	Template         string            `yaml:"-"`
	Expanded         bool              `yaml:"expanded"`
}

func (t formField) Copy() formField {
	return formField{
		Name:             t.Name,
		Description:      t.Description,
		Type:             t.Type,
		ValueType:        t.ValueType,
		Label:            t.Label,
		ItemLabel:        t.ItemLabel,
		DeleteBtnCaption: t.DeleteBtnCaption,
		AddBtnCaption:    t.AddBtnCaption,
		Disabled:         t.Disabled,
		Rows:             t.Rows,
		Placeholder:      t.Placeholder,
		Readonly:         t.Readonly,
		Options:          t.Options,
		Length:           t.Length,
	}
}
