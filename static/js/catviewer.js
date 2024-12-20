document.getElementById('favs-tab').addEventListener('click', function() {
    showFavorites(); // Fetch and display the favorites when the tab is clicked
    document.getElementById('favs-section').style.display = 'block';
    document.getElementById('voting-section').style.display = 'none';
});



