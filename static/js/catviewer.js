// Ensure that showFavorites is defined and accessible when the user clicks
function showFavorites() {
    fetch('/cat/favorites')
        .then(response => response.json()) // Parse the JSON response
        .then(favorites => {
            if (favorites.error) {
                console.error("Error:", favorites.error);
                alert(favorites.error); // Display an alert if there's an error
                return;
            }
            renderFavorites(favorites); // Render the favorites
        })
        .catch(error => {
            console.error("Error fetching favorites:", error);
        });
}



// Render the fetched favorites in the grid
function renderFavorites(favorites) {
    const favoritesGrid = document.querySelector('.favorites-grid');
    const favoritesContainer = document.getElementById('favorites-container');

    if (!favoritesGrid) {
        console.error("Favorites grid element not found");
        return;
    }

    // Clear the grid before rendering
    favoritesGrid.innerHTML = '';

    if (favorites.length === 0) {
        favoritesGrid.innerHTML = '<p>No favorites found. Add some from the voting section!</p>';
    } else {
        favorites.forEach(favorite => {
            const favoriteItem = document.createElement('div');
            favoriteItem.classList.add('favorite-item');

            const img = document.createElement('img');
            img.src = favorite.Image.URL;
            img.alt = 'Favorite Cat';

            favoriteItem.appendChild(img);
            favoritesGrid.appendChild(favoriteItem);
        });
    }

    // Display the favorites container
    favoritesContainer.style.display = 'block';
}


// Grid view toggle functionality
document.addEventListener('DOMContentLoaded', function () {
    document.querySelector('.grid-view').addEventListener('click', function () {
        const grid = document.querySelector('.favorites-grid');
        grid.classList.add('grid-view');
        grid.classList.remove('bar-view');
    });
    
    document.querySelector('.bar-view').addEventListener('click', function () {
        const grid = document.querySelector('.favorites-grid');
        grid.classList.add('bar-view');
        grid.classList.remove('grid-view');
    });
});



