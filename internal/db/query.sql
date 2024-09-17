-- roles table
DROP TABLE IF EXISTS roles;
CREATE TABLE roles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE (name)
);

-- users table
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    firstName VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    lastName VARCHAR(50) COLLATE utf8mb4_unicode_ci NOT NULL,
    email VARCHAR(100) COLLATE utf8mb4_unicode_ci NOT NULL,
    dateOfBirth DATE NOT NULL,
    sex ENUM('Male', 'Female', 'Other') COLLATE utf8mb4_unicode_ci NOT NULL,
    height DECIMAL(5,2), -- Nullable, allowing for null values
    password VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    role_id INT NOT NULL DEFAULT 1, -- Foreign key to roles table
    isActive TINYINT(1) NULL DEFAULT 1, -- Nullable, with default value 1 (active)
    lastLogin DATETIME, -- Nullable
    refreshToken VARCHAR(255) COLLATE utf8mb4_unicode_ci, -- Nullable
    tokenExpiration DATETIME, -- Nullable
    isEmailVerified TINYINT(1) NULL DEFAULT 0, -- Nullable, default value 0 (not verified)
    emailVerificationToken VARCHAR(100) COLLATE utf8mb4_unicode_ci, -- Nullable
    emailVerificationTokenExpires DATETIME, -- Nullable
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY (email), -- Ensures no duplicate emails
    FOREIGN KEY (role_id) REFERENCES roles(id) -- Foreign key constraint
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


DROP TABLE IF EXISTS social_logins;
CREATE TABLE social_logins (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL, -- Foreign key referencing the users table
    provider ENUM('Google', 'Facebook', 'Twitter', 'Apple', 'GitHub') NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    provider_token TEXT,
    token_expires_at DATETIME,
    email VARCHAR(100),
    name VARCHAR(100),
    avatar_url VARCHAR(255),
    created_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Hospitals table
DROP TABLE IF EXISTS hospitals;
CREATE TABLE hospitals (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  address TEXT,
  phone VARCHAR(20),
  email VARCHAR(100)
);

-- Laboratory tests table
DROP TABLE IF EXISTS lab_tests;
CREATE TABLE lab_tests (
  id INT PRIMARY KEY AUTO_INCREMENT,
  hospital_id INT,
  name VARCHAR(100),
  description TEXT,
  price DECIMAL(10, 2),
  FOREIGN KEY (hospital_id) REFERENCES hospitals(id)
);


-- doctors table
DROP TABLE IF EXISTS doctors;
CREATE TABLE doctors (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hospital_id INT, -- Foreign key linking to the hospital table
    name VARCHAR(100) NOT NULL,
    specialization VARCHAR(100), -- Doctor's field of specialization
    email VARCHAR(100),
    phone VARCHAR(20),
    available_from TIME, -- Available start time for appointments
    available_to TIME, -- Available end time for appointments
    FOREIGN KEY (hospital_id) REFERENCES hospitals(id)
);

-- appointments table
DROP TABLE IF EXISTS appointments;
CREATE TABLE appointments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL, -- Foreign key linking to the user table
    doctor_id INT, -- Foreign key linking to the doctor table, NULL if lab test
    lab_test_id INT, -- Foreign key linking to the lab_test table, NULL if doctor appointment
    appointment_type ENUM('doctor', 'lab_test') NOT NULL, -- Type of appointment
    appointment_date DATE NOT NULL,
    appointment_time TIME NOT NULL,
    status ENUM('pending','scheduled', 'completed', 'canceled') DEFAULT 'pending', -- Status of the appointment
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (doctor_id) REFERENCES doctors(id),
    FOREIGN KEY (lab_test_id) REFERENCES lab_tests(id)
);

-- lab test appointment details table
DROP TABLE IF EXISTS lab_test_appointment_details;
CREATE TABLE lab_test_appointment_details (
    id INT PRIMARY KEY AUTO_INCREMENT,
    appointment_id INT NOT NULL,  -- Foreign key linking to the appointment table
    pickup_type ENUM('home', 'hospital') NOT NULL,  -- Type of pickup
    home_location TEXT, -- Address details for home pickup
    test_type_id INT NOT NULL,  -- Test type details
    hospital_id INT, -- hospital id for hospital type of tests
    additional_instructions TEXT,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id),
    FOREIGN KEY (test_type_id) REFERENCES lab_tests(id),
    FOREIGN KEY (hospital_id) REFERENCES hospitals(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- doctor appointment details table
DROP TABLE IF EXISTS doctor_appointment_details;
CREATE TABLE doctor_appointment_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT NOT NULL, -- Foreign key linking to the appointment table
    reason_for_visit TEXT, -- Reason for the appointment
    symptoms TEXT, -- Symptoms or other relevant patient information
    additional_notes TEXT, -- Any other additional notes for the doctor
    FOREIGN KEY (appointment_id) REFERENCES appointments(id)
);


-- Inserting data to hospital
INSERT INTO hospitals (name, address, phone, email)
VALUES (
    'Miracle Hospital',
    'Prof. Hilmi Forward Street, No: 24, NICOSIA',
    '0392 444 67 25',
    'info@wellcarelaborators.com'
);


-- Inserting data to doctors
INSERT INTO roles ( name, description) VALUES
( 'user', 'Regular user with standard privileges'),
('admin', 'Administrator with full system access'),
('doctor', 'Medical professional with access to patient data');

-- Assuming Miracle Hospital has ID 1
INSERT INTO lab_tests (hospital_id, name, description, price) VALUES
(1, 'Complete Blood Count (CBC)', 'Measures different components of blood including red and white blood cells, hemoglobin, and platelets.', 50.00),
(1, 'Lipid Panel', 'Measures cholesterol levels to assess risk of cardiovascular disease.', 65.00),
(1, 'Thyroid Function Test', 'Checks the function of the thyroid gland by measuring hormone levels.', 80.00),
(1, 'Urinalysis', 'Analyzes urine sample for various health indicators.', 30.00),
(1, 'Hemoglobin A1C', 'Measures average blood sugar levels over the past 2-3 months.', 70.00),
(1, 'Vitamin D Test', 'Measures the level of Vitamin D in the blood.', 90.00),
(1, 'Liver Function Test', 'Assesses the health and function of the liver.', 75.00),
(1, 'COVID-19 PCR Test', 'Detects genetic material of the SARS-CoV-2 virus.', 120.00);
