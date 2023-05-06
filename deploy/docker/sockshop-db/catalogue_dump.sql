CREATE USER IF NOT EXISTS 'sockshop' IDENTIFIED BY 'password';

GRANT ALL ON socksdb.* TO 'sockshop';

CREATE TABLE IF NOT EXISTS sock (
	id varchar(40) NOT NULL, 
	name varchar(20), 
	description varchar(200), 
	price float, 
	count int, 
	image_urls varchar(100),  
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS tag (
	id MEDIUMINT NOT NULL AUTO_INCREMENT, 
	name varchar(20), 
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS sock_tag (
	sock_id varchar(40), 
	tag_id MEDIUMINT NOT NULL, 
	FOREIGN KEY (sock_id) 
		REFERENCES sock(id), 
	FOREIGN KEY(tag_id)
		REFERENCES tag(id)
);

INSERT INTO sock VALUES ("6d62d909-f957-430e-8689-b5129c0bb75e", "Weave special", "Limited issue Weave socks.", 17.15, 33, "/catalogue/images/weave1.jpg,/catalogue/images/weave2.jpg");
INSERT INTO sock VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "Nerd leg", "For all those leg lovers out there. A perfect example of a swivel chair trained calf. Meticulously trained on a diet of sitting and Pina Coladas. Phwarr...", 7.99, 115, "/catalogue/images/bit_of_leg_1.jpeg,/catalogue/images/bit_of_leg_2.jpeg");
INSERT INTO sock VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "Crossed", "A mature sock, crossed, with an air of nonchalance.",  17.32, 738, "/catalogue/images/cross_1.jpeg,/catalogue/images/cross_2.jpeg");
INSERT INTO sock VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "SuperSport XL", "Ready for action. Engineers: be ready to smash that next bug! Be ready, with these super-action-sport-masterpieces. This particular engineer was chased away from the office with a stick.",  15.00, 820, "/catalogue/images/puma_1.jpeg,/catalogue/images/puma_2.jpeg");
INSERT INTO sock VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "Holy", "Socks fit for a Messiah. You too can experience walking in water with these special edition beauties. Each hole is lovingly proggled to leave smooth edges. The only sock approved by a higher power.",  99.99, 1, "/catalogue/images/holy_1.jpeg,/catalogue/images/holy_2.jpeg");
INSERT INTO sock VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "YouTube.sock", "We were not paid to sell this sock. It's just a bit geeky.",  10.99, 801, "/catalogue/images/youtube_1.jpeg,/catalogue/images/youtube_2.jpeg");
INSERT INTO sock VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "Figueroa", "enim officia aliqua excepteur esse deserunt quis aliquip nostrud anim",  14, 808, "/catalogue/images/WAT.jpg,/catalogue/images/WAT2.jpg");
INSERT INTO sock VALUES ("647d653f-14ed-4391-bc00-ed6ee2687b51", "Classic", "Keep it simple.",  12, 127, "/catalogue/images/classic.jpg,/catalogue/images/classic2.jpg");
INSERT INTO sock VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "Colourful", "proident occaecat irure et excepteur labore minim nisi amet irure",  18, 438, "/catalogue/images/colourful_socks.jpg,/catalogue/images/colourful_socks.jpg");
INSERT INTO sock VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "Cat socks", "consequat amet cupidatat minim laborum tempor elit ex consequat in",  15, 175, "/catalogue/images/catsocks.jpg,/catalogue/images/catsocks2.jpg");

INSERT INTO tag (name) VALUES ("brown");
INSERT INTO tag (name) VALUES ("geek");
INSERT INTO tag (name) VALUES ("formal");
INSERT INTO tag (name) VALUES ("blue");
INSERT INTO tag (name) VALUES ("skin");
INSERT INTO tag (name) VALUES ("red");
INSERT INTO tag (name) VALUES ("action");
INSERT INTO tag (name) VALUES ("sport");
INSERT INTO tag (name) VALUES ("black");
INSERT INTO tag (name) VALUES ("magic");
INSERT INTO tag (name) VALUES ("green");

INSERT INTO sock_tag VALUES ("6d62d909-f957-430e-8689-b5129c0bb75e", "2");
INSERT INTO sock_tag VALUES ("6d62d909-f957-430e-8689-b5129c0bb75e", "9");
INSERT INTO sock_tag VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "4");
INSERT INTO sock_tag VALUES ("a0a4f044-b040-410d-8ead-4de0446aec7e", "5");
INSERT INTO sock_tag VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "4");
INSERT INTO sock_tag VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "6");
INSERT INTO sock_tag VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "7");
INSERT INTO sock_tag VALUES ("808a2de1-1aaa-4c25-a9b9-6612e8f29a38", "3");
INSERT INTO sock_tag VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "8");
INSERT INTO sock_tag VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "9");
INSERT INTO sock_tag VALUES ("510a0d7e-8e83-4193-b483-e27e09ddc34d", "3");
INSERT INTO sock_tag VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "10");
INSERT INTO sock_tag VALUES ("03fef6ac-1896-4ce8-bd69-b798f85c6e0b", "7");
INSERT INTO sock_tag VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "2");
INSERT INTO sock_tag VALUES ("d3588630-ad8e-49df-bbd7-3167f7efb246", "3");
INSERT INTO sock_tag VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "3");
INSERT INTO sock_tag VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "11");
INSERT INTO sock_tag VALUES ("819e1fbf-8b7e-4f6d-811f-693534916a8b", "4");
INSERT INTO sock_tag VALUES ("647d653f-14ed-4391-bc00-ed6ee2687b51", "1");
INSERT INTO sock_tag VALUES ("647d653f-14ed-4391-bc00-ed6ee2687b51", "11");
INSERT INTO sock_tag VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "1");
INSERT INTO sock_tag VALUES ("3395a43e-2d88-40de-b95f-e00e1502085b", "4");
INSERT INTO sock_tag VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "1");
INSERT INTO sock_tag VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "11");
INSERT INTO sock_tag VALUES ("837ab141-399e-4c1f-9abc-bace40296bac", "3");


CREATE TABLE IF NOT EXISTS customer (
	id varchar(40) NOT NULL, 
	first_name varchar(20), 
	last_name varchar(20), 
	email varchar(40), 
	username varchar(20), 
	password varchar(40), 
	salt varchar(40),
	PRIMARY KEY(id),
	UNIQUE (email)
);

CREATE UNIQUE INDEX customer_username_uq ON customer (username);

CREATE TABLE IF NOT EXISTS address (
	id varchar(40) NOT NULL, 
	street varchar(40), 
	number varchar(40), 
	country varchar(20), 
	city varchar(20), 
	postcode varchar(20), 
	PRIMARY KEY(id)
);


CREATE TABLE IF NOT EXISTS card (
	id varchar(40) NOT NULL, 
	long_num varchar(60), 
	expires varchar(20), 
	ccv varchar(20), 
	PRIMARY KEY(id),
	UNIQUE (long_num)
);

CREATE TABLE IF NOT EXISTS customer_address (
	customer_id varchar(40), 
	address_id varchar(40), 
	FOREIGN KEY (customer_id) 
		REFERENCES customer(id), 
	FOREIGN KEY(address_id)
		REFERENCES address(id)
);

CREATE INDEX customer_address_customer_id ON customer_address (customer_id);

CREATE TABLE IF NOT EXISTS customer_card (
	customer_id varchar(40), 
	card_id varchar(40), 
	FOREIGN KEY (customer_id) 
		REFERENCES customer(id), 
	FOREIGN KEY(card_id)
		REFERENCES card(id)
);

CREATE INDEX customer_card_customer_id ON customer_card (customer_id);