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