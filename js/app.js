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
    "esri/geometry/Mesh"
], function (Map, Basemap, SceneView, Graphic, GraphicsLayer, Point, Mesh) {
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

    // Empty Basemap
    const emptyBasemap = new Basemap({
        baseLayers: [],
        title: "Empty",
        id: "empty"
    });

    const map = new Map({
        basemap: emptyBasemap,
        ground: {
            surfaceColor: [180, 188, 196, 1]
        }
    });

    // SceneView
    const homePos = localToLngLat(220, -280);
    const view = new SceneView({
        container: "viewDiv",
        map: map,
        qualityProfile: "high",
        camera: {
            position: {
                longitude: homePos.longitude,
                latitude: homePos.latitude,
                z: 240
            },
            heading: 322,
            tilt: 62
        },
        environment: {
            background: { type: "color", color: [222, 232, 242, 1] },
            starsEnabled: false,
            atmosphereEnabled: false,
            lighting: {
                date: new Date("2026-05-06T09:30:00+07:00"),
                directShadowsEnabled: true,
                cameraTrackingEnabled: false
            }
        },
        ui: { components: ["zoom", "compass", "navigation-toggle"] }
    });

    view.when(
        () => console.log("[VWP] LocalScene ready"),
        (err) => console.error("[VWP] SceneView error:", err)
    );

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

        addBox({
            x: p.center.x, y: p.center.y, z: 0,
            width: p.width, depth: p.depth, height: 1.2,
            color: p.baseColor, roughness: 0.85
        });

        addBox({
            x: p.center.x, y: p.center.y, z: 1.2,
            width: p.width - 2, depth: p.depth - 2, height: totalH - 1.2,
            color: p.glassColor, metallic: 0.4, roughness: 0.18,
            attributes: { id: p.id, name: p.name, floors: p.floors },
            popup: { title: "{name}", content: "Số tầng: <b>{floors}</b>" }
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
        const totalH = t.floors * t.floorHeight;
        const baseZ = 18;

        addBox({
            x: cx, y: cy, z: baseZ,
            width: t.width - 1.4, depth: t.depth - 1.4, height: totalH,
            color: t.glassColor, metallic: 0.55, roughness: 0.15,
            attributes: {
                id: t.id, name: t.name,
                floors: t.floors, height: totalH
            },
            popup: {
                title: "{name}",
                content: [{
                    type: "fields",
                    fieldInfos: [
                        { fieldName: "id", label: "Mã toà" },
                        { fieldName: "floors", label: "Số tầng" },
                        { fieldName: "height", label: "Chiều cao thân (m)" }
                    ]
                }]
            }
        });

        // Slab
        for (let f = 1; f <= t.floors; f++) {
            const z = baseZ + f * t.floorHeight - 0.18;
            addBox({
                x: cx, y: cy, z: z,
                width: t.width, depth: t.depth, height: 0.36,
                color: t.slabColor, roughness: 0.85
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
                color: t.pillarColor, roughness: 0.75
            });
        });

        // Mullion
        const mc = 5;
        for (let i = 1; i < mc; i++) {
            const fx = -t.width / 2 + (i * t.width) / mc;
            [-1, 1].forEach((sy) => {
                addBox({
                    x: cx + fx,
                    y: cy + (sy * t.depth) / 2,
                    z: baseZ,
                    width: 0.6, depth: 0.4, height: totalH,
                    color: t.pillarColor, roughness: 0.75
                });
            });
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
