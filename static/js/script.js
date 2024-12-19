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

// Array to hold the breed list
let breeds = [];
let selectedBreed = null;

        
// Function to display breed details
function displayBreed(breed) {
            const imageContainer = document.getElementById('image-container');
            const footer = document.querySelector('.footer');
            footer.style.display = 'none'; // Hide the footer

            // Ensure breed.image is defined before accessing breed.image.url
            const breedImageUrl = breed.image ? breed.image.url : 'default-image-url'; // Use a fallback image URL if image is not available

            imageContainer.innerHTML = `
                <section id="breed-details">
                    <div class="breed-container">
                        <div class="breed-select">
                            <img src="${breedImageUrl}" alt="${breed.name}">
                            <h1>${breed.name}</h1>
                            <span>(${breed.origin})</span>
                            <span>${breed.id}</span>
                            <p>${breed.description}</p>
                            <a href="${breed.wikipedia_url}" target="_blank">Wikipedia</a>
                        </div>
                    </div>
                </section>
            `;
    }

// Function to display selected breed details
function selectBreed(breed) {
            selectedBreed = breed; // Set the selected breed
            displayBreed(breed); // Display the breed's details
            const breedsList = document.getElementById('breeds-list');
            //breedsList.innerHTML = ''; // Clear the breed list after selection
        }
        // Function to display the breed list and search bar
function displayBreeds(breedsList) {
            const imageContainer = document.getElementById('image-container');
            const footer = document.querySelector('.footer');
            footer.style.display = 'none'; // Hide the footer when displaying breeds
            

            imageContainer.innerHTML = `
                <section id="breeds">
                    <div class="breed-container">
                        <div class="breed-select">
                            <span class="breed-text">Search for a breed</span>
                            <div class="value-container">
                                <input type="text" id="search-breed-input" placeholder="Search breeds" class="search-breed-input">
                                <button id="close-button" class="close-button">×</button>
                            </div>
                        </div>
                        <ul id="breeds-list" class="breed-list"></ul>
                        <div class="breed-details" id="breed-details">
                            
                        </div>
                    </div>
                </section>
            `;
            
            const breedsListElement = document.getElementById('breeds-list');
            breedsList.forEach(breed => {
                const li = document.createElement('li');
                li.textContent = breed.name;
                li.onclick = () => selectBreed(breed); // Handle breed selection
                breedsListElement.appendChild(li);
            });
        }

// Function to set up the search bar functionality
function setupSearchBar() {
            const searchInput = document.getElementById('search-breed-input');
            const breedsListElement = document.getElementById('breeds-list');
            const closeButton = document.getElementById('close-button');

            // Show the breed list when the search input is clicked
            searchInput.addEventListener('focus', () => {
                breedsListElement.innerHTML = ''; // Clear any previous list
                breedsListElement.style.display = 'block'; // Show the breed list

                breeds.forEach(breed => {
                    const li = document.createElement('li');
                    li.textContent = breed.name;
                    li.onclick = () => selectBreed(breed); // Select the breed when clicked
                    breedsListElement.appendChild(li);
                });
            });

            // Filter the list based on search input
            searchInput.addEventListener('input', () => {
                const query = searchInput.value.toLowerCase();
                const filteredBreeds = breeds.filter(breed =>
                    breed.name.toLowerCase().includes(query)
                );
                breedsListElement.innerHTML = ''; // Clear the current list
                filteredBreeds.forEach(breed => {
                    const li = document.createElement('li');
                    li.textContent = breed.name;
                    li.onclick = () => selectBreed(breed); // Select the breed when clicked
                    breedsListElement.appendChild(li);
                });
            });

            // Add functionality to the close button
            closeButton.addEventListener('click', () => {
                searchInput.value = ''; // Clear the search input
                breedsListElement.style.display = 'none'; // Hide the breed list
            });

            // Hide the breed list when clicking outside the input or list
            document.addEventListener('click', (e) => {
                if (!searchInput.contains(e.target) && !closeButton.contains(e.target)) {
                    breedsListElement.style.display = 'none'; // Hide the list
                }
            });
        }


// Function to fetch breeds from The Cat API
function fetchBreeds() {
            fetch('https://api.thecatapi.com/v1/breeds', {
                headers: { 'x-api-key': 'YOUR_API_KEY' } // Replace 'YOUR_API_KEY' with your actual API key
            })
            .then(response => {
                if (!response.ok) {
                    console.error('Error fetching breeds:', response.statusText);
                    alert('Error fetching breeds.');
                }
                return response.json(); // Parse the JSON response
            })
            .then(data => {
                breeds = data; // Store the list of breeds
                displayBreeds(breeds); // Display the breed list
                setupSearchBar(); // Set up search functionality for breeds
            })
            .catch(error => {
                console.error('Error fetching breeds:', error);
                alert('Error fetching breeds.');
            });
        }

        
// Show the breed details when clicking on the "Breeds" link
function showBreeds() {
            fetchBreeds(); // Fetch and display the list of breeds when the "Breeds" link is clicked
        }

         

window.onload = fetchRandomImage;