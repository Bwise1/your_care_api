-- Migration script to enhance appointment system
-- Run this script to add new fields and tables for improved appointment management

-- Add new columns to appointments table
ALTER TABLE appointments 
ADD COLUMN admin_notes TEXT AFTER status,
ADD COLUMN user_notes TEXT AFTER admin_notes,
ADD COLUMN rejection_reason TEXT AFTER user_notes,
ADD COLUMN provider_id INT AFTER lab_test_id;

-- Update appointment status enum to include new statuses
ALTER TABLE appointments 
MODIFY COLUMN status ENUM(
    'pending',
    'admin_review', 
    'confirmed',
    'scheduled',
    'reschedule_offered',
    'reschedule_accepted',
    'in_progress',
    'completed',
    'canceled',
    'rejected',
    'no_show'
) DEFAULT 'pending';

-- Update appointment_type enum to include IVF
ALTER TABLE appointments 
MODIFY COLUMN appointment_type ENUM('doctor', 'lab_test', 'ivf') NOT NULL;

-- Create reschedule offers table
CREATE TABLE reschedule_offers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT NOT NULL,
    proposed_date DATE NOT NULL,
    proposed_time TIME NOT NULL,
    admin_notes TEXT,
    status ENUM('pending', 'accepted', 'rejected') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
);

-- Create appointment status history table for tracking status changes
CREATE TABLE appointment_status_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    notes TEXT,
    changed_by_user_id INT, -- Who made the change (admin or user)
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by_user_id) REFERENCES users(id)
);

-- Create IVF appointment details table (for future use)
CREATE TABLE ivf_appointment_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    appointment_id INT NOT NULL,
    treatment_type VARCHAR(100),
    cycle_day INT,
    special_instructions TEXT,
    preparation_notes TEXT,
    FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_appointments_status ON appointments(status);
CREATE INDEX idx_appointments_type ON appointments(appointment_type);
CREATE INDEX idx_appointments_date ON appointments(appointment_date);
CREATE INDEX idx_appointments_user_id ON appointments(user_id);
CREATE INDEX idx_appointments_provider_id ON appointments(provider_id);
CREATE INDEX idx_reschedule_offers_appointment_id ON reschedule_offers(appointment_id);
CREATE INDEX idx_status_history_appointment_id ON appointment_status_history(appointment_id);

-- Insert initial status history for existing appointments
INSERT INTO appointment_status_history (appointment_id, status, notes, changed_at)
SELECT id, status, 'Initial status from migration', created_at
FROM appointments;