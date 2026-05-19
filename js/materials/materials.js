// js/materials/materials.js
window.WPMaterials = {
    // 1. Facade căn hộ cao cấp (Towers)
    residential_glass: { color: [40, 60, 80, 0.95], metallic: 0.7, roughness: 0.1 },     // Kính cửa sổ xanh xám
    white_aluminum: { color: [250, 252, 255, 1], roughness: 0.9 },                       // Khung nhôm trắng sáng
    light_concrete: { color: [240, 240, 240, 1], roughness: 0.95 },                      // Bê tông sáng
    balcony_recess: { color: [20, 25, 30, 1], roughness: 0.8 },                          // Phần lõm ban công (đổ bóng tối)
    parapet: { color: [220, 225, 230, 1], roughness: 0.8 },                              // Lan can mái

    // 2. Khối đế thương mại (Podium)
    commercial_glass: { color: [180, 210, 240, 0.6], metallic: 0.8, roughness: 0.1 },    // Kính trong suốt có phản chiếu
    podium_base: { color: [230, 235, 240, 1], roughness: 0.85 },                         // Bê tông khối đế

    // 3. Landscape & Rooftop Amenities
    pavement: { color: [200, 205, 210, 1], roughness: 0.95 },                            // Đá lát quảng trường
    wood_deck: { color: [160, 110, 70, 1], roughness: 0.85 },                            // Sàn gỗ hồ bơi
    pool_water: { color: [30, 100, 160, 0.8], metallic: 0.9, roughness: 0.05 },          // Màu nước hồ bơi cơ bản
    grass: { color: [90, 150, 80, 1], roughness: 0.95 },                                 // Thảm cỏ
    rubber_floor: { color: [200, 80, 80, 1], roughness: 0.9 },                           // Sàn cao su khu vui chơi
    pathway: { color: [220, 215, 205, 1], roughness: 0.9 }                               // Đường dạo
};
