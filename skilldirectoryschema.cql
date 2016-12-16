CREATE KEYSPACE IF NOT EXISTS skills2 with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS skills2.skills (
	name varchar,
	id uuid PRIMARY KEY,
	skilltype varchar
);

CREATE TABLE IF NOT EXISTS skills2.links (
	id uuid,
	name varchar,
	url varchar,
	skill_id uuid,
	linktype varchar,
	PRIMARY KEY (id, skill_id)
);

CREATE INDEX IF NOT EXISTS ON skills2.skills ( skilltype );
CREATE INDEX IF NOT EXISTS ON skills2.links ( skill_id );
