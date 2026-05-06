# DoAn_3D — Vinhomes West Point (ArcGIS Maps SDK for JavaScript)

Mô hình 3D minh hoạ khu đô thị **Vinhomes West Point** (Phạm Hùng, Mễ Trì,
Nam Từ Liêm, Hà Nội) gồm 3 toà tháp **W1 / W2 / W3** và khối đế (podium),
dựng bằng `SceneView` của **ArcGIS Maps SDK for JavaScript 4.30**.

## Cấu trúc thư mục

```
DoAn_3D/
├── index.html          # Trang chính, nhúng ArcGIS SDK
├── css/style.css       # Style cho panel & UI
├── js/
│   |
│   ├── buildings.js    # Toạ độ + kích thước các toà nhà (xấp xỉ)
│   └── app.js          # Tạo Map / SceneView / GraphicsLayer extrude
└── README.md
```

## Cách chạy

### Bằng VS Code

Cài extension **Live Server** → chuột phải `index.html` → *Open with Live
Server*.
