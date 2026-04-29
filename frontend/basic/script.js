fetch("http://localhost:8080/products")
    .then(response => response.json())
    .then(data => {
        console.log("Products from API:", data);

        const container = document.getElementById("product-list");
        container.innerHTML = "";

        data.forEach((product, index) => {
            // Create tile
            const tile = document.createElement("div");
            tile.className = "product-tile";

            tile.style.animationDelay = `${index * 0.2}s`;

            // Add content
            tile.innerHTML = `
                <h2>${product.name}</h2>
                <p><strong>Brand:</strong> ${product.brand}</p>
                <p><strong>Price:</strong> ₹${product.price}</p>
                <p><strong>Quantity:</strong> ${product.quantity}</p>
            `;

            // Add tile to container
            container.appendChild(tile);
        });
    })
    .catch(error => {
        console.error("Error fetching products:", error);
    });