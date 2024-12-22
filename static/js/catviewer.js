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
                document.getElementById("voted-images-section").style.display = "none"; // Hide the voted images section
                document.getElementById("favs-section").style.display = "block"; // Show the favorites section

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






