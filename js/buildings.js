window.SITE_DATA = {
towers: [
    {
        id: "W3",
        name: "Tòa West 3 - Vinhomes West Point",
        center: { x: -135, y: 0 },
        width: 38,
        depth: 152,
        mainFace: "right",
        floors: 35,
        modeledHeight: 140,
        floorHeight: 140 / 35,
        basementFloors: 3,

        floorInfo: "Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-35: Căn hộ chung cư",
        units: 442,
        density: "14 căn hộ/mặt sàn",
        elevators: "7 thang thường + 1 thang hàng",
        apartmentTypes: "Studio, 2PN, 3PN, 4PN",
        areaRange: "35.6 - 139.1 m²",
        handoverStandard: "Liền tường",
        deliveryTime: "Tháng 4/2020",

        mechanicalHeight: 6,
        antennaHeight: 10,
        glassColor: [60, 80, 110, 0.78],
        slabColor: [230, 226, 218, 1],
        pillarColor: [200, 196, 188, 1],
        crownColor: [170, 174, 180, 1]
    },
    {
        id: "W2",
        name: "Tòa West 2 - Vinhomes West Point",
        center: { x: 20, y: 55 },
        width: 152,
        depth: 38,
        mainFace: "bottom",
        floors: 39,
        modeledHeight: 140,
        floorHeight: 140 / 39,
        basementFloors: 3,

        floorInfo: "Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-4: Trung tâm thương mại & Officetel; Tầng 5-37: Căn hộ chung cư; Tầng 38-39: Căn hộ Duplex",
        units: 452,
        density: "Tầng 5-18: 12 căn hộ/mặt sàn; Tầng 19-39: 13 căn hộ/mặt sàn",
        elevators: "7 thang thường + 1 thang hàng",
        apartmentTypes: "Studio, 2PN, 3PN, 4PN, Duplex",
        areaRange: "39.7 - 285.2 m²",
        handoverStandard: "Liền tường",
        deliveryTime: "Tháng 4/2020",

        mechanicalHeight: 6,
        antennaHeight: 10,
        glassColor: [70, 95, 130, 0.78],
        slabColor: [230, 226, 218, 1],
        pillarColor: [200, 196, 188, 1],
        crownColor: [170, 174, 180, 1]
    },
    {
        id: "W1",
        name: "Tòa West 1 - Vinhomes West Point",
        center: { x: 114, y: -55 },
        width: 152,
        depth: 38,
        mainFace: "top",
        floors: 39,
        modeledHeight: 140,
        floorHeight: 140 / 39,
        basementFloors: 3,

        floorInfo: "Tầng 1: Sảnh Lounge & shop thương mại; Tầng 2-4: Trung tâm thương mại & Officetel; Tầng 5-37: Officetel; Tầng 38-39: Căn hộ Duplex",
        units: 407,
        density: "Tầng 5-37: 11 căn hộ/mặt sàn; Tầng 38-39: 13 căn hộ/mặt sàn",
        elevators: "7 thang thường + 1 thang hàng",
        apartmentTypes: "Studio, 2PN, 3PN, 4PN, Duplex",
        areaRange: "35.9 - 221.8 m²",
        handoverStandard: "Liền tường",
        deliveryTime: "Tháng 4/2020",

        mechanicalHeight: 6,
        antennaHeight: 10,
        glassColor: [55, 75, 105, 0.78],
        slabColor: [230, 226, 218, 1],
        pillarColor: [200, 196, 188, 1],
        crownColor: [170, 174, 180, 1]
    }
],
    // Khối đế bao quanh 3 tòa
    podium: {
        id: "PODIUM",
        name: "Khối đế thương mại - tiện ích Vinhomes West Point",
        center: { x: 18, y: -20 },
        width: 500,
        depth: 250,

        floors: 4,
        floorHeight: 4.5,
        modeledHeight: 18,

        usage: "Khối đế thương mại, dịch vụ và tiện ích nội khu",
        floorInfo: "Tầng 1: Shophouse, sảnh Lounge và shop thương mại; Tầng 2: Khu mua sắm, Vinmart/Vinmart+ và tiện ích giáo dục; Tầng 3-4: Trung tâm thương mại, dịch vụ, officetel và khu tiện ích",
        mainFacilities: "Trung tâm thương mại Vincom, siêu thị Vinmart/Vinmart+, shophouse chân đế, sảnh Lounge, phòng gym, sân chơi trẻ em, vườn cảnh quan",
        rooftopFacilities: "Bể bơi bốn mùa mái kính, vườn nổi, khu thư giãn và không gian cảnh quan trên khối đế",
        parking: "Kết nối 3 tầng hầm để xe của dự án",
        serviceRole: "Phục vụ nhu cầu mua sắm, sinh hoạt, giải trí và kết nối tiện ích cho cư dân W1, W2, W3",
        note: "Thông tin tiện ích được mô phỏng từ các nguồn công khai; chi tiết mặt bằng từng tầng có thể khác theo bản vẽ vận hành thực tế.",

        glassColor: [110, 140, 170, 0.7],
        slabColor: [215, 175, 120, 1],
        baseColor: [195, 155, 105, 1]
    },

    // Sân nội khu / mặt bằng tổng thể
    plaza: {
        width: 315,
        depth: 255,
        thickness: 0.4,
        color: [205, 210, 218, 1],

        walkways: [
            // trục đường giữa W3 và cụm W1/W2
            { x: -5, y: -15, w: 18, d: 175, color: [180, 185, 190, 1] },

            // đường phía dưới gần Đỗ Đức Dục
            { x: 25, y: -120, w: 230, d: 20, color: [170, 175, 180, 1] },

            // đường phía phải gần Phạm Hùng
            { x: 145, y: 5, w: 18, d: 215, color: [170, 175, 180, 1] },

            // khu hồ bơi ở giữa
            { x: -8, y: -20, w: 42, d: 32, color: [80, 170, 210, 0.85] },

            // mảng sân vườn trung tâm
            { x: 55, y: -15, w: 85, d: 60, color: [120, 180, 110, 1] },

            // sân/lối đi trước W1
            { x: 73, y: 8, w: 105, d: 16, color: [225, 215, 195, 1] },

            // sân/lối đi trước W2
            { x: -36, y: 8, w: 110, d: 16, color: [225, 215, 195, 1] },

            // sân/lối đi cạnh W3
            { x: -70, y: -66, w: 16, d: 110, color: [225, 215, 195, 1] }
        ]
    }
};