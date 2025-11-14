-- =========================
-- СЖА 1102: Типы участков
-- =========================
CREATE TABLE IF NOT EXISTS tip_uch (
  id          integer      PRIMARY KEY,                 -- Идентификатор
  kod         integer      NOT NULL,                    -- Код
  name        varchar(255) NOT NULL,                    -- Наименование
  tip_corr    varchar(50)  NOT NULL,                    -- Тип коррекции
  date_start  date         NOT NULL,                    -- Дата начала действия
  date_end    date         NOT NULL,                    -- Дата конца действия
  upd_time    timestamp    NOT NULL,                    -- Время обновления
  CONSTRAINT tip_uch_date_chk CHECK (date_end >= date_start)
);

-- =========================
-- ДОП 1: Участки
-- =========================
CREATE TABLE IF NOT EXISTS uch (
  id        integer      PRIMARY KEY,       -- Идентификатор
  uch_num   integer      NOT NULL,          -- Номер участка
  uch_name  varchar(255)  NOT NULL,          -- Наименование
  uch_tip   integer      NOT NULL,          -- Тип участка -> tip_uch.id
  CONSTRAINT uch_num_uk UNIQUE (uch_num),
  CONSTRAINT uch_tip_fk FOREIGN KEY (uch_tip)
    REFERENCES tip_uch(id) ON UPDATE CASCADE ON DELETE RESTRICT
);

-- =========================
-- ДОП 2: Станции (справочник)
-- =========================
CREATE TABLE IF NOT EXISTS stan (
  id         integer      PRIMARY KEY,      -- Идентификатор
  stan_kod   integer      NOT NULL,         -- Код станции
  stan_name  varchar(50)  NOT NULL,         -- Наименование
  stan_priznak integer NOT NULL,
  paragraph  varchar(50),         -- Параграф
  CONSTRAINT stan_kod_uk UNIQUE (stan_kod)
);

