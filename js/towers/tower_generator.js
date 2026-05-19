// js/towers/tower_generator.js
window.WPModules = window.WPModules || {};

window.WPModules.buildTower = function(t, baseZ) {
    const engine = window.WPEngine;
    const mat = window.WPMaterials;

    // --- KIỂM TRA BẬT TẮT THÁP CỦA TỪNG PHÂN KHU (W1, W2, W3) ---
    if (t.id === "W1" && !window.WPConfig.towerW1) return;
    if (t.id === "W2" && !window.WPConfig.towerW2) return;
    if (t.id === "W3" && !window.WPConfig.towerW3) return;

    const cx = t.center.x;
    const cy = t.center.y;
    const totalH = t.modeledHeight || (t.floors * t.floorHeight);

    const towerAttributes = {
        id: t.id, name: t.name, floors: t.floors,
        modeledHeight: t.modeledHeight, floorInfo: t.floorInfo,
        units: t.units, density: t.density
    };

    const towerPopup = {
        title: "{name}",
        content: [{
            type: "fields",
            fieldInfos: [
                { fieldName: "id", label: "Mã tòa" },
                { fieldName: "floors", label: "Số tầng nổi" },
                { fieldName: "modeledHeight", label: "Chiều cao (m)" },
                { fieldName: "floorInfo", label: "Phân bổ tầng" },
                { fieldName: "units", label: "Tổng số căn" }
            ]
        }]
    };

    const R = 1.5; // Bo góc nhẹ cho luxury tower

    // --- 0.9 LÕI BÊ TÔNG TRUNG TÂM (Central Concrete Core) ---
    // Trục vách kỹ thuật và lõi thang máy bằng bê tông đặc ở trung tâm tòa nhà,
    // ngăn chặn hoàn toàn việc nhìn xuyên qua tòa tháp từ bên ngoài.
    engine.addBox({
        x: cx, y: cy, z: baseZ,
        width: t.width * 0.45, depth: t.depth * 0.45, height: totalH,
        color: mat.light_concrete.color, roughness: mat.light_concrete.roughness,
        attributes: { name: `Vách kỹ thuật tháp ${t.id}` },
        type: "tower_core"
    });

    // --- 1. THÂN KÍNH CHÍNH (Core Glass Structure) ---
    // Chỉ vẽ kính chính nếu lớp tháp tổng được bật
    engine.addBox({
        x: cx, y: cy, z: baseZ,
        width: t.width - 2*R - 0.4, depth: t.depth - 0.4, height: totalH,
        color: mat.residential_glass.color, 
        metallic: mat.residential_glass.metallic, 
        roughness: mat.residential_glass.roughness,
        attributes: towerAttributes, popup: towerPopup,
        type: "tower_glass"
    });
    engine.addBox({
        x: cx, y: cy, z: baseZ,
        width: t.width - 0.4, depth: t.depth - 2*R - 0.4, height: totalH,
        color: mat.residential_glass.color, 
        metallic: mat.residential_glass.metallic, 
        roughness: mat.residential_glass.roughness,
        attributes: towerAttributes, popup: towerPopup,
        type: "tower_glass"
    });

    // Bo góc bằng Cylinder bằng nhôm
    if (engine.addCylinder && window.WPConfig.towerAluminum) {
        [[-1, -1], [1, -1], [-1, 1], [1, 1]].forEach(([sx, sy]) => {
            engine.addCylinder({
                x: cx + sx * (t.width / 2 - R),
                y: cy + sy * (t.depth / 2 - R),
                z: baseZ,
                radius: R, height: totalH,
                color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
                attributes: towerAttributes, popup: towerPopup,
                type: "tower_aluminum"
            });
        });
    }

    // --- 2. VIỀN NGANG CHIA TẦNG (Slabs) ---
    for (let f = 1; f <= t.floors; f++) {
        const z = baseZ + f * t.floorHeight;
        engine.addBox({
            x: cx, y: cy, z: z,
            width: t.width, depth: t.depth, height: 0.3,
            color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
            type: "tower_slab"
        });
    }

    // --- 3. PHÂN MẢNG MẶT ĐỨNG VÀ LOGIA (Mullions & Balconies) ---
    const mx = Math.floor(t.width / 6); 
    const my = Math.floor(t.depth / 6); 

    for(let i=1; i<mx; i++) {
        const fx = -t.width/2 + (i*t.width)/mx;
        
        // Sọc dọc thanh nhôm trang trí (Mullions)
        if (window.WPConfig.towerAluminum) {
            engine.addBox({
                x: cx + fx, y: cy + t.depth/2, z: baseZ,
                width: 1.2, depth: 0.6, height: totalH,
                color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
                type: "tower_aluminum"
            });
            engine.addBox({
                x: cx + fx, y: cy - t.depth/2, z: baseZ,
                width: 1.2, depth: 0.6, height: totalH,
                color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
                type: "tower_aluminum"
            });
        }
        
        // Mô phỏng Balcony lõm (Logia xen kẽ)
        if (i % 2 === 0 && window.WPConfig.towerBalcony) {
            engine.addBox({
                x: cx + fx - (t.width/mx)/2, y: cy + t.depth/2 - 0.2, z: baseZ,
                width: t.width/mx - 1.2, depth: 0.8, height: totalH,
                color: mat.balcony_recess.color, roughness: mat.balcony_recess.roughness,
                type: "tower_balcony"
            });
        }
    }

    for(let j=1; j<my; j++) {
        const fy = -t.depth/2 + (j*t.depth)/my;
        
        if (window.WPConfig.towerAluminum) {
            engine.addBox({
                x: cx + t.width/2, y: cy + fy, z: baseZ,
                width: 0.6, depth: 1.2, height: totalH,
                color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
                type: "tower_aluminum"
            });
            engine.addBox({
                x: cx - t.width/2, y: cy + fy, z: baseZ,
                width: 0.6, depth: 1.2, height: totalH,
                color: mat.white_aluminum.color, roughness: mat.white_aluminum.roughness,
                type: "tower_aluminum"
            });
        }
    }

    // --- 4. MÁI VÀ KHỐI KỸ THUẬT (Mechanical Penthouse & Antenna) ---
    const crownZ = baseZ + totalH;
    
    // Lan can mái (Parapet)
    engine.addBox({
        x: cx, y: cy, z: crownZ,
        width: t.width, depth: t.depth, height: 1.5,
        color: mat.parapet.color, roughness: mat.parapet.roughness,
        type: "tower_slab"
    });

    // Lõi kỹ thuật thang máy (Service Core)
    engine.addBox({
        x: cx, y: cy, z: crownZ + 1.5,
        width: t.width - 10, depth: t.depth - 6, height: t.mechanicalHeight,
        color: mat.light_concrete.color, roughness: 0.8,
        type: "tower_slab"
    });

    // Ăng-ten (Antenna)
    engine.addBox({
        x: cx, y: cy, z: crownZ + 1.5 + t.mechanicalHeight,
        width: 1.0, depth: 1.0, height: t.antennaHeight,
        color: [200, 205, 210, 1], metallic: 0.8, roughness: 0.2,
        type: "tower_slab"
    });
};
