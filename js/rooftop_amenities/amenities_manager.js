// js/rooftop_amenities/amenities_manager.js
window.WPModules = window.WPModules || {};

window.WPModules.buildRooftopAmenities = function(pData, roofZ) {
    const engine = window.WPEngine;
    const mat = window.WPMaterials;

    // Kích thước khu vực Red Box (Tầng mái khối đế) trong mặt bằng:
    // Từ W2 (trên, y:20) tới W1 (dưới, y:-20).
    // Bắt đầu ngang W2 (x:-30) đến ngang W1 (x:70).
    
    // 1. Tấm đệm base của lõi mái (Grey deck)
    // Tăng độ dày lên 0.2 để dễ thấy, không bị engine cull vì quá mỏng
    engine.addBox({
        x: 20, y: 0, z: roofZ,
        width: 100, depth: 40, height: 0.2,
        color: [210, 215, 220, 1], // Xám nhạt
        attributes: { name: "Đường dạo chính tầng mái" }
    });

    // ==========================================
    // KHU BÊN TRÁI: BỂ BƠI BỐN MÙA (3)
    // Tâm X=-5, rộng 35, sâu 24
    // ==========================================
    
    // Viền bờ hồ màu sáng (Pool Deck border)
    engine.addBox({
        x: -5, y: 0, z: roofZ + 0.2,
        width: 37, depth: 26, height: 0.2,
        color: [240, 240, 240, 1], 
        attributes: { name: "Bờ hồ" }
    });
    
    // Lòng bể bơi (Pool Basin - 3)
    engine.addBox({
        x: -5, y: 0, z: roofZ + 0.25,
        width: 35, depth: 24, height: 0.2,
        color: [80, 160, 210, 0.9], // Xanh dương
        metallic: 0.2, roughness: 0.2,
        attributes: { name: "Bể bơi bốn mùa" }
    });

    // Mái kính che bể bơi
    engine.addBox({
        x: -5, y: 0, z: roofZ + 3.0,
        width: 39, depth: 28, height: 0.3,
        color: [180, 210, 240, 0.4], // Kính trong suốt màu xanh nhạt
        metallic: 0.8, roughness: 0.1,
        attributes: { name: "Mái kính bể bơi" }
    });
    
    // Khung mái che (Trụ đỡ kính)
    [[-23, -13], [-23, 13], [13, -13], [13, 13]].forEach(pos => {
        engine.addBox({
            x: pos[0], y: pos[1], z: roofZ + 0.4,
            width: 0.8, depth: 0.8, height: 2.6,
            color: [200, 200, 200, 1]
        });
    });

    // ==========================================
    // KHU BÊN PHẢI: THẢM CỎ & CÁC BỆ (7, 11, 13, 14)
    // ==========================================
    
    // Thảm cỏ hữu cơ lớn (Organic Lawn - 7)
    // Dùng nhiều cylinder đan xen để tạo khối liền mạch (từ x=15 đến x=55)
    if (engine.addCylinder) {
        engine.addCylinder({ x: 20, y: 0, z: roofZ + 0.25, radius: 10, height: 0.2, color: mat.grass.color, roughness: mat.grass.roughness });
        engine.addCylinder({ x: 30, y: 0, z: roofZ + 0.25, radius: 14, height: 0.2, color: mat.grass.color, roughness: mat.grass.roughness, attributes: { name: "Thảm cỏ trung tâm" }});
        engine.addCylinder({ x: 38, y: 8, z: roofZ + 0.25, radius: 10, height: 0.2, color: mat.grass.color, roughness: mat.grass.roughness });
        engine.addCylinder({ x: 40, y: -6, z: roofZ + 0.25, radius: 10, height: 0.2, color: mat.grass.color, roughness: mat.grass.roughness });
        engine.addCylinder({ x: 48, y: 2, z: roofZ + 0.25, radius: 8, height: 0.2, color: mat.grass.color, roughness: mat.grass.roughness });
    }

    // Đường dạo cong (8) - Dải trắng uốn lượn cắt ngang
    engine.addBox({ x: 35, y: 0, z: roofZ + 0.28, width: 3, depth: 28, height: 0.2, color: [230, 230, 230, 1] });

    // Khu thiền yoga (11) - Vòng tròn xám (top left cỏ)
    if (engine.addCylinder) {
        engine.addCylinder({ x: 18, y: 15, z: roofZ + 0.3, radius: 3, height: 0.2, color: [180, 190, 190, 1], attributes: { name: "Khu thiền yoga" } });
    }

    // Khu vực ngắm cảnh thành phố (14) - Mặt sàn gỗ hình tròn góc trên phải
    if (engine.addCylinder) {
        engine.addCylinder({ x: 55, y: 12, z: roofZ + 0.3, radius: 4, height: 0.2, color: [180, 130, 80, 1], roughness: mat.wood_deck.roughness, attributes: { name: "Khu ngắm cảnh thành phố" } });
    }

    // Sân chơi trẻ em (13) - Lót cao su màu đỏ/đỏ đất góc dưới phải
    if (engine.addCylinder) {
        engine.addCylinder({ x: 50, y: -5, z: roofZ + 0.3, radius: 3, height: 0.2, color: [220, 80, 80, 1], roughness: mat.rubber_floor.roughness, attributes: { name: "Sân chơi trẻ em" } });
        // Đồ chơi (màu vàng)
        engine.addBox({ x: 50, y: -5, z: roofZ + 0.5, width: 1.5, depth: 1.5, height: 1.2, color: [255, 200, 50, 1] });
    }

    // Các băng ghế gỗ (Ghế nghỉ - 9)
    [ {x: 28, y: 16}, {x: 35, y: -16}, {x: 8, y: 12} ].forEach(pos => {
        engine.addBox({ x: pos.x, y: pos.y, z: roofZ + 0.5, width: 1.5, depth: 0.6, height: 0.4, color: [150, 100, 50, 1], attributes: { name: "Ghế nghỉ" } });
    });

    // ==========================================
    // CÂY XANH TẦNG MÁI (TIỂU CẢNH ĐỎ / XANH LÁ)
    // ==========================================
    const treesData = (pData && pData.plaza) ? pData.plaza.rooftopTrees : null;
    if (treesData && engine.WebStyleSymbol && engine.Graphic && engine.Point) {
        const greenShrub = new engine.WebStyleSymbol({
            name: "Rhus", // Bụi cây xanh
            styleName: "EsriRealisticTreesStyle"
        });
        const redShrub = new engine.WebStyleSymbol({
            name: "Prunus", // Cây có sắc hoa/lá đỏ
            styleName: "EsriRealisticTreesStyle"
        });

        treesData.forEach(t => {
            const isRed = t.type === "red";
            const { longitude, latitude } = window.localToLngLat(t.x, t.y);
            engine.layer.add(new engine.Graphic({
                geometry: new engine.Point({ longitude, latitude, z: Number(roofZ) + 0.2 }),
                symbol: isRed ? redShrub : greenShrub,
                elevationInfo: { mode: "relative-to-ground" }
            }));
        });
    }
};
