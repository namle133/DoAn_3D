// js/roads/outer_roads.js
window.WPModules = window.WPModules || {};

window.WPModules.buildPlaza = function(p) {
    const engine = window.WPEngine;
    const mat = window.WPMaterials;

    const cx = p.center ? p.center.x : 0;
    const cy = p.center ? p.center.y : 0;
    
    // ==========================================
    // 1. NỀN DỰ ÁN & ĐƯỜNG GIAO THÔNG CHÍNH (Base Roads) - LỚP 0 & 1
    // ==========================================
    if (window.WPConfig.ENABLE_BASE_ROAD_PAVEMENT) {
        // Nền sân dự án (màu xám đậm để tạo tương phản cực cao cho đường trắng) - Layer 0 (Z: 2.10, cao 20cm)
        engine.addBox({
            x: cx, y: cy, z: 2.10,
            width: 172, depth: 114, height: 0.20,
            color: [100, 105, 110, 1], roughness: 0.9,
            attributes: { name: "Nền sân dự án (Xám đậm)" }
        });

        // Đường Đỗ Đức Dục (Phía dưới) - Layer 0 (Z: 2.10, cao 20cm)
        engine.addBox({
            x: cx, y: -70, z: 2.10,
            width: p.width, depth: 10, height: 0.20,
            color: [45, 48, 50, 1], roughness: 0.9,
            attributes: { name: "Đường Đỗ Đức Dục" }
        });

        // Đường Phạm Hùng (Phía phải) - Layer 0 (Z: 2.10, cao 20cm)
        engine.addBox({
            x: 88, y: cy, z: 2.10,
            width: 14, depth: p.depth, height: 0.20,
            color: [45, 48, 50, 1], roughness: 0.9,
            attributes: { name: "Đường Phạm Hùng" }
        });

        // Vỉa hè & Lối đi nội bộ bao quanh dự án (Màu xám vừa) - Đặt nổi ở Z: 2.30, cao 20cm
        const roadColor = [75, 80, 85, 1];
        
        // Lối đi bên trái (Tây)
        engine.addBox({ x: -100, y: 0, z: 2.30, width: 12, depth: 130, height: 0.20, color: roadColor });
        // Lối đi phía trên (Bắc)
        engine.addBox({ x: -10, y: 60, z: 2.30, width: 170, depth: 10, height: 0.20, color: roadColor });
        // Lối đi bên phải (Đông - song song Phạm Hùng)
        engine.addBox({ x: 74, y: 0, z: 2.30, width: 14, depth: 130, height: 0.20, color: roadColor });
        // Lối đi phía dưới (Nam)
        engine.addBox({ x: -10, y: -60, z: 2.30, width: 170, depth: 10, height: 0.20, color: roadColor });
    }

    // ==========================================
    // 2. HỆ THỐNG ĐƯỜNG NỘI BỘ TRẮNG SÁNG (White Internal Lanes) - LỚP 2
    // ==========================================
    const whiteRoadColor = [255, 255, 255, 1]; // Màu trắng tinh khiết nổi bật cực kỳ đẹp mắt
    const rw = 2; // Chiều rộng tiêu chuẩn làn đường (2m) - cho nét vẽ mượt mà, thanh mảnh giống hình mẫu
    const rh = 0.40; // Chiều cao 40cm để chống Z-fighting triệt để trong chế độ Global
    const rz = 2.50; // Cao độ Z an toàn chống Z-fighting tuyệt đối (Layer 2)
    
    // --- 2a. ĐƯỜNG VÒNG BIÊN NGOÀI (Outer Loop Lanes & Corner Chamfers) ---
    if (window.WPConfig.ENABLE_WHITE_OUTER_LOOP) {
        // Đường biên ngang phía Nam (chạy ngang từ x: -90 đến x: 76)
        engine.addBox({ x: -7, y: -48, z: rz, width: 166, depth: rw, height: rh, color: whiteRoadColor });
        
        // Đường biên ngang phía Bắc (chạy ngang từ x: -90 đến x: 76)
        engine.addBox({ x: -7, y: 48, z: rz, width: 166, depth: rw, height: rh, color: whiteRoadColor });
        
        // Đường biên dọc phía Tây (chạy dọc từ y: -48 đến y: 48)
        engine.addBox({ x: -90, y: 0, z: rz, width: rw, depth: 96, height: rh, color: whiteRoadColor });
        
        // --- BÁN NGUYỆT UỐN LỚN PHÍA TÂY (West U-Turn Loop) ---
        // Nhánh ngang trên của U-turn: y: 12, từ x: -90 sang trái đến x: -100
        engine.addBox({ x: -95, y: 12, z: rz, width: 10, depth: rw, height: rh, color: whiteRoadColor });
        // Nhánh dọc của U-turn: x: -100, từ y: -12 đến y: 12
        engine.addBox({ x: -100, y: 0, z: rz, width: rw, depth: 24, height: rh, color: whiteRoadColor });
        // Nhánh ngang dưới của U-turn: y: -12, từ x: -90 sang trái đến x: -100
        engine.addBox({ x: -95, y: -12, z: rz, width: 10, depth: rw, height: rh, color: whiteRoadColor });

        // --- ĐƯỜNG CONG MỊN LÕM VÀO PHÍA ĐÔNG (Smooth Inward-Curving Eastern Bypass) ---
        // Thiết kế cong cosine lõm sâu vào phía khối đế để uốn lượn tuyệt đẹp đúng tỉ lệ
        const numSegments = 16;
        for (let i = 0; i < numSegments; i++) {
            const y1 = -48 + (96 / numSegments) * i;
            const y2 = -48 + (96 / numSegments) * (i + 1);
            
            // Công thức cong cosine uốn lượn mượt mà tránh hoàn toàn khối đế (y: 0 -> x: 70, y: 48 -> x: 76)
            const x1 = 76 - 6 * Math.cos(Math.PI * y1 / 96);
            const x2 = 76 - 6 * Math.cos(Math.PI * y2 / 96);
            
            const cx = (x1 + x2) / 2;
            const cy = (y1 + y2) / 2;
            const dx = x2 - x1;
            const dy = y2 - y1;
            const len = Math.sqrt(dx * dx + dy * dy);
            const angle = Math.atan2(dx, dy) * 180 / Math.PI;
            
            engine.addBox({
                x: cx, y: cy, z: rz,
                width: rw, depth: len + 0.2, height: rh,
                color: whiteRoadColor,
                heading: angle
            });
        }

        // --- CÁC ĐƯỜNG NHÁNH CHÉO ĐÂM RA NGOÀI DỰ ÁN (Diagonal Exit Paths) ---
        // Nhánh đâm ra góc dưới-trái (Bottom-Left Exit) - Kết nối mượt mà từ (-90, -48) xuống dưới-trái
        engine.addBox({ x: -96, y: -53, z: rz, width: rw, depth: 15.6, height: rh, color: whiteRoadColor, heading: 130 });
        
        // Nhánh đâm ra góc dưới-phải (Bottom-Right Exit) - Kết nối mượt mà từ (76, -48) song song nhánh trái
        engine.addBox({ x: 70, y: -53, z: rz, width: rw, depth: 15.6, height: rh, color: whiteRoadColor, heading: 130 });

        // Đường phụ nhỏ cụt đâm ra lề trái của tháp West 3
        engine.addBox({ x: -94, y: 0, z: rz, width: 8, depth: rw, height: rh, color: whiteRoadColor });
    }

    // --- 2b. ĐƯỜNG TRONG COURTYARD PHÂN NHÁNH THẲNG HÀNG TUYỆT ĐỐI (Courtyard Straight Connectors) ---
    if (window.WPConfig.ENABLE_WHITE_COURTYARD_BRANCHES) {
        // Đường trục dọc trung tâm giữa West 3 và khối đế: x: -48, chạy suốt từ Nam ra Bắc
        engine.addBox({ x: -48, y: 0, z: rz, width: rw, depth: 96, height: rh, color: whiteRoadColor });

        // Nhánh kết nối ngang phía Bắc (bao quanh đầu tháp West 3): y: 35, chạy ngang từ x: -90 sang x: -48
        engine.addBox({ x: -69, y: 35, z: rz, width: 42, depth: rw, height: rh, color: whiteRoadColor });

        // Nhánh kết nối ngang phía Nam (bao quanh chân tháp West 3): y: -35, chạy ngang từ x: -90 sang x: -48
        engine.addBox({ x: -69, y: -35, z: rz, width: 42, depth: rw, height: rh, color: whiteRoadColor });
    }
 
    // --- 3. LỐI VÀO PHÍA NAM & ĐƯỜNG HẦM ĐỖ XE (Do Duc Duc Entrances & Parking Lane) ---
    if (window.WPConfig.ENABLE_WHITE_SOUTH_ENTRANCES) {
        // Dốc đi xuống của hầm đỗ xe góc Tây Nam (SW Parking Ramp): x: -54, từ y: -35 xuống y: -10
        engine.addBox({ x: -54, y: -22.5, z: rz, width: rw, depth: 25, height: rh, color: [200, 200, 200, 1] });

        // Dốc đi xuống của hầm đỗ xe góc Đông Nam (SE Parking Ramp): x: -15, chạy từ y: -48 đến y: -38
        engine.addBox({ x: -15, y: -43, z: rz, width: rw, depth: 10, height: rh, color: whiteRoadColor });

        // Cổng chào chính từ Đỗ Đức Dục (Double-lane main entrance with green divider strip)
        // Làn trái vào: x: -37, chạy từ y: -58 đến y: -48
        engine.addBox({ x: -37, y: -53, z: rz, width: rw, depth: 10, height: rh, color: whiteRoadColor });
        // Làn phải ra: x: -33, chạy từ y: -58 đến y: -48
        engine.addBox({ x: -33, y: -53, z: rz, width: rw, depth: 10, height: rh, color: whiteRoadColor });
    }

    // ==========================================
    // 3. SẢNH ĐÓN PHẠM HÙNG & BỂ CẢNH QUAN (Plaza Lobby & Fountains) - LỚP 2b
    // ==========================================
    if (window.WPConfig.ENABLE_PLAZA_LOBBY) {
        // Sân đón trả khách chính mặt Phạm Hùng (Granite sáng)
        engine.addBox({
            x: 65, y: 0, z: 2.50, // Layer 2 (Z: 2.50, cao 20cm)
            width: 16, depth: 60, height: 0.20,
            color: [220, 225, 230, 1], roughness: 0.4,
            attributes: { name: "Sảnh đón trả khách Phạm Hùng" }
        });

        // Đường vát chéo lối vào từ Phạm Hùng (Chevron driveways)
        engine.addBox({ x: 77, y: 35, z: 2.50, width: 14, depth: 14, height: 0.20, color: [220, 225, 230, 1] });
        engine.addBox({ x: 77, y: -35, z: 2.50, width: 14, depth: 14, height: 0.20, color: [220, 225, 230, 1] });

        // Dãy bể phun nước dọc Phạm Hùng
        const basins = [
            { y: 35, d: 35 },
            { y: -35, d: 35 }
        ];

        basins.forEach((b, idx) => {
            // Viền đá granite đen bao quanh bể
            engine.addBox({
                x: 81, y: b.y, z: 2.50, // Layer 2 (Z: 2.50, cao 20cm)
                width: 5, depth: b.d, height: 0.20,
                color: [50, 55, 60, 1], roughness: 0.2
            });
            // Lòng hồ nước xanh
            engine.addBox({
                x: 81, y: b.y, z: 2.60, // Đặt hơi nổi hẳn để nhìn thấy nước (Z: 2.60, cao 10cm)
                width: 4, depth: b.d - 1, height: 0.10,
                color: [40, 130, 190, 0.95], roughness: 0.05, metallic: 0.3,
                attributes: { name: `Bể phun nước tiểu cảnh ${idx + 1}` }
            });
        });
    }

    // ==========================================
    // 4. TIỂU CẢNH SÂN CHƠI & THẢM CỎ NỘI KHU (Amenities & Playgrounds) - LỚP 3
    // ==========================================
    if (window.WPConfig.ENABLE_AMENITIES_GROUND) {
        // Sân chơi trẻ em (Item 25) - Giữa W3 và W1 (Đặt chuẩn bên trái đường dọc x: -48)
        engine.addBox({
            x: -54.5, y: -12, z: 2.70, // Layer 3 (Z: 2.70, cao 20cm)
            width: 9.5, depth: 18, height: 0.20,
            color: [235, 140, 60, 1], // Cao su đàn hồi màu cam đất ấm
            roughness: 0.8,
            attributes: { name: "Sân chơi trẻ em trung tâm" }
        });
        
        // Các mảng hoa văn màu sắc trên sân chơi (Xanh dương, vàng)
        engine.addBox({ x: -56.5, y: -15, z: 2.90, width: 3.5, depth: 6, height: 0.02, color: [40, 120, 180, 1] });
        engine.addBox({ x: -52.5, y: -8, z: 2.90, width: 3.5, depth: 6, height: 0.02, color: [230, 190, 50, 1] });

        // Thảm cỏ bao quanh tháp W3 (Góc trái) - Shifted to x: -96 to avoid overlapping white road at x: -90
        engine.addBox({ x: -96, y: 0, z: 2.70, width: 4, depth: 85, height: 0.20, color: mat.grass.color, roughness: mat.grass.roughness });
        
        // Vẽ bồn cỏ dưới dạng thảm cỏ tròn phẳng hoàn toàn trên mặt đất - Shifted to x: -96
        const planters = [
            { x: -96, y: 40, r: 2 },
            { x: -96, y: 0, r: 2 },
            { x: -96, y: -40, r: 2 }
        ];
        planters.forEach((p, idx) => {
            engine.addCylinder({
                x: p.x, y: p.y, z: 2.70, // Layer 3 (Z: 2.70, cao 20cm)
                radius: p.r, height: 0.20,
                color: mat.grass.color, roughness: mat.grass.roughness,
                attributes: { name: `Thảm cỏ tròn tiểu cảnh ${idx + 1}` }
            });
        });

        // Đường dạo đá sa thạch vàng nhẹ ấm áp (Walkways) - Repositioned inside courtyard
        const walkwayColor = [225, 215, 195, 1];
        engine.addBox({ x: -48, y: 15, z: 2.70, width: 10, depth: 3, height: 0.20, color: walkwayColor });
        engine.addBox({ x: -45, y: 0, z: 2.70, width: 3, depth: 20, height: 0.20, color: walkwayColor });
    }

    // ==========================================
    // 7. HÀNG RÀO CÂY XANH TẦNG 1 (Ground Trees)
    // ==========================================
    if (window.WPConfig.ENABLE_VEGETATION && p.largeTrees && engine.WebStyleSymbol && engine.Graphic && engine.Point) {
        const largeTree = new engine.WebStyleSymbol({
            name: "Quercus",
            styleName: "EsriRealisticTreesStyle"
        });
        p.largeTrees.forEach(t => {
            const { longitude, latitude } = window.localToLngLat(t.x, t.y);
            engine.layer.add(new engine.Graphic({
                geometry: new engine.Point({ longitude, latitude, z: 0 }),
                symbol: largeTree,
                elevationInfo: { mode: "on-the-ground" }
            }));
        });
    }

    if (window.WPConfig.ENABLE_VEGETATION && p.mediumTrees && engine.WebStyleSymbol && engine.Graphic && engine.Point) {
        const mediumTree = new engine.WebStyleSymbol({
            name: "Citrus",
            styleName: "EsriRealisticTreesStyle"
        });
        p.mediumTrees.forEach(t => {
            const { longitude, latitude } = window.localToLngLat(t.x, t.y);
            engine.layer.add(new engine.Graphic({
                geometry: new engine.Point({ longitude, latitude, z: 0 }),
                symbol: mediumTree,
                elevationInfo: { mode: "on-the-ground" }
            }));
        });
    }
};

