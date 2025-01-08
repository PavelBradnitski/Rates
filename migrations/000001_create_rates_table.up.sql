CREATE TABLE rates (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cur_id INT NOT NULL,
    date DATE NOT NULL,
    cur_abbreviation VARCHAR(10) NOT NULL,
    cur_scale INT NOT NULL,
    cur_official_rate FLOAT NOT NULL
);