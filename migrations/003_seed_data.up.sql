-- Insert tribun data
INSERT INTO tribun (nama) VALUES 
('Tribun Utara'),
('Tribun Selatan'),
('Tribun Timur'),
('Tribun Barat'),
('VIP Utara'),
('VIP Selatan');

-- Insert kursi for regular tribun (10 baris x 20 kursi = 200 per tribun)
INSERT INTO kursi (baris, nomor_kursi, tribun_id)
SELECT baris, nomor_kursi, t.id
FROM tribun t
CROSS JOIN generate_series(1, 10) AS baris
CROSS JOIN generate_series(1, 20) AS nomor_kursi
WHERE t.nama IN ('Tribun Utara', 'Tribun Selatan', 'Tribun Timur', 'Tribun Barat');

-- Insert kursi for VIP tribun (5 baris x 50 kursi = 250 per tribun)
INSERT INTO kursi (baris, nomor_kursi, tribun_id)
SELECT baris, nomor_kursi, t.id
FROM tribun t
CROSS JOIN generate_series(1, 5) AS baris
CROSS JOIN generate_series(1, 50) AS nomor_kursi
WHERE t.nama IN ('VIP Utara', 'VIP Selatan');