window.WPModules.buildStreetlights = function(p) {
    const engine = window.WPEngine;
    const mat = window.WPMaterials;

    if (!engine || !mat) return;

    // Chỉ vẽ nếu được bật trong cấu hình
    if (!window.WPConfig.ENABLE_ROADS) return;

    // Kiểm tra chế độ Ban đêm từ UI
    const isNight = document.getElementById("chkNightMode") && document.getElementById("chkNightMode").checked;

    const streetlightPositions = [];

    // 1. Dọc đường Phạm Hùng (Phía Đông, x = 79.5)
    for (let y = -50; y <= 55; y += 15) {
        // Tránh đặt đè lên lối vào chéo
        if (Math.abs(y - 35) < 5 || Math.abs(y + 35) < 5) continue;
        streetlightPositions.push({ x: 79.5, y: y, z: 2.50, armDir: -1, angle: 0 }); // Hướng cần đèn sang Tây (trái)
    }

    // 2. Dọc đường Đỗ Đức Dục (Phía Nam, y = -52.5)
    for (let x = -80; x <= 65; x += 15) {
        // Tránh hầm đỗ xe
        if (Math.abs(x + 54) < 5 || Math.abs(x + 15) < 5 || Math.abs(x + 35) < 5) continue;
        streetlightPositions.push({ x: x, y: -52.5, z: 2.50, armDir: 1, angle: 90 }); // Hướng cần đèn sang Bắc (lên)
    }

    // 3. Dọc đường phía Tây (x = -93)
    for (let y = -45; y <= 45; y += 15) {
        streetlightPositions.push({ x: -93, y: y, z: 2.50, armDir: 1, angle: 0 }); // Hướng cần đèn sang Đông (phải)
    }

    streetlightPositions.forEach((pos, idx) => {
        // --- 3a. CỘT THÉP ĐÈN (Steel Pole) ---
        engine.addBox({
            x: pos.x, y: pos.y, z: pos.z,
            width: 0.16, depth: 0.16, height: 5.5,
            color: [70, 75, 80, 1], roughness: 0.7,
            attributes: { name: `Cột đèn thép ${idx + 1}` },
            type: "street_light"
        });

        // --- 3b. CẦN ĐÈN CHÉO (Diagonal Arm) ---
        // Vươn ngang ra đường
        const armLength = 1.6;
        const armX = pos.x + (pos.angle === 0 ? pos.armDir * (armLength / 2) : 0);
        const armY = pos.y + (pos.angle === 90 ? pos.armDir * (armLength / 2) : 0);
        
        engine.addBox({
            x: armX, y: armY, z: pos.z + 5.5,
            width: pos.angle === 0 ? armLength : 0.08,
            depth: pos.angle === 90 ? armLength : 0.08,
            height: 0.08,
            color: [70, 75, 80, 1], roughness: 0.7,
            attributes: { name: `Cần đèn ${idx + 1}` },
            type: "street_light"
        });

        // --- 3c. BÓNG ĐÈN PHÁT SÁNG (Illuminated Lamp Head) ---
        const lampX = pos.x + (pos.angle === 0 ? pos.armDir * armLength : 0);
        const lampY = pos.y + (pos.angle === 90 ? pos.armDir * armLength : 0);
        const lampColor = isNight ? [255, 255, 200, 1] : [220, 220, 220, 1]; // Phát sáng vàng ấm ban đêm

        engine.addBox({
            x: lampX, y: lampY, z: pos.z + 5.3,
            width: 0.3, depth: 0.3, height: 0.15,
            color: lampColor, metallic: 0.2, roughness: 0.1,
            attributes: { name: `Bóng đèn ${idx + 1}` },
            type: "street_light"
        });
    });
};

