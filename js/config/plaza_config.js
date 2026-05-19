// js/config/plaza_config.js
window.WP_PLAZA_CONFIG = {
    center: { x: -10, y: 0 },
    width: 210,
    depth: 150,
    thickness: 0.4,
    color: [205, 210, 218, 1],

    walkways: [
        { x: -5, y: -15, w: 18, d: 175, color: [180, 185, 190, 1] },
        { x: 25, y: -100, w: 230, d: 20, color: [170, 175, 180, 1] },
        { x: 145, y: 5, w: 18, d: 215, color: [170, 175, 180, 1] },
        { x: 55, y: -15, w: 85, d: 60, color: [120, 180, 110, 1] },
        { x: 73, y: 8, w: 105, d: 16, color: [225, 215, 195, 1] },
        { x: -36, y: 8, w: 110, d: 16, color: [225, 215, 195, 1] },
        { x: -70, y: -66, w: 16, d: 110, color: [225, 215, 195, 1] }
    ],
    
    pools: [],

    canopies: [
        // Mái sảnh các tòa
        { x: -60, y: 0, w: 6, d: 15, h: 5 }, // W3 canopy
        { x: 15, y: 20, w: 18, d: 6, h: 5 }, // W2 canopy
        { x: 27.5, y: -20, w: 18, d: 6, h: 5 } // W1 canopy
    ],

    largeTrees: [
        // Dọc Đỗ Đức Dục (Dưới cùng, Y: -65), 14 cây
        {x:-85, y:-65}, {x:-75, y:-65}, {x:-65, y:-65}, {x:-55, y:-65}, {x:-45, y:-65},
        {x:-35, y:-65}, {x:-25, y:-65}, {x:-15, y:-65}, {x:-5, y:-65}, {x:5, y:-65},
        {x:15, y:-65}, {x:25, y:-65}, {x:35, y:-65}, {x:45, y:-65},
        
        // Dọc viền trái (W3, X: -95), 9 cây
        {x:-95, y:-45}, {x:-95, y:-30}, {x:-95, y:-15}, {x:-95, y:0}, {x:-95, y:15},
        {x:-95, y:30}, {x:-95, y:45}, {x:-95, y:60}, {x:-95, y:75},

        // Cụm góc ngã tư Đỗ Đức Dục - Phạm Hùng
        {x:85, y:-65}, {x:80, y:-55}, {x:90, y:-55}, {x:85, y:-45},

        // Dọc đường nội khu phía trên
        {x:-50, y:65}, {x:-30, y:65}, {x:-10, y:65}, {x:10, y:65}, {x:30, y:65}, {x:50, y:65}, {x:70, y:65}
    ],

    mediumTrees: [
        // Hàng cây nhỏ dọc đài phun nước (Phạm Hùng, X: 75)
        {x:75, y:-30}, {x:75, y:-20}, {x:75, y:-10}, {x:75, y:0}, {x:75, y:10}, {x:75, y:20}, {x:75, y:30}, {x:75, y:40},
        // Dọc quanh toà W1, W2
        {x:27.5, y:-10}, {x:10, y:-34}, {x:45, y:-34}, // Quanh W1
        {x:15, y:10}, {x:-5, y:34}, {x:35, y:34} // Quanh W2
    ],

    rooftopTrees: [
        // Bờ trên bể bơi
        {x: -20, y: 14, type: "red"}, {x: -15, y: 14, type: "red"}, {x: -10, y: 14, type: "red"}, {x: -5, y: 14, type: "green"},
        // Bờ dưới bể bơi
        {x: -20, y: -14, type: "red"}, {x: -15, y: -14, type: "red"}, {x: -10, y: -14, type: "red"}, {x: -5, y: -14, type: "green"},
        // Cạnh trái bể bơi
        {x: -26, y: 0, type: "red"},
        
        // Quanh viền thảm cỏ (Trên)
        {x: 20, y: 16, type: "green"}, {x: 25, y: 18, type: "red"}, {x: 32, y: 16, type: "green"}, {x: 38, y: 14, type: "green"}, {x: 45, y: 15, type: "red"},
        // Quanh viền thảm cỏ (Dưới)
        {x: 20, y: -16, type: "red"}, {x: 28, y: -14, type: "green"}, {x: 35, y: -16, type: "green"}, {x: 45, y: -14, type: "red"},
        // Nằm giữa thảm cỏ
        {x: 25, y: 5, type: "green"}, {x: 30, y: -5, type: "red"}, {x: 38, y: 0, type: "green"}, {x: 42, y: 5, type: "red"},
        // Cạnh phải ngắm cảnh
        {x: 58, y: 5, type: "red"}, {x: 55, y: -5, type: "green"}
    ]
};
