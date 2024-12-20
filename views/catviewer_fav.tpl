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
                <a href="#" class="active">⬆️⬇️ Voting</a>
                <a href="#" >🔍 Breeds</a>
                <a href="#" >💖 Favs</a>
            </div>

            <!-- Image Section -->
            <div class="image-container" id="image-container">
                <!-- Favorites Section -->
                <section class="favorites-container" id="favorites-container">
                    <div class="view-icons">
                        <div class="grid-view">
                            <i class="fa-solid fa-th"></i>
                        </div>
                        <div class="bar-view">
                            <i class="fa-solid fa-bars"></i>
                        </div>
                    </div>
                    <div class="favorites-grid">
                        {{range .favorites}}
                        <div class="favorite-item">
                            <img src="{{.Image.URL}}" alt="Favorite Cat" />
                        </div>
                        {{end}}
                    </div>
                </section>

            </div>
        </div>

        <script src="/static/js/catviewer.js"></script>
    </body>
</html>
