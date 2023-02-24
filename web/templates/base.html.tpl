{{ define "base" }}
<!doctype html>
<html>
    <head>
        <meta charset="utf-8">
        <meta http-equiv="refresh" content="3">
    </head>
    <body>
        <main>
            {{ .Content }}
        </main>
    </body>
</html>
{{ end }}
