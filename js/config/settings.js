// js/config/settings.js
window.WPConfig = {
    // Toggles cho từng lớp sàn dưới mặt đất (Ground Layers) để tách biệt xác nhận
    ENABLE_BASE_ROAD_PAVEMENT: true, // Lớp 0 & 1: Nền sân xám nhạt, đường Đỗ Đức Dục/Phạm Hùng, vỉa hè xám sẫm
    ENABLE_WHITE_OUTER_LOOP: true,          // Lớp 2a: Đường vòng biên ngoài trắng & các góc vát chéo
    ENABLE_WHITE_COURTYARD_BRANCHES: true,  // Lớp 2b: Các đường nhánh nội khu trắng (sát W2, W3)
    ENABLE_WHITE_SOUTH_ENTRANCES: true,     // Lớp 2c: Lối vào phía Nam (Đỗ Đức Dục) & dốc hầm trắng
    ENABLE_PLAZA_LOBBY: false,         // Lớp 2b: Sảnh đón trả khách Phạm Hùng & các bể phun nước tiểu cảnh (Tắt để tránh màu sắc mờ)
    ENABLE_AMENITIES_GROUND: false,   // Lớp 3: Sân chơi trẻ em màu cam, thảm cỏ West 3, lối dạo sa thạch vàng (Tắt để tránh màu sắc mờ)

    // Toggles cho các cấu trúc phía trên
    ENABLE_TERRAIN: true,
    ENABLE_ROADS: true,             // Nút tổng để bật/tắt toàn bộ cảnh quan mặt sàn
    ENABLE_PODIUM: true,
    ENABLE_TOWERS: false,
    ENABLE_ROOFTOP_AMENITIES: false,
    ENABLE_VEGETATION: false,

    // Debug settings
    DEBUG_MODE: false
};
