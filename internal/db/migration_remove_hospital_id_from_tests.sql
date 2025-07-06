-- Migration to remove hospital_id from lab_tests table
-- This change is needed because hospital-specific test offerings are now handled
-- through the hospital_lab_tests junction table

-- First drop the foreign key constraint
ALTER TABLE lab_tests DROP FOREIGN KEY lab_tests_ibfk_1;

-- Then drop the column
ALTER TABLE lab_tests DROP COLUMN hospital_id;
