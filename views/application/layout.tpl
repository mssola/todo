<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8" />
    <title>todo</title>
    <link rel="icon" sizes="196x196" href="/images/icon.png?v=0">
    <link href="/stylesheets/{{ view }}.css" rel="stylesheet" type="text/css" />
    {{if .JS}}
        <script src="/javascripts/jquery/jquery.min.js"></script>
        <script src="/javascripts/{{ .JS }}.js"></script>
    {{end}}
</head>
<body>
    <article>
        {{ yield }}
    </article>
</body>
</html>
