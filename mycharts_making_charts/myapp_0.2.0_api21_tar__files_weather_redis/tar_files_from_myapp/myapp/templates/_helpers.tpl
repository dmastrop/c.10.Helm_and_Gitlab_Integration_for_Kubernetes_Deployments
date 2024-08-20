{{/*
Expand the name of the chart.
myapp.name is the default chart name .Chart.Name or if the nameOverride is -set command then use that.
*/}}
{{- define "myapp.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
DM: if no fullnameOverride specified in -set then check the .Values.nameOverride. NOTE that $name will be either the .Chart.Name
or the .Values.nameOverride if it is present.
If that name (for example Chart.Name is contained in the .Release.Name then can use the Release.Name as the fully qualified name
Contains function example: .Release.Name = myapp01 and $name = myapp then the function returns true since myapp
is a subset of myapp01.  myapp01 will be used at the myapp.fullname 
If this is not true then use Release.Name-$name as the fully qualified myapp.fullname
Example: the release name is test: helm install test . --dry-run, and the $name is myapp
The final fully qualified name myapp.fullname is test-myapp
*/}}
{{- define "myapp.fullname" -}}
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
replace function will replace any + in the version with underscores for readability
*/}}
{{- define "myapp.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
myapp.lables is the name and version of the chart and the generation service which will always be helm
Note that selectorLabels is another named chart and it is injected into the myapp.labels block below.
The complete output of both of these blocks will look somtning like this: for a helm install test . --dry-run
helm.sh/chart: myapp-0.1.1
app.kubernetes.io/name: myapp
app.kubernetes.io/instance: test
app.kubernetes.io/version: "2.0.0"
app.kubernetes.io/managed-by: Helm

All of these together will uniquely identify an deployed app for a given service in the helm list.
*/}}
{{- define "myapp.labels" -}}
helm.sh/chart: {{ include "myapp.chart" . }}
{{ include "myapp.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
This is injected into the myapp.labels named template above
*/}}
{{- define "myapp.selectorLabels" -}}
app.kubernetes.io/name: {{ include "myapp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
If .Values.serviceAccount.name is not specified it iwll use myapp.fullname or the fully qualified name
*/}}
{{- define "myapp.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "myapp.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Custom named template: Change the image pull policy according to the type of envronment
For production use Always to force authentiation for all new containers
For development can use ifNotPresent
If the environment is not in .Values.environment then use default "production"
If the environment is not production then use IfNotPresent (development)
If it is a production environment (by default this is set) the else to Always for the imagePullPolicy
for authentication/security reasons.
IfNotPresent will only be used if explicitly set to development. By default it will be Always imagePullPolicy.
NOTE that .Values.environment variable has been added to the values.yaml file and removed the pullPolicy
In deployment.yaml replaced the .Values.image.pullPolicy with the imagePullPolicy: {{ include "myapp.imagePullPolicy" . }}
to this custom named template below.
*/}}
{{- define "myapp.imagePullPolicy" -}}
    {{- $environment := default "production" .Values.environment }}
    {{- if not (eq $environment "production") }}
        {{- "IfNotPresent" }}
    {{- else }}
        {{- "Always" }}
    {{- end }}
{{- end }}    
