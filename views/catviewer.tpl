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
            <a href="http://localhost:8080/cat/vote" id="voting-tab" class="active">⬆️⬇️ Voting</a>
            <a href="#" id="breeds-tab">🔍 Breeds</a>
            <a href="#" id="favs-tab">💖 Favs</a>
        </div>

        <!-- Dynamic Content Sections -->
        <div id="voting-section" class="section active">
            <div class="image-container" id="voting-image-container">
                <img id="cat-image" src="{{.ImageURL}}" alt="Cat Image" data-image-id="{{.ImageID}}">
            </div>
            <div class="footer">

                <!-- Favorite Button (Heart Icon) -->
                <form method="POST" action="/cat/favorite">
                    <input type="hidden" name="image_id" value="{{.ImageID}}">
                    <button type="submit" name="fav" class="fav-icon">💖</button> 
                </form>

                <!-- Voting Buttons -->
                <form method="POST" action="/cat/vote">
                    <input type="hidden" name="image_id" value="{{.ImageID}}">
                    <button type="submit" name="vote" value="1" class="vote-icon" id="upvote">&#128077;</button>
                    <button type="submit" name="vote" value="-1" class="vote-icon" id="downvote">&#128078;</button>
                    <button type="button" id="votedPicsBtn" class="vote-icon">Voted Pics</button>
                </form>

            </div>
        </div>
        <!-- Add a Voted Images Container -->
        <div id="voted-images-section" class="voted-image-container" style="display: none;">
            <h2>Voted Cat Images</h2>
            <div id="voted-images-grid" class="voted-grid"></div>
        </div> 

        <!-- favs -->
        <div id="favs-section" class="section" style="display: none;">
            <div class="favorites-container" id="favorites-container">
                <div class="view-icons">
                    <div class="grid-view active">
                        <i class="fa-solid fa-th"></i>
                    </div>
                    <div class="bar-view">
                        <i class="fa-solid fa-bars"></i>
                    </div>
                </div>
                <div id="favorite-images-grid" class="favorite-images-grid">
                    <!-- The images will be shown here -->
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
