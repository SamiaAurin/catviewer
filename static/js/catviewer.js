// Fetch and display voted images  STARTS //
document.getElementById("votedPicsBtn").addEventListener("click", function () {
    fetch('/cat/voted_pics') // Fetch voted images and values
        .then(response => response.json())
        .then(data => {
            // Hide the voting section and show the voted images section
            document.getElementById("voting-section").style.display = "none";
            document.getElementById("voted-images-section").style.display = "block";

            const grid = document.getElementById("voted-images-grid");
            grid.innerHTML = ""; // Clear previous images

            data.forEach(vote => {
                const container = document.createElement("div"); // Container for image + vote info
                container.className = "voted-item";

                // Create and append the image
                const img = document.createElement("img");
                img.src = vote.image.url;
                img.alt = "Voted Cat";
                container.appendChild(img);

                // Create and append the vote info
                const voteInfo = document.createElement("p");
                voteInfo.textContent = vote.value === 1 ? "Upvoted" : "Downvoted"; // Display vote type
                voteInfo.className = vote.value === 1 ? "upvote" : "downvote"; // Optional for styling
                container.appendChild(voteInfo);

                // Add the container to the grid
                grid.appendChild(container);
            });
        })
        .catch(error => console.error("Error fetching voted images:", error));
});
// Fetch and display voted images  ENDS //

////////////////////// JS for BREEDS ////////////////////////////
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
                populateBreedsDropdown(breeds);
                displayBreedDetails(breeds[0]); // Display details of the first breed
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



/////////////////// JS for FAVS /////////////////

document.addEventListener('DOMContentLoaded', function () {
    // Fetch favorite images when the favorites tab is clicked
    document.getElementById("favs-tab").addEventListener("click", function () {
        console.log("Favorites tab clicked!");

        fetch('/cat/fav_pics') // Fetch favorite images
            .then(response => response.json())
            .then(data => {
                console.log("Fetched favorite images:", data);

                // Hide other sections and show the favorites section
                document.getElementById("voting-section").style.display = "none";
                document.getElementById("voted-images-section").style.display = "none"; 
                document.getElementById("breeds-section").style.display = "none";
                document.getElementById("favs-section").style.display = "block"; 

                const grid = document.getElementById("favorite-images-grid");
                grid.innerHTML = ""; // Clear previous images

                data.forEach(fav => {
                    const container = document.createElement("div"); // Container for image
                    container.className = "favorite-item";

                    // Create and append the image
                    const img = document.createElement("img");
                    img.src = fav.image.url; // Use the image URL from the API response
                    img.alt = "Favorite Cat";
                    img.className = "favorite-image"; // Add class for styling if needed
                    container.appendChild(img);

                    // Optionally, add information about the favorite, e.g., date or ID
                    const info = document.createElement("p");
                    info.textContent = `Favorited on: ${new Date(fav.created_at).toLocaleDateString()}`;
                    container.appendChild(info);

                    // Add the container to the grid
                    grid.appendChild(container);
                });
            })
            .catch(error => console.error("Error fetching favorite images:", error));
    });

    // Add event listeners to toggle between grid and bar view
    const gridView = document.querySelector(".grid-view");
    const barView = document.querySelector(".bar-view");

    if (gridView && barView) {
        gridView.addEventListener("click", function () {
            const grid = document.getElementById("favorite-images-grid");
            grid.classList.remove("bar-view");
            grid.classList.add("grid-view");

            // Update active state on icons
            gridView.classList.add("active");
            barView.classList.remove("active");
        });

        barView.addEventListener("click", function () {
            const grid = document.getElementById("favorite-images-grid");
            grid.classList.remove("grid-view");
            grid.classList.add("bar-view");

            // Update active state on icons
            barView.classList.add("active");
            gridView.classList.remove("active");
        });
    } else {
        console.error("Grid or Bar view icons not found!");
    }
});






