CREATE TABLE public.Mirror
(
  id   bigserial PRIMARY KEY NOT NULL,
  name varchar(128)          NOT NULL
);
CREATE UNIQUE INDEX Mirror_id_uindex
  ON public.Mirror (id);
CREATE UNIQUE INDEX Mirror_name_uindex
  ON public.Mirror (name);

CREATE TABLE public.Area
(
  id   bigserial PRIMARY KEY NOT NULL,
  name varchar(128)          NOT NULL
);
CREATE UNIQUE INDEX Area_name_uindex
  ON public.Area (name);
CREATE UNIQUE INDEX Area_id_uindex
  ON public.Area (id);

CREATE TABLE public.ISP
(
  id   bigserial PRIMARY KEY NOT NULL,
  name varchar(128)          NOT NULL
);
CREATE UNIQUE INDEX ISP_id_uindex
  ON public.ISP (id);
CREATE UNIQUE INDEX ISP_name_uindex
  ON public.ISP (name);

CREATE TABLE public.MirrorStation
(
  id   bigserial PRIMARY KEY NOT NULL,
  name varchar(128),
  url  text                  NOT NULL
);
CREATE UNIQUE INDEX MirrorStation_id_uindex
  ON public.MirrorStation (id);
CREATE UNIQUE INDEX MirrorStation_url_uindex
  ON public.MirrorStation (url);

CREATE TABLE public.WebIndexedMirrorStation
(
  selector text NOT NULL
)
  INHERITS (MirrorStation);

CREATE TABLE public.JsonIndexedMirrorStation
(
  indexUrl text NOT NULL
)
  INHERITS (MirrorStation);


CREATE TABLE public.MirrorStation_Mirror
(
  mirrorstation_id bigint NOT NULL,
  mirror_id        bigint NOT NULL,
  CONSTRAINT MirrorStation_Mirror_pk PRIMARY KEY (mirrorstation_id, mirror_id),
  CONSTRAINT MirrorStation_Mirror_mirrorstation_id_fk FOREIGN KEY (mirrorstation_id) REFERENCES public.mirrorstation (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT MirrorStation_Mirror_mirror_id_fk FOREIGN KEY (mirror_id) REFERENCES public.mirror (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE public.Area_Area
(
  from_id  bigint NOT NULL,
  to_id    bigint NOT NULL,
  distance int    NOT NULL,
  CONSTRAINT Area_Area_pk PRIMARY KEY (from_id, to_id)
);

CREATE TABLE public.IPRange
(
  id bigserial PRIMARY KEY NOT NULL,
  IP cidr                  NOT NULL,
  EXCLUDE USING SPGIST (IP with &&)
);
CREATE UNIQUE INDEX IPRange_id_uindex
  ON public.IPRange (id);
CREATE UNIQUE INDEX IPRange_IP_uindex
  ON public.IPRange (IP);

CREATE TABLE public.IPRange_Area_ISP
(
  IPRange_id bigint PRIMARY KEY NOT NULL,
  Area_id    bigint             NOT NULL,
  ISP_id     bigint             NOT NULL
);
CREATE UNIQUE INDEX IPRange_Area_ISP_IPRange_id_uindex
  ON public.IPRange_Area_ISP (IPRange_id);
