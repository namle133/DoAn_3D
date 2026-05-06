
window.SITE_DATA = {
    // 3 Towers: W1 / W2 / W3
    towers: [
        {
            id: "W1",
            name: "Tháp W1 - Vinhomes West Point",
            center: { x: -55, y: 25 },
            width: 46,
            depth: 38,
            floors: 40,
            floorHeight: 3.5,
            mechanicalHeight: 6,
            antennaHeight: 14,
            glassColor: [60, 80, 110, 0.78],
            slabColor: [230, 226, 218, 1],
            pillarColor: [200, 196, 188, 1],
            crownColor: [170, 174, 180, 1]
        },
        {
            id: "W2",
            name: "Tháp W2 - Vinhomes West Point",
            center: { x: 55, y: 25 },
            width: 46,
            depth: 38,
            floors: 40,
            floorHeight: 3.5,
            mechanicalHeight: 6,
            antennaHeight: 14,
            glassColor: [70, 95, 130, 0.78],
            slabColor: [230, 226, 218, 1],
            pillarColor: [200, 196, 188, 1],
            crownColor: [170, 174, 180, 1]
        },
        {
            id: "W3",
            name: "Tháp W3 - Vinhomes West Point",
            center: { x: 0, y: -45 },
            width: 50,
            depth: 40,
            floors: 42,
            floorHeight: 3.5,
            mechanicalHeight: 6,
            antennaHeight: 14,
            glassColor: [55, 75, 105, 0.78],
            slabColor: [230, 226, 218, 1],
            pillarColor: [200, 196, 188, 1],
            crownColor: [170, 174, 180, 1]
        }
    ],

    // Podium: Surrounding 3 towers
    podium: {
        id: "PODIUM",
        name: "Khối đế TTTM Vincom",
        center: { x: 0, y: -5 },
        width: 180,
        depth: 130,
        floors: 4,
        floorHeight: 4.5,
        glassColor: [110, 140, 170, 0.7],
        slabColor: [215, 175, 120, 1],
        baseColor: [195, 155, 105, 1]
    },

    // Plaza: Surrounding 3 towers
    plaza: {
        width: 260,
        depth: 220,
        thickness: 0.4,
        color: [205, 210, 218, 1],
        // Walkways: Small rectangles on the plaza
        walkways: [
            { x: 0, y: 60, w: 30, d: 60, color: [225, 215, 195, 1] },
            { x: 0, y: -90, w: 30, d: 40, color: [225, 215, 195, 1] }
        ]
    }
};
