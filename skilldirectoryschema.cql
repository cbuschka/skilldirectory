CREATE KEYSPACE IF NOT EXISTS skill_directory_keyspace with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

CREATE TABLE IF NOT EXISTS skill_directory_keyspace.skills (
	name varchar,
	id uuid PRIMARY KEY,
	skilltype varchar
);

CREATE TABLE IF NOT EXISTS skill_directory_keyspace.links (
	id uuid,
	name varchar,
	url varchar,
	skill_id uuid,
	linktype varchar,
	PRIMARY KEY (id, skill_id)
);

CREATE INDEX IF NOT EXISTS ON skill_directory_keyspace.skills ( skilltype );
CREATE INDEX IF NOT EXISTS ON skill_directory_keyspace.links ( skill_id );
