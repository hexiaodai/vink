{{/* Hook annotations */}}
{{- define "vink.hook.annotations" -}}
  annotations:
    "helm.sh/hook": {{ .hookType }}
    "helm.sh/hook-delete-policy": before-hook-creation,hook-succeeded
    "helm.sh/hook-weight": {{ .hookWeight | quote }}
{{- end -}}

{{/* Namespace modifying hook annotations */}}
{{- define "vink.namespaceHook.annotations" -}}
{{ template "vink.hook.annotations" merge (dict "hookType" "pre-install") . }}
{{- end -}}

{{/* CRD upgrading hook annotations */}}
{{- define "vink.crdUpgradeHook.annotations" -}}
{{ template "vink.hook.annotations" merge (dict "hookType" "pre-upgrade") . }}
{{- end -}}

{{/* Custom resource uninstalling hook annotations */}}
{{- define "vink.crUninstallHook.annotations" -}}
{{ template "vink.hook.annotations" merge (dict "hookType" "pre-delete") . }}
{{- end -}}

{{/* CRD uninstalling hook annotations */}}
{{- define "vink.crdUninstallHook.annotations" -}}
{{ template "vink.hook.annotations" merge (dict "hookType" "post-delete") . }}
{{- end -}}

{{/* Namespace modifying hook name */}}
{{- define "vink.namespaceHook.name" -}}
{{ include "vink.fullname" . }}-namespace-modify
{{- end }}

{{/* CRD upgrading hook name */}}
{{- define "vink.crdUpgradeHook.name" -}}
{{ include "vink.fullname" . }}-crd-upgrade
{{- end }}

{{/* Custom resource uninstalling hook name */}}
{{- define "vink.crUninstallHook.name" -}}
{{ include "vink.fullname" . }}-uninstall
{{- end }}

{{/* CRD uninstalling hook name */}}
{{- define "vink.crdUninstallHook.name" -}}
{{ include "vink.fullname" . }}-crd-uninstall
{{- end }}
