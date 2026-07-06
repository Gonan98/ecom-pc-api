-- +goose Up
INSERT INTO roles (name, description) VALUES ('admin', 'User authorized for create, update, delete resources'), ('customer', 'User authorized for see and buy products');

INSERT INTO brands (id, name, website) VALUES
(1, 'AMD', 'https://www.amd.com'),
(2, 'Intel', 'https://www.intel.com'),
(3, 'NVIDIA', 'https://www.nvidia.com'),
(4, 'ASUS', 'https://www.asus.com'),
(5, 'MSI', 'https://www.msi.com'),
(6, 'Gigabyte', 'https://www.gigabyte.com'),
(7, 'ASRock', 'https://www.asrock.com'),
(8, 'Corsair', 'https://www.corsair.com'),
(9, 'Kingston', 'https://www.kingston.com'),
(10, 'G.Skill', 'https://www.gskill.com'),
(11, 'Samsung', 'https://www.samsung.com'),
(12, 'Western Digital', 'https://www.westerndigital.com'),
(13, 'Seagate', 'https://www.seagate.com'),
(14, 'Cooler Master', 'https://www.coolermaster.com'),
(15, 'Noctua', 'https://www.noctua.at'),
(16, 'BeQuiet', 'https://www.bequiet.com'),
(17, 'Thermaltake', 'https://www.thermaltake.com'),
(18, 'EVGA', 'https://www.evga.com'),
(19, 'DeepCool', 'https://www.deepcool.com'),
(20, 'NZXT', 'https://nzxt.com');

INSERT INTO categories (id, name, description) VALUES
(1, 'Processors', 'Desktop CPUs from AMD and Intel.'),
(2, 'Graphics Cards', 'Dedicated GPUs for gaming and professional workloads.'),
(3, 'Motherboards', 'ATX, Micro-ATX and Mini-ITX motherboards.'),
(4, 'RAM', 'DDR4 and DDR5 memory modules.'),
(5, 'Solid State Drives (SSD)', 'SATA and NVMe solid-state storage devices.'),
(6, 'Hard Disk Drives (HDD)', 'Mechanical hard drives for mass storage.'),
(7, 'Power Supplies', 'ATX power supply units with various efficiency ratings.'),
(8, 'PC Cases', 'Computer chassis for desktop builds.'),
(9, 'CPU Coolers', 'Air and liquid cooling solutions for processors.'),
(10, 'Case Fans', 'Cooling fans for improving airflow inside the case.'),
(11, 'Thermal Paste', 'Thermal compounds for CPU and GPU cooling.'),
(12, 'Monitors', 'Gaming and productivity monitors.'),
(13, 'Keyboards', 'Mechanical and membrane keyboards.'),
(14, 'Mice', 'Gaming and productivity mice.'),
(15, 'Headsets', 'Gaming headsets with microphone.'),
(16, 'Speakers', 'Desktop speakers and sound systems.'),
(17, 'Networking', 'Network adapters, Wi-Fi cards and routers.'),
(18, 'Accessories', 'Cables, adapters and other PC accessories.'),
(19, 'Storage Accessories', 'External enclosures, heatsinks and storage accessories.'),
(20, 'Operating Systems', 'Operating system licenses and installation media.');

-- +goose Down
DELETE FROM roles;
DELETE FROM brands;
DELETE FROM categories;
