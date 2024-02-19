{{/*
Expand the name of the chart.
*/}}
{{- define "naaprs.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "naaprs.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "naaprs.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "naaprs.labels" -}}
helm.sh/chart: {{ include "naaprs.chart" . }}
{{ include "naaprs.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "naaprs.selectorLabels" -}}
app.kubernetes.io/name: {{ include "naaprs.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "naaprs.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "naaprs.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Calculate schedule from interval in seconds.
*/}}
{{- define "naaprs.scheduleFromInterval" -}}
{{- $seconds := default 60 .Values.env.INTERVAL -}}
{{- if lt (int $seconds) 60 -}}
{{- print "* * * * *" -}} {{/* Every minute */}}
{{- else -}}
{{- $minutes := div (int $seconds) 60 -}}
{{- if eq $minutes 60 -}}
{{- print "0 * * * *" -}} {{/* Every hour */}}
{{- else if lt $minutes 1440 -}}
{{- printf "*/%d * * * *" $minutes -}} {{/* Every N minutes */}}
{{- else -}}
{{- print "Unsupported interval. Must be less than 86400 seconds (24 hours)." -}}
{{- end -}}
{{- end -}}
{{- end -}}

