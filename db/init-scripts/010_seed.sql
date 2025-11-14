SET datestyle TO 'ISO, DMY';

\copy stan (id, stan_kod, stan_name, stan_priznak, paragraph) FROM '/tmp/data/stations.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

\copy tip_uch (id, kod, name, tip_corr, date_start, date_end, upd_time) FROM '/tmp/data/stations02.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

\copy uch (id, uch_num, uch_name, uch_tip) FROM '/tmp/data/areas.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

\copy tr_punkt (id, tr, updated_at, tip_corr, date_start) FROM '/tmp/data/stations16.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

\copy kom_oper (id, place_code, opr_code, opr_name, tip_corr, date_start, date_end, updated_at) FROM '/tmp/data/stations08.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

\copy uch_uzel (id, uch, stan1, stan2, stan3, updated_at, tip_corr, date_start) FROM '/tmp/data/stations03.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8', NULL '0');

\copy tarif_stan (id, uch, stan, stan1, dist1, stan2, dist2, stan3, dist3, updated_at, tip_corr, date_start) FROM '/tmp/data/stations04.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8', NULL '0');

\copy stan_zd (id, stan, tr1, dist1, tr2, dist2, tr3, dist3, tr4, dist4, tip_corr, date_start, updated_at) FROM '/tmp/data/stations09.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8', NULL '0');

\copy tarif_tr (id, tr_start, tr_end, dist_tr, updated_at, tip_corr, date_start) FROM '/tmp/data/stations17.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF-8');

