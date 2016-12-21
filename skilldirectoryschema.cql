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
	PRIMARY KEY (skill_id)
);

CREATE TABLE IF NOT EXISTS skill_directory_keyspace.teammembers (
    name varchar,
    id uuid PRIMARY KEY,
    title varchar
);

CREATE INDEX IF NOT EXISTS ON skill_directory_keyspace.skills ( skilltype );
