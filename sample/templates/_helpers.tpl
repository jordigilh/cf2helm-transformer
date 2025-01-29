{{- define "metadata" -}}
  {{- if .Values.space }}
{{- print "namespace: " .Values.space }}
  {{- end }}
{{ print "name: " .Values.name }}
labels:
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
