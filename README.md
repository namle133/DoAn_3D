# 3D Geographic Information System — Vinhomes West Point

A 3D visualization of the **Vinhomes West Point** urban complex located in Pham Hung, Me Tri, Nam Tu Liem, Hanoi, Vietnam.

The project models the three main towers **W1, W2, and W3**, together with the podium, using `SceneView` from the **ArcGIS Maps SDK for JavaScript 4.30**.

This project was developed as an academic project for the **3D Geographic Information Systems (3D GIS)** course.

## Project Structure

```text
3d-geographic-information-system/
├── index.html          # Main page containing the ArcGIS SDK and UI overlays
├── css/
│   └── style.css       # Styles for the title bar and north arrow
├── docs/
│   └── diagrams/       # Architecture diagrams (PNG/SVG)
├── js/
│   ├── buildings.js    # Approximate building coordinates and dimensions
│   └── app.js          # Map, SceneView, Mesh objects, and widgets
└── README.md
```

## Features

### Interactive 3D Map

The application uses `SceneView` from the ArcGIS Maps SDK for JavaScript to provide an interactive 3D visualization of the Vinhomes West Point complex.

### 3D Building Models

The scene contains:

- Tower W1
- Tower W2
- Tower W3
- Podium

The buildings are constructed using `Mesh.createBox` and include several architectural components such as:

- Glass building bodies
- Floor bands
- Corner columns
- Mullions
- Technical structures
- Antennas

Building coordinates and dimensions are approximated for visualization purposes.

### Basemap Selection

The default basemap is **OpenStreetMap**.

Users can switch between several basemap styles:

- OpenStreetMap
- Satellite
- Hybrid
- Topographic
- Streets
- Dark Gray

The basemap selector is located in the upper-right corner of the application.

### Compass and North Arrow

The application provides two orientation tools:

- The **Esri Compass widget** in the upper-left corner. Clicking it resets the scene orientation toward North.
- A custom **North Arrow** in the lower-right corner displaying the four cardinal directions: N, E, S, and W.

The red arrow dynamically remains pointed toward geographic North regardless of camera rotation.

### Map Navigation Tools

The application includes several tools for interacting with the 3D scene:

- Home
- Fullscreen
- Scale Bar
- Coordinate Readout
- Zoom
- Navigation Toggle
- Search

## Technologies

- HTML5
- CSS3
- JavaScript
- ArcGIS Maps SDK for JavaScript 4.30
- ArcGIS `SceneView`
- ArcGIS `Mesh`

## How to Run

### Using Visual Studio Code

1. Clone or download this repository.
2. Open the project folder in **Visual Studio Code**.
3. Install the **Live Server** extension.
4. Right-click `index.html`.
5. Select **Open with Live Server**.
6. The application will open in your web browser.

## Academic Context

This project was developed as an individual academic project for the **3D Geographic Information Systems (3D GIS)** course.

The project demonstrates the use of web-based GIS technologies to create and visualize a three-dimensional urban environment, including 3D building geometry, spatial navigation, basemap switching, coordinate visualization, and orientation tools.

## Author

**Le Thanh Nam** - **Nguyen Huu Duy Dat** - **Nguyen Tan Hoang**

Bachelor of Science in Information Technology  
VNUHCM - University of Information Technology (UIT)
