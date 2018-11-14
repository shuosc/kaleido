INSERT INTO public.jsonindexedmirrorstation (id, name, url, indexurl)
VALUES (1, '上海大学', 'https://mirrors.shu.edu.cn', '/data/mirrors.json');
INSERT INTO public.jsonindexedmirrorstation (id, name, url, indexurl)
VALUES (2, '清华大学', 'https://mirrors.tuna.tsinghua.edu.cn', '/static/tunasync.json');
INSERT INTO public.jsonindexedmirrorstation (id, name, url, indexurl)
VALUES (11, '阿里巴巴', 'https://mirrors.aliyun.com', 'https://opsx.alibaba.com/api/v1/repo?_input_charset=utf-8&lang=chs');

INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (3, '中国科学技术大学', 'https://mirrors.ustc.edu.cn', '.filelist td.filename a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (4, '上海交通大学', 'https://ftp.sjtu.edu.cn', 'a:not(a:first-of-type)');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (5, '上海科技大学', 'https://mirrors.geekpie.club', '.indexcolicon a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (6, '浙江大学', 'https://mirrors.zju.edu.cn', '.zju-mirrors-body dl a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (7, '兰州大学', 'https://mirror.lzu.edu.cn', '.mirror_item a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (8, '华中科技大学', 'http://mirrors.hust.edu.cn', '#mirror-tbody tr td:first-of-type a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (9, '东软信息学院', 'https://mirrors.neusoft.edu.cn', '.container section .table-mirror tbody td:first-of-type a');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (10, '北京理工大学', 'http://mirror.bit.edu.cn', 'ul>li>a:first-of-type');
INSERT INTO public.webindexedmirrorstation (id, name, url, selector)
VALUES (12, '网易', 'http://mirrors.163.com', '#distro-table > tbody > tr td:first-of-type a');

INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (1, 10, 'download.html');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (2, 4, 'favicon.ico');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (3, 4, 'index2.html');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (4, 4, 'robots.txt');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (5, 4, 'sjtu.edu.cn.html');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (6, 4, 'test');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (7, 10, 'Downloads');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (8, 10, 'linux');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (9, 2, '",');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (10, 11, '系统');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (11, 11, '容器');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (12, 11, '存储');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (13, 11, '语言');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (14, 11, '其他');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (15, 11, '网络');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (16, 11, '内核');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (17, 11, '运维');
INSERT INTO public.mirrorignore (id, mirrorstation_id, name)
VALUES (18, 11, '驱动');