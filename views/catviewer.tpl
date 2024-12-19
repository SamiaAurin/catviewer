<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Cat Viewer</title>
        <!-- Font Awesome -->
        <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
        <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css" rel="stylesheet">
        <!-- Custom CSS : static/css/styles.css-->
        <link href="/static/css/styles.css" rel="stylesheet">
                
    </head>
    <body>
        <h1>{{.message}}</h1>
        <div class="container">
            <div class="header">
                <a href="#" class="active">⬆️⬇️ Voting</a>
                <a href="#" onclick="showBreeds()">🔍 Breeds</a>
                <a href="#" onclick="showFavorites()"><i>&#9825;</i> Favs</a>
            </div>
            <div class="image-container" id="image-container">
                <img src="https://cdn2.thecatapi.com/images/0XYvRd7oD.jpg" alt="Cat Image" data-image-id="0XYvRd7oD">
            </div>
            <div class="footer">
                <i onclick="saveFavorite()" class="fav-icon">&#9825;</i>
                <div>
                    <i id="thumbs-up" class="vote-icon" onclick="vote(1)">&#128077;</i>
                    <i id="thumbs-down" class="vote-icon" onclick="vote(-1)">&#128078;</i>
                </div>
            </div>
        </div>
        <script src="/static/js/catviewer.js"></script>
    </body>
</html>
