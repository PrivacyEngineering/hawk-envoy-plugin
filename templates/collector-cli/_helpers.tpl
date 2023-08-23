{{/*
Expand the name of the chart.
*/}}
{{- define "collector-cli.name" -}}
{{ $name := default .Chart.Name .Values.nameOverride }}
{{- printf "%s-%s" $name "collector-cli" | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "collector-cli.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- printf "%s-%s" .Release.Name "collector-cli" | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" (printf "%s-%s" .Release.Name $name | trimSuffix "-") "collector-cli" | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "collector-cli.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "collector-cli.labels" -}}
helm.sh/chart: {{ include "collector-cli.chart" . }}
{{ include "collector-cli.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{/*
Selector labels
*/}}
{{- define "collector-cli.selectorLabels" -}}
app.kubernetes.io/name: {{ include "collector-cli.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
