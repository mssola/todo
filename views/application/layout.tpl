<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>todo</title>
    <link rel="icon" sizes="196x196" href="/images/icon.png?v=0">
    <link href="{{ .Prefix }}/stylesheets/vendor/normalize.min.css" rel="stylesheet" type="text/css" />
    <link href="{{ .Prefix }}/stylesheets/{{ view }}.css" rel="stylesheet" type="text/css" />
    {{if .Print}}
        <script src="{{ .Prefix }}/javascripts/print.js"></script>
    {{else}}
        {{if .JS}}
            <script src="{{ .Prefix }}/javascripts/{{ .JS }}.js"></script>
        {{end}}
    {{end}}
</head>
<body>
    {{if .Print}}
    <article style="background-color: white">
    {{else}}
    <article>
    {{end}}
        {{ yield }}
    </article>
</body>
</html>
