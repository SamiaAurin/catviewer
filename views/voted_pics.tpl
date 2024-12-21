<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Viewer</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="/static/css/styles.css" rel="stylesheet">
</head>
<body>
    
    <div class="container">
            <h1>Voted Images</h1>

            {{if .Error}}
                <p class="error">{{.Error}}</p>
            {{else}}
                {{range .Votes}}
                <div class="voted-image">
                    <img src="{{index .image "url"}}" alt="Cat Image" />
                    <p>Vote Value: {{.value}}</p>
                    <p>Image ID: {{.image_id}}</p>
                    <p>Vote ID: {{.id}}</p>
                </div>
                {{end}}
            {{end}}
    </div>


    <script src="/static/js/catviewer.js"></script>
</body>
</html>
