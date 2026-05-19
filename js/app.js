// =====================================================================
// Vinhomes West Point - LocalScene 3D detailed model
// ArcGIS Maps SDK for JavaScript 4.30
// =====================================================================

require([
    "esri/Map",
    "esri/Basemap",
    "esri/views/SceneView",
    "esri/Graphic",
    "esri/layers/GraphicsLayer",
    "esri/geometry/Point",
    "esri/geometry/Mesh",
    "esri/symbols/WebStyleSymbol",
    "esri/widgets/Compass",
    "esri/widgets/ScaleBar",
    "esri/widgets/Expand",
    "esri/widgets/BasemapGallery",
    "esri/widgets/Home",
    "esri/widgets/Fullscreen"
], function (
    Map, Basemap, SceneView, Graphic, GraphicsLayer, Point, Mesh, WebStyleSymbol,
    Compass, ScaleBar, Expand, BasemapGallery, Home, Fullscreen
) {
    const data = window.SITE_DATA;

    const SITE_LNG = 105.78375;
    const SITE_LAT = 21.01266;
    const M_PER_DEG_LAT = 110940;
    const M_PER_DEG_LNG = 111320 * Math.cos((SITE_LAT * Math.PI) / 180);

    function localToLngLat(x, y) {
        return {
            longitude: SITE_LNG + x / M_PER_DEG_LNG,
            latitude: SITE_LAT + y / M_PER_DEG_LAT
        };
    }

    // Map: OpenStreetMap (không cần API key)
    const map = new Map({
        basemap: "osm",
        ground: { surfaceColor: [220, 225, 230, 1] }
    });

    // SceneView - Căn chỉnh góc nhìn camera gần hơn, tập trung chính xác vào tâm dự án (0, 0)
    const homePos = localToLngLat(140, -140); // Đưa vị trí camera gần hơn nhiều
    const HOME_VIEWPOINT = {
        position: {
            longitude: homePos.longitude,
            latitude: homePos.latitude,
            z: 220 // Độ cao Z thấp hơn (220m thay vì 750m) để nhìn cận cảnh chi tiết sắc nét cực đẹp
        },
        heading: 315, // Hướng nhìn từ Đông Nam chếch sang Tây Bắc
        tilt: 60 // Độ nghiêng 60 độ lý tưởng để ngắm toàn cảnh 3D nổi bật
    };

    const view = new SceneView({
        container: "viewDiv",
        map: map,
        qualityProfile: "high",
        camera: HOME_VIEWPOINT,
        environment: {
            atmosphere: { quality: "high" },
            starsEnabled: true,
            atmosphereEnabled: true,
            lighting: {
                date: new Date("2026-05-06T09:30:00+07:00"),
                directShadowsEnabled: true,
                cameraTrackingEnabled: false
            }
        },
        ui: { components: ["zoom", "navigation-toggle"] }
    });

    view.when(
        () => {
            console.log("[VWP] SceneView ready");
            setupWidgets();
            setupNorthArrow();
        },
        (err) => console.error("[VWP] SceneView error:", err)
    );

    // -----------------------------------------------------------------
    // WIDGETS: Compass, ScaleBar, BasemapGallery, Home, Fullscreen
    // -----------------------------------------------------------------
    function setupWidgets() {
        // Home: bay về vị trí mặc định
        const homeBtn = new Home({ view: view, viewpoint: { targetGeometry: new Point({ longitude: SITE_LNG, latitude: SITE_LAT }), scale: 5000 } });
        homeBtn.goToOverride = function (view, goToParams) {
            return view.goTo(HOME_VIEWPOINT, goToParams.options);
        };
        view.ui.add(homeBtn, "top-left");

        // Compass: la bàn (nhấn để đưa heading về bắc = 0°)
        const compass = new Compass({ view: view });
        view.ui.add(compass, "top-left");

        // Fullscreen
        view.ui.add(new Fullscreen({ view: view }), "top-left");

        // ScaleBar: thước tỉ lệ (km/m)
        view.ui.add(new ScaleBar({ view: view, unit: "metric" }), "bottom-left");

        // BasemapGallery: chọn loại bản đồ (đặt trong Expand cho gọn)
        const basemapGallery = new BasemapGallery({
            view: view,
            source: [
                Basemap.fromId("osm"),
                Basemap.fromId("satellite"),
                Basemap.fromId("hybrid"),
                Basemap.fromId("topo-vector"),
                Basemap.fromId("streets-vector"),
                Basemap.fromId("dark-gray-vector")
            ]
        });
        view.ui.add(
            new Expand({
                view: view,
                content: basemapGallery,
                expandIcon: "basemap",
                expandTooltip: "Chọn bản đồ nền"
            }),
            "top-right"
        );

        // Toạ độ con trỏ (DIY): hiện kinh độ + vĩ độ ở góc dưới-phải
        const coordEl = document.createElement("div");
        coordEl.className = "coords-readout";
        coordEl.textContent = "—";
        view.ui.add(coordEl, "bottom-right");
        view.on("pointer-move", (event) => {
            const pt = view.toMap({ x: event.x, y: event.y });
            if (pt) {
                coordEl.textContent =
                    `${pt.latitude.toFixed(5)}°N, ${pt.longitude.toFixed(5)}°E`;
            }
        });
    }

    // -----------------------------------------------------------------
    // CUSTOM NORTH ARROW (mũi tên la bàn lớn ở góc dưới phải)
    // -----------------------------------------------------------------
    function setupNorthArrow() {
        const arrow = document.getElementById("northArrow");
        if (!arrow) return;
        const inner = arrow.querySelector(".compass-inner");

        // Khi camera xoay (heading thay đổi) → xoay ngược kim la bàn
        view.watch("camera.heading", (heading) => {
            if (inner) inner.style.transform = `rotate(${-heading}deg)`;
        });

        // Click vào la bàn → quay heading về 0 (bắc lên trên)
        arrow.addEventListener("click", () => {
            const cam = view.camera.clone();
            cam.heading = 0;
            view.goTo(cam, { duration: 800, easing: "in-out-cubic" });
        });
    }

    // GraphicsLayer - Sử dụng relative-to-ground để đồng bộ các Mesh 3D sát khít với mặt đất thực tế
    const layer = new GraphicsLayer({
        title: "Vinhomes West Point",
        elevationInfo: { mode: "relative-to-ground" }
    });
    map.add(layer);

    // Khởi tạo các cấu hình mặc định bổ sung cho bảng điều khiển BIM
    window.WPConfig = window.WPConfig || {};
    window.WPConfig.currentHeightLimit = 160;
    window.WPConfig.terrainZOffset = -0.30;
    window.WPConfig.roadsZOffset = 0.00;

    // Các cờ cấu hình chính bật mặc định
    window.WPConfig.ENABLE_TERRAIN = true;
    window.WPConfig.ENABLE_ROADS = true;
    window.WPConfig.ENABLE_AMENITIES_GROUND = true;
    window.WPConfig.ENABLE_PODIUM = true;
    window.WPConfig.ENABLE_TOWERS = true;
    window.WPConfig.ENABLE_ROOFTOP_AMENITIES = true;
    window.WPConfig.ENABLE_VEGETATION = true;

    window.WPConfig.podiumStructure = true;
    window.WPConfig.podiumGlass = true;
    window.WPConfig.podiumRoof = true;

    window.WPConfig.towerW1 = true;
    window.WPConfig.towerW2 = true;
    window.WPConfig.towerW3 = true;
    window.WPConfig.towerBalcony = true;
    window.WPConfig.towerAluminum = true;

    // Helper: addBox có hỗ trợ lọc cắt lớp cao độ & dịch chuyển Z-Offset động
    function addBox(opts) {
        const {
            x, y, z, width, depth, height,
            color, metallic, roughness,
            heading,
            attributes, popup,
            type
        } = opts;

        let drawZ = z;
        let drawHeight = height;

        // Áp dụng dịch chuyển độ cao Z-Offset theo tính chất cấu trúc
        if (attributes && attributes.name) {
            if (attributes.name.includes("Nền") || attributes.name.includes("đất") || attributes.name.includes("cỏ") || attributes.name.includes("Cột") || attributes.name.includes("Lõi")) {
                drawZ += window.WPConfig.terrainZOffset;
            } else if (attributes.name.includes("Đường") || attributes.name.includes("Vỉa hè") || attributes.name.includes("Lối đi") || attributes.name.includes("Dốc")) {
                drawZ += window.WPConfig.roadsZOffset;
            }
        }

        // Lọc cắt lớp cao độ thi công BIM (Capping)
        const heightLimit = window.WPConfig.currentHeightLimit;
        if (drawZ > heightLimit) {
            return null; // Nằm hoàn toàn phía trên đường cắt giới hạn
        }
        if (drawZ + drawHeight > heightLimit) {
            drawHeight = heightLimit - drawZ; // Cắt ngắn khối hộp
            if (drawHeight < 0.05) return null; // Quá mỏng để hiển thị sạch sẽ
        }

        const { longitude, latitude } = localToLngLat(x, y);
        const center = new Point({
            longitude: longitude,
            latitude: latitude,
            z: drawZ, // Sử dụng drawZ làm gốc vì Mesh.createBox neo ở đáy của hộp (bottom center)
            spatialReference: { wkid: 4326 }
        });

        const mesh = Mesh.createBox(center, {
            size: { width: width, depth: depth, height: drawHeight },
            unit: "meters"
        });

        if (heading !== undefined && heading !== 0) {
            mesh.rotate(0, 0, -heading);
        }

        const symbol = {
            type: "mesh-3d",
            symbolLayers: [
                {
                    type: "fill",
                    material: {
                        color: color,
                        colorMixMode: "replace"
                    },
                    edges: null
                }
            ]
        };

        if (metallic !== undefined || roughness !== undefined) {
            symbol.symbolLayers[0].material.metallic = metallic ?? 0;
            symbol.symbolLayers[0].material.roughness = roughness ?? 0.8;
        }

        const g = new Graphic({
            geometry: mesh,
            symbol: symbol,
            attributes: attributes || {},
            popupTemplate: popup || null
        });

        layer.add(g);
        return g;
    }

    // Helper: addCylinder có hỗ trợ lọc cắt lớp cao độ thi công BIM (Capping)
    function addCylinder(opts) {
        const {
            x, y, z, radius, height,
            color, metallic, roughness,
            attributes, popup,
            type
        } = opts;

        let drawZ = z;
        let drawHeight = height;

        // Lọc cắt lớp cao độ thi công BIM (Capping)
        const heightLimit = window.WPConfig.currentHeightLimit;
        if (drawZ > heightLimit) {
            return null;
        }
        if (drawZ + drawHeight > heightLimit) {
            drawHeight = heightLimit - drawZ;
            if (drawHeight < 0.05) return null;
        }

        const { longitude, latitude } = localToLngLat(x, y);
        const center = new Point({
            longitude: longitude,
            latitude: latitude,
            z: drawZ, // Sử dụng drawZ làm gốc vì Mesh.createCylinder neo ở đáy của hình trụ (bottom center)
            spatialReference: { wkid: 4326 }
        });

        const mesh = Mesh.createCylinder(center, {
            size: { width: radius * 2, depth: radius * 2, height: drawHeight },
            unit: "meters"
        });

        const symbol = {
            type: "mesh-3d",
            symbolLayers: [
                {
                    type: "fill",
                    material: {
                        color: color,
                        colorMixMode: "replace"
                    },
                    edges: null
                }
            ]
        };

        if (metallic !== undefined || roughness !== undefined) {
            symbol.symbolLayers[0].material.metallic = metallic ?? 0;
            symbol.symbolLayers[0].material.roughness = roughness ?? 0.8;
        }

        const g = new Graphic({
            geometry: mesh,
            symbol: symbol,
            attributes: attributes || {},
            popupTemplate: popup || null
        });

        layer.add(g);
        return g;
    }

    // Hàm vẽ lại toàn bộ mô hình (Reactive Redraw)
    function redrawAll() {
        layer.removeAll();
        try {
            if (window.WPModules && window.WPConfig) {
                // Địa hình & Đường xá
                if (window.WPConfig.ENABLE_ROADS && window.WPModules.buildPlaza) {
                    window.WPModules.buildPlaza(data.plaza);
                    if (window.WPModules.buildStreetlights) {
                        window.WPModules.buildStreetlights(data.plaza);
                    }
                }
                // Khối đế thương mại
                if (window.WPConfig.ENABLE_PODIUM && window.WPModules.buildPodium) {
                    window.WPModules.buildPodium(data.podium);
                }
                // Tiện ích tầng mái khối đế
                if (window.WPConfig.ENABLE_ROOFTOP_AMENITIES && window.WPModules.buildRooftopAmenities) {
                    const roofZ = window.WPConfig.ENABLE_PODIUM ? 25.20 : 2.90;
                    window.WPModules.buildRooftopAmenities(data, roofZ);
                }
                // 3 Tòa tháp căn hộ
                if (window.WPConfig.ENABLE_TOWERS && window.WPModules.buildTower) {
                    data.towers.forEach(t => {
                        let baseZ = 3.50; // Tòa West 3 độc lập ngồi trên footprint tầng 1 (Z: 3.50)
                        if (t.id === "W1" || t.id === "W2") {
                            baseZ = 2.70 + (data.podium.modeledHeight || 22.5); // Tòa West 1 & 2 ngồi trên mái khối đế chính (Z: 25.20)
                        }
                        window.WPModules.buildTower(t, baseZ);
                    });
                }
                console.log("[VWP] Redraw success. Graphics count:", layer.graphics.length);
            }
        } catch (err) {
            console.error("[VWP] Redraw error:", err);
        }
    }

    window.redrawAll = redrawAll;
    window.localToLngLat = localToLngLat;

    window.WPEngine = {
        addBox: addBox,
        addCylinder: addCylinder,
        layer: layer,
        Mesh: Mesh,
        Point: Point,
        Graphic: Graphic,
        WebStyleSymbol: WebStyleSymbol
    };

    // Vẽ mô hình lần đầu tiên
    redrawAll();

    // -----------------------------------------------------------------
    // KHỞI TẠO BỘ LẮNG NGHE SỰ KIỆN GIAO DIỆN (UI Controller Events)
    // -----------------------------------------------------------------
    setupUIControls();

    function setupUIControls() {
        const heightSlider = document.getElementById("heightSlider");
        const heightValue = document.getElementById("heightValue");

        // 1. Slider điều chỉnh Mặt cắt Chiều cao
        if (heightSlider && heightValue) {
            heightSlider.addEventListener("input", function() {
                const val = parseFloat(this.value);
                heightValue.textContent = val + "m";
                window.WPConfig.currentHeightLimit = val;
                redrawAll();
            });
        }

        // 2. Nút Quét thi công tự động (Auto construction animation sweep)
        let playInterval = null;
        const btnAutoPlay = document.getElementById("btnAutoPlay");
        if (btnAutoPlay) {
            btnAutoPlay.addEventListener("click", function() {
                if (playInterval) {
                    // Đang chạy -> Tạm dừng
                    clearInterval(playInterval);
                    playInterval = null;
                    this.textContent = "► Quét thi công tự động";
                    this.classList.remove("playing");
                } else {
                    // Đang dừng -> Chạy
                    this.textContent = "⏸ Tạm dừng quét";
                    this.classList.add("playing");
                    if (parseFloat(heightSlider.value) >= 160) {
                        heightSlider.value = 0;
                    }
                    playInterval = setInterval(() => {
                        let val = parseFloat(heightSlider.value);
                        val += 2;
                        if (val > 160) {
                            val = 160;
                            clearInterval(playInterval);
                            playInterval = null;
                            btnAutoPlay.textContent = "► Quét thi công tự động";
                            btnAutoPlay.classList.remove("playing");
                        }
                        heightSlider.value = val;
                        heightValue.textContent = val + "m";
                        window.WPConfig.currentHeightLimit = val;
                        redrawAll();
                    }, 45);
                }
            });
        }

        // 3. Nút Reset về trạng thái mặc định đầy đủ
        const btnReset = document.getElementById("btnReset");
        if (btnReset) {
            btnReset.addEventListener("click", function() {
                if (playInterval) {
                    clearInterval(playInterval);
                    playInterval = null;
                    btnAutoPlay.textContent = "► Quét thi công tự động";
                    btnAutoPlay.classList.remove("playing");
                }
                heightSlider.value = 160;
                heightValue.textContent = "160m";
                window.WPConfig.currentHeightLimit = 160;
                
                // Trả toàn bộ checkbox về trạng thái checked
                document.querySelectorAll("#controlPanel input[type=checkbox]").forEach(chk => {
                    if (chk.id !== "chkNightMode") {
                        chk.checked = true;
                        window.WPConfig[getCheckboxConfigKey(chk.id)] = true;
                    } else {
                        chk.checked = false;
                    }
                });

                // Reset Z-Offset
                window.WPConfig.terrainZOffset = -0.30;
                window.WPConfig.roadsZOffset = 0.00;
                document.getElementById("valTerrainOffset").textContent = "-0.30m";
                document.getElementById("valRoadOffset").textContent = "0.00m";

                // Đưa hướng camera ban ngày về mặc định
                view.environment.lighting.date = new Date("2026-05-06T09:30:00+07:00");
                view.environment.lighting.directShadowsEnabled = true;

                redrawAll();
            });
        }

        // Helper maps checkbox id with WPConfig key
        function getCheckboxConfigKey(id) {
            const keys = {
                chkTerrain: "ENABLE_TERRAIN",
                chkRoads: "ENABLE_ROADS",
                chkAmenitiesGround: "ENABLE_AMENITIES_GROUND",
                chkPodium: "ENABLE_PODIUM",
                chkPodiumStructure: "podiumStructure",
                chkPodiumGlass: "podiumGlass",
                chkPodiumRoof: "podiumRoof",
                chkTowers: "ENABLE_TOWERS",
                chkTowerW1: "towerW1",
                chkTowerW2: "towerW2",
                chkTowerW3: "towerW3",
                chkTowerBalcony: "towerBalcony",
                chkTowerAluminum: "towerAluminum",
                chkRooftopAmenities: "ENABLE_ROOFTOP_AMENITIES",
                chkVegetation: "ENABLE_VEGETATION"
            };
            return keys[id] || null;
        }

        // Gắn sự kiện Change cho tất cả Checkboxes và đồng bộ ban đầu
        const checkboxes = document.querySelectorAll("#controlPanel input[type=checkbox]");
        checkboxes.forEach(chk => {
            const key = getCheckboxConfigKey(chk.id);
            if (key) {
                // Đồng bộ trạng thái UI với WPConfig ban đầu
                if (window.WPConfig[key] !== undefined) {
                    chk.checked = window.WPConfig[key];
                }
                chk.addEventListener("change", function() {
                    window.WPConfig[key] = this.checked;
                    redrawAll();
                });
            }
        });

        // 4. Chế độ Ban đêm (Night mode - Lighting)
        const chkNightMode = document.getElementById("chkNightMode");
        if (chkNightMode) {
            // Lưu trữ màu sắc ban ngày ban đầu
            const originalColors = {
                residential_glass: [...window.WPMaterials.residential_glass.color],
                commercial_glass: [...window.WPMaterials.commercial_glass.color],
                pool_water: [...window.WPMaterials.pool_water.color],
                wood_deck: [...window.WPMaterials.wood_deck.color],
                white_aluminum: [...window.WPMaterials.white_aluminum.color],
                light_concrete: [...window.WPMaterials.light_concrete.color]
            };

            chkNightMode.addEventListener("change", function() {
                // Giữ nguyên Bản đồ nền OSM sáng rõ để xung quanh không bị tối đen
                view.map.basemap = "osm";

                if (this.checked) {
                    // Chuyển thời gian về đêm hẳn (7:30 Tối - 19:30) để trời tối đen có sao, đèn đường BẬT SÁNG RỰC RỠ
                    view.environment.lighting.date = new Date("2026-05-06T19:30:00+07:00"); 
                    
                    // Kích hoạt cameraTrackingEnabled = true để biến camera thành một đèn pha (headlight) khổng lồ chiếu trực diện
                    // vào tòa nhà, giúp tòa nhà luôn sáng bừng rực rỡ từ mọi góc nhìn, triệt tiêu mọi bóng tối u ám!
                    view.environment.lighting.cameraTrackingEnabled = true;
                    view.environment.lighting.directShadowsEnabled = false; // Loại bỏ bóng đổ tối đen rườm rà

                    // Màu kính và tiện ích phát sáng cực đại
                    window.WPMaterials.residential_glass.color = [255, 235, 10, 1.0];   // Kính vàng ròng sáng rực
                    window.WPMaterials.commercial_glass.color = [255, 250, 150, 1.0];  // Khối đế bật đèn thương mại sáng trưng
                    window.WPMaterials.pool_water.color = [0, 255, 255, 1.0];           // Bể bơi xanh neon cực rực rỡ
                    window.WPMaterials.wood_deck.color = [100, 70, 50, 1];              // Sàn gỗ sáng ấm áp
                    
                    window.WPMaterials.white_aluminum.color = [255, 255, 255, 1.0];     // Viền nhôm sáng trắng phản chiếu đèn cực mạnh
                    window.WPMaterials.light_concrete.color = [255, 255, 255, 1.0];     // Bê tông cột sáng sủa rõ ràng
                } else {
                    view.environment.lighting.date = new Date("2026-05-06T09:30:00+07:00"); // 9:30 Sáng (Ánh nắng ban mai)
                    view.environment.lighting.cameraTrackingEnabled = false; // Trả về góc chiếu mặt trời thực tế ban ngày
                    view.environment.lighting.directShadowsEnabled = true;

                    // Khôi phục màu sắc ban ngày gốc
                    window.WPMaterials.residential_glass.color = [...originalColors.residential_glass];
                    window.WPMaterials.commercial_glass.color = [...originalColors.commercial_glass];
                    window.WPMaterials.pool_water.color = [...originalColors.pool_water];
                    window.WPMaterials.wood_deck.color = [...originalColors.wood_deck];
                    window.WPMaterials.white_aluminum.color = [...originalColors.white_aluminum];
                    window.WPMaterials.light_concrete.color = [...originalColors.light_concrete];
                }
                // Vẽ lại toàn bộ mô hình để cập nhật hiệu ứng ánh sáng đèn đêm rực rỡ
                redrawAll();
            });
        }

        // 5. Stepper Z-Offset chỉnh khoảng cách cao độ
        setupStepper("btnTerrainMinus", "btnTerrainPlus", "valTerrainOffset", "terrainZOffset", -0.05);
        setupStepper("btnRoadMinus", "btnRoadPlus", "valRoadOffset", "roadsZOffset", 0.05);

        function setupStepper(minusId, plusId, valueId, configKey, step) {
            const btnMinus = document.getElementById(minusId);
            const btnPlus = document.getElementById(plusId);
            const valEl = document.getElementById(valueId);

            if (btnMinus && btnPlus && valEl) {
                btnMinus.addEventListener("click", function() {
                    window.WPConfig[configKey] = parseFloat((window.WPConfig[configKey] - Math.abs(step)).toFixed(2));
                    valEl.textContent = (window.WPConfig[configKey] >= 0 ? "+" : "") + window.WPConfig[configKey].toFixed(2) + "m";
                    redrawAll();
                });
                btnPlus.addEventListener("click", function() {
                    window.WPConfig[configKey] = parseFloat((window.WPConfig[configKey] + Math.abs(step)).toFixed(2));
                    valEl.textContent = (window.WPConfig[configKey] >= 0 ? "+" : "") + window.WPConfig[configKey].toFixed(2) + "m";
                    redrawAll();
                });
            }
        }

        // 6. BIM Construction Phase Buttons (Phân giai đoạn thi công nhanh)
        const phaseBtns = document.querySelectorAll(".phase-btn");
        phaseBtns.forEach(btn => {
            btn.addEventListener("click", function() {
                phaseBtns.forEach(b => b.classList.remove("highlighted"));
                this.classList.add("highlighted");

                const phase = parseInt(this.getAttribute("data-phase"));
                
                if (phase === 1) {
                    // 1. Nền & Đường (Chỉ hiện địa hình, đường xá)
                    setBIMStates({
                        currentHeightLimit: 5,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: false,
                        ENABLE_PODIUM: false, ENABLE_TOWERS: false, ENABLE_ROOFTOP_AMENITIES: false, ENABLE_VEGETATION: false
                    });
                } else if (phase === 2) {
                    // 2. Khối đế (Cắt ở tầng 4 khối đế Z: 23m)
                    setBIMStates({
                        currentHeightLimit: 23,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: true, podiumRoof: true,
                        ENABLE_TOWERS: false, ENABLE_ROOFTOP_AMENITIES: false, ENABLE_VEGETATION: false
                    });
                } else if (phase === 3) {
                    // 3. Khung tháp (Vẽ tháp không nhôm/không logia kính để tạo bộ xương bê tông cốt thép cốt lõi)
                    setBIMStates({
                        currentHeightLimit: 110,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: false, podiumRoof: true,
                        ENABLE_TOWERS: true, towerW1: true, towerW2: true, towerW3: true,
                        towerBalcony: false, towerAluminum: false, ENABLE_ROOFTOP_AMENITIES: false, ENABLE_VEGETATION: false
                    });
                } else if (phase === 4) {
                    // 4. Vỏ kính (Vách kính & thanh nhôm nổi bật)
                    setBIMStates({
                        currentHeightLimit: 110,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: true, podiumRoof: true,
                        ENABLE_TOWERS: true, towerW1: true, towerW2: true, towerW3: true,
                        towerBalcony: true, towerAluminum: true, ENABLE_ROOFTOP_AMENITIES: false, ENABLE_VEGETATION: false
                    });
                } else if (phase === 5) {
                    // 5. Đỉnh tháp (Mái, thang máy, ăng-ten Z: 140m)
                    setBIMStates({
                        currentHeightLimit: 140,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: true, podiumRoof: true,
                        ENABLE_TOWERS: true, towerW1: true, towerW2: true, towerW3: true,
                        towerBalcony: true, towerAluminum: true, ENABLE_ROOFTOP_AMENITIES: false, ENABLE_VEGETATION: false
                    });
                } else if (phase === 6) {
                    // 6. Tiện ích mái (Rooftop pool & gardens)
                    setBIMStates({
                        currentHeightLimit: 160,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: true, podiumRoof: true,
                        ENABLE_TOWERS: true, towerW1: true, towerW2: true, towerW3: true,
                        towerBalcony: true, towerAluminum: true, ENABLE_ROOFTOP_AMENITIES: true, ENABLE_VEGETATION: false
                    });
                } else if (phase === 7) {
                    // 7. Cảnh quan (Hàng cây tầng trệt)
                    setBIMStates({
                        currentHeightLimit: 160,
                        ENABLE_TERRAIN: true, ENABLE_ROADS: true, ENABLE_AMENITIES_GROUND: true,
                        ENABLE_PODIUM: true, podiumStructure: true, podiumGlass: true, podiumRoof: true,
                        ENABLE_TOWERS: true, towerW1: true, towerW2: true, towerW3: true,
                        towerBalcony: true, towerAluminum: true, ENABLE_ROOFTOP_AMENITIES: true, ENABLE_VEGETATION: true
                    });
                } else if (phase === 8) {
                    // 8. Ánh sáng (Night mode!)
                    chkNightMode.checked = !chkNightMode.checked;
                    chkNightMode.dispatchEvent(new Event("change"));
                }
            });
        });

        function setBIMStates(states) {
            for (let key in states) {
                window.WPConfig[key] = states[key];
                
                // Đồng bộ check box trực quan
                const chkId = getCheckboxIdFromKey(key);
                if (chkId) {
                    const chk = document.getElementById(chkId);
                    if (chk) chk.checked = states[key];
                }
            }

            // Cập nhật thanh trượt slider chiều cao
            if (states.currentHeightLimit !== undefined) {
                heightSlider.value = states.currentHeightLimit;
                heightValue.textContent = states.currentHeightLimit + "m";
            }

            redrawAll();
        }

        function getCheckboxIdFromKey(key) {
            const mapKeys = {
                ENABLE_TERRAIN: "chkTerrain",
                ENABLE_ROADS: "chkRoads",
                ENABLE_AMENITIES_GROUND: "chkAmenitiesGround",
                ENABLE_PODIUM: "chkPodium",
                podiumStructure: "chkPodiumStructure",
                podiumGlass: "chkPodiumGlass",
                podiumRoof: "chkPodiumRoof",
                ENABLE_TOWERS: "chkTowers",
                towerW1: "chkTowerW1",
                towerW2: "chkTowerW2",
                towerW3: "chkTowerW3",
                towerBalcony: "chkTowerBalcony",
                towerAluminum: "chkTowerAluminum",
                ENABLE_ROOFTOP_AMENITIES: "chkRooftopAmenities",
                ENABLE_VEGETATION: "chkVegetation"
            };
            return mapKeys[key] || null;
        }

        // 7. ArcGIS Click Event & BIM Selection Readout (Xem dữ liệu kỹ thuật công trình khi click vào mô hình 3D)
        view.on("click", (event) => {
            view.hitTest(event).then((response) => {
                const results = response.results;
                const selectionBox = document.getElementById("selectionReadout");
                
                if (results.length > 0 && results[0].graphic) {
                    const graphic = results[0].graphic;
                    const attr = graphic.attributes;
                    
                    if (attr && attr.name) {
                        selectionBox.classList.add("active");
                        
                        let html = `<div class="readout-field"><span class="readout-label">Tên phần tử:</span><span class="readout-value" style="color:#00f2fe">${attr.name}</span></div>`;
                        
                        if (attr.id) {
                            html += `<div class="readout-field"><span class="readout-label">Mã phân khu:</span><span class="readout-value">${attr.id}</span></div>`;
                        }
                        if (attr.floors) {
                            html += `<div class="readout-field"><span class="readout-label">Số tầng nổi:</span><span class="readout-value">${attr.floors} tầng</span></div>`;
                        }
                        if (attr.modeledHeight) {
                            html += `<div class="readout-field"><span class="readout-label">Chiều cao:</span><span class="readout-value">${attr.modeledHeight}m</span></div>`;
                        }
                        if (attr.floorInfo) {
                            html += `<div class="readout-field"><span class="readout-label">Phân bổ:</span><span class="readout-value" style="font-size:10px">${attr.floorInfo}</span></div>`;
                        }
                        if (attr.units) {
                            html += `<div class="readout-field"><span class="readout-label">Tổng số căn:</span><span class="readout-value">${attr.units} căn hộ</span></div>`;
                        }
                        
                        selectionBox.innerHTML = html;
                    } else {
                        selectionBox.classList.remove("active");
                        selectionBox.textContent = "Click vào mô hình 3D để quét và hiển thị dữ liệu kỹ thuật công trình";
                    }
                } else {
                    selectionBox.classList.remove("active");
                    selectionBox.textContent = "Click vào mô hình 3D để quét và hiển thị dữ liệu kỹ thuật công trình";
                }
            });
        });
    }
});

