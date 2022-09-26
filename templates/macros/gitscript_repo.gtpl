#!/bin/bash
# This file is generated, do not update manually
{{ $model := . }}
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
ME="$(basename "$0")"

pushd "$DIR" || exit 1

echo "git {{ .gitActionLabel }} - group {{ .groupNode.group.full_path }}"

{{- if .groupNode.group.projects }}
  {{- $projectsLength := len .groupNode.group.projects }}
echo "git {{ .gitActionLabel }} {{ $projectsLength }} projects of {{ .groupNode.group.full_path }}"
  {{ range $idx, $project := .groupNode.group.projects }}
echo "git {{ $model.gitActionLabel }} {{ $project.path_with_namespace }}"
git -C {{ $project.path }} {{ $model.gitAction }}
  {{ end }}
{{- end -}}

{{- if .groupNode.children -}}
  {{ $childrenLength := len .groupNode.children }}
echo "git {{ .gitActionLabel }} {{ $childrenLength }} sub-groups of {{ .groupNode.group.full_path }}"
  {{ range $idx, $subGroup := .groupNode.children }}
"{{ $subGroup.group.path }}/$ME"
  {{- end -}}
{{ end }}

popd || exit 1