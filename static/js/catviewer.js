// To toggle the active class dynamically between tabs when clicked
document.addEventListener("DOMContentLoaded", function () {
    
    const tabs = document.querySelectorAll("#voting-tab, #breeds-tab, #favs-tab");

    tabs.forEach(tab => {
        tab.addEventListener("click", function (event) {
        
            if (this.id === "voting-tab") {
                // Let it navigate to "/cat/vote"
                return; // Skip toggling for navigation
            } else {
                event.preventDefault(); 
            }

            tabs.forEach(t => t.classList.remove("active"));

            this.classList.add("active");
        });
    });
});

/////////////////// JS for VOTES STARTS //////////////////////////////////

// Fetch and display voted images   //
document.getElementById("votedPicsBtn").addEventListener("click", function () {
    fetch('/cat/voted_pics') // Fetch voted images and values
        .then(response => response.json())
        .then(data => {
            // Hide the voting section and show the voted images section
            document.getElementById("voting-section").style.display = "none";
            document.getElementById("voted-images-section").style.display = "block";

            const grid = document.getElementById("voted-images-grid");
            grid.innerHTML = ""; 

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
                voteInfo.textContent = vote.value === 1 ? "😺 Upvoted" : "😿 Downvoted"; // Display vote type
                voteInfo.className = vote.value === 1 ? "upvote" : "downvote"; // Optional for styling
                container.appendChild(voteInfo);

                // Add the container to the grid
                grid.appendChild(container);
            });
        })
        .catch(error => console.error("Error fetching voted images:", error));
});

/////////////////// JS for VOTES  ENDS //////////////////////////////////

////////////////////// JS for BREEDS STARTS //////////////////////////////

document.addEventListener('DOMContentLoaded', function () {
    const elements = {
        dropdown: document.getElementById('search-breed-dropdown'),
        breedName: document.getElementById('breed-name'),
        breedOrigin: document.getElementById('breed-origin'),
        breedId: document.getElementById('breed-id'),
        breedDescription: document.getElementById('breed-description'),
        wikiLink: document.getElementById('wiki-link'),
        breedImagesContainer: document.querySelector('.slider-images'),
        sliderDotsContainer: document.querySelector('.slider-dots'),
        sections: {
            breeds: document.getElementById("breeds-section"),
            voting: document.getElementById("voting-section"),
            votedImages: document.getElementById("voted-images-section"),
            favorites: document.getElementById("favs-section")
        }
    };

    if (Object.values(elements).includes(null)) {
        console.error("One or more elements are missing in the DOM.");
        return;
    }

    const { dropdown, breedName, breedOrigin, breedId, breedDescription, wikiLink, breedImagesContainer, sliderDotsContainer, sections } = elements;

    // Switch section visibility
    function showSection(activeSection) {
        Object.values(sections).forEach(section => {
            section.style.display = section === activeSection ? 'block' : 'none';
        });
    }

    // Fetch and display breed data
    function fetchAndDisplayBreedData(breedId) {
        fetch(`/cat/fetch_breeds?id=${breedId}`)
            .then(response => response.json())
            .then(data => {
                const breed = data.BreedDetails;
                const breedImages = data.BreedImages;

                if (breed && breedImages) {
                    displayBreedDetails(breed);
                    displayBreedImages(breedImages);
                } else {
                    console.error("Breed or images data is missing.");
                }
            })
            .catch(error => console.error("Error fetching breed data:", error));
    }

    // Populate dropdown with breeds
    function populateBreedsDropdown(breeds) {
        dropdown.innerHTML = `<option value="" disabled selected>Please select</option>`;
        breeds.forEach(breed => {
            const option = document.createElement("option");
            option.value = breed.id;
            option.textContent = breed.name;
            dropdown.appendChild(option);
        });
    }

    // Display breed details
    function displayBreedDetails(breed) {
        breedName.textContent = breed.name || "Unknown Breed";
        breedOrigin.textContent = breed.origin ? `(${breed.origin})` : "";
        breedId.textContent = breed.id || "";
        breedDescription.textContent = breed.description || "No description available.";
        wikiLink.href = breed.wikipedia_url || "#";
    }

    // Display breed images
    function displayBreedImages(images) {
        breedImagesContainer.innerHTML = "";
        sliderDotsContainer.innerHTML = "";

        images.forEach((image, index) => {
            const imgElement = document.createElement('img');
            imgElement.src = image.url;
            imgElement.alt = "Breed Image";
            imgElement.classList.add('slider-img');
            breedImagesContainer.appendChild(imgElement);

            const dotElement = document.createElement('span');
            dotElement.classList.add('dot');
            dotElement.addEventListener('click', () => {
                stopSlideShow();
                showSlide(index);
                startSlideShow();
            });
            sliderDotsContainer.appendChild(dotElement);
        });

        // Reinitialize the slider after adding new slides
        initializeSlider();
    }

    let currentIndex = 0;
    let slideInterval;

    // Show the slide at the given index
    function showSlide(index) {
        const slides = document.querySelectorAll('.slider-img');
        const dots = document.querySelectorAll('.dot');

        if (slides.length === 0) return;

        currentIndex = index < 0 ? slides.length - 1 : index % slides.length;

        const offset = -currentIndex * 100; // Calculate the translateX value
        breedImagesContainer.style.transform = `translateX(${offset}%)`;
        breedImagesContainer.style.transition = 'transform 0.5s ease-in-out'; // Smooth transition

        // Update active dot
        dots.forEach((dot, idx) => {
            dot.classList.toggle('active', idx === currentIndex);
        });
    }

    // Start automatic sliding
    function startSlideShow() {
        slideInterval = setInterval(() => {
            showSlide(currentIndex + 1);
        }, 5000); // Change slide every 3 seconds
    }

    // Stop the automatic sliding
    function stopSlideShow() {
        clearInterval(slideInterval);
    }

    // Initialize the slider
    function initializeSlider() {
        currentIndex = 0;
        showSlide(0);
        startSlideShow();
    }

    // Event listener for breeds tab click
    document.getElementById("breeds-tab").addEventListener("click", function () {
        fetch("/cat/fetch_breeds")
            .then(response => response.json())
            .then(breeds => {
                if (breeds && breeds.length > 0) {
                    const firstBreed = breeds[0];
                    populateBreedsDropdown(breeds);
                    displayBreedDetails(firstBreed);

                    fetchAndDisplayBreedData(firstBreed.id);
                } else {
                    console.error("No breeds found.");
                }

                showSection(sections.breeds);
            })
            .catch(error => console.error("Error fetching breeds:", error));
    });

    // Event listener for dropdown change
    dropdown.addEventListener('change', function () {
        const selectedBreedId = dropdown.value;
        if (selectedBreedId) {
            fetchAndDisplayBreedData(selectedBreedId);
        }
    });

});

