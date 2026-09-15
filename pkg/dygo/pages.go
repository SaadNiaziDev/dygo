package dygo

// Page is the public contract for a metadata-driven app Page.
//
// Name is the app-scoped identity (app.key), and Key is the page bundle key.
// Path is the public URL path. Renderer is entity-index or vue. A vue Page
// loads the sibling <page>.vue file from the Page bundle. Options are
// renderer-specific metadata passed to the Page component.
type Page struct {
	App         string         `json:"app,omitempty" yaml:"-"`
	Name        string         `json:"name,omitempty" yaml:"-"`
	Key         string         `json:"key,omitempty" yaml:"-"`
	Source      string         `json:"source,omitempty" yaml:"-"`
	Label       string         `json:"label" yaml:"label"`
	Description string         `json:"description,omitempty" yaml:"description,omitempty"`
	Icon        string         `json:"icon,omitempty" yaml:"icon,omitempty"`
	Path        string         `json:"path" yaml:"path"`
	Renderer    string         `json:"renderer" yaml:"renderer"`
	Options     map[string]any `json:"options,omitempty" yaml:"options,omitempty"`
}
