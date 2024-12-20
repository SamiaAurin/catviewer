<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Viewer</title>
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css" rel="stylesheet">
    <link href="/static/css/styles.css" rel="stylesheet">
</head>
<body>
    <div class="container">
        <div class="header">
            <a href="#" class="active">⬆️⬇️ Voting</a>
            <a href="#" onclick="showBreeds()">🔍 Breeds</a>
            <a href="#" onclick="showFavorites()">💖 Favs</a>
        </div>

        <div class="image-container" id="image-container">
            <img id="cat-image" src="{{.ImageURL}}" alt="Cat Image" data-image-id="{{.ImageID}}">
        </div>

        <div class="footer">
            <form method="POST" action="/cat/favorite">
                <input type="hidden" name="image_id" value="{{.ImageID}}">
                <button type="submit" class="fav-icon">&#9825;</button>
            </form>
            <form method="POST" action="/cat/vote">
                <input type="hidden" name="image_id" value="{{.ImageID}}">
                <button type="submit" name="vote" value="1" class="vote-icon">&#128077;</button>
                <button type="submit" name="vote" value="-1" class="vote-icon">&#128078;</button>
            </form>
        </div>
    </div>

    <!-- JS -->
    <script src="/static/js/catviewer.js"></script>
</body>
</html>
