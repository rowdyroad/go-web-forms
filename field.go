package forms

type formField struct {
	ID               string            `yaml:"-"`
	Name             string            `yaml:"name"`
	Description      string            `yaml:"description"`
	Type             string            `yaml:"type"`
	Label            string            `yaml:"label"`
	ItemLabel        string            `yaml:"itemLabel"`
	DeleteBtnCaption string            `yaml:"deleteBtnCaption"`
	AddBtnCaption    string            `yaml:"addBtnCaption"`
	Value            interface{}       `yaml:"-"`
	Disabled         bool              `yaml:"disabled"`
	Rows             int               `yaml:"rows"`
	Placeholder      string            `yaml:"placeholder"`
	Indent           int               `yaml:"-"`
	Readonly         bool              `yaml:"readonly"`
	Options          map[string]string `yaml:"options"`
	Skip             bool              `yaml:"-"`
	IsArrayItem      bool              `yaml:"-"`
	Length           int               `yaml:"-"`
}

func (t formField) Copy() formField {
	return formField{
		Name:             t.Name,
		Description:      t.Description,
		Type:             t.Type,
		Label:            t.Label,
		ItemLabel:        t.ItemLabel,
		DeleteBtnCaption: t.DeleteBtnCaption,
		AddBtnCaption:    t.AddBtnCaption,
		Disabled:         t.Disabled,
		Rows:             t.Rows,
		Placeholder:      t.Placeholder,
		Readonly:         t.Readonly,
		Options:          t.Options,
	}
}