////////////////////// JS for BREEDS ENDS //////////////////////////////

/////////////////// JS for FAVS STARTS   //////////////////////////////
/*
document.addEventListener('DOMContentLoaded', function () {

    // Fetch favorite images when the favorites tab is clicked
    document.getElementById("favs-tab").addEventListener("click", function () {
        console.log("Favorites tab clicked!");
        
        // Add event listeners to toggle between grid and bar view
        const gridView = document.querySelector(".grid-view");
        const barView = document.querySelector(".bar-view");
        const grid = document.getElementById("favorite-images-grid");

        if (gridView && barView && grid) {
            // Ensure grid view is applied by default when the Favorites tab is clicked
            grid.classList.remove("bar-view");
            grid.classList.add("grid-view"); // Default to grid view

            // Set the grid-view button as active by default
            gridView.classList.add("active");
            barView.classList.remove("active");

            gridView.addEventListener("click", function () {
                grid.classList.remove("bar-view");
                grid.classList.add("grid-view");

                // Update active state on icons
                gridView.classList.add("active");
                barView.classList.remove("active");
            });

            barView.addEventListener("click", function () {
                grid.classList.remove("grid-view");
                grid.classList.add("bar-view");

                // Update active state on icons
                barView.classList.add("active");
                gridView.classList.remove("active");
            });
        } else {
            console.error("Grid or Bar view icons not found!");
        }

        fetch('/cat/fav_pics') // Fetch favorite images
            .then(response => response.json())
            .then(data => {
                console.log("Fetched favorite images:", data);

                // Hide other sections and show the favorites section
                document.getElementById("voting-section").style.display = "none";
                document.getElementById("voted-images-section").style.display = "none"; 
                document.getElementById("breeds-section").style.display = "none";
                document.getElementById("favs-section").style.display = "block"; 

                grid.innerHTML = ""; 

                data.forEach(fav => {
                    const container = document.createElement("div"); 
                    container.className = "favorite-item";

                    // Create and append the image
                    const img = document.createElement("img");
                    img.src = fav.image.url; 
                    img.alt = "Favorite Cat";
                    img.className = "favorite-image"; 
                    container.appendChild(img);

                    
                    const info = document.createElement("p");
                    info.textContent = `😻 Added: ${new Date(fav.created_at).toLocaleDateString()}`;
                    container.appendChild(info);

                    // Add the container to the grid
                    grid.appendChild(container);
                });
                // Scroll the container to the top after adding images
                grid.scrollTop = 0;
            })
            .catch(error => console.error("Error fetching favorite images:", error));
    });
});
*/

