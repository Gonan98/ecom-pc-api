-- +goose Up
INSERT INTO products (
    category_id,
    brand_id,
    name,
    description,
    image_url,
    price,
    stock
) VALUES

-- Processors
(1, 1, 'AMD Ryzen 5 9600X', '6-core 12-thread desktop gaming processor.', NULL, 279.99, 25),
(1, 1, 'AMD Ryzen 7 9800X3D', 'High-end gaming processor with 3D V-Cache.', NULL, 519.99, 15),
(1, 2, 'Intel Core i5-14600KF', '14th Gen unlocked desktop processor.', NULL, 309.99, 20),
(1, 2, 'Intel Core i7-14700K', '20-core unlocked gaming processor.', NULL, 449.99, 12),

-- Graphics Cards
(2, 3, 'NVIDIA GeForce RTX 5070', 'Next-generation gaming graphics card.', NULL, 699.99, 10),
(2, 3, 'NVIDIA GeForce RTX 5080', 'High-performance graphics card for 4K gaming.', NULL, 1099.99, 6),
(2, 4, 'ASUS TUF Gaming Radeon RX 7800 XT', 'Factory-overclocked AMD graphics card.', NULL, 579.99, 8),
(2, 5, 'MSI GeForce RTX 4060 Ti Ventus', 'Efficient graphics card for 1080p gaming.', NULL, 449.99, 18),

-- Motherboards
(3, 4, 'ASUS ROG Strix B650-A Gaming WiFi', 'AM5 ATX gaming motherboard.', NULL, 259.99, 14),
(3, 5, 'MSI MAG B760 Tomahawk WiFi', 'Intel LGA1700 gaming motherboard.', NULL, 219.99, 16),
(3, 6, 'Gigabyte B650 Aorus Elite AX', 'AMD AM5 motherboard with WiFi 6.', NULL, 239.99, 12),
(3, 7, 'ASRock Z790 Steel Legend', 'Intel Z790 motherboard for enthusiasts.', NULL, 299.99, 8),

-- RAM
(4, 10, 'G.Skill Trident Z5 RGB 32GB DDR5-6000', '32GB DDR5 dual-channel memory kit.', NULL, 169.99, 20),
(4, 9, 'Kingston Fury Beast 32GB DDR5-5600', 'High-speed DDR5 gaming memory.', NULL, 149.99, 22),
(4, 8, 'Corsair Vengeance RGB 32GB DDR5-6400', 'Premium RGB DDR5 memory kit.', NULL, 189.99, 18),

-- SSD
(5, 11, 'Samsung 990 PRO 2TB NVMe SSD', 'PCIe 4.0 high-performance SSD.', NULL, 209.99, 18),
(5, 12, 'WD Black SN850X 2TB', 'Gaming NVMe SSD with high sequential speeds.', NULL, 199.99, 20),
(5, 9, 'Kingston KC3000 1TB', 'PCIe Gen4 NVMe SSD.', NULL, 109.99, 25),

-- HDD
(6, 13, 'Seagate Barracuda 2TB', '7200 RPM desktop hard drive.', NULL, 69.99, 30),
(6, 12, 'WD Blue 4TB', 'Reliable storage for desktop PCs.', NULL, 94.99, 18),

-- Power Supplies
(7, 8, 'Corsair RM850x 850W Gold', '80 Plus Gold fully modular PSU.', NULL, 159.99, 15),
(7, 16, 'be quiet! Pure Power 12M 750W', 'ATX 3.0 modular power supply.', NULL, 139.99, 12),
(7, 14, 'Cooler Master MWE Gold 850 V2', '850W fully modular power supply.', NULL, 149.99, 10),

-- Cases
(8, 20, 'NZXT H7 Flow', 'Mid-tower case optimized for airflow.', NULL, 139.99, 10),
(8, 17, 'Thermaltake View 270 TG', 'Tempered glass gaming case.', NULL, 119.99, 9),

-- CPU Coolers
(9, 15, 'Noctua NH-D15 chromax.black', 'Premium dual-tower air cooler.', NULL, 109.99, 15),
(9, 19, 'DeepCool AK620', 'High-performance dual-tower CPU cooler.', NULL, 69.99, 18),
(9, 14, 'Cooler Master MasterLiquid 360L', '360mm all-in-one liquid cooler.', NULL, 129.99, 12),

-- Case Fans
(10, 15, 'Noctua NF-A12x25 PWM', '120mm premium PWM cooling fan.', NULL, 34.99, 40),
(10, 8, 'Corsair iCUE AF120 RGB Elite', '120mm RGB PWM case fan.', NULL, 29.99, 35);

-- +goose Down
DELETE FROM products;
