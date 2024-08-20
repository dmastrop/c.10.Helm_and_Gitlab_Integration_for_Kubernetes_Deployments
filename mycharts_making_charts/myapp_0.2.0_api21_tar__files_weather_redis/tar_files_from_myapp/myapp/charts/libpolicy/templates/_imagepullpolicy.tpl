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

NOTE: this is the same as the original code in the _helpers.tpl parent helper file except change .myapp.imagePullPolicy name to 
libpolicy.imagePullPolicy since this is being used as part of the shared libpolicy helper files
*/}}
{{- define "libpolicy.imagePullPolicy" -}}
    {{- $environment := default "production" .Values.environment }}
    {{- if not (eq $environment "production") }}
        {{- "IfNotPresent" }}
    {{- else }}
        {{- "Always" }}
    {{- end }}
{{- end }}    