document.addEventListener('DOMContentLoaded', function () {

    // Fetch favorite images when the favorites tab is clicked
    document.getElementById("favs-tab").addEventListener("click", function () {
        console.log("Favorites tab clicked!");
        
        // Add event listeners to toggle between grid and bar view
        const gridView = document.querySelector(".grid-view");
        const barView = document.querySelector(".bar-view");
        const grid = document.getElementById("favorite-images-grid");

        if (gridView && barView && grid) {
            // Ensure grid view is applied by default when the Favorites tab is clicked
            grid.classList.remove("bar-view");
            grid.classList.add("grid-view"); // Default to grid view

            // Set the grid-view button as active by default
            gridView.classList.add("active");
            barView.classList.remove("active");

            gridView.addEventListener("click", function () {
                grid.classList.remove("bar-view");
                grid.classList.add("grid-view");

                // Update active state on icons
                gridView.classList.add("active");
                barView.classList.remove("active");
            });

            barView.addEventListener("click", function () {
                grid.classList.remove("grid-view");
                grid.classList.add("bar-view");

                // Update active state on icons
                barView.classList.add("active");
                gridView.classList.remove("active");
            });
        } else {
            console.error("Grid or Bar view icons not found!");
        }

        fetch('/cat/fav_pics') // Fetch favorite images
            .then(response => response.json())
            .then(data => {
                console.log("Fetched favorite images:", data);

                // Hide other sections and show the favorites section
                document.getElementById("voting-section").style.display = "none";
                document.getElementById("voted-images-section").style.display = "none"; 
                document.getElementById("breeds-section").style.display = "none";
                document.getElementById("favs-section").style.display = "block"; 

                grid.innerHTML = ""; 

                data.forEach(fav => {
                    const container = document.createElement("div"); 
                    container.className = "favorite-item";
                
                    // Create and append the image
                    const img = document.createElement("img");
                    img.src = fav.image.url; 
                    img.alt = "Favorite Cat";
                    img.className = "favorite-image"; 
                    container.appendChild(img);
                
                    // Create and append the information
                    const info = document.createElement("p");
                    info.textContent = `😻 Added: ${new Date(fav.created_at).toLocaleDateString()}`;
                    container.appendChild(info);
                
                    // Create and append the delete button
                    const deleteButton = document.createElement("button");
                    deleteButton.className = "delete-btn"; // Add class for styling
                
                    // Use an icon for the delete button (you can use a font awesome icon, for example)
                    deleteButton.innerHTML = "🗑️"; // You can use any icon here, like a trash bin
                
                    // Add event listener for delete functionality
                    deleteButton.addEventListener("click", function () {
                        // Optionally, remove the image from the DOM
                        container.remove();
                
                        // You can also trigger a server-side request to delete the favorite image from the database
                        fetch(`/cat/delete_fav/${fav.id}`, { method: "DELETE" })
                            .then(response => {
                                if (response.ok) {
                                    console.log("Image deleted successfully!");
                                } else {
                                    console.log("Error deleting image.");
                                }
                            })
                            .catch(error => {
                                console.error("Error:", error);
                            });
                    });
                
                    container.appendChild(deleteButton); // Add the delete button to the container
                
                    // Add the container to the grid
                    grid.appendChild(container);
                });
                
                // Scroll the container to the top after adding images
                grid.scrollTop = 0;
            })
            .catch(error => console.error("Error fetching favorite images:", error));
    });
});




