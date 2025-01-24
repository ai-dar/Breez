--
-- PostgreSQL database dump
--

-- Dumped from database version 16.2
-- Dumped by pg_dump version 16.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: likes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.likes (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    tweet_id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);


ALTER TABLE public.likes OWNER TO postgres;

--
-- Name: likes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.likes_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.likes_id_seq OWNER TO postgres;

--
-- Name: likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.likes_id_seq OWNED BY public.likes.id;


--
-- Name: tweets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.tweets (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    content text NOT NULL,
    user_id bigint NOT NULL
);


ALTER TABLE public.tweets OWNER TO postgres;

--
-- Name: tweets_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.tweets_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.tweets_id_seq OWNER TO postgres;

--
-- Name: tweets_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.tweets_id_seq OWNED BY public.tweets.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: breez
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name text,
    email text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    password text,
    role text DEFAULT 'user'::text,
    is_verified boolean DEFAULT false
);


ALTER TABLE public.users OWNER TO breez;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: breez
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO breez;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: breez
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: likes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.likes ALTER COLUMN id SET DEFAULT nextval('public.likes_id_seq'::regclass);


--
-- Name: tweets id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tweets ALTER COLUMN id SET DEFAULT nextval('public.tweets_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: breez
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: likes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.likes (id, user_id, tweet_id, created_at, updated_at, deleted_at) FROM stdin;
1	0	0	\N	\N	\N
2	0	9	\N	\N	\N
3	0	9	\N	\N	\N
4	0	9	\N	\N	\N
5	0	9	\N	\N	\N
6	0	9	\N	\N	\N
7	0	9	\N	\N	\N
8	0	8	\N	\N	\N
9	0	8	\N	\N	\N
10	0	8	\N	\N	\N
11	0	10	\N	\N	\N
12	0	10	\N	\N	\N
13	0	10	\N	\N	\N
14	0	10	\N	\N	\N
15	0	10	\N	\N	\N
16	0	10	\N	\N	\N
17	0	10	\N	\N	\N
18	0	9	\N	\N	\N
19	0	9	\N	\N	\N
20	0	9	\N	\N	\N
21	0	9	\N	\N	\N
22	0	9	\N	\N	\N
23	0	9	\N	\N	\N
24	0	10	\N	\N	\N
25	0	10	\N	\N	\N
26	0	10	\N	\N	\N
27	0	10	\N	\N	\N
28	0	10	\N	\N	\N
29	0	10	\N	\N	\N
30	0	10	\N	\N	\N
31	0	10	\N	\N	\N
32	0	10	\N	\N	\N
33	0	10	\N	\N	\N
34	0	10	\N	\N	\N
35	0	10	\N	\N	\N
36	0	9	\N	\N	\N
37	0	9	\N	\N	\N
38	0	9	\N	\N	\N
39	0	10	\N	\N	\N
40	0	9	\N	\N	\N
41	0	5	\N	\N	\N
42	15	0	2024-12-26 19:37:20.121305+05	2024-12-26 19:37:20.121305+05	2024-12-26 19:37:21.007012+05
43	15	0	2024-12-26 19:37:22.376401+05	2024-12-26 19:37:22.376401+05	2024-12-26 19:37:31.195772+05
44	15	0	2024-12-26 19:37:31.458278+05	2024-12-26 19:37:31.458278+05	2024-12-26 19:37:31.60786+05
45	15	0	2024-12-26 19:37:40.106898+05	2024-12-26 19:37:40.106898+05	2024-12-26 19:37:42.443008+05
46	15	0	2024-12-26 19:37:43.599867+05	2024-12-26 19:37:43.599867+05	2024-12-26 19:40:33.618748+05
47	16	4	2025-01-06 19:12:33.804419+05	2025-01-06 19:12:33.804419+05	2025-01-06 19:12:35.396376+05
48	16	4	2025-01-06 19:12:36.075223+05	2025-01-06 19:12:36.075223+05	2025-01-06 19:12:36.768393+05
49	16	4	2025-01-06 19:12:38.280751+05	2025-01-06 19:12:38.280751+05	2025-01-06 19:12:38.80012+05
50	16	4	2025-01-06 19:12:39.300204+05	2025-01-06 19:12:39.300204+05	\N
51	20	22	2025-01-07 18:41:13.893301+05	2025-01-07 18:41:13.893301+05	2025-01-07 18:41:14.732545+05
52	20	22	2025-01-07 18:41:16.600608+05	2025-01-07 18:41:16.600608+05	\N
53	21	23	2025-01-07 21:43:18.021054+05	2025-01-07 21:43:18.021054+05	\N
54	21	22	2025-01-07 21:43:21.085371+05	2025-01-07 21:43:21.085371+05	\N
55	21	5	2025-01-07 21:43:26.266622+05	2025-01-07 21:43:26.266622+05	2025-01-07 21:43:33.315185+05
56	21	5	2025-01-07 21:43:33.781137+05	2025-01-07 21:43:33.781137+05	2025-01-07 21:43:35.985711+05
57	21	5	2025-01-07 21:43:37.217789+05	2025-01-07 21:43:37.217789+05	2025-01-07 21:43:46.261426+05
58	21	5	2025-01-07 21:43:47.984132+05	2025-01-07 21:43:47.984132+05	\N
59	19	8	2025-01-07 21:46:36.320743+05	2025-01-07 21:46:36.320743+05	2025-01-07 21:46:37.237482+05
60	15	8	2025-01-07 21:46:46.983896+05	2025-01-07 21:46:46.983896+05	\N
61	19	23	2025-01-07 21:47:22.881669+05	2025-01-07 21:47:22.881669+05	\N
62	22	9	2025-01-08 00:28:59.765868+05	2025-01-08 00:28:59.765868+05	2025-01-08 00:29:00.172952+05
63	22	9	2025-01-08 00:29:02.352155+05	2025-01-08 00:29:02.352155+05	2025-01-08 00:29:03.503369+05
64	22	9	2025-01-08 00:29:08.717568+05	2025-01-08 00:29:08.717568+05	2025-01-08 00:29:24.771785+05
65	22	9	2025-01-08 00:29:27.297923+05	2025-01-08 00:29:27.297923+05	2025-01-08 00:29:29.412584+05
66	22	24	2025-01-08 00:30:31.667203+05	2025-01-08 00:30:31.667203+05	\N
67	19	19	2025-01-08 02:57:37.09667+05	2025-01-08 02:57:37.09667+05	2025-01-08 02:57:37.581695+05
69	19	25	2025-01-08 03:28:27.045933+05	2025-01-08 03:28:27.045933+05	2025-01-08 03:28:30.083122+05
68	22	25	2025-01-08 03:28:26.181575+05	2025-01-08 03:28:26.181575+05	2025-01-08 03:28:46.2559+05
70	22	25	2025-01-08 03:28:47.155628+05	2025-01-08 03:28:47.155628+05	2025-01-08 03:28:51.435162+05
71	19	25	2025-01-08 03:28:48.972284+05	2025-01-08 03:28:48.972284+05	2025-01-08 03:29:13.816732+05
73	22	25	2025-01-08 03:29:16.49787+05	2025-01-08 03:29:16.49787+05	\N
72	19	25	2025-01-08 03:29:14.632073+05	2025-01-08 03:29:14.632073+05	2025-01-08 03:29:17.637789+05
74	19	25	2025-01-08 03:29:22.728228+05	2025-01-08 03:29:22.728228+05	\N
75	19	8	2025-01-08 03:46:09.108407+05	2025-01-08 03:46:09.108407+05	\N
76	24	29	2025-01-08 16:01:43.578055+05	2025-01-08 16:01:43.578055+05	\N
77	24	28	2025-01-08 16:01:44.796557+05	2025-01-08 16:01:44.796557+05	\N
78	24	18	2025-01-08 16:01:51.325351+05	2025-01-08 16:01:51.325351+05	2025-01-08 16:01:54.061326+05
79	24	6	2025-01-08 16:05:15.584643+05	2025-01-08 16:05:15.584643+05	\N
80	19	7	2025-01-08 16:39:27.443587+05	2025-01-08 16:39:27.443587+05	\N
82	26	6	2025-01-17 14:48:40.634941+05	2025-01-17 14:48:40.634941+05	\N
83	20	30	2025-01-23 00:48:17.779019+05	2025-01-23 00:48:17.779019+05	2025-01-23 00:48:18.975567+05
81	19	4	2025-01-08 17:59:34.070849+05	2025-01-08 17:59:34.070849+05	2025-01-23 01:59:57.628451+05
84	19	4	2025-01-23 01:59:58.247809+05	2025-01-23 01:59:58.247809+05	2025-01-23 02:00:00.674228+05
\.


--
-- Data for Name: tweets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.tweets (id, created_at, updated_at, deleted_at, content, user_id) FROM stdin;
6	2024-12-26 19:09:18.07858+05	2024-12-26 19:09:18.07858+05	\N	asd	15
7	2024-12-26 19:09:45.277053+05	2024-12-26 19:09:45.277053+05	\N	He,,	15
8	2024-12-26 19:11:57.293767+05	2024-12-26 19:11:57.293767+05	\N	Hello	15
9	2024-12-26 19:12:10.22252+05	2024-12-26 19:12:10.22252+05	\N	sad	15
10	2024-12-26 19:32:44.724771+05	2024-12-26 19:32:44.724771+05	\N	A	15
11	2025-01-05 14:37:51.38351+05	2025-01-05 14:37:51.38351+05	\N	hbhb	15
12	2025-01-06 16:24:02.731081+05	2025-01-06 16:24:02.731081+05	\N	фыыы	15
13	2025-01-06 16:25:32.12034+05	2025-01-06 16:25:32.12034+05	\N	фы	16
14	2025-01-06 16:31:33.494492+05	2025-01-06 16:31:33.494492+05	\N	Как дела	16
15	2025-01-06 16:31:56.747166+05	2025-01-06 16:31:56.747166+05	\N	ааа	16
16	2025-01-06 16:33:39.44184+05	2025-01-06 16:33:39.44184+05	\N	Hello	16
17	2025-01-06 16:33:50.33876+05	2025-01-06 16:33:50.33876+05	\N	YES	16
18	2025-01-06 17:31:00.931263+05	2025-01-06 17:31:00.931263+05	\N	GO I LIKE BREEZ!!!!!	17
19	2025-01-06 17:31:10.549075+05	2025-01-06 17:31:10.549075+05	\N	GO	17
20	2025-01-06 17:44:31.991396+05	2025-01-06 17:44:31.991396+05	\N	LETSGOOO	17
21	2025-01-07 18:36:28.162624+05	2025-01-07 18:36:28.162624+05	\N	Hello ALL	19
22	2025-01-07 18:40:33.250487+05	2025-01-07 18:40:33.250487+05	\N	Hi	20
23	2025-01-07 21:43:06.396781+05	2025-01-07 21:43:06.396781+05	\N	LETSGO MU	21
24	2025-01-08 00:30:23.592427+05	2025-01-08 00:30:23.592427+05	\N	Hello All	22
25	2025-01-08 03:27:20.521879+05	2025-01-08 03:27:20.521879+05	\N	NewCastle 2:0 Arsenal ARSENAL BAD GAME LOLOLOLOL	22
26	2025-01-08 03:28:11.977538+05	2025-01-08 03:28:11.977538+05	\N	No TROLLING!!!	19
28	2025-01-08 15:54:32.674278+05	2025-01-08 15:54:32.674278+05	\N	I like Breez!	23
29	2025-01-08 16:01:17.085637+05	2025-01-08 16:01:17.085637+05	\N	Hello all	24
30	2025-01-23 00:48:08.459484+05	2025-01-23 00:48:08.459484+05	\N	Hello	20
31	2025-01-23 00:48:23.529983+05	2025-01-23 00:48:23.529983+05	\N	Hi	20
32	2025-01-23 00:55:46.334297+05	2025-01-23 00:55:46.334297+05	\N	Hello	20
33	2025-01-23 01:00:33.95262+05	2025-01-23 01:00:33.95262+05	\N	!!!	20
34	2025-01-23 01:55:13.791015+05	2025-01-23 01:55:13.791015+05	\N	Hello	19
4	2024-12-26 18:58:01.466726+05	2025-01-23 03:21:43.245202+05	\N	asdasd	15
5	2024-12-26 19:05:38.444888+05	2024-12-26 19:05:38.444888+05	2025-01-23 03:21:48.22344+05	Hello	15
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: breez
--

COPY public.users (id, name, email, created_at, updated_at, deleted_at, password, role, is_verified) FROM stdin;
14	Valid User	validuser@example.com	\N	\N	\N	$2a$10$.PwlXMGUR70Tw.1fXk1fZ.oLugiG1R3T1INY7xTX7ywyBr/lmA2Hm	user	f
15	Aidar	sabirgali05@mail.ru	\N	\N	\N	$2a$10$baTy4b.OzhOqLU05rWfbVuKkIKaglmBmFMjD9kO8mVIVEFow5KDE2	user	f
16	Sanzhar	sancho@mail.ru	\N	\N	\N	$2a$10$AvF60z0kqsi65Zo64a98SOyTf0XJH48TYWNinmNdhXn6xQvJP5TIm	user	f
17	Диас	dias@mail.ru	\N	\N	\N	$2a$10$IvHaNnzfo2jYFetB/lijRurep2Cx8FRSrRvFA0BI3YIB3XlVizIEy	user	f
20	Айдар	aidarsoulsov@mail.ru	\N	\N	\N	$2a$10$T5RVDgAJTTt5FK3.hjyqJeU4RTwbaHGG8fesGhFaedxt8Aq0pXEb6	user	f
21	Danik	bananamaninmoon@gmail.com	\N	\N	\N	$2a$10$VUrzrjUmTYBEkFXd.ZSSyuVpwn/23tbqp2yIq5rZUnD19EF6fI.ju	user	f
22	SanzharNurmukhanbet	n.sanzhar06@gmail.com	\N	\N	\N	$2a$10$RubzkXGiR1Xe20VZriPFB./Ilt8EMz9DNss2KFuF72EzUwO8m2E3.	user	f
23	Miras	zhumaseytov06@mail.ru	\N	\N	\N	$2a$10$0XN7NYfNRu6zYangN82cFul5ay.54q1QywRk.oolBZBKi2kMMfhGy	user	f
24	Test	test@mail.ru	\N	\N	\N	$2a$10$XyzlO2zzM1WuCYA3ShiRQ.ajcSLyjfJXL3.GWogJh6nYZGNBMJKF.	user	f
25	Inkar	inkar.usurbaeva@inbox.ru	\N	\N	\N	$2a$10$vF27E1H.9oaxVOOFRnWV4uyVhBxeogVPybQs.nd.k9o.clfaAcQ7a	user	f
26	Айдар	4engineersaitu@gmail.com	\N	\N	\N	$2a$10$A/c.lOtIS5wjL27cZJhyBOFqKTHGbBmRyYfVWwVwO/0n2wK2FgouC	user	t
28	Айдар	4engineersatu@gmail.com	\N	\N	\N	$2a$10$gZvEL7Ws1y41k4ue0bSpS.kuplZEWeet5saj8GqRswN6enl8TzJRi	user	f
19	Breez	admbreez@gmail.com	\N	\N	\N	$2a$10$Gaw5/TnukZenRWfA1KF1MOdToKCjY2/.0AkRP6Hpqx0W5PRxAkAI.	admin	t
\.


--
-- Name: likes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.likes_id_seq', 84, true);


--
-- Name: tweets_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.tweets_id_seq', 34, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: breez
--

SELECT pg_catalog.setval('public.users_id_seq', 29, true);


--
-- Name: likes likes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.likes
    ADD CONSTRAINT likes_pkey PRIMARY KEY (id);


--
-- Name: tweets tweets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tweets
    ADD CONSTRAINT tweets_pkey PRIMARY KEY (id);


--
-- Name: users unique_email; Type: CONSTRAINT; Schema: public; Owner: breez
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT unique_email UNIQUE (email);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: breez
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: breez
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_tweets_deleted_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_tweets_deleted_at ON public.tweets USING btree (deleted_at);


--
-- Name: idx_users_deleted_at; Type: INDEX; Schema: public; Owner: breez
--

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);


--
-- Name: tweets fk_tweets_user; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.tweets
    ADD CONSTRAINT fk_tweets_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO breez;


--
-- PostgreSQL database dump complete
--

