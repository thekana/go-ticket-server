INSERT INTO users(username) values
    ('john'),('jessy'),('jay');
INSERT INTO roles(role) values
    ('admin'),('organizer'),('customer');
INSERT INTO user_roles(user_id,role_id) values
    (1,1),(2,2),(3,3),(1,3);