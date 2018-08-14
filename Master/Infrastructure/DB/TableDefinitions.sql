CREATE TABLE public.area
(
  id       serial PRIMARY KEY NOT NULL,
  province varchar(16)        NOT NULL
);
CREATE UNIQUE INDEX area_id_uindex
  ON public.area (id);
CREATE UNIQUE INDEX area_province_uindex
  ON public.area (province);

CREATE TABLE public.ip
(
  id   bigserial PRIMARY KEY NOT NULL,
  data cidr                  NOT NULL
);
CREATE UNIQUE INDEX ip_id_uindex
  ON public.ip (id);
CREATE UNIQUE INDEX ip_data_uindex
  ON public.ip (data);
CREATE INDEX ip_data_excl
  ON public.ip (data);

CREATE TABLE public.mirrorstation
(
  id  serial PRIMARY KEY NOT NULL,
  url varchar(128)       NOT NULL
);
CREATE UNIQUE INDEX mirrorstation_id_uindex
  ON public.mirrorstation (id);
CREATE UNIQUE INDEX mirrorstation_url_uindex
  ON public.mirrorstation (url);

CREATE TABLE public.area_mirrorstation
(
  area_id          integer PRIMARY KEY NOT NULL,
  mirrorstation_id integer             NOT NULL,
  CONSTRAINT area_mirrorstation_area_id_fk FOREIGN KEY (area_id) REFERENCES public.area (id),
  CONSTRAINT area_mirrorstation_mirrorstation_id_fk FOREIGN KEY (mirrorstation_id) REFERENCES public.mirrorstation (id) ON DELETE SET DEFAULT
);
CREATE UNIQUE INDEX area_mirrorstation_area_id_uindex
  ON public.area_mirrorstation (area_id);

CREATE TABLE public.ip_area
(
  ip_id   bigint PRIMARY KEY NOT NULL,
  area_id integer            NOT NULL,
  CONSTRAINT ip_area_ip_id_fk FOREIGN KEY (ip_id) REFERENCES public.ip (id) ON DELETE CASCADE,
  CONSTRAINT ip_area_area_id_fk FOREIGN KEY (area_id) REFERENCES public.area (id)
);
CREATE UNIQUE INDEX ip_area_ip_id_uindex
  ON public.ip_area (ip_id);