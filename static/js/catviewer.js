// Array to store favorite images
let favoriteImages = [];

function saveFavorite() {
    const currentImage = document.querySelector('#image-container img');
    
    if (currentImage) {
        const imageId = currentImage.dataset.imageId;

        // Call the API to save the image as a favorite
        const body = JSON.stringify({
            image_id: imageId,
            sub_id: "user-123" // Unique user ID (can be dynamic or stored)
        });
        

        fetch("/add_to_favorites", {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: body
        })
        .then(response => response.json())
        .then(data => {
            if (data.id) {
                // Image added to favorites successfully
                favoriteImages.push({ id: imageId, url: currentImage.src });
                alert('Image added to favorites!');
                fetchRandomImage(); // Fetch a new image after favouriting
            } else {
                alert('Failed to add image to favorites.');
            }
        })
        .catch(error => {
            console.error("Error saving favorite:", error);
            alert('Failed to save image as favorite.');
        });
    }
}

function showFavorites() {
    const favoriteContainer = document.getElementById('favorite-container');
    favoriteContainer.innerHTML = '';

    if (favoriteImages.length > 0) {
        favoriteImages.forEach(image => {
            const imgElement = document.createElement('img');
            imgElement.src = image.url;
            imgElement.alt = 'Favorite Cat Image';
            favoriteContainer.appendChild(imgElement);
        });
    } else {
        favoriteContainer.innerHTML = 'No favorites yet.';
    }
}

function vote(action) {
    console.log(`Voted ${action}`);
    
    // Perform the vote (thumbs up or down)
    const imageId = document.querySelector("#image-container img").getAttribute("data-image-id");
    const voteData = {
        image_id: imageId,
        value: action,
    };

    fetch("/cast_vote", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(voteData),
    })
    .then(response => response.json())
    .then(data => {
        // After voting, fetch a new image
        fetchRandomImage();
    })
    .catch(error => {
        console.error("Error casting vote:", error);
        alert("Failed to cast vote.");
    });
}

function fetchRandomImage() {
    fetch("/random_cat_image", {
        method: "GET",
    })
    .then(response => response.json())
    .then(data => {
        if (data && data.url) {
            const imgElement = document.createElement('img');
            imgElement.src = data.url;
            imgElement.alt = 'Cat Image';
            imgElement.dataset.imageId = data.id;

            const imageContainer = document.getElementById('image-container');
            imageContainer.innerHTML = '';
            imageContainer.appendChild(imgElement);
        } else {
            alert("Error fetching new image.");
        }
    })
    .catch(() => alert('Error fetching a new image.'));
}
