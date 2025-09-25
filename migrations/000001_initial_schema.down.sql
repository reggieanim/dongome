-- Drop triggers first
DROP TRIGGER IF EXISTS update_listings_updated_at ON listings;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_seller_profiles_updated_at ON seller_profiles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_listings_search;
DROP INDEX IF EXISTS idx_listing_attributes_key;
DROP INDEX IF EXISTS idx_listing_attributes_listing_id;
DROP INDEX IF EXISTS idx_listing_images_listing_id;
DROP INDEX IF EXISTS idx_listings_promoted;
DROP INDEX IF EXISTS idx_listings_location;
DROP INDEX IF EXISTS idx_listings_price;
DROP INDEX IF EXISTS idx_listings_created_at;
DROP INDEX IF EXISTS idx_listings_status;
DROP INDEX IF EXISTS idx_listings_category_id;
DROP INDEX IF EXISTS idx_listings_seller_id;
DROP INDEX IF EXISTS idx_categories_is_active;
DROP INDEX IF EXISTS idx_categories_parent_id;
DROP INDEX IF EXISTS idx_seller_profiles_verification_status;
DROP INDEX IF EXISTS idx_seller_profiles_user_id;
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS listing_tag_relations;
DROP TABLE IF EXISTS listing_tags;
DROP TABLE IF EXISTS listing_attributes;
DROP TABLE IF EXISTS listing_images;
DROP TABLE IF EXISTS listings;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS seller_profiles;
DROP TABLE IF EXISTS users;