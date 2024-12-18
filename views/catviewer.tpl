<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Cat Viewer</title>
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
    </style>
</head>
<body>
    <div class="container">
        <!-- Header -->
        <div class="header">
            <a href="#" class="active">⬆️⬇️ Voting</a>
            <a href="#">🔍 Breeds</a>
            <a href="#"><i>&#9825;</i> <!-- Heart -->Favs</a>
        </div>

        <!-- Image -->
        <div class="image-container" id="image-container">
            <img src="https://cdn2.thecatapi.com/images/0XYvRd7oD.jpg" alt="Cat Image">
        </div>

        <!-- Footer -->
        <div class="footer">
            <i>&#9825;</i> <!-- Heart -->
            <div>
                <i id="thumbs-up" class="vote-icon" onclick="vote('up')">&#128077;</i> <!-- Thumbs Up -->
                <i id="thumbs-down" class="vote-icon" onclick="vote('down')">&#128078;</i> <!-- Thumbs Down -->
            </div>
        </div>

        
    </div>


    <!-- Add a script to handle the voting actions -->
    <!-- At the end of the body -->
    <script>
        // Function to fetch random cat images
        function fetchImages() {
            // Fetch images from the API (replace YOUR_API_KEY with your actual key)
            fetch('https://api.thecatapi.com/v1/images/search?limit=10', {
                headers: {
                    'x-api-key': 'YOUR_API_KEY'  // Replace with your actual API key
                }
            })
            .then(response => response.json())
            .then(images => {
                const imageContainer = document.getElementById('image-container');
                if (!imageContainer) {
                    console.error('Image container not found');
                    return;
                }
                imageContainer.innerHTML = '';  // Clear previous images

                // Loop through images and display them
                images.forEach(image => {
                    const imgElement = document.createElement('img');
                    imgElement.src = image.url;
                    imgElement.alt = 'Cat Image';
                    imgElement.classList.add('cat-image');
                    imgElement.dataset.imageId = image.id;  // Store image ID for voting
                    imageContainer.appendChild(imgElement);
                });
            })
            .catch(error => {
                console.error('Error fetching images:', error);
                alert('Error fetching images.');
            });
        }

        // Load images when the page loads
        window.onload = function() {
            fetchImages();  // Fetch images on page load
        };
    </script>




</body>
</html>
