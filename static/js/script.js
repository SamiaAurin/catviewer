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

/*
document.addEventListener('DOMContentLoaded', function () {
    const dropdown = document.getElementById('search-breed-dropdown');
    const breedName = document.getElementById('breed-name');
    const breedOrigin = document.getElementById('breed-origin');
    const breedId = document.getElementById('breed-id');
    const breedDescription = document.getElementById('breed-description');
    const wikiLink = document.getElementById('wiki-link');
    const breedImagesContainer = document.getElementById('slider-images');
    const sliderDotsContainer = document.querySelector('.slider-dots');
    
    if (!dropdown || !breedName || !breedOrigin || !breedDescription || !wikiLink || !breedImagesContainer || !sliderDotsContainer) {
        console.error("One or more elements are missing in the DOM.");
        return;
    }


    document.getElementById("breeds-tab").addEventListener("click", function () {
        console.log("Breeds tab clicked!");
    
        fetch("/cat/fetch_breeds")
            .then(response => response.json())
            .then(breeds => {
                console.log("Fetched breeds:", breeds);
    
                if (breeds && breeds.length > 0) {
                    // Display the first breed details
                    const firstBreed = breeds[0]; // Get the first breed
                    populateBreedsDropdown(breeds); // Populate the dropdown with breed details
                    displayBreedDetails(firstBreed); // Display details of the first breed
    
                    // Fetch images for the first breed dynamically
                    fetch(`/cat/fetch_breeds?id=${firstBreed.id}`)
                        .then(response => response.json())
                        .then(data => {
                            if (data.BreedImages) {
                                displayBreedImages(data.BreedImages); // Display images for the first breed
                            } else {
                                console.error("No images found for the first breed.");
                            }
                        })
                        .catch(error => console.error("Error fetching images for the first breed:", error));
                } else {
                    console.error("No breeds found.");
                }
    
                // Show the breeds section and hide other sections
                document.getElementById("voting-section").style.display = "none";
                document.getElementById("voted-images-section").style.display = "none";
                document.getElementById("favs-section").style.display = "none";
                document.getElementById("breeds-section").style.display = "block";
            })
            .catch(error => console.error("Error fetching breeds:", error));
    });
    
    // Helper function to display breed images
    function displayBreedImages(images) {
        // Clear previous images and dots
        breedImagesContainer.innerHTML = "";
        sliderDotsContainer.innerHTML = "";

        // Display breed images and create dots dynamically
        images.forEach((image, index) => {
            const imgElement = document.createElement('img');
            imgElement.src = image.url;
            imgElement.alt = "Breed Image";
            imgElement.classList.add('slider-img');
            breedImagesContainer.appendChild(imgElement);

            const dotElement = document.createElement('span');
            dotElement.classList.add('dot');
            dotElement.addEventListener('click', () => showSlide(index));
            sliderDotsContainer.appendChild(dotElement);
        });

        // Initialize slider to show the first image
        showSlide(0);
    }
            
    // Function to handle selecting a breed from the dropdown
    dropdown.addEventListener('change', function () {
        
        const selectedBreedId = dropdown.value;

        fetch(`/cat/fetch_breeds?id=${selectedBreedId}`)
            .then(response => response.json())
            .then(data => {
                // Check if response contains breed details and images
                const breed = data.BreedDetails;
                const breedImages = data.BreedImages;
    
                if (breed && breedImages) {
                    // Display breed details
                    breedName.textContent = breed.name;
                    breedOrigin.textContent = breed.origin ? `(${breed.origin})` : "";
                    breedId.textContent = breed.id ? `${breed.id}` : "";
                    breedDescription.textContent = breed.description || "No description available.";
                    wikiLink.href = breed.wikipedia_url || "#";
    
                    // Clear previous images and dots
                    breedImagesContainer.innerHTML = "";
                    sliderDotsContainer.innerHTML = "";
    
                    // Display breed images and create dots dynamically
                    breedImages.forEach((image, index) => {
                        const imgElement = document.createElement('img');
                        imgElement.src = image.url;
                        imgElement.alt = "Breed Image";
                        imgElement.classList.add('slider-img');
                        breedImagesContainer.appendChild(imgElement);
    
                        const dotElement = document.createElement('span');
                        dotElement.classList.add('dot');
                        dotElement.addEventListener('click', () => showSlide(index));
                        sliderDotsContainer.appendChild(dotElement);
                    });
    
                    // Initialize slider to show the first image
                    showSlide(0);
                } else {
                    console.error('Breed or images data is missing.');
                }
            })
            .catch(error => console.error('Error fetching breed:', error));
    });
    

    function showSlide(index) {
        const slides = document.querySelectorAll('.slider-img');
        const dots = document.querySelectorAll('.dot');
        
        // Hide all images and deactivate all dots
        slides.forEach((slide, idx) => {
            slide.style.display = idx === index ? 'block' : 'none';
        });
        dots.forEach((dot, idx) => {
            dot.classList.toggle('active', idx === index);
        });
    }

    function populateBreedsDropdown(breeds) {
        // Clear the existing options in the dropdown
        dropdown.innerHTML = "<option value=''>Select a breed</option>"; // Optional placeholder for dropdown
        breeds.forEach(breed => {
            const option = document.createElement("option");
            option.value = breed.id;
            option.textContent = breed.name;
            dropdown.appendChild(option);
        });
    }

    function displayBreedDetails(breed) {
        //breedImage.src = breed.image?.url || "placeholder.jpg";
        breedName.textContent = breed.name || "Unknown Breed";
        breedOrigin.textContent = breed.origin ? `(${breed.origin})` : "";
        breedId.textContent = breed.id ? `${breed.id}` : "";
        breedDescription.textContent = breed.description || "No description available.";
        wikiLink.href = breed.wikipedia_url || "#";
    }
});

*/