-- ============================================
-- СЖА 1116: Транзитные пункты администраций
-- ============================================
CREATE TABLE IF NOT EXISTS tr_punkt (
  id         integer     PRIMARY KEY,       -- Идентификатор
  tr         integer,                        -- Код транзитного пункта
  updated_at timestamp,                      -- Время обновления
  tip_corr   varchar(50),                    -- Тип коррекции
  date_start date,                             -- Дата начала действия
  CONSTRAINT tr_punkt_tr_fk   FOREIGN KEY (tr) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE UNIQUE INDEX tr_punkt_tr_uk ON tr_punkt(tr);

-- ==================================================
-- СЖА 1108: Условные обозначения коммерч. операций
--  NOTE: расширяем длину OPR_CODE до 50, чтобы
--  согласовать FK-совместимость с СЖА 1109.
-- ==================================================
CREATE TABLE IF NOT EXISTS kom_oper (
  id         integer       PRIMARY KEY,      -- Идентификатор
  place_code integer,                        -- Код места проведения операций
  opr_code   varchar(50)   NOT NULL,         -- Условные обозначения операций
  opr_name   varchar(255)  NOT NULL,         -- Наименования операций
  tip_corr   varchar(50),                    -- Тип коррекции
  date_start date,                           -- Дата начала действия
  date_end   date,                           -- Дата конца действия
  updated_at timestamp,                      -- Время обновления
  CONSTRAINT kom_oper_dates_chk CHECK (date_end IS NULL OR date_start IS NULL OR date_end >= date_start)
);

-- =====================================================
-- СЖА 1103: Участки и узловые станции железных дорог
-- =====================================================
CREATE TABLE IF NOT EXISTS uch_uzel (
  id         integer     PRIMARY KEY,        -- Идентификатор
  uch        integer,                        -- Номер участка -> uch.uch_num
  stan1      integer,                        -- Код 1-й узловой станции -> stan.stan_kod
  stan2      integer,                        -- Код 2-й узловой станции -> stan.stan_kod
  stan3      integer,                        -- Код 3-й узловой станции -> stan.stan_kod
  updated_at timestamp,                      -- Время обновления
  tip_corr   varchar(50),                    -- Тип коррекции
  date_start date,                           -- Дата начала действия
  CONSTRAINT uch_uzel_uch_fk  FOREIGN KEY (uch)   REFERENCES uch(uch_num)   ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT uch_uzel_s1_fk   FOREIGN KEY (stan1) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT uch_uzel_s2_fk   FOREIGN KEY (stan2) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT uch_uzel_s3_fk   FOREIGN KEY (stan3) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL
);

-- ======================================================================
-- СЖА 1104: Тарифные расстояния между станциями участков железных дорог
-- ======================================================================
CREATE TABLE IF NOT EXISTS tarif_stan (
  id         integer      PRIMARY KEY,        -- Идентификатор
  uch        integer,                         -- Номер участка -> uch.uch_num
  stan       integer,                         -- Код станции -> stan.stan_kod
  stan1      integer,                         -- Код 1-го узла -> stan.stan_kod
  dist1      numeric(10,3),                   -- Расстояние до 1-го узла (км)
  stan2      integer,
  dist2      numeric(10,3),
  stan3      integer,
  dist3      numeric(10,3),
  updated_at timestamp,                       -- Время обновления
  tip_corr   varchar(50),                     -- Тип коррекции
  date_start date,                            -- Дата начала действия
  CONSTRAINT tarif_stan_uch_fk  FOREIGN KEY (uch)   REFERENCES uch(uch_num)   ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT tarif_stan_s_fk    FOREIGN KEY (stan)  REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT tarif_stan_s1_fk   FOREIGN KEY (stan1) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT tarif_stan_s2_fk   FOREIGN KEY (stan2) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT tarif_stan_s3_fk   FOREIGN KEY (stan3) REFERENCES stan(stan_kod) ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX tarif_stan_uch_idx  ON tarif_stan(uch);
CREATE INDEX tarif_stan_stan_idx ON tarif_stan(stan);

-- ======================================================================================
-- СЖА 1109: Железнодорожные станции (раздельные и пассажирские остановочные пункты)
-- ======================================================================================
CREATE TABLE IF NOT EXISTS stan_zd (
  id         integer      PRIMARY KEY,        -- Идентификатор
  stan       integer,                         -- Код станции -> stan.stan_kod
  tr1        integer,                          -- Код 1-го транзитного пункта -> tr_punkt.tr
  dist1      numeric(10,3),                    -- Расстояние до 1-го ТП (км)
  tr2        integer,
  dist2      numeric(10,3),
  tr3        integer,
  dist3      numeric(10,3),
  tr4        integer,
  dist4      numeric(10,3),
  tip_corr   varchar(50),                      -- Тип коррекции
  date_start date,                             -- Дата начала действия
  updated_at timestamp,                        -- Время обновления
  CONSTRAINT stan_zd_stan_fk  FOREIGN KEY (stan)     REFERENCES stan(stan_kod)    ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT stan_zd_tr1_fk   FOREIGN KEY (tr1)      REFERENCES tr_punkt(tr)      ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT stan_zd_tr2_fk   FOREIGN KEY (tr2)      REFERENCES tr_punkt(tr)      ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT stan_zd_tr3_fk   FOREIGN KEY (tr3)      REFERENCES tr_punkt(tr)      ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT stan_zd_tr4_fk   FOREIGN KEY (tr4)      REFERENCES tr_punkt(tr)      ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX stan_zd_stan_idx ON stan_zd(stan);

-- ===========================================================
-- СЖА 1117: Тарифные расстояния между транзитными пунктами
-- ===========================================================
CREATE TABLE IF NOT EXISTS tarif_tr (
  id         integer      PRIMARY KEY,        -- Идентификатор
  tr_start   integer,                         -- Код начального ТП -> tr_punkt.tr
  tr_end     integer,                         -- Код конечного ТП  -> tr_punkt.tr
  dist_tr    numeric(10,3),                   -- Расстояние (км)
  updated_at timestamp,                       -- Время обновления
  tip_corr   varchar(50),                     -- Тип коррекции
  date_start date,                            -- Дата начала действия
  CONSTRAINT tarif_tr_s_fk FOREIGN KEY (tr_start) REFERENCES tr_punkt(tr) ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT tarif_tr_e_fk FOREIGN KEY (tr_end)   REFERENCES tr_punkt(tr) ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE INDEX tarif_tr_pair_idx ON tarif_tr(tr_start, tr_end);
