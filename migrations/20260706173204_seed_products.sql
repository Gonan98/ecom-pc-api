-- +goose Up
INSERT INTO products (
    id,
    category_id,
    brand_id,
    name,
    description,
    image_url,
    price,
    stock
) VALUES

-- Processors
(1, 1, 1, 'AMD Ryzen 5 9600X', '6-core 12-thread desktop gaming processor.', NULL, 279.90, 25),
(2, 1, 1, 'AMD Ryzen 7 9800X3D', 'High-end gaming processor with 3D V-Cache.', NULL, 519.90, 15),
(3, 1, 2, 'Intel Core i5-14600KF', '14th Gen unlocked desktop processor.', NULL, 309.90, 20),
(4, 1, 2, 'Intel Core i7-14700K', '20-core unlocked gaming processor.', NULL, 449.90, 12),

-- Graphics Cards
(5, 2, 3, 'NVIDIA GeForce RTX 5070', 'Next-generation gaming graphics card.', NULL, 699.90, 10),
(6, 2, 3, 'NVIDIA GeForce RTX 5080', 'High-performance graphics card for 4K gaming.', NULL, 1099.90, 6),
(7, 2, 4, 'ASUS TUF Gaming Radeon RX 7800 XT', 'Factory-overclocked AMD graphics card.', NULL, 579.90, 8),
(8, 2, 5, 'MSI GeForce RTX 4060 Ti Ventus', 'Efficient graphics card for 1080p gaming.', NULL, 449.90, 18),

-- Motherboards
(9, 3, 4, 'ASUS ROG Strix B650-A Gaming WiFi', 'AM5 ATX gaming motherboard.', NULL, 259.90, 14),
(10, 3, 5, 'MSI MAG B760 Tomahawk WiFi', 'Intel LGA1700 gaming motherboard.', NULL, 219.90, 16),
(11, 3, 6, 'Gigabyte B650 Aorus Elite AX', 'AMD AM5 motherboard with WiFi 6.', NULL, 239.90, 12),
(12, 3, 7, 'ASRock Z790 Steel Legend', 'Intel Z790 motherboard for enthusiasts.', NULL, 299.90, 8),

-- RAM
(13, 4, 10, 'G.Skill Trident Z5 RGB 32GB DDR5-6000', '32GB DDR5 dual-channel memory kit.', NULL, 169.90, 20),
(14, 4, 9, 'Kingston Fury Beast 32GB DDR5-5600', 'High-speed DDR5 gaming memory.', NULL, 149.90, 22),
(15, 4, 8, 'Corsair Vengeance RGB 32GB DDR5-6400', 'Premium RGB DDR5 memory kit.', NULL, 189.90, 18),

-- SSD
(16, 5, 11, 'Samsung 990 PRO 2TB NVMe SSD', 'PCIe 4.0 high-performance SSD.', NULL, 209.90, 18),
(17, 5, 12, 'WD Black SN850X 2TB', 'Gaming NVMe SSD with high sequential speeds.', NULL, 199.90, 20),
(18, 5, 9, 'Kingston KC3000 1TB', 'PCIe Gen4 NVMe SSD.', NULL, 109.90, 25),

-- HDD
(19, 6, 13, 'Seagate Barracuda 2TB', '7200 RPM desktop hard drive.', NULL, 69.90, 30),
(20, 6, 12, 'WD Blue 4TB', 'Reliable storage for desktop PCs.', NULL, 94.90, 18),

-- Power Supplies
(21, 7, 8, 'Corsair RM850x 850W Gold', '80 Plus Gold fully modular PSU.', NULL, 159.90, 15),
(22, 7, 16, 'be quiet! Pure Power 12M 750W', 'ATX 3.0 modular power supply.', NULL, 139.90, 12),
(23, 7, 14, 'Cooler Master MWE Gold 850 V2', '850W fully modular power supply.', NULL, 149.90, 10),

-- Cases
(24, 8, 20, 'NZXT H7 Flow', 'Mid-tower case optimized for airflow.', NULL, 139.90, 10),
(25, 8, 17, 'Thermaltake View 270 TG', 'Tempered glass gaming case.', NULL, 119.90, 9),

-- CPU Coolers
(26, 9, 15, 'Noctua NH-D15 chromax.black', 'Premium dual-tower air cooler.', NULL, 109.90, 15),
(27, 9, 19, 'DeepCool AK620', 'High-performance dual-tower CPU cooler.', NULL, 69.90, 18),
(28, 9, 14, 'Cooler Master MasterLiquid 360L', '360mm all-in-one liquid cooler.', NULL, 129.90, 12),

-- Case Fans
(29, 10, 15, 'Noctua NF-A12x25 PWM', '120mm premium PWM cooling fan.', NULL, 34.90, 40),
(30, 10, 8, 'Corsair iCUE AF120 RGB Elite', '120mm RGB PWM case fan.', NULL, 29.90, 35);

-- +goose Down
DELETE FROM products;
