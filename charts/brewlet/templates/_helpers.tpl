{{/* Common helpers for the brewlet chart. */}}

{{- define "brewlet.namespace" -}}
{{- default "brewlet" .Values.namespace -}}
{{- end -}}

{{/* Standard labels applied to every rendered object. */}}
{{- define "brewlet.labels" -}}
app.kubernetes.io/name: brewlet
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: brewlet-{{ .Chart.Version }}
{{- end -}}

{{/* Fully-qualified DNS name of the admission webhook Service. */}}
{{- define "brewlet.admission.serviceName" -}}
brewlet-admission
{{- end -}}

{{/*
brewlet.jdkItems renders a comma-separated "<dist>-<feature>" JDK string (e.g.
"temurin-21,microsoft-25") as the NodeProfile spec.jdks list. Call with a dict:
  {{ include "brewlet.jdkItems" (dict "value" "temurin-21,microsoft-25") }}
*/}}
{{- define "brewlet.jdkItems" -}}
{{- $spec := trim .value -}}
{{- if not $spec -}}{{- fail "provisioner.jdks must list at least one <dist>-<feature> JDK" -}}{{- end -}}
{{- range $tok := splitList "," $spec -}}
{{- $t := trim $tok -}}
{{- if $t -}}
{{- $parts := splitList "-" $t -}}
{{- if ne (len $parts) 2 -}}{{- fail (printf "invalid JDK token %q; want <distribution>-<feature>" $t) -}}{{- end -}}
- distribution: {{ index $parts 0 | quote }}
  feature: {{ index $parts 1 }}
{{ end -}}
{{- end -}}
{{- end -}}

{{/*
brewlet.launcherItems renders a comma-separated launcher string as a YAML list.
Emits nothing when empty.
*/}}
{{- define "brewlet.launcherItems" -}}
{{- $spec := trim .value -}}
{{- range $tok := splitList "," $spec -}}
{{- $t := trim $tok -}}
{{- if $t }}
- {{ $t | quote }}
{{- end -}}
{{- end -}}
{{- end -}}
