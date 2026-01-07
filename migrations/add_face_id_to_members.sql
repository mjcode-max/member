-- 添加 face_id 字段到 members 表
-- 如果字段已存在，此SQL会报错，可以忽略

ALTER TABLE members ADD COLUMN face_id VARCHAR(255) DEFAULT NULL COMMENT '华为云人脸ID';

