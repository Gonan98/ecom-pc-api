-- +goose Up
INSERT INTO roles (name, description) VALUES ('admin', 'User authorized to create, update, delete resources'), ('customer', 'User authorized to consult and buy products');

INSERT INTO brands (name, website) VALUES
('AMD', 'https://www.amd.com'),
('Intel', 'https://www.intel.com'),
('NVIDIA', 'https://www.nvidia.com'),
('ASUS', 'https://www.asus.com'),
('MSI', 'https://www.msi.com'),
('Gigabyte', 'https://www.gigabyte.com'),
('ASRock', 'https://www.asrock.com'),
('Corsair', 'https://www.corsair.com'),
('Kingston', 'https://www.kingston.com'),
('G.Skill', 'https://www.gskill.com'),
('Samsung', 'https://www.samsung.com'),
('Western Digital', 'https://www.westerndigital.com'),
('Seagate', 'https://www.seagate.com'),
('Cooler Master', 'https://www.coolermaster.com'),
('Noctua', 'https://www.noctua.at'),
('BeQuiet', 'https://www.bequiet.com'),
('Thermaltake', 'https://www.thermaltake.com'),
('EVGA', 'https://www.evga.com'),
('DeepCool', 'https://www.deepcool.com'),
('NZXT', 'https://nzxt.com');

INSERT INTO categories (name, description) VALUES
('Processors', 'Desktop CPUs from AMD and Intel.'),
('Graphics Cards', 'Dedicated GPUs for gaming and professional workloads.'),
('Motherboards', 'ATX, Micro-ATX and Mini-ITX motherboards.'),
('RAM', 'DDR4 and DDR5 memory modules.'),
('Solid State Drives (SSD)', 'SATA and NVMe solid-state storage devices.'),
('Hard Disk Drives (HDD)', 'Mechanical hard drives for mass storage.'),
('Power Supplies', 'ATX power supply units with various efficiency ratings.'),
('PC Cases', 'Computer chassis for desktop builds.'),
('CPU Coolers', 'Air and liquid cooling solutions for processors.'),
('Case Fans', 'Cooling fans for improving airflow inside the case.'),
('Thermal Paste', 'Thermal compounds for CPU and GPU cooling.'),
('Monitors', 'Gaming and productivity monitors.'),
('Keyboards', 'Mechanical and membrane keyboards.'),
('Mice', 'Gaming and productivity mice.'),
('Headsets', 'Gaming headsets with microphone.'),
('Speakers', 'Desktop speakers and sound systems.'),
('Networking', 'Network adapters, Wi-Fi cards and routers.'),
('Accessories', 'Cables, adapters and other PC accessories.'),
('Storage Accessories', 'External enclosures, heatsinks and storage accessories.'),
('Operating Systems', 'Operating system licenses and installation media.');

-- +goose Down
DELETE FROM roles;
DELETE FROM brands;
DELETE FROM categories;
