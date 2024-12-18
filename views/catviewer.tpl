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
        <div class="image-container">
            <img src="https://cdn2.thecatapi.com/images/0XYvRd7oD.jpg" alt="Cat Image">
        </div>

        <!-- Footer -->
        <div class="footer">
            <i>&#9825;</i> <!-- Heart -->
            <div>
                <i>&#128077;</i> <!-- Thumbs Up -->
                <i>&#128078;</i> <!-- Thumbs Down -->
            </div>
        </div>

        
    </div>
</body>
</html>
