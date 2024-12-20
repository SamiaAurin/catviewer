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
        <div class="header">
            <a href="#" id="voting-tab" class="active">⬆️⬇️ Voting</a>
            <a href="#" id="breeds-tab">🔍 Breeds</a>
            <a href="#" id="favs-tab">💖 Favs</a>
        </div>

        <!-- Dynamic Content Sections -->
        <div id="voting-section" class="section active">
            <div class="image-container" id="image-container">
                <img id="cat-image" src="{{.ImageURL}}" alt="Cat Image" data-image-id="{{.ImageID}}">
            </div>
            <div class="footer">
                <form method="POST" action="/cat/vote">
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

        <div id="favs-section" class="section" style="display: none;">
            <div class="favorites-container" id="favorites-container">
                <div class="view-icons">
                    <div class="grid-view">
                        <i class="fa-solid fa-th"></i>
                    </div>
                    <div class="bar-view">
                        <i class="fa-solid fa-bars"></i>
                    </div>
                </div>
                <div class="favorites-grid">
                    {{range .Favorites}}
                    <div class="favorite-item">
                        <img src="{{.Image.URL}}" alt="Favorite Cat" />
                    </div>
                    {{end}}
                </div>
            </div>
        </div>
        

        <div id="breeds-section" class="section" style="display: none;">
            <!-- You can add breeds content here as needed -->
            <h2>Breeds Section</h2>
        </div>
        
        
        
    </div>

    <script src="/static/js/catviewer.js"></script>
</body>
</html>
