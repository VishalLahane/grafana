package channels

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/prometheus/alertmanager/template"
	"github.com/stretchr/testify/require"
)

const DefaultTemplateString = `
{{ define "__subject" }}[{{ .Status | toUpper }}{{ if eq .Status "firing" }}:{{ .Alerts.Firing | len }}{{ end }}] {{ .GroupLabels.SortedPairs.Values | join " " }} {{ if gt (len .CommonLabels) (len .GroupLabels) }}({{ with .CommonLabels.Remove .GroupLabels.Names }}{{ .Values | join " " }}{{ end }}){{ end }}{{ end }}

{{ define "__text_alert_list" }}{{ range . }}Labels:
{{ range .Labels.SortedPairs }} - {{ .Name }} = {{ .Value }}
{{ end }}Annotations:
{{ range .Annotations.SortedPairs }} - {{ .Name }} = {{ .Value }}
{{ end }}Source: {{ .GeneratorURL }}
{{ end }}{{ end }}

{{ define "default.title" }}{{ template "__subject" . }}{{ end }}

{{ define "default.message" }}{{ if gt (len .Alerts.Firing) 0 }}
**Firing**
{{ template "__text_alert_list" .Alerts.Firing }}

{{ end }}
{{ if gt (len .Alerts.Resolved) 0 }}
**Resolved**
{{ template "__text_alert_list" .Alerts.Resolved }}
{{ end }}
{{ end }}
`

func templateForTests(t *testing.T) *template.Template {
	f, err := ioutil.TempFile("/tmp", "template")
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, os.RemoveAll(f.Name()))
	})

	_, err = f.WriteString(DefaultTemplateString)
	require.NoError(t, err)

	tmpl, err := template.FromGlobs(f.Name())
	require.NoError(t, err)

	return tmpl
}
