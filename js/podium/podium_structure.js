// js/podium/podium_structure.js
window.WPModules = window.WPModules || {};

window.WPModules.buildPodium = function(pData) {
    const engine = window.WPEngine;
    const mat = window.WPMaterials;

    if (!engine || !mat) return;

    const baseZ = 2.70;
    const slabHeight = 0.80;
    const modeledHeight = pData.modeledHeight || 22.5;

    // --- 1. SÀN MÓNG / FOOTPRINTS (Layer 3b) ---
    if (window.WPConfig.ENABLE_TERRAIN) {
        // Footprint Khối đế chính West 1 & 2
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: baseZ,
            width: pData.width, depth: pData.depth, height: slabHeight,
            color: [195, 180, 165, 1], roughness: 0.9,
            attributes: { name: "Sàn Footprint Khối đế chính West 1 & 2" },
            type: "podium_structure"
        });

        // Footprint West 3
        engine.addBox({
            x: -69, y: 0, z: baseZ,
            width: 18, depth: 50, height: slabHeight,
            color: [185, 175, 165, 1], roughness: 0.9,
            attributes: { name: "Sàn Footprint West 3" },
            type: "podium_structure"
        });
    }

    // --- 2. KẾT CẤU CỐT ĐẾ (Podium Concrete Columns & Shear Walls) ---
    if (window.WPConfig.podiumStructure) {
        // Cột chịu lực góc Đông Bắc
        engine.addBox({
            x: pData.center.x + pData.width/2 - 1.5, y: pData.center.y + pData.depth/2 - 1.5, z: baseZ + slabHeight,
            width: 3, depth: 3, height: modeledHeight - slabHeight,
            color: mat.light_concrete.color, roughness: 0.8,
            attributes: { name: "Cột chịu lực Đông Bắc" },
            type: "podium_structure"
        });
        // Cột chịu lực góc Tây Bắc
        engine.addBox({
            x: pData.center.x - pData.width/2 + 1.5, y: pData.center.y + pData.depth/2 - 1.5, z: baseZ + slabHeight,
            width: 3, depth: 3, height: modeledHeight - slabHeight,
            color: mat.light_concrete.color, roughness: 0.8,
            attributes: { name: "Cột chịu lực Tây Bắc" },
            type: "podium_structure"
        });
        // Cột chịu lực góc Đông Nam
        engine.addBox({
            x: pData.center.x + pData.width/2 - 1.5, y: pData.center.y - pData.depth/2 + 1.5, z: baseZ + slabHeight,
            width: 3, depth: 3, height: modeledHeight - slabHeight,
            color: mat.light_concrete.color, roughness: 0.8,
            attributes: { name: "Cột chịu lực Đông Nam" },
            type: "podium_structure"
        });
        // Cột chịu lực góc Tây Nam
        engine.addBox({
            x: pData.center.x - pData.width/2 + 1.5, y: pData.center.y - pData.depth/2 + 1.5, z: baseZ + slabHeight,
            width: 3, depth: 3, height: modeledHeight - slabHeight,
            color: mat.light_concrete.color, roughness: 0.8,
            attributes: { name: "Cột chịu lực Tây Nam" },
            type: "podium_structure"
        });

        // Lõi kỹ thuật thang máy trung tâm khối đế
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: baseZ + slabHeight,
            width: 12, depth: 10, height: modeledHeight - slabHeight,
            color: mat.light_concrete.color, roughness: 0.8,
            attributes: { name: "Lõi thang máy trung tâm khối đế" },
            type: "podium_structure"
        });
    }

    // --- 3. MẶT VÁCH KÍNH ĐẾ (Podium Glass Facade) ---
    if (window.WPConfig.podiumGlass) {
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: baseZ + slabHeight,
            width: pData.width - 2.0, depth: pData.depth - 2.0, height: modeledHeight - slabHeight,
            color: mat.commercial_glass.color,
            metallic: mat.commercial_glass.metallic,
            roughness: mat.commercial_glass.roughness,
            attributes: { name: "Mặt vách kính khối đế thương mại" },
            type: "podium_glass"
        });
    }

    // --- 4. SÀN MÁI TẦNG 5 (Podium Roof & Parapet) ---
    if (window.WPConfig.podiumRoof) {
        // Tấm sàn mái chính ở độ cao 22.5m (dày 1.5m)
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: baseZ + modeledHeight,
            width: pData.width + 1.0, depth: pData.depth + 1.0, height: 1.5,
            color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
            attributes: { name: "Sàn mái khối đế" },
            type: "podium_roof"
        });

        // Viền lan can an toàn quanh mái (Parapet)
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: baseZ + modeledHeight + 1.5,
            width: pData.width + 1.0, depth: pData.depth + 1.0, height: 0.8,
            color: [240, 240, 240, 1], roughness: 0.9,
            attributes: { name: "Lan can mái khối đế" },
            type: "podium_roof"
        });
    }

    // --- 5. HOẠT CẢNH HOẠT ĐỘNG BÊN TRONG TRUNG TÂM THƯƠNG MẠI (Interior Shops & Shoppers) ---
    // Khối đế có vách kính trong suốt (opacity 0.6) nên có thể nhìn thấy hoạt cảnh bên trong cực kỳ sinh động!
    
    // 5a. Sàn phân tầng nội bộ (Interior Floor Slabs)
    const floorSlabs = [
        { z: 8.0, name: "Sàn tầng 2" },
        { z: 12.5, name: "Sàn tầng 3" },
        { z: 17.0, name: "Sàn tầng 4" },
        { z: 21.5, name: "Sàn tầng 5" }
    ];
    floorSlabs.forEach(slab => {
        engine.addBox({
            x: pData.center.x, y: pData.center.y, z: slab.z,
            width: pData.width - 6.0, depth: pData.depth - 6.0, height: 0.20,
            color: [235, 235, 235, 1], roughness: 0.9,
            attributes: { name: `${slab.name} (Khối đế)` },
            type: "podium_interior"
        });
    });

    // 5b. Gian hàng thương mại (Retail Shops)
    const shops = [
        // Tầng 1 (z = 3.50)
        { name: "Siêu thị WinMart", x: pData.center.x - 30, y: pData.center.y - 18, z: 3.50, w: 22, d: 15, h: 4.0, color: [230, 40, 40, 1] }, // Đỏ rực
        { name: "Starbucks Coffee", x: pData.center.x + 25, y: pData.center.y + 15, z: 3.50, w: 14, d: 12, h: 4.0, color: [10, 110, 60, 1] }, // Xanh lá Starbucks
        { name: "Sảnh Lounge Khối đế", x: pData.center.x - 5, y: pData.center.y + 20, z: 3.50, w: 20, d: 10, h: 4.0, color: [230, 210, 180, 1] }, // Vàng kem lịch lãm
        { name: "Highlands Coffee", x: pData.center.x - 20, y: pData.center.y + 18, z: 3.50, w: 14, d: 12, h: 4.0, color: [120, 30, 30, 1] }, // Đỏ đô Highlands

        // Tầng 2 (z = 8.20)
        { name: "Cửa hàng Thời trang Zara", x: pData.center.x - 28, y: pData.center.y + 10, z: 8.20, w: 18, d: 14, h: 4.0, color: [30, 30, 30, 1] }, // Đen Zara sang trọng
        { name: "Cửa hàng Apple Store", x: pData.center.x + 22, y: pData.center.y - 12, z: 8.20, w: 16, d: 12, h: 4.0, color: [240, 240, 240, 1] }, // Trắng Apple
        { name: "Nike Premium Store", x: pData.center.x + 30, y: pData.center.y + 12, z: 8.20, w: 15, d: 12, h: 4.0, color: [20, 20, 20, 1] }, // Đen Nike
        { name: "Uniqlo West Point", x: pData.center.x - 5, y: pData.center.y - 20, z: 8.20, w: 22, d: 14, h: 4.0, color: [235, 10, 30, 1] }, // Đỏ Uniqlo

        // Tầng 3 (z = 12.70)
        { name: "Nhà hàng Haidilao Hotpot", x: pData.center.x - 25, y: pData.center.y - 12, z: 12.70, w: 18, d: 16, h: 4.0, color: [245, 50, 40, 1] }, // Đỏ lẩu Haidilao
        { name: "Nhà hàng Pizza 4P's", x: pData.center.x + 24, y: pData.center.y + 14, z: 12.70, w: 16, d: 14, h: 4.0, color: [160, 140, 100, 1] }, // Nâu ấm gỗ
        { name: "Khu vui chơi trẻ em TiniWorld", x: pData.center.x + 10, y: pData.center.y - 18, z: 12.70, w: 24, d: 15, h: 4.0, color: [30, 150, 240, 1] }, // Xanh cyan vui vẻ

        // Tầng 4 (z = 17.20)
        { name: "Rạp chiếu phim CGV Cinema", x: pData.center.x - 20, y: pData.center.y, z: 17.20, w: 35, d: 24, h: 4.0, color: [180, 20, 20, 1] }, // Đỏ rạp phim
        { name: "Trung tâm California Fitness", x: pData.center.x + 22, y: pData.center.y - 10, z: 17.20, w: 28, d: 20, h: 4.0, color: [40, 40, 45, 1] } // Xám đen Gym
    ];

    shops.forEach(s => {
        engine.addBox({
            x: s.x, y: s.y, z: s.z,
            width: s.w, depth: s.d, height: s.h,
            color: s.color, roughness: 0.6,
            attributes: { 
                name: s.name, 
                info: `Gian hàng tại tầng ${Math.ceil((s.z - baseZ) / 4.5)} của Trung tâm thương mại`
            },
            popupTemplate: {
                title: "{name}",
                content: "{info}"
            },
            type: "podium_shop"
        });
    });

    // 5c. Khách mua sắm / Người đi bộ (Shoppers / People)
    // Tạo ra khoảng 35 người sinh động phân bổ ngẫu nhiên trên các tầng
    const shoppersData = [
        // Tầng 1
        { name: "Khách mua sắm", x: pData.center.x - 15, y: pData.center.y - 5, z: 3.50, color: [255, 80, 80] },
        { name: "Khách mua sắm", x: pData.center.x - 8, y: pData.center.y + 12, z: 3.50, color: [80, 120, 255] },
        { name: "Nhân viên an ninh", x: pData.center.x + 5, y: pData.center.y + 3, z: 3.50, color: [40, 40, 50] },
        { name: "Khách mua sắm", x: pData.center.x + 15, y: pData.center.y + 8, z: 3.50, color: [255, 200, 0] },
        { name: "Nhóm bạn trẻ", x: pData.center.x + 12, y: pData.center.y - 12, z: 3.50, color: [255, 80, 255] },
        { name: "Khách mua sắm", x: pData.center.x - 22, y: pData.center.y - 2, z: 3.50, color: [50, 220, 150] },

        // Tầng 2
        { name: "Khách xem đồ Zara", x: pData.center.x - 18, y: pData.center.y + 10, z: 8.20, color: [255, 100, 100] },
        { name: "Khách trải nghiệm iPhone", x: pData.center.x + 15, y: pData.center.y - 10, z: 8.20, color: [240, 240, 240] },
        { name: "Khách mua sắm", x: pData.center.x + 2, y: pData.center.y + 14, z: 8.20, color: [255, 180, 50] },
        { name: "Cư dân Officetel", x: pData.center.x - 5, y: pData.center.y - 8, z: 8.20, color: [100, 255, 100] },
        { name: "Khách dạo chơi", x: pData.center.x + 24, y: pData.center.y + 5, z: 8.20, color: [100, 180, 255] },

        // Tầng 3
        { name: "Khách chờ bàn Haidilao", x: pData.center.x - 14, y: pData.center.y - 8, z: 12.70, color: [255, 50, 50] },
        { name: "Khách ăn Pizza", x: pData.center.x + 16, y: pData.center.y + 10, z: 12.70, color: [200, 160, 80] },
        { name: "Gia đình đi chơi", x: pData.center.x - 2, y: pData.center.y - 10, z: 12.70, color: [0, 200, 255] },
        { name: "Nhóm trẻ em", x: pData.center.x + 5, y: pData.center.y - 12, z: 12.70, color: [255, 100, 255] },

        // Tầng 4
        { name: "Khách mua vé xem phim", x: pData.center.x - 10, y: pData.center.y + 2, z: 17.20, color: [255, 220, 50] },
        { name: "Khách tập Gym", x: pData.center.x + 12, y: pData.center.y - 5, z: 17.20, color: [100, 100, 120] }
    ];

    // Thêm các shopper ngẫu nhiên để tăng tính sinh động
    for (let i = 0; i < 20; i++) {
        const floorIdx = Math.floor(Math.random() * 4); // 0, 1, 2, 3
        const floorZ = floorIdx === 0 ? 3.50 : 3.50 + floorIdx * 4.5;
        const rx = pData.center.x + (Math.random() - 0.5) * (pData.width - 20);
        const ry = pData.center.y + (Math.random() - 0.5) * (pData.depth - 20);
        
        shoppersData.push({
            name: "Khách tham quan",
            x: rx, y: ry, z: floorZ,
            color: [
                Math.floor(Math.random() * 155) + 100,
                Math.floor(Math.random() * 155) + 100,
                Math.floor(Math.random() * 155) + 100
            ]
        });
    }

    shoppersData.forEach((shop, sIdx) => {
        // Thang người đi bộ (Cylinder nhỏ thẳng đứng tượng trưng cho 1 người)
        engine.addCylinder({
            x: shop.x, y: shop.y, z: shop.z,
            radius: 0.35, height: 1.6, // Bán kính 35cm, chiều cao 1.6 mét
            color: [...shop.color, 1.0],
            attributes: { 
                name: `${shop.name} #${sIdx + 1}`,
                info: `${shop.name} đang tham quan mua sắm tại Vinhomes West Point.`
            },
            popupTemplate: {
                title: "{name}",
                content: "{info}"
            },
            type: "podium_shopper"
        });
    });

    // 5d. Cầu thang cuốn liên tầng (Retail Escalators)
    // Tạo 3 cụm thang cuốn đan chéo ở trung tâm khoảng thông tầng của TTTM
    const escalators = [
        { name: "Thang cuốn Tầng 1 -> Tầng 2", x: pData.center.x, y: -2, z: 3.50 + 2.0, w: 2.2, d: 15, h: 0.2 },
        { name: "Thang cuốn Tầng 2 -> Tầng 3", x: pData.center.x - 4, y: 2, z: 8.20 + 2.0, w: 2.2, d: 15, h: 0.2 },
        { name: "Thang cuốn Tầng 3 -> Tầng 4", x: pData.center.x + 4, y: -2, z: 12.70 + 2.0, w: 2.2, d: 15, h: 0.2 }
    ];

    escalators.forEach(esc => {
        engine.addBox({
            x: esc.x, y: esc.y, z: esc.z,
            width: esc.w, depth: esc.d, height: esc.h,
            color: [70, 80, 95, 1], roughness: 0.3, // Màu thép xám kim loại
            attributes: { name: esc.name },
            type: "podium_escalator"
        });
    });
};
