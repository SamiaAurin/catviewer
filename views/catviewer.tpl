<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Viewer</title>
    <!-- Font Awesome -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css" rel="stylesheet">
    <link href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css" rel="stylesheet">

    <style>
        /* Container Styling */
        .container {
            width: 100%;
            max-width: 600px;
            margin: auto;
            border: 1px solid #ddd;
            border-radius: 10px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
            overflow: hidden;
            font-family: Arial, sans-serif;
        }

        /* Header Navigation */
        .header {
            display: flex;
            justify-content: space-around;
            border-bottom: 1px solid #eee;
            padding: 10px 0;
        }

        .header a {
            text-decoration: none;
            font-weight: bold;
            color: #666;
        }

        .header a.active {
            color: red;
        }

        /* Image Styling */
        .image-container {
            display: flex;
            justify-content: center;
            background-color: #f9f9f9;
        }

        .image-container img {
            max-width: 100%;
            height: auto;
            display: block;
        }

        /* Footer Icons */
        .footer {
            display: flex;
            justify-content: space-between;
            padding: 10px 15px;
        }

        .footer i {
            font-size: 1.2em;
            color: #666;
            cursor: pointer;
        }

        .footer i:hover {
            color: red;
        }

        /* General container styling */
        .favorites-container {
            padding: 1rem;
            background-color: #f9f9f9;
        }

        /* Flex container for the icons */
        .view-icons {
            display: flex;
            justify-content: flex-start;
            gap: 1rem;
            padding: 1rem 0;
        }

        /* Styling for the icon containers */
        .grid-view, .bar-view {
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 0.5rem;
            border-radius: 5px;
            transition: background-color 0.3s ease;
        }

        .grid-view:hover, .bar-view:hover {
            background-color: #e0e0e0;
        }

        /* Icons styling */
        i {
            font-size: 1.5rem; /* Adjust the size of icons */
        }

        /* Grid view layout for the favorites */
        .favorites-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
            gap: 1rem;
            padding: 1rem;
            background-color: #ffffff;
            transition: all 0.3s ease;
        }

        /* Responsive grid for smaller screens */
        @media (max-width: 768px) {
            .favorites-grid {
                grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); /* More responsive grid */
            }
        }

        /* Bar (scrolling) view layout */
        .favorites-grid.scroll-view {
            display: flex;
            flex-direction: column; /* Make it a column layout */
            overflow-y: auto; /* Enable vertical scroll */
            gap: 1rem;
            padding: 1rem;
            background-color: #ffffff;
            height: 400px; /* Set a fixed height for scrolling */
        }

        /* Make the images in scroll view behave responsively */
        .favorites-grid.scroll-view img {
            width: 100%;
            max-width: 200px;
            object-fit: cover;
            height: 150px; /* Adjust according to your needs */
        }

        /* Hover state for scroll images */
        .favorites-grid.scroll-view img:hover {
            transform: scale(1.05);
        }



        /* Active icon color change */
        .grid-view.active, .bar-view.active {
            background-color: #ff6841;
            color: #ffffff;
        }

        /* Button focus states */
        [role="button"]:focus {
            outline: 2px solid #ff6841;
        }

    </style>
</head>
<body>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <a href="#" class="active">⬆️⬇️ Voting</a>
            <a href="#">🔍 Breeds</a>
            <a href="#" onclick="showFavorites()"><i>&#9825;</i> Favs</a> <!-- Heart -->
        </div>

        <!-- Image -->
        <div class="image-container" id="image-container">
            <img src="https://cdn2.thecatapi.com/images/0XYvRd7oD.jpg" alt="Cat Image">
        </div>

        <!-- Footer -->
        <div class="footer">
            <i onclick="saveFavorite()" class="fav-icon">&#9825;</i> <!-- Heart -->
            <div>
                <i id="thumbs-up" class="vote-icon" onclick="vote('up')">&#128077;</i> <!-- Thumbs Up -->
                <i id="thumbs-down" class="vote-icon" onclick="vote('down')">&#128078;</i> <!-- Thumbs Down -->
            </div>
        </div>
    </div>

    <script>
        let favoriteImages = [];

        function fetchRandomImage() {
            fetch('https://api.thecatapi.com/v1/images/search', {
                headers: { 'x-api-key': 'YOUR_API_KEY' }
            })
            .then(response => response.json())
            .then(images => {
                if (images.length > 0) {
                    const imgElement = document.createElement('img');
                    imgElement.src = images[0].url;
                    imgElement.alt = 'Cat Image';
                    imgElement.dataset.imageId = images[0].id;

                    const imageContainer = document.getElementById('image-container');
                    imageContainer.innerHTML = '';
                    imageContainer.appendChild(imgElement);
                }
            })
            .catch(() => alert('Error fetching an image.'));
        }

        function saveFavorite() {
            const currentImage = document.querySelector('#image-container img');
            if (currentImage) {
                favoriteImages.push({ id: currentImage.dataset.imageId, url: currentImage.src });
                alert('Image added to favorites!');
                fetchRandomImage();
            }
        }

        function showFavorites() {
            const imageContainer = document.getElementById('image-container');
            const footer = document.querySelector('.footer');
            footer.style.display = 'none';

            imageContainer.innerHTML = `
                <section class="favorites-container">
                    <div class="view-icons">
                        <!-- Grid View Icon -->
                        <div class="grid-view">
                            <i class="fa-solid fa-th"></i>
                        </div>
                        <!-- Bar View Icon -->
                        <div class="bar-view">
                            <i class="fa-solid fa-bars"></i>
                        </div>
                    </div>
                    <!-- Grid for displaying favorites -->
                    <div class="favorites-grid"></div>
                </section>
            `;

            const favoritesGrid = document.querySelector('.favorites-grid');
            if (favoriteImages.length === 0) {
                favoritesGrid.innerHTML = '<p>No favorite images to display.</p>';
            } else {
                favoriteImages.forEach(image => {
                    const imgDiv = document.createElement('div');
                    imgDiv.classList.add('aspect-square');
                    imgDiv.innerHTML = `<img src="${image.url}" alt="Favorite Cat Image">`;
                    favoritesGrid.appendChild(imgDiv);
                });
            }

                // Add event listeners for view icons inside the showFavorites() function
                document.querySelector('.grid-view').addEventListener('click', function() {
                // Remove active class from all icons
                document.querySelector('.grid-view').classList.add('active');
                document.querySelector('.bar-view').classList.remove('active');
                
                // Switch to grid layout
                document.querySelector('.favorites-grid').classList.remove('scroll-view');
                document.querySelector('.favorites-grid').classList.add('grid-layout');
                });

                document.querySelector('.bar-view').addEventListener('click', function() {
                    // Remove active class from all icons
                    document.querySelector('.bar-view').classList.add('active');
                    document.querySelector('.grid-view').classList.remove('active');
                    
                    // Switch to scroll layout
                    document.querySelector('.favorites-grid').classList.remove('grid-layout');
                    document.querySelector('.favorites-grid').classList.add('scroll-view');
                });
        }



        function vote(action) {
            console.log(`Voted ${action}`);
            fetchRandomImage();
        }

        



        window.onload = fetchRandomImage;
    </script>
</body>
</html>
