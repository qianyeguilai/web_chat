--create  tables   in  database  web_chat_go
CREATE  TABLE  IF NOT  EXISTS  `user_base_info` (
   `user_id`  SMALLINT  NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
   `cname`    VARCHAR(255)   NOT  NULL,
   `password` VARCHAR(255)    DEFAULT  'hello1234'
) ENGINE=InnoDB  DEFAULT  CHARSET=utf8;