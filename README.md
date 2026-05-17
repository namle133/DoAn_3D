# DoAn_3D — Vinhomes West Point (ArcGIS Maps SDK for JavaScript)

Mô hình 3D minh hoạ khu đô thị **Vinhomes West Point** (Phạm Hùng, Mễ Trì,
Nam Từ Liêm, Hà Nội) gồm 3 toà tháp **W1 / W2 / W3** và khối đế (podium),
dựng bằng `SceneView` của **ArcGIS Maps SDK for JavaScript 4.30**.

## Cấu trúc thư mục

```
DoAn_3D/
├── index.html          # Trang chính, nhúng ArcGIS SDK + UI overlay
├── css/style.css       # Style cho title bar, north arrow
├── js/
│   ├── buildings.js    # Toạ độ + kích thước các toà nhà (xấp xỉ)
│   └── app.js          # Map / SceneView / Mesh + widgets (Compass, Search, ...)
└── README.md
```

## Tính năng

- **Bản đồ nền**: OpenStreetMap mặc định, có thể đổi sang Satellite, Hybrid,
  Topo, Streets, Dark Gray qua nút chọn basemap (góc trên-phải).
- **Mô hình 3D**: 3 toà tháp W1/W2/W3 + khối đế podium, dựng bằng `Mesh.createBox`
  (thân kính + dải sàn + cột góc + mullion + khối kỹ thuật + antenna).
- **Định hướng (la bàn)**:
  - Compass widget của Esri ở góc trên-trái (nhấn để quay về hướng Bắc).
  - **North arrow lớn** ở góc dưới-phải hiển thị 4 hướng N/E/S/W, mũi tên đỏ
    luôn chỉ về phía Bắc bất kể camera xoay.
- **Công cụ**: Home, Fullscreen, ScaleBar (thước tỉ lệ), Coordinates readout
  (toạ độ con trỏ), Zoom, Navigation Toggle.

## Cách chạy

### Bằng VS Code

Cài extension **Live Server** → chuột phải `index.html` → *Open with Live
Server*.

<img width="800" height="450" alt="vinhomes-west-point" src="https://github.com/user-attachments/assets/e795ace0-209c-4398-a8d8-5b3290aabc66" />

