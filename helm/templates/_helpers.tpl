{{- define "rathole-operator.created-meta" }}
app.kubernetes.io/created-by: rathole-operator
app.kubernetes.io/part-of: rathole-operator
{{- end }}

{{- define "rathole-operator.rbac-labels" }}
app.kubernetes.io/component: rbac
{{- include "rathole-operator.created-meta" . }}
{{- end }}

{{- define "rathole-operator.selector-labels" -}}
control-plane: rathole-operator
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "rathole-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "rathole-operator.fullname" -}}
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
{{- define "rathole-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}
