DELIMITER //

-- Drop the procedures if they exist
DROP PROCEDURE IF EXISTS messages_create//
DROP PROCEDURE IF EXISTS messages_read_by_id//
DROP PROCEDURE IF EXISTS messages_read_all//
DROP PROCEDURE IF EXISTS messages_update//
DROP PROCEDURE IF EXISTS messages_delete//

-- Create
CREATE PROCEDURE messages_create(_content varchar(300))
BEGIN
  INSERT INTO messages (content) VALUES (_content);
  -- SELECT LAST_INSERT_ID() as id; Return only inserted id
  SELECT id, content, createdAt FROM messages WHERE id = LAST_INSERT_ID() LIMIT 1;
END //

-- Get By ID
CREATE PROCEDURE messages_read_by_id(_id int)
BEGIN
  SELECT id, content, createdAt FROM messages WHERE id = _id;
END //

-- Get All
CREATE PROCEDURE messages_read_all()
BEGIN
  SELECT * FROM messages ORDER BY id;
END //

-- Update
CREATE PROCEDURE messages_update(_id int, _content varchar(300))
BEGIN
  UPDATE messages SET content = _content WHERE id = _id;
END //

-- Delete
CREATE PROCEDURE messages_delete(_id int)
BEGIN
  DELETE FROM messages WHERE id = _id;
END //

DELIMITER ;


-- CALL messages_create('starting docker container');