{{- define "metadata" -}}
  {{- if .Values.space }}
{{- print "namespace: " .Values.space }}
  {{- end }}
{{ print "name: " .Values.name }}
labels:
{{- print "app: " .Values.name | nindent 2 }}
  {{-  range $k,$v:= .Values.labels -}}
{{ print $k ": " $v |nindent 2 }}
  {{- end }}
annotations:
  {{-  range $k,$v:= .Values.annotations -}}
{{ print $k ": " $v |nindent 2 -}}
  {{- end -}}
    {{- if .Values.stack }}
{{- print "stack: " .Values.stack |nindent 2 -}}
    {{- end }}
{{- end -}}
