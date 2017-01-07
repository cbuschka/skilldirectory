// Create new superuser
CREATE USER IF NOT EXISTS star WITH PASSWORD 'wars' SUPERUSER;

// Use starwars superuser to "disable" default user w/ ridiculous password:
//ALTER USER cassandra WITH PASSWORD 'k/xT<SEUAUZ/kV4k>6pA[VPvGk}&X`)+' NOSUPERUSER;
