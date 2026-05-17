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
    "esri/widgets/Compass",
    "esri/widgets/ScaleBar",
    "esri/widgets/Expand",
    "esri/widgets/BasemapGallery",
    "esri/widgets/Home",
    "esri/widgets/Fullscreen"
], function (
    Map, Basemap, SceneView, Graphic, GraphicsLayer, Point, Mesh,
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

    // SceneView
    const homePos = localToLngLat(450, -450);
    const HOME_VIEWPOINT = {
        position: {
            longitude: homePos.longitude,
            latitude: homePos.latitude,
            z: 750
        },
        heading: 315,
        tilt: 55
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

    // GraphicsLayer
    const layer = new GraphicsLayer({
        title: "Vinhomes West Point",
        elevationInfo: { mode: "absolute-height" }
    });
    map.add(layer);

    // Helper
    function addBox(opts) {
        const {
            x, y, z, width, depth, height,
            color, metallic, roughness,
            attributes, popup
        } = opts;

        const { longitude, latitude } = localToLngLat(x, y);
        const center = new Point({
            longitude: longitude,
            latitude: latitude,
            z: z + height / 2,
            spatialReference: { wkid: 4326 }
        });

        const mesh = Mesh.createBox(center, {
            size: { width: width, depth: depth, height: height },
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

        // PBR (metallic/roughness)
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

    // Plaza
    function buildPlaza(p) {
        addBox({
            x: 0, y: 0, z: -p.thickness,
            width: p.width, depth: p.depth, height: p.thickness,
            color: p.color, roughness: 0.95,
            attributes: { name: "Sân plaza" }
        });

        p.walkways.forEach((w, i) => {
            addBox({
                x: w.x, y: w.y, z: 0.01,
                width: w.w, depth: w.d, height: 0.05,
                color: w.color, roughness: 0.9,
                attributes: { name: "Lối đi " + (i + 1) }
            });
        });
    }

    // Podium
    function buildPodium(p) {
        const totalH = p.floors * p.floorHeight;
        const podiumAttributes = {
        id: p.id,
        name: p.name,
        floors: p.floors,
        modeledHeight: p.modeledHeight || totalH,
        usage: p.usage,
        floorInfo: p.floorInfo,
        mainFacilities: p.mainFacilities,
        rooftopFacilities: p.rooftopFacilities,
        parking: p.parking,
        serviceRole: p.serviceRole,
        note: p.note
    };

    const podiumPopup = {
        title: "{name}",
        content: [{
            type: "fields",
            fieldInfos: [
                { fieldName: "id", label: "Mã khu" },
                { fieldName: "floors", label: "Số tầng khối đế" },
                { fieldName: "modeledHeight", label: "Chiều cao mô phỏng (m)" },
                { fieldName: "usage", label: "Chức năng" },
                { fieldName: "floorInfo", label: "Phân bổ tầng" },
                { fieldName: "mainFacilities", label: "Tiện ích chính" },
                { fieldName: "rooftopFacilities", label: "Tiện ích trên mái khối đế" },
                { fieldName: "parking", label: "Bãi đỗ xe" },
                { fieldName: "serviceRole", label: "Vai trò phục vụ" },
                { fieldName: "note", label: "Ghi chú dữ liệu" }
            ]
        }]
    };

        addBox({
            x: p.center.x, y: p.center.y, z: 0,
            width: p.width, depth: p.depth, height: 1.2,
            color: p.baseColor, roughness: 0.85,
            attributes: podiumAttributes,
            popup: podiumPopup
        });

        addBox({
            x: p.center.x, y: p.center.y, z: 1.2,
            width: p.width - 2, depth: p.depth - 2, height: totalH - 1.2,
            color: p.glassColor, metallic: 0.4, roughness: 0.18,
            attributes: podiumAttributes,
            popup: podiumPopup
        });

        for (let f = 1; f <= p.floors; f++) {
            const z = f * p.floorHeight - 0.3;
            addBox({
                x: p.center.x, y: p.center.y, z: z,
                width: p.width + 0.6, depth: p.depth + 0.6, height: 0.6,
                color: p.slabColor, roughness: 0.8
            });
        }
    }

    // Tower
    function buildTower(t) {
        const cx = t.center.x;
        const cy = t.center.y;
        const totalH = t.modeledHeight || (t.floors * t.floorHeight);
        const baseZ = 18;
        const towerAttributes = {
            id: t.id,
            name: t.name,
            floors: t.floors,
            modeledHeight: t.modeledHeight,
            basementFloors: t.basementFloors,
            floorInfo: t.floorInfo,
            units: t.units,
            density: t.density,
            elevators: t.elevators,
            apartmentTypes: t.apartmentTypes,
            areaRange: t.areaRange,
            handoverStandard: t.handoverStandard,
            deliveryTime: t.deliveryTime,
            mainFace: t.mainFace
        };

        const towerPopup = {
            title: "{name}",
            content: [{
                type: "fields",
                fieldInfos: [
                    { fieldName: "id", label: "Mã tòa" },
                    { fieldName: "floors", label: "Số tầng nổi" },
                    { fieldName: "modeledHeight", label: "Chiều cao mô phỏng (m)" },
                    { fieldName: "basementFloors", label: "Số tầng hầm" },
                    { fieldName: "floorInfo", label: "Phân bổ tầng" },
                    { fieldName: "units", label: "Tổng số căn" },
                    { fieldName: "density", label: "Mật độ căn hộ" },
                    { fieldName: "elevators", label: "Thang máy" },
                    { fieldName: "apartmentTypes", label: "Loại hình căn hộ" },
                    { fieldName: "areaRange", label: "Diện tích căn hộ" },
                    { fieldName: "handoverStandard", label: "Tiêu chuẩn bàn giao" },
                    { fieldName: "deliveryTime", label: "Thời điểm bàn giao" },
                    { fieldName: "mainFace", label: "Mặt chính" }
                ]
            }]
        };
        addBox({
            x: cx, y: cy, z: baseZ,
            width: t.width - 1.4, depth: t.depth - 1.4, height: totalH,
            color: t.glassColor, metallic: 0.55, roughness: 0.15,
            attributes: towerAttributes,
            popup: towerPopup
        });

        // Slab
        for (let f = 1; f <= t.floors; f++) {
            const z = baseZ + f * t.floorHeight - 0.18;
            addBox({
                x: cx, y: cy, z: z,
                width: t.width, depth: t.depth, height: 0.36,
                color: t.slabColor, roughness: 0.85,
                attributes: towerAttributes,
                popup: towerPopup
            });
        }

        // Pillar
        const ps = 1.6;
        [[-1, -1], [1, -1], [-1, 1], [1, 1]].forEach(([sx, sy]) => {
            addBox({
                x: cx + (sx * (t.width - ps)) / 2,
                y: cy + (sy * (t.depth - ps)) / 2,
                z: baseZ,
                width: ps, depth: ps, height: totalH,
                color: t.pillarColor, roughness: 0.75,
                attributes: towerAttributes,
                popup: towerPopup
            });
        });

        // Mullion
        // Mullion theo mặt chính của từng tòa
        const mc = 8;
        const mainFace = t.mainFace || "bottom";

        if (mainFace === "top" || mainFace === "bottom") {
            const sideY = mainFace === "top" ? 1 : -1;

            for (let i = 1; i < mc; i++) {
                const fx = -t.width / 2 + (i * t.width) / mc;

                addBox({
                    x: cx + fx,
                    y: cy + sideY * t.depth / 2,
                    z: baseZ,
                    width: 0.7,
                    depth: 0.45,
                    height: totalH,
                    color: t.pillarColor,
                    roughness: 0.75,
                    attributes: {
                        ...towerAttributes,
                        partName: `Mullion mặt chính ${t.id}`,
                        partType: "Mullion"
                    },
                    popup: towerPopup
                });
            }
        }

        if (mainFace === "left" || mainFace === "right") {
            const sideX = mainFace === "right" ? 1 : -1;

            for (let i = 1; i < mc; i++) {
                const fy = -t.depth / 2 + (i * t.depth) / mc;

                addBox({
                    x: cx + sideX * t.width / 2,
                    y: cy + fy,
                    z: baseZ,
                    width: 0.45,
                    depth: 0.7,
                    height: totalH,
                    color: t.pillarColor,
                    roughness: 0.75,
                    attributes: {
                        ...towerAttributes,
                        partName: `Mullion mặt chính ${t.id}`,
                        partType: "Mullion"
                    },
                    popup: towerPopup
                });
            }
        }

        // Crown
        addBox({
            x: cx, y: cy, z: baseZ + totalH,
            width: t.width - 6, depth: t.depth - 6, height: t.mechanicalHeight,
            color: t.crownColor, roughness: 0.7
        });
        addBox({
            x: cx, y: cy, z: baseZ + totalH + t.mechanicalHeight - 0.5,
            width: t.width - 4, depth: t.depth - 4, height: 0.6,
            color: [120, 125, 132, 1], roughness: 0.6
        });

        // Antenna
        addBox({
            x: cx, y: cy, z: baseZ + totalH + t.mechanicalHeight,
            width: 1.0, depth: 1.0, height: t.antennaHeight,
            color: [80, 85, 92, 1], metallic: 0.7, roughness: 0.3
        });
        addBox({
            x: cx, y: cy,
            z: baseZ + totalH + t.mechanicalHeight + t.antennaHeight,
            width: 0.4, depth: 0.4, height: 8,
            color: [50, 55, 60, 1], metallic: 0.8, roughness: 0.2
        });
    }

    // Build all
    try {
        buildPlaza(data.plaza);
        buildPodium(data.podium);
        data.towers.forEach(buildTower);
        console.log("[VWP] Built", layer.graphics.length, "graphics");
    } catch (err) {
        console.error("[VWP] Build error:", err);
    }

